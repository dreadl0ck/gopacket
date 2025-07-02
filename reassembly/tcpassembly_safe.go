// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package reassembly

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// SafeAssemblerOptions provides validated options for the assembler
type SafeAssemblerOptions struct {
	// MaxBufferedPagesTotal is the maximum number of pages to buffer across all connections
	MaxBufferedPagesTotal int
	// MaxBufferedPagesPerConnection is the maximum number of pages per connection
	MaxBufferedPagesPerConnection int
	// ConnectionTimeout is how long to keep idle connections
	ConnectionTimeout time.Duration
	// FlushInterval is how often to flush old connections
	FlushInterval time.Duration
}

// Validate checks that options are sensible
func (o *SafeAssemblerOptions) Validate() error {
	if o.MaxBufferedPagesTotal < 0 {
		return fmt.Errorf("%w: MaxBufferedPagesTotal must be >= 0", ErrInvalidOptions)
	}
	if o.MaxBufferedPagesPerConnection < 0 {
		return fmt.Errorf("%w: MaxBufferedPagesPerConnection must be >= 0", ErrInvalidOptions)
	}
	if o.MaxBufferedPagesTotal == 0 {
		o.MaxBufferedPagesTotal = 50000 // Default
	}
	if o.MaxBufferedPagesPerConnection == 0 {
		o.MaxBufferedPagesPerConnection = 1000 // Default  
	}
	if o.ConnectionTimeout == 0 {
		o.ConnectionTimeout = 2 * time.Minute
	}
	if o.FlushInterval == 0 {
		o.FlushInterval = 10 * time.Second
	}
	return nil
}

// SafeAssembler provides resilient TCP stream reassembly
type SafeAssembler struct {
	SafeAssemblerOptions
	ret      []byteContainer
	pc       *SafePageCache
	connPool *SafeStreamPool
	cacheLP  livePacket
	cacheSG  reassemblyObject
	
	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
	
	// Background flush goroutine
	flushTicker *time.Ticker
	wg          sync.WaitGroup
	
	// Statistics
	stats AssemblerStats
}

// AssemblerStats tracks assembler statistics
type AssemblerStats struct {
	PacketsProcessed uint64
	BytesProcessed   uint64
	Errors           uint64
	ConnectionsFlush uint64
}

// NewSafeAssembler creates a new assembler with safety improvements
func NewSafeAssembler(pool *SafeStreamPool, options SafeAssemblerOptions) (*SafeAssembler, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	a := &SafeAssembler{
		SafeAssemblerOptions: options,
		ret:                  make([]byteContainer, 0, assemblerReturnValueInitialSize),
		pc:                   newSafePageCache(int64(options.MaxBufferedPagesTotal)),
		connPool:             pool,
		ctx:                  ctx,
		cancel:               cancel,
		flushTicker:          time.NewTicker(options.FlushInterval),
	}
	
	// Start background flush goroutine
	a.wg.Add(1)
	go a.backgroundFlush()
	
	return a, nil
}

// backgroundFlush periodically flushes old connections
func (a *SafeAssembler) backgroundFlush() {
	defer a.wg.Done()
	
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-a.flushTicker.C:
			cutoff := time.Now().Add(-a.ConnectionTimeout)
			flushed, closed := a.FlushWithOptions(FlushOptions{T: cutoff, TC: cutoff})
			if flushed > 0 || closed > 0 {
				atomic.AddUint64(&a.stats.ConnectionsFlush, uint64(flushed))
				if *debugLog {
					log.Printf("Background flush: flushed=%d, closed=%d", flushed, closed)
				}
			}
		}
	}
}

