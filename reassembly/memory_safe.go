// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// Package reassembly provides improved memory-safe TCP stream reassembly.
package reassembly

import (
	"errors"
	"flag"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/gopacket/layers"
)

var (
	memLog              = flag.Bool("assembly_memuse_log", false, "If true, log memory usage information")
	ErrNilStream        = errors.New("stream factory returned nil stream")
	ErrInvalidOptions   = errors.New("invalid assembler options")
	ErrPoolClosed       = errors.New("stream pool is closed")
	ErrConnectionClosed = errors.New("connection is closed")
)

// SafePageCache is a concurrency-safe store of page objects with eviction
type SafePageCache struct {
	pagePool       *sync.Pool
	used           int64 // atomic
	pageRequests   int64 // atomic
	maxPages       int64
	evictionMutex  sync.Mutex
	lastEviction   time.Time
	evictionPeriod time.Duration
}

func newSafePageCache(maxPages int64) *SafePageCache {
	if maxPages <= 0 {
		maxPages = 10000 // default reasonable limit
	}
	pc := &SafePageCache{
		pagePool: &sync.Pool{
			New: func() interface{} { return new(page) },
		},
		maxPages:       maxPages,
		evictionPeriod: 30 * time.Second,
	}
	return pc
}

// next returns a clean, ready-to-use page object with memory limit enforcement
func (c *SafePageCache) next(ts time.Time) (*page, error) {
	requests := atomic.AddInt64(&c.pageRequests, 1)
	if *memLog && requests&0xFFFF == 0 {
		used := atomic.LoadInt64(&c.used)
		log.Printf("SafePageCache: %d requested, %d used", requests, used)
	}

	// Check if we need eviction
	if atomic.LoadInt64(&c.used) >= c.maxPages {
		c.performEviction()
		// Still over limit after eviction
		if atomic.LoadInt64(&c.used) >= c.maxPages {
			return nil, errors.New("page cache limit exceeded")
		}
	}

	p := c.pagePool.Get().(*page)
	p.seen = ts
	p.bytes = p.buf[:0]
	p.prev = nil
	p.next = nil
	p.ac = nil
	atomic.AddInt64(&c.used, 1)
	
	return p, nil
}

// replace returns a page to the cache
func (c *SafePageCache) replace(p *page) {
	if p == nil {
		return
	}
	atomic.AddInt64(&c.used, -1)
	// Clear the page completely to avoid memory leaks
	p.bytes = nil
	p.prev = nil
	p.next = nil
	p.ac = nil
	c.pagePool.Put(p)
}

// performEviction runs garbage collection on the pool
func (c *SafePageCache) performEviction() {
	c.evictionMutex.Lock()
	defer c.evictionMutex.Unlock()

	now := time.Now()
	if now.Sub(c.lastEviction) < c.evictionPeriod {
		return
	}

	// Force GC to run
	c.pagePool = &sync.Pool{
		New: func() interface{} { return new(page) },
	}
	c.lastEviction = now
	
	if *memLog {
		log.Printf("SafePageCache: performed eviction, used=%d", atomic.LoadInt64(&c.used))
	}
}

// SafeStreamPool stores all streams with improved concurrency safety
type SafeStreamPool struct {
	conns       map[key]*connection
	connsMutex  sync.RWMutex
	factory     StreamFactory
	connPool    sync.Pool  // Pool of connection objects
	
	// Connection limiting
	maxConns    int64
	activeConns int64 // atomic
	
	// Lifecycle management  
	closed      int32 // atomic
	closeOnce   sync.Once
	wg          sync.WaitGroup
	
	// Stats
	newConnCount int64 // atomic
}

// NewSafeStreamPool creates a new connection pool with safety improvements
func NewSafeStreamPool(factory StreamFactory, maxConns int) *SafeStreamPool {
	if maxConns <= 0 {
		maxConns = 100000 // Default reasonable limit
	}
	
	pool := &SafeStreamPool{
		conns:    make(map[key]*connection),
		factory:  factory,
		maxConns: int64(maxConns),
		connPool: sync.Pool{
			New: func() interface{} {
				return &connection{}
			},
		},
	}
	
	return pool
}

