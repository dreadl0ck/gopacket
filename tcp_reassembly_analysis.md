# TCP Stream Reassembly Analysis

## Executive Summary

This document analyzes the TCP stream reassembly package for potential bugs in memory allocation, locking/concurrency, and logic errors that could cause crashes. Multiple critical and moderate issues were identified that could lead to memory leaks, race conditions, panics, and denial of service vulnerabilities.

## Critical Issues Found

### 1. **Race Condition in StreamPool.getConnection()**

**Location**: `tcpassembly/assembly.go:513-522` and `reassembly/memory.go:176-204`

**Issue**: In both implementations, there's a time-of-check-time-of-use (TOCTOU) race condition:
```go
p.mu.RLock()
conn := p.conns[k]
p.mu.RUnlock()
if end || conn != nil {
    return conn
}
// Race window here - another goroutine could create the same connection
s := p.factory.New(k[0], k[1])
p.mu.Lock()
// ...
```

**Impact**: Could lead to:
- Multiple Stream objects created for the same connection
- Memory leaks (orphaned Stream objects)
- Inconsistent state

### 2. **Nil Pointer Dereference in sendToConnection()**

**Location**: `tcpassembly/assembly.go:624-632`

**Issue**: The function checks if `conn.stream == nil` and panics:
```go
func (a *Assembler) sendToConnection(conn *connection) {
    a.addContiguous(conn)
    if conn.stream == nil {
        panic("why?")  // This is not a proper error handling
    }
    conn.stream.Reassembled(a.ret)
    // ...
}
```

**Impact**: Application crash in production

### 3. **Integer Overflow in Sequence Arithmetic**

**Location**: Multiple locations using `Sequence.Add()` and `Sequence.Difference()`

**Issue**: The sequence number handling doesn't properly validate against integer overflow in all cases:
```go
func (s Sequence) Add(t int) Sequence {
    return (s + Sequence(t)) & (uint32Size - 1)
}
```

**Impact**: Incorrect sequence tracking, potential infinite loops

### 4. **Memory Leak in Page Cache**

**Location**: `tcpassembly/assembly.go` - pageCache implementation

**Issue**: The page cache grows unboundedly and never shrinks:
```go
// TODO: The page caches used by an Assembler will grow to the size necessary
// to handle a workload, and currently will never shrink.
```

**Impact**: Memory exhaustion under varying traffic patterns

### 5. **Unsafe Concurrent Access to Connection State**

**Location**: `reassembly/tcpassembly.go:1269-1286` in FlushWithOptions

**Issue**: The connection's half-connections are accessed after unlocking:
```go
conn.mu.Unlock()
if remove {
    a.connPool.remove(conn)  // conn might be accessed by another thread
}
```

**Impact**: Use-after-free, crashes

### 6. **Improper Slice Handling in checkOverlap()**

**Location**: `reassembly/tcpassembly.go:818-820`

**Issue**: Negative slice indices without bounds checking:
```go
cur.bytes = cur.bytes[:-start.Difference(cur.seq)]
```

**Impact**: Panic if the difference calculation is incorrect

### 7. **Missing Validation in Assembler Options**

**Location**: Both implementations

**Issue**: No validation of AssemblerOptions values:
```go
type AssemblerOptions struct {
    MaxBufferedPagesTotal         int  // No validation
    MaxBufferedPagesPerConnection int  // No validation
}
```

**Impact**: Negative values could cause unexpected behavior

### 8. **Deadlock Risk in Flush Operations**

**Location**: `reassembly/tcpassembly.go` - FlushAll()

**Issue**: Calling flush operations while holding locks could deadlock if Stream callbacks try to access the assembler

## Moderate Issues

### 1. **Inefficient Connection Lookup**

The connection pool uses a map with composite keys, causing unnecessary allocations for each lookup.

### 2. **No Backpressure Mechanism**

The assembler can accept unlimited packets, leading to memory exhaustion under load.

### 3. **Poor Error Propagation**

Many functions use panics instead of returning errors, making recovery difficult.

### 4. **Missing Context Cancellation**

No way to cleanly shut down the assembler or cancel operations.

## Recommendations for Improvement

### 1. **Fix Race Conditions**
- Use double-checked locking pattern for connection creation
- Add connection state validation
- Use atomic operations for state changes

### 2. **Implement Proper Memory Management**
- Add page cache eviction policy
- Implement connection limits
- Add memory usage monitoring

### 3. **Improve Error Handling**
- Replace panics with error returns
- Add context support for cancellation
- Implement graceful degradation

### 4. **Add Input Validation**
- Validate all user inputs
- Add bounds checking for sequence operations
- Implement safe slice operations

### 5. **Enhance Concurrency Safety**
- Reduce lock granularity
- Use RWMutex where appropriate
- Implement lock-free data structures where possible

### 6. **Add Monitoring and Metrics**
- Connection count tracking
- Memory usage metrics
- Performance counters
- Error rate monitoring

## Security Considerations

1. **DoS Vulnerability**: Attackers can exhaust memory by:
   - Sending many partial connections
   - Sending highly fragmented streams
   - Exploiting the unbounded page cache

2. **Resource Exhaustion**: No limits on:
   - Number of concurrent connections
   - Memory per connection
   - Total memory usage

3. **Crash Vulnerabilities**: Multiple panic conditions can be triggered by malformed input

## Performance Implications

1. **Lock Contention**: Global locks on connection pool limit scalability
2. **Memory Fragmentation**: Frequent allocation/deallocation of pages
3. **Cache Inefficiency**: No locality optimization for hot connections