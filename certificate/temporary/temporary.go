// Copyright (c) Safer.Place and contributors. All rights reserved
// Licensed under the MIT license. See LICENSE file in the project root for details.

// Package temporary provides generated and temporary certificate.
package temporary

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"
)

// Provider generates a temporary certificate for the given domains. It
// implements the certificate.Provider interface.
type Provider struct {
	Config
}

type Config struct {
	ValidFor   time.Duration
	UseED25519 bool
}

// NewProvider returns a new temporary certificate provider.
func NewProvider(c Config) *Provider {
	return &Provider{c}
}

// Provide generates the certificate and returns the TLS config for the
// provided domains.
func (p *Provider) Provide(ctx context.Context, domains []string) (*tls.Config, error) {
	var err error

	var pub, priv any
	if p.UseED25519 {
		pub, priv, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("unable to generate key: %w", err)
		}
	} else {
		priv, err = rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, fmt.Errorf("unable to generate key: %w", err)
		}
		pub = &priv.(*rsa.PrivateKey).PublicKey
	}

	serial := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, _ = rand.Int(rand.Reader, serial)

	cert := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"Messr"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(p.ValidFor),
		PublicKey: pub,
		IsCA:      true,
		DNSNames:  domains,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, pub, priv)
	if err != nil {
		return nil, fmt.Errorf("unable to generate certificate: %w", err)
	}

	return &tls.Config{
		ClientAuth:               tls.NoClientCert,
		InsecureSkipVerify:       true,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{certBytes},
			PrivateKey:  priv,
			Leaf:        cert,
		}},
	}, nil
}
