package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// HTTPClient handles HTTP/HTTPS requests through proxies
type HTTPClient struct {
	client  *http.Client
	timeout time.Duration
}

// NewHTTPClient creates a new HTTP client with proxy support
func NewHTTPClient(proxy *Proxy, timeout time.Duration) (*HTTPClient, error) {
	proxyURL, err := url.Parse(fmt.Sprintf("http://%s:%s@%s", proxy.Username, proxy.Password, proxy.Address()))
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	return &HTTPClient{
		client:  client,
		timeout: timeout,
	}, nil
}

// MakeRequest performs an HTTP request and returns the time taken
func (h *HTTPClient) MakeRequest(ctx context.Context, targetURL string) (time.Duration, error) {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return 0, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	return duration, nil
}
