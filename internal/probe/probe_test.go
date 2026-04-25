package probe_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/probe"
)

func startTCPServer(t *testing.T) (string, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()
	_, port, _ := net.SplitHostPort(ln.Addr().String())
	return port, func() { _ = ln.Close() }
}

func TestProbe_Reachable(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	p := probe.New(2*time.Second, port)
	res := p.Probe(context.Background(), "127.0.0.1")

	if !res.Reachable {
		t.Fatalf("expected reachable, got error: %v", res.Error)
	}
	if res.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestProbe_Unreachable(t *testing.T) {
	p := probe.New(300*time.Millisecond, "9")
	res := p.Probe(context.Background(), "127.0.0.1")

	if res.Reachable {
		t.Fatal("expected unreachable")
	}
	if res.Error == nil {
		t.Error("expected non-nil error")
	}
}

func TestProbeAll_ReturnsAllResults(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	p := probe.New(2*time.Second, port)
	hosts := []string{"127.0.0.1", "127.0.0.1"}
	results := p.ProbeAll(context.Background(), hosts)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Reachable {
			t.Errorf("expected reachable for host %s", r.Host)
		}
	}
}

func TestProbe_DefaultTimeout(t *testing.T) {
	p := probe.New(0, "")
	if p == nil {
		t.Fatal("expected non-nil prober")
	}
}
