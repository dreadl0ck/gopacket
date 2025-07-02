# TCP Reassembly Implementation - Detailed Code Analysis

## 1. Sequence Number Arithmetic Implementation

### RFC 793 Section 3.3 / RFC 9293 Section 3.4 Requirements:
> "TCP must perform comparisons on sequence numbers in a way that accounts for 
> the possibility of sequence number wrap-around"

### Implementation Analysis:

**reassembly/tcpassembly.go:**
```go
const uint32Max = 0xFFFFFFFF

// Sequence is a TCP sequence number with wrap-around handling
type Sequence int64

func (s Sequence) Difference(t Sequence) int {
    // Handle wrap-around by checking if sequences are in different quarters
    if s > uint32Max-uint32Max/4 && t < uint32Max/4 {
        t += uint32Max
    } else if t > uint32Max-uint32Max/4 && s < uint32Max/4 {
        s += uint32Max
    }
    return int(t - s)
}

func (s Sequence) Add(t int) Sequence {
    return (s + Sequence(t)) & uint32Max
}
```

**Compliance:** ✅ FULLY COMPLIANT
- Correctly implements modulo 2^32 arithmetic
- Handles wrap-around cases properly
- Uses the "quarter space" approach for sequence comparison

## 2. Out-of-Order Segment Handling

### RFC Requirements:
> "A TCP receiver SHOULD accept out-of-order segments"

### Implementation:

**reassembly/tcpassembly.go - Page-based buffering system:**
```go
type page struct {
    bytes      []byte
    seq        Sequence
    prev, next *page      // Doubly-linked list
    buf        [pageBytes]byte
    ac         AssemblerContext
    seen       time.Time
    start, end bool
}

type halfconnection struct {
    first, last       *page    // Out-of-order pages (seq > nextSeq)
    saved             *page    // In-order pages already delivered
    nextSeq           Sequence // Next expected sequence
    // ... other fields
}
```

**Segment insertion logic:**
```go
func (a *Assembler) handleBytes(bytes []byte, seq Sequence, half *halfconnection, 
    start bool, end bool, action assemblerAction, ac AssemblerContext) assemblerAction {
    
    if half.nextSeq == invalidSequence {
        if start {
            half.nextSeq = seq.Add(len(bytes))
            // ... deliver data
        }
    } else {
        diff := half.nextSeq.Difference(seq)
        if diff > 0 {
            // Future packet - buffer it
            a.insertIntoConn(half, seq, bytes, ac)
        } else if diff < 0 {
            // Past packet - check for overlap
            bytes, seq = a.overlapExisting(half, seq, seq.Add(len(bytes)), bytes)
        }
    }
    // ...
}
```

**Compliance:** ✅ FULLY COMPLIANT
- Maintains ordered doubly-linked list of out-of-order segments
- Correctly inserts segments based on sequence numbers
- Handles gaps and overlaps appropriately

## 3. Duplicate and Overlapping Segment Handling

### RFC 793: 
> "TCP must be prepared to handle duplicate segments"

### Implementation:

**reassembly/tcpassembly.go:**
```go
func (a *Assembler) overlapExisting(half *halfconnection, start, end Sequence, 
    bytes []byte) ([]byte, Sequence) {
    
    if half.nextSeq == invalidSequence || start.Difference(half.nextSeq) < 0 {
        return bytes, start
    }
    
    overlap := half.nextSeq.Difference(start)
    if overlap > len(bytes) {
        // Completely included in what we've already seen
        return nil, start
    }
    
    // Remove overlap
    return bytes[overlap:], start.Add(overlap)
}
```

**Compliance:** ✅ COMPLIANT
- Detects overlapping segments
- Removes duplicate data
- Preserves non-duplicate portions

## 4. TCP State Machine Implementation

### RFC 793 Connection State Diagram

### Simplified Implementation:

**reassembly/tcpcheck.go:**
```go
const (
    TCPStateClosed      = 0
    TCPStateSynSent     = 1
    TCPStateEstablished = 2
    TCPStateCloseWait   = 3
    TCPStateLastAck     = 4
    TCPStateReset       = 5
)

func (t *TCPSimpleFSM) CheckState(tcp *layers.TCP, dir TCPFlowDirection) bool {
    switch t.state {
    case TCPStateClosed:
        if tcp.SYN && !tcp.ACK {
            t.state = TCPStateSynSent
            return true
        }
    case TCPStateSynSent:
        if tcp.SYN && tcp.ACK {
            t.state = TCPStateEstablished
            return true
        }
    // ... other states
    }
}
```

