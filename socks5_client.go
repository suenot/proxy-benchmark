package main

import (
	"context"
	"golang.org/x/net/proxy"
	"net/http"
	"time"
)

// SOCKS5Client handles SOCKS5 requests through proxies
type SOCKS5Client struct {
	client  *http.Client
	timeout time.Duration
}

// NewSOCKS5Client creates a new SOCKS5 client with proxy support
func NewSOCKS5Client(p *Proxy, timeout time.Duration) (*SOCKS5Client, error) {
	auth := &proxy.Auth{
		User:     p.Username,
		Password: p.Password,
	}

	dialer, err := proxy.SOCKS5("tcp", p.Address(), auth, proxy.Direct)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	return &SOCKS5Client{
		client:  client,
		timeout: timeout,
	}, nil
}

// MakeRequest performs an HTTP request through SOCKS5 proxy and returns the time taken
func (s *SOCKS5Client) MakeRequest(ctx context.Context, targetURL string) (time.Duration, error) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	return duration, nil
}