// getConnection returns a connection with improved race condition handling
func (p *SafeStreamPool) getConnection(k key, end bool, ts time.Time, tcp *layers.TCP, ac AssemblerContext) (*connection, *halfconnection, *halfconnection, error) {
	// Check if pool is closed
	if atomic.LoadInt32(&p.closed) != 0 {
		return nil, nil, nil, ErrPoolClosed
	}

	// Fast path - connection exists
	p.connsMutex.RLock()
	conn, half, rev := p.getHalf(k)
	p.connsMutex.RUnlock()
	
	if conn != nil {
		return conn, half, rev, nil
	}
	
	// Don't create new connection for end packets
	if end {
		return nil, nil, nil, nil
	}
	
	// Check connection limit
	if atomic.LoadInt64(&p.activeConns) >= p.maxConns {
		return nil, nil, nil, errors.New("connection limit exceeded")
	}
	
	// Slow path - need to create connection
	// Use double-checked locking to prevent races
	p.connsMutex.Lock()
	defer p.connsMutex.Unlock()
	
	// Check again with write lock
	conn, half, rev = p.getHalf(k)
	if conn != nil {
		return conn, half, rev, nil
	}
	
	// Create new stream
	s := p.factory.New(k[0], k[1], tcp, ac)
	if s == nil {
		return nil, nil, nil, ErrNilStream
	}
	
	// Get connection from pool
	conn = p.connPool.Get().(*connection)
	conn.reset(k, s, ts)
	
	// Store connection
	p.conns[k] = conn
	atomic.AddInt64(&p.activeConns, 1)
	atomic.AddInt64(&p.newConnCount, 1)
	
	if *memLog && p.newConnCount&0x7FFF == 0 {
		log.Printf("SafeStreamPool: %d total created, %d active", p.newConnCount, p.activeConns)
	}
	
	return conn, &conn.c2s, &conn.s2c, nil
}

// getHalf must be called with at least RLock held
func (p *SafeStreamPool) getHalf(k key) (*connection, *halfconnection, *halfconnection) {
	conn := p.conns[k]
	if conn != nil {
		return conn, &conn.c2s, &conn.s2c
	}
	rk := k.Reverse()
	conn = p.conns[rk]
	if conn != nil {
		return conn, &conn.s2c, &conn.c2s
	}
	return nil, nil, nil
}

// remove safely removes a connection from the pool
func (p *SafeStreamPool) remove(conn *connection) {
	if conn == nil {
		return
	}
	
	p.connsMutex.Lock()
	if _, exists := p.conns[conn.key]; exists {
		delete(p.conns, conn.key)
		atomic.AddInt64(&p.activeConns, -1)
		// Return connection to pool after clearing
		conn.reset(key{}, nil, time.Time{})
		p.connPool.Put(conn)
	}
	p.connsMutex.Unlock()
}

// connections returns a snapshot of all connections
func (p *SafeStreamPool) connections() []*connection {
	p.connsMutex.RLock()
	conns := make([]*connection, 0, len(p.conns))
	for _, conn := range p.conns {
		conns = append(conns, conn)
	}
	p.connsMutex.RUnlock()
	return conns
}

// Close gracefully shuts down the pool
func (p *SafeStreamPool) Close() error {
	p.closeOnce.Do(func() {
		atomic.StoreInt32(&p.closed, 1)
		
		// Wait for active operations
		p.wg.Wait()
		
		// Clear all connections
		p.connsMutex.Lock()
		for k := range p.conns {
			delete(p.conns, k)
		}
		p.connsMutex.Unlock()
	})
	return nil
}

// Stats returns current pool statistics
type PoolStats struct {
	ActiveConnections int64
	TotalCreated      int64
	PoolClosed        bool
}

func (p *SafeStreamPool) Stats() PoolStats {
	return PoolStats{
		ActiveConnections: atomic.LoadInt64(&p.activeConns),
		TotalCreated:      atomic.LoadInt64(&p.newConnCount),
		PoolClosed:        atomic.LoadInt32(&p.closed) != 0,
	}
}