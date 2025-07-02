# Safe TCP Stream Reassembly

This is an improved, production-ready implementation of TCP stream reassembly that addresses multiple critical issues found in the original gopacket implementation.

## Key Features

- **Thread-Safe**: Fixed race conditions in connection management
- **Memory-Safe**: Enforces memory limits and prevents unbounded growth  
- **Crash-Resistant**: No panics - all errors handled gracefully
- **Resource Management**: Automatic cleanup of idle connections
- **Observable**: Built-in statistics and monitoring
- **Configurable**: Flexible options with validation

## Quick Start

```go
import "github.com/yourorg/gopacket/reassembly"

// Create a stream factory
factory := &MyStreamFactory{}

// Create a connection pool with max 10000 connections
pool := reassembly.NewSafeStreamPool(factory, 10000)
defer pool.Close()

// Configure assembler options
opts := reassembly.SafeAssemblerOptions{
    MaxBufferedPagesTotal:         50000,  // Max pages across all connections
    MaxBufferedPagesPerConnection: 1000,   // Max pages per connection
    ConnectionTimeout:             2 * time.Minute,
    FlushInterval:                 10 * time.Second,
}

// Create assembler
assembler, err := reassembly.NewSafeAssembler(pool, opts)
if err != nil {
    log.Fatal(err)
}
defer assembler.Close()

// Process packets
err = assembler.AssembleWithContext(netFlow, tcpLayer, assemblerContext)
if err != nil {
    log.Printf("Assembly error: %v", err)
}

// Get statistics
stats := assembler.Stats()
log.Printf("Processed %d packets, %d bytes, %d errors", 
    stats.PacketsProcessed, stats.BytesProcessed, stats.Errors)
```

## Improvements Over Original

### 1. Fixed Critical Bugs

- **Race Condition**: Connection creation is now thread-safe with double-checked locking
- **Memory Leaks**: Page cache now has eviction and proper cleanup
- **Nil Panics**: All nil pointers checked before use
- **Slice Panics**: All slice operations have bounds checking

### 2. Added Safety Features  

- **Resource Limits**: Configurable limits on connections and memory
- **Input Validation**: All options validated with clear errors
- **Error Handling**: Errors returned instead of panics
- **Protected Callbacks**: Stream handlers wrapped in recover()

### 3. Enhanced Functionality

- **Context Support**: Clean cancellation and shutdown
- **Background Cleanup**: Automatic flushing of old connections  
- **Statistics**: Track packets, bytes, errors, and flushes
- **Object Pooling**: Reuse connections and pages for efficiency

## Configuration Options

```go
type SafeAssemblerOptions struct {
    // Maximum total pages buffered across all connections
    // Default: 50000
    MaxBufferedPagesTotal int
    
    // Maximum pages buffered per connection
    // Default: 1000  
    MaxBufferedPagesPerConnection int
    
    // How long to keep idle connections
    // Default: 2 minutes
    ConnectionTimeout time.Duration
    
    // How often to flush old connections
    // Default: 10 seconds
    FlushInterval time.Duration
}
```

## Error Handling

The safe implementation returns errors instead of panicking:

```go
var (
    ErrNilStream        = errors.New("stream factory returned nil stream")
    ErrInvalidOptions   = errors.New("invalid assembler options")  
    ErrPoolClosed       = errors.New("stream pool is closed")
    ErrConnectionClosed = errors.New("connection is closed")
)
```

## Monitoring

Get real-time statistics:

```go
// Assembler statistics
stats := assembler.Stats()
// AssemblerStats{
//     PacketsProcessed: 1234567,
//     BytesProcessed:   987654321,
//     Errors:           12,
//     ConnectionsFlush: 456,
// }

// Pool statistics  
poolStats := pool.Stats()
// PoolStats{
//     ActiveConnections: 234,
//     TotalCreated:      5678,
//     PoolClosed:        false,
// }
```

## Testing

Comprehensive test coverage including:

- Concurrent connection creation
- Memory limit enforcement
- Nil pointer handling
- Slice bounds checking
- Graceful shutdown
- Option validation

Run tests:
```bash
go test -race ./reassembly/...
```

## Performance

- **Lock Optimization**: RWMutex for read-heavy operations
- **Atomic Counters**: Lock-free statistics
- **Object Pooling**: Reduced allocations
- **Early Returns**: Optimized hot paths

## Security

- **DoS Protection**: Connection and memory limits
- **Input Validation**: All inputs sanitized
- **Crash Prevention**: No panics in production
- **Resource Limits**: Prevents exhaustion attacks

## Migration from Original

See [tcp_reassembly_improvements.md](../tcp_reassembly_improvements.md) for detailed migration guide.

## License

Same as original gopacket - BSD-style license.