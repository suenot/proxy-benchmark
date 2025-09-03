package main

import (
	"fmt"
	"strings"
)

// Proxy represents a proxy server with all its details
type Proxy struct {
	Protocol string
	Host     string
	Port     string
	Username string
	Password string
	Status   string
}

// ParseProxy parses a proxy string into a Proxy struct
func ParseProxy(proxyString string) (*Proxy, error) {
	parts := strings.Split(proxyString, ":")
	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid proxy format: %s", proxyString)
	}

	return &Proxy{
		Protocol: parts[0],
		Host:     parts[1],
		Port:     parts[2],
		Username: parts[3],
		Password: parts[4],
		Status:   parts[5],
	}, nil
}

// IsValid checks if the proxy is valid for use
func (p *Proxy) IsValid() bool {
	return p.Status == "enabled"
}

// Address returns the host:port combination
func (p *Proxy) Address() string {
	return fmt.Sprintf("%s:%s", p.Host, p.Port)
}

// String returns the proxy as a string representation
func (p *Proxy) String() string {
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", p.Protocol, p.Host, p.Port, p.Username, p.Password, p.Status)
}
