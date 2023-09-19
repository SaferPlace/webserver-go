package insecure

import (
	"context"
	"crypto/tls"
)

// Provider provides noop [certificate.Provider], which doesn't return a new [tls.Config].
type Provider struct {
}

// NewProvider provides noop [certificate.Provider], which doesn't return a new [tls.Config].
func NewProvider() *Provider {
	return &Provider{}
}

// Provide no security.
func (*Provider) Provide(_ context.Context, _ []string) (*tls.Config, error) {
	return nil, nil
}
