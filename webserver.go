// Copyright 2022 SaferPlace

package webserver

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rs/cors"
	"go.uber.org/zap"
)

// Server hosts the connect service.
type Server struct {
	services           map[string]http.Handler
	logger             *zap.Logger
	server             *http.Server
	allowedCORSDomains []string
}

// New creates a new connect server.
func New(opts ...Option) (*Server, error) {
	logger, _ := zap.NewDevelopment()

	s := &Server{
		services: make(map[string]http.Handler),
		logger:   logger,
		server: &http.Server{
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
	}

	// Apply all the options.
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	mux := http.NewServeMux()

	// Register all the handlers
	for path, handler := range s.services {
		s.logger.Info("registering handler", zap.String("path", path))
		mux.Handle(path, handler)
	}

	corsHandler := cors.New(cors.Options{
		AllowedMethods: []string{
			// CORS preflight
			http.MethodOptions,
			// Metrics
			http.MethodGet,
			// connect RPCs
			http.MethodPost,
		},
		// Mirror the `Origin` header value in the `Access-Control-Allow-Origin`
		// preflight response header.
		// This is equivalent to `Access-Control-Allow-Origin: *`, but allows
		// for requests with credentials.
		// Note that this effectively disables CORS and is not safe for use in
		// production environments.
		AllowOriginFunc: func(origin string) bool {
			// Disable CORS when the domain list is not specified.
			// This might be a security issue long term.
			if len(s.allowedCORSDomains) == 0 {
				return true
			}
			return inSlice(s.allowedCORSDomains, origin)
		},
		// Note that rs/cors does not return `Access-Control-Allow-Headers: *`
		// in response to preflight requests with the following configuration.
		// It simply mirrors all headers listed in the `Access-Control-Request-Headers`
		// preflight request header.
		AllowedHeaders: []string{"*"},
		// We explicitly set the exposed header names instead of using the wildcard *,
		// because in requests with credentials, it is treated as the literal header
		// name "*" without special semantics.
		ExposedHeaders: []string{
			"Grpc-Status", "Grpc-Message", "Grpc-Status-Details-Bin", "X-Grpc-Test-Echo-Initial",
			"Trailer-X-Grpc-Test-Echo-Trailing-Bin"},
	}).Handler(mux)

	s.server.Handler = corsHandler

	return s, nil
}

// Run the server. If the TLSConfiguration is provided, the server will run securely,
// otherwise it will run insecurely.
func (s *Server) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("unable to listen on port %d: %w", port, err)
	}

	// Start insecurely if the TLSConfig is missing
	if s.server.TLSConfig == nil {
		s.logger.Info("starting server",
			zap.Int("port", port),
			zap.Bool("tls", false),
		)
		if err := s.server.Serve(lis); err != nil {
			return fmt.Errorf("unable to listed on port %d: %w", port, err)
		}
	}

	// Start a secure server if the TLSConfig is provided.
	s.logger.Info("starting server",
		zap.Int("port", port),
		zap.Bool("tls", true),
	)
	if err := s.server.ServeTLS(lis, "", ""); err != nil {
		return fmt.Errorf("unable to listed on port %d: %w", port, err)
	}

	return nil
}

func inSlice[T comparable](ss []T, x T) bool {
	for _, s := range ss {
		if s == x {
			return true
		}
	}
	return false
}