// AssembleWithContext safely reassembles TCP packets
func (a *SafeAssembler) AssembleWithContext(netFlow gopacket.Flow, t *layers.TCP, ac AssemblerContext) error {
	// Check context
	if err := a.ctx.Err(); err != nil {
		return fmt.Errorf("assembler closed: %w", err)
	}
	
	// Update stats
	atomic.AddUint64(&a.stats.PacketsProcessed, 1)
	atomic.AddUint64(&a.stats.BytesProcessed, uint64(len(t.Payload)))
	
	// Reset return buffer
	a.ret = a.ret[:0]
	
	// Validate inputs
	if t == nil {
		return errors.New("nil TCP layer")
	}
	
	key := key{netFlow, t.TransportFlow()}
	ci := ac.GetCaptureInfo()
	timestamp := ci.Timestamp
	
	// Get connection
	conn, half, _, err := a.connPool.getConnection(key, false, timestamp, t, ac)
	if err != nil {
		atomic.AddUint64(&a.stats.Errors, 1)
		return fmt.Errorf("failed to get connection: %w", err)
	}
	if conn == nil {
		if *debugLog {
			log.Printf("%v got empty packet on empty connection", key)
		}
		return nil
	}
	
	// Process under connection lock
	conn.mu.Lock()
	defer conn.mu.Unlock()
	
	// Update timestamp
	if half.lastSeen.Before(timestamp) {
		half.lastSeen = timestamp
	}
	
	// Check if stream accepts packet
	start := half.nextSeq == invalidSequence && t.SYN
	if !half.stream.Accept(t, ci, half.dir, half.nextSeq, &start, ac) {
		if *debugLog {
			log.Printf("Stream rejected packet")
		}
		return nil
	}
	
	if half.closed {
		if *debugLog {
			log.Printf("%v got packet on closed half", key)
		}
		return nil
	}
	
	// Process packet
	seq, bytes := Sequence(t.Seq), t.Payload
	if t.ACK {
		half.ackSeq = Sequence(t.Ack)
	}
	
	action := assemblerAction{
		nextSeq: invalidSequence,
		queue:   true,
	}
	
	// Handle sequence processing
	if half.nextSeq == invalidSequence {
		if t.SYN || start {
			half.nextSeq = seq
			if t.SYN {
				half.nextSeq = half.nextSeq.Add(1)
			}
			action.queue = false
		}
	} else {
		diff := half.nextSeq.Difference(seq)
		if diff <= 0 {
			action.queue = false
		}
	}
	
	// Process bytes safely
	if err := a.handleBytesSafe(bytes, seq, half, t.SYN, t.RST || t.FIN, action, ac); err != nil {
		atomic.AddUint64(&a.stats.Errors, 1)
		return fmt.Errorf("failed to handle bytes: %w", err)
	}
	
	// Send to stream if we have data
	if len(a.ret) > 0 {
		nextSeq, err := a.sendToConnectionSafe(conn, half, ac)
		if err != nil {
			atomic.AddUint64(&a.stats.Errors, 1)
			return fmt.Errorf("failed to send to connection: %w", err)
		}
		if nextSeq != invalidSequence {
			half.nextSeq = nextSeq
			if t.FIN {
				half.nextSeq = half.nextSeq.Add(1)
			}
		}
	}
	
	return nil
}

// handleBytesSafe processes bytes with bounds checking
func (a *SafeAssembler) handleBytesSafe(bytes []byte, seq Sequence, half *halfconnection, start, end bool, action assemblerAction, ac AssemblerContext) error {
	a.cacheLP.bytes = bytes
	a.cacheLP.start = start
	a.cacheLP.end = end
	a.cacheLP.seq = seq
	a.cacheLP.ac = ac
	
	if action.queue {
		if err := a.checkOverlapSafe(half, true, ac); err != nil {
			return err
		}
		
		// Check buffer limits
		if (a.MaxBufferedPagesPerConnection > 0 && half.pages >= a.MaxBufferedPagesPerConnection) ||
			(a.MaxBufferedPagesTotal > 0 && atomic.LoadInt64(&a.pc.used) >= int64(a.MaxBufferedPagesTotal)) {
			if *debugLog {
				log.Printf("Hit buffer limit, flushing")
			}
			a.addNextFromConn(half)
			action.queue = false
		}
	}
	
	if !action.queue {
		// Process overlap for immediate data
		a.cacheLP.bytes, a.cacheLP.seq = a.overlapExistingSafe(half, seq, seq.Add(len(bytes)), a.cacheLP.bytes)
		if err := a.checkOverlapSafe(half, false, ac); err != nil {
			return err
		}
		if len(a.cacheLP.bytes) != 0 || end || start {
			a.ret = append(a.ret, &a.cacheLP)
		}
	}
	
	return nil
}

// checkOverlapSafe safely handles overlapping segments
func (a *SafeAssembler) checkOverlapSafe(half *halfconnection, queue bool, ac AssemblerContext) error {
	if len(a.cacheLP.bytes) == 0 && !queue {
		return nil
	}
	
	var next *page
	cur := half.last
	bytes := a.cacheLP.bytes
	start := a.cacheLP.seq
	end := start.Add(len(bytes))
	
	for cur != nil {
		curEnd := cur.seq.Add(len(cur.bytes))
		
		// Calculate differences safely
		diffStart := start.Difference(cur.seq)
		diffEnd := end.Difference(curEnd)
		
		// Handle different overlap cases
		if diffEnd <= 0 && diffStart >= 0 {
			// Complete overlap - remove current
			if cur.prev != nil {
				cur.prev.next = cur.next
			} else {
				half.first = cur.next
			}
			if cur.next != nil {
				cur.next.prev = cur.prev
			} else {
				half.last = cur.prev
			}
			tmp := cur.prev
			half.pages -= cur.release(a.pc)
			cur = tmp
			continue
		}
		
		// Partial overlaps - safe slice operations
		if diffEnd < 0 && start.Difference(curEnd) > 0 {
			// Trim end of current
			trimLen := start.Difference(cur.seq)
			if trimLen > 0 && trimLen < len(cur.bytes) {
				cur.bytes = cur.bytes[:trimLen]
			}
			break
		}
		
		if diffStart > 0 && end.Difference(cur.seq) < 0 {
			// Trim start of current
			trimLen := end.Difference(cur.seq)
			if trimLen > 0 && trimLen < len(cur.bytes) {
				cur.bytes = cur.bytes[trimLen:]
				cur.seq = cur.seq.Add(trimLen)
			}
			next = cur
		}
		
		cur = cur.prev
	}
	
	// Queue remaining bytes if needed
	if len(bytes) > 0 && queue {
		p, p2, numPages, err := a.cacheLP.convertToPagesSafe(a.pc, 0, ac)
		if err != nil {
			return err
		}
		
		half.queuedPackets++
		half.queuedBytes += len(bytes)
		half.pages += numPages
		
		// Link pages
		if cur != nil {
			cur.next = p
			p.prev = cur
		} else {
			half.first = p
		}
		
		if next != nil {
			p2.next = next
			next.prev = p2
		} else {
			half.last = p2
		}
	}
	
	return nil
}

// overlapExistingSafe safely handles overlap with existing data
func (a *SafeAssembler) overlapExistingSafe(half *halfconnection, start, end Sequence, bytes []byte) ([]byte, Sequence) {
	if half.nextSeq == invalidSequence || len(bytes) == 0 {
		return bytes, start
	}
	
	diff := start.Difference(half.nextSeq)
	if diff <= 0 {
		return bytes, start
	}
	
	// Safe slice operation
	if diff >= len(bytes) {
		return []byte{}, half.nextSeq
	}
	
	half.overlapPackets++
	half.overlapBytes += diff
	
	return bytes[diff:], half.nextSeq
}

// sendToConnectionSafe safely sends data to the stream
func (a *SafeAssembler) sendToConnectionSafe(conn *connection, half *halfconnection, ac AssemblerContext) (Sequence, error) {
	if half.stream == nil {
		return invalidSequence, errors.New("nil stream in connection")
	}
	
	end, nextSeq := a.buildSG(half)
	
	// Call stream handler - this could panic, so we protect it
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Stream handler panicked: %v", r)
			}
		}()
		half.stream.ReassembledSG(&a.cacheSG, ac)
	}()
	
	a.cleanSG(half, ac)
	
	if end {
		a.closeHalfConnection(conn, half)
	}
	
	return nextSeq, nil
}

// convertToPagesSafe safely converts live packet to pages
func (lp *livePacket) convertToPagesSafe(pc *SafePageCache, skip int, ac AssemblerContext) (*page, *page, int, error) {
	if skip < 0 || skip > len(lp.bytes) {
		return nil, nil, 0, errors.New("invalid skip value")
	}
	
	ts := lp.captureInfo().Timestamp
	first, err := pc.next(ts)
	if err != nil {
		return nil, nil, 0, err
	}
	
	current := first
	current.prev = nil
	first.ac = ac
	first.start = lp.start && skip == 0
	numPages := 1
	
	seq, bytes := lp.seq.Add(skip), lp.bytes[skip:]
	
	for {
		length := min(len(bytes), pageBytes)
		current.bytes = current.buf[:length]
		copy(current.bytes, bytes)
		current.seq = seq
		
		bytes = bytes[length:]
		if len(bytes) == 0 {
			current.end = lp.end
			current.next = nil
			break
		}
		
		seq = seq.Add(length)
		next, err := pc.next(ts)
		if err != nil {
			// Clean up allocated pages
			for p := first; p != nil; {
				tmp := p.next
				pc.replace(p)
				p = tmp
			}
			return nil, nil, 0, err
		}
		
		current.next = next
		next.prev = current
		current = next
		current.ac = nil
		numPages++
	}
	
	return first, current, numPages, nil
}

// Close gracefully shuts down the assembler
func (a *SafeAssembler) Close() error {
	a.cancel()
	a.flushTicker.Stop()
	a.wg.Wait()
	
	// Flush all remaining connections
	a.FlushAll()
	
	return nil
}

// Stats returns current assembler statistics
func (a *SafeAssembler) Stats() AssemblerStats {
	return AssemblerStats{
		PacketsProcessed: atomic.LoadUint64(&a.stats.PacketsProcessed),
		BytesProcessed:   atomic.LoadUint64(&a.stats.BytesProcessed),
		Errors:           atomic.LoadUint64(&a.stats.Errors),
		ConnectionsFlush: atomic.LoadUint64(&a.stats.ConnectionsFlush),
	}
}