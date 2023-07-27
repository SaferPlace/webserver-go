// Copyright (c) Safer.Place and contributors. All rights reserved
// Licensed under the MIT license. See LICENSE file in the project root for details.

// Package certificate provides configuration loader.
package certificate

import (
	"context"
	"crypto/tls"
)

// Provider allows to load TLS configuration.
type Provider interface {
	// Provide TLS config for the given domains.
	Provide(ctx context.Context, domains []string) (*tls.Config, error)
}
