package reassembly

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// TestSafeAssemblerRaceCondition tests that the race condition in connection creation is fixed
func TestSafeAssemblerRaceCondition(t *testing.T) {
	factory := &testStreamFactory{}
	pool := NewSafeStreamPool(factory, 1000)
	defer pool.Close()

	opts := SafeAssemblerOptions{
		MaxBufferedPagesTotal:         100,
		MaxBufferedPagesPerConnection: 10,
	}
	
	assembler, err := NewSafeAssembler(pool, opts)
	if err != nil {
		t.Fatalf("Failed to create assembler: %v", err)
	}
	defer assembler.Close()

	// Create concurrent connections with same key
	var wg sync.WaitGroup
	key := key{
		gopacket.NewFlow(gopacket.EndpointIPv4, []byte{1, 2, 3, 4}, []byte{5, 6, 7, 8}),
		gopacket.NewFlow(gopacket.EndpointTCPPort, []byte{0, 80}, []byte{0, 8080}),
	}

	// Try to create same connection from multiple goroutines
	created := 0
	var mu sync.Mutex
	
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			tcp := &layers.TCP{
				SrcPort: 80,
				DstPort: 8080,
				SYN:     true,
				Seq:     1000,
			}
			
			ctx := &assemblerSimpleContext{Timestamp: time.Now()}
			conn, _, _, err := pool.getConnection(key, false, time.Now(), tcp, ctx)
			if err == nil && conn != nil {
				mu.Lock()
				created++
				mu.Unlock()
			}
		}()
	}
	
	wg.Wait()
	
	// Should only create one connection despite race attempts
	if created != 1 {
		t.Errorf("Expected 1 connection created, got %d", created)
	}
	
	// Verify pool stats
	stats := pool.Stats()
	if stats.ActiveConnections != 1 {
		t.Errorf("Expected 1 active connection, got %d", stats.ActiveConnections)
	}
}

// TestSafeAssemblerMemoryLimits tests memory limit enforcement
func TestSafeAssemblerMemoryLimits(t *testing.T) {
	factory := &testStreamFactory{}
	pool := NewSafeStreamPool(factory, 10) // Low connection limit
	defer pool.Close()

	opts := SafeAssemblerOptions{
		MaxBufferedPagesTotal:         5, // Very low page limit
		MaxBufferedPagesPerConnection: 2,
	}
	
	assembler, err := NewSafeAssembler(pool, opts)
	if err != nil {
		t.Fatalf("Failed to create assembler: %v", err)
	}
	defer assembler.Close()

	// Try to create more connections than allowed
	errors := 0
	for i := 0; i < 20; i++ {
		flow := gopacket.NewFlow(gopacket.EndpointIPv4, 
			[]byte{1, 2, 3, byte(i)}, 
			[]byte{5, 6, 7, 8})
		
		tcp := &layers.TCP{
			SrcPort:   80,
			DstPort:   uint16(8080 + i),
			SYN:       true,
			Seq:       1000,
			BaseLayer: layers.BaseLayer{Payload: []byte{1, 2, 3}},
		}
		
		err := assembler.AssembleWithContext(flow, tcp, &assemblerSimpleContext{Timestamp: time.Now()})
		if err != nil {
			errors++
		}
	}
	
	// Should have rejected some connections
	if errors == 0 {
		t.Error("Expected some connection rejections due to limit")
	}
	
	stats := pool.Stats()
	if stats.ActiveConnections > 10 {
		t.Errorf("Active connections %d exceeded limit of 10", stats.ActiveConnections)
	}
}

// TestSafeAssemblerNilHandling tests nil pointer protection
func TestSafeAssemblerNilHandling(t *testing.T) {
	factory := &testStreamFactory{returnNil: true}
	pool := NewSafeStreamPool(factory, 100)
	defer pool.Close()

	opts := SafeAssemblerOptions{}
	assembler, err := NewSafeAssembler(pool, opts)
	if err != nil {
		t.Fatalf("Failed to create assembler: %v", err)
	}
	defer assembler.Close()

	flow := gopacket.NewFlow(gopacket.EndpointIPv4, []byte{1, 2, 3, 4}, []byte{5, 6, 7, 8})
	tcp := &layers.TCP{
		SrcPort:   80,
		DstPort:   8080,
		SYN:       true,
		Seq:       1000,
		BaseLayer: layers.BaseLayer{Payload: []byte{1, 2, 3}},
	}
	
	// Should not panic even if factory returns nil
	err = assembler.AssembleWithContext(flow, tcp, &assemblerSimpleContext{Timestamp: time.Now()})
	if err == nil {
		t.Error("Expected error when factory returns nil")
	}
}

// TestSafeAssemblerSliceBounds tests safe slice operations
func TestSafeAssemblerSliceBounds(t *testing.T) {
	factory := &testStreamFactory{}
	pool := NewSafeStreamPool(factory, 100)
	defer pool.Close()

	opts := SafeAssemblerOptions{}
	assembler, err := NewSafeAssembler(pool, opts)
	if err != nil {
		t.Fatalf("Failed to create assembler: %v", err)
	}
	defer assembler.Close()

	flow := gopacket.NewFlow(gopacket.EndpointIPv4, []byte{1, 2, 3, 4}, []byte{5, 6, 7, 8})
	
	// Send packets that would cause slice bounds issues in original code
	packets := []layers.TCP{
		{
			SrcPort:   80,
			DstPort:   8080,
			SYN:       true,
			Seq:       4294967290, // Near uint32 max
			BaseLayer: layers.BaseLayer{Payload: []byte{1, 2, 3, 4, 5}},
		},
		{
			SrcPort:   80,
			DstPort:   8080,
			Seq:       10, // Wrapped around
			BaseLayer: layers.BaseLayer{Payload: []byte{6, 7, 8, 9, 10}},
		},
		{
			SrcPort:   80,
			DstPort:   8080,
			Seq:       5, // Overlap
			BaseLayer: layers.BaseLayer{Payload: []byte{11, 12, 13, 14, 15}},
		},
	}
	
	// None of these should panic
	for _, tcp := range packets {
		err := assembler.AssembleWithContext(flow, &tcp, &assemblerSimpleContext{Timestamp: time.Now()})
		if err != nil {
			t.Logf("Packet processing error (expected): %v", err)
		}
	}
}

// TestSafeAssemblerGracefulShutdown tests clean shutdown
func TestSafeAssemblerGracefulShutdown(t *testing.T) {
	factory := &testStreamFactory{}
	pool := NewSafeStreamPool(factory, 100)
	
	opts := SafeAssemblerOptions{
		FlushInterval: 100 * time.Millisecond,
	}
	
	assembler, err := NewSafeAssembler(pool, opts)
	if err != nil {
		t.Fatalf("Failed to create assembler: %v", err)
	}
	
	// Create some connections
	for i := 0; i < 5; i++ {
		flow := gopacket.NewFlow(gopacket.EndpointIPv4, 
			[]byte{1, 2, 3, byte(i)}, 
			[]byte{5, 6, 7, 8})
		
		tcp := &layers.TCP{
			SrcPort:   80,
			DstPort:   uint16(8080 + i),
			SYN:       true,
			Seq:       1000,
			BaseLayer: layers.BaseLayer{Payload: []byte{1, 2, 3}},
		}
		
		assembler.AssembleWithContext(flow, tcp, &assemblerSimpleContext{Timestamp: time.Now()})
	}
	
	// Close should not hang
	done := make(chan bool)
	go func() {
		assembler.Close()
		pool.Close()
		done <- true
	}()
	
	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Error("Shutdown hung")
	}
}

// TestSafeAssemblerInvalidOptions tests option validation
func TestSafeAssemblerInvalidOptions(t *testing.T) {
	factory := &testStreamFactory{}
	pool := NewSafeStreamPool(factory, 100)
	defer pool.Close()

	testCases := []struct {
		name string
		opts SafeAssemblerOptions
		valid bool
	}{
		{
			name: "negative total pages",
			opts: SafeAssemblerOptions{
				MaxBufferedPagesTotal: -1,
			},
			valid: false,
		},
		{
			name: "negative per connection pages",
			opts: SafeAssemblerOptions{
				MaxBufferedPagesPerConnection: -1,
			},
			valid: false,
		},
		{
			name: "valid options",
			opts: SafeAssemblerOptions{
				MaxBufferedPagesTotal:         100,
				MaxBufferedPagesPerConnection: 10,
			},
			valid: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewSafeAssembler(pool, tc.opts)
			if tc.valid && err != nil {
				t.Errorf("Expected valid options, got error: %v", err)
			}
			if !tc.valid && err == nil {
				t.Error("Expected error for invalid options")
			}
		})
	}
}

// testStreamFactory for testing
type testStreamFactory struct {
	streams   []Stream
	returnNil bool
	mu        sync.Mutex
}

func (f *testStreamFactory) New(net, transport gopacket.Flow, tcp *layers.TCP, ac AssemblerContext) Stream {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	if f.returnNil {
		return nil
	}
	
	s := &testStream{
		net:       net,
		transport: transport,
	}
	f.streams = append(f.streams, s)
	return s
}

type testStream struct {
	net       gopacket.Flow
	transport gopacket.Flow
	data      []byte
	mu        sync.Mutex
}

func (s *testStream) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir TCPFlowDirection, nextSeq Sequence, start *bool, ac AssemblerContext) bool {
	return true
}

func (s *testStream) ReassembledSG(sg ScatterGather, ac AssemblerContext) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	length, _ := sg.Lengths()
	if length > 0 {
		data := sg.Fetch(length)
		s.data = append(s.data, data...)
	}
}

func (s *testStream) ReassemblyComplete(ac AssemblerContext) bool {
	return true
}