package insecure

import (
	"context"
	"crypto/tls"
)

// Provider provides no TLSConfig
type Provider struct {
}

// NewProvider provies not security
func NewProvider() *Provider {
	return &Provider{}
}

// Provide no security.
func (*Provider) Provide(_ context.Context, _ []string) (*tls.Config, error) {
	return nil, nil
}
