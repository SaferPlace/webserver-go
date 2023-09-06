// Copyright 2022 SaferPlace

package webserver

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/saferplace/webserver-go/middleware"
)

// Server hosts the connect service.
type Server struct {
	services   map[string]http.Handler
	middleware []middleware.Middleware
	logger     *zap.Logger
	server     *http.Server
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
		middleware: []middleware.Middleware{
			// By default enable cors for all.
			middleware.Cors(nil),
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

	var handler http.Handler = mux

	// Register all middleware
	for _, middleware := range s.middleware {
		s.logger.Info("using middleware", zap.String("type", fmt.Sprintf("%T", middleware)))
		handler = middleware(handler)
	}

	s.server.Handler = handler

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