**Compliance:** ⚠️ PARTIALLY COMPLIANT
- Missing states: LISTEN, SYN-RECEIVED, FIN-WAIT-1, FIN-WAIT-2, CLOSING, TIME-WAIT
- Simplified for passive monitoring use case
- Handles basic connection lifecycle

## 5. Window Management and Flow Control

### RFC 793/9293 Window Requirements

### Implementation:

**reassembly/tcpcheck.go:**
```go
type tcpStreamOptions struct {
    mss           int
    scale         int
    receiveWindow uint
}

func (t *TCPOptionCheck) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, 
    dir TCPFlowDirection, nextSeq Sequence, start *bool) error {
    
    // Parse window scale option from SYN
    if tcp.SYN {
        for _, o := range tcp.Options {
            if o.OptionType == 3 { // Window scaling
                scale = int(o.OptionData[0])
            }
        }
    }
    
    // Calculate actual window
    options.receiveWindow = uint(tcp.Window)
    if options.scale > 0 {
        options.receiveWindow = options.receiveWindow << uint(options.scale)
    }
    
    // Check if data fits in window
    if revOptions.receiveWindow != 0 && diff > int(revOptions.receiveWindow) {
        return fmt.Errorf("%d > receiveWindow(%d)", diff, revOptions.receiveWindow)
    }
}
```

**Compliance:** ⚠️ BASIC COMPLIANCE
- Tracks advertised window
- Supports window scaling (RFC 7323)
- Limited enforcement (passive monitoring focus)

## 6. TCP Options Parsing

### Implementation:

**layers/tcp.go:**
```go
func (tcp *TCP) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
    // ... header parsing ...
    
    // Parse options
    data = data[20:dataStart]
    for len(data) > 0 {
        tcp.Options = append(tcp.Options, TCPOption{
            OptionType: TCPOptionKind(data[0])
        })
        opt := &tcp.Options[len(tcp.Options)-1]
        
        switch opt.OptionType {
        case TCPOptionKindEndList:
            tcp.Padding = data[1:]
            break OPTIONS
        case TCPOptionKindNop:
            opt.OptionLength = 1
        default:
            opt.OptionLength = data[1]
            if int(opt.OptionLength) > len(data) {
                return fmt.Errorf("Invalid TCP option length")
            }
            opt.OptionData = data[2:opt.OptionLength]
        }
        data = data[opt.OptionLength:]
    }
}
```

**Compliance:** ✅ FULLY COMPLIANT
- Correctly parses all option formats
- Handles variable-length options
- Validates option lengths

## 7. Memory Management and Performance

### Implementation Features:

**Page Cache System:**
```go
type pageCache struct {
    free         []*page
    pcSize       int
    size, used   int
    pages        [][]page
}

func (c *pageCache) next(ts time.Time) (p *page) {
    if len(c.free) == 0 {
        c.grow()  // Exponential growth
    }
    i := len(c.free) - 1
    p, c.free = c.free[i], c.free[:i]
    p.Reassembly = Reassembly{Bytes: p.buf[:0], Seen: ts}
    c.used++
    return p
}
```

**Design Benefits:**
- Object pooling reduces GC pressure
- Pre-allocated buffers (1900 bytes per page)
- Exponential growth strategy
- Reuses page objects

## 8. Key Non-Compliance Areas

### 1. Missing Congestion Control
```go
// No implementation of:
// - Slow start
// - Congestion avoidance
// - Fast retransmit
// - Fast recovery
```

### 2. No Retransmission Timer
```go
// No implementation of:
// - RTO calculation
// - Retransmission queue
// - Karn's algorithm
```

### 3. Simplified Connection Management
```go
// Missing:
// - TIME-WAIT state and 2MSL timer
// - Simultaneous open handling
// - Half-closed connection support
```

## Summary

The implementation is **well-suited for its intended purpose** as a passive TCP stream reassembly library. It correctly implements the core algorithms needed for:

1. **Sequence number handling** - Full compliance
2. **Out-of-order reassembly** - Full compliance  
3. **Duplicate detection** - Full compliance
4. **Basic state tracking** - Partial compliance (simplified)
5. **Options parsing** - Full compliance

The omissions (congestion control, retransmission, full state machine) are **appropriate design decisions** for a packet capture library that only needs to reconstruct streams, not actively participate in TCP connections.