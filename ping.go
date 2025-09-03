package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

// PingClient measures direct connection latency to proxy servers
type PingClient struct {
	timeout time.Duration
}

// NewPingClient creates a new ping client
func NewPingClient(timeout time.Duration) *PingClient {
	return &PingClient{
		timeout: timeout,
	}
}

// PingProxy measures the round-trip time to establish a TCP connection to the proxy
func (p *PingClient) PingProxy(ctx context.Context, proxy *Proxy) (time.Duration, error) {
	address := fmt.Sprintf("%s:%s", proxy.Host, proxy.Port)
	
	dialer := &net.Dialer{
		Timeout: p.timeout,
	}
	
	start := time.Now()
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to proxy %s: %w", address, err)
	}
	duration := time.Since(start)
	
	// Close the connection immediately after measuring connection time
	conn.Close()
	
	return duration, nil
}