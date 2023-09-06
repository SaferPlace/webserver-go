// Copyright 2022 SaferPlace

package webserver

import (
	"crypto/tls"
	"time"

	"go.uber.org/zap"
	"safer.place/webserver/middleware"
)

// Option allows to override the default behaviour of the Server.
type Option func(*Server) error

// Logger overrides the default logger in the Server.
func Logger(logger *zap.Logger) Option {
	return func(s *Server) error {
		s.logger = logger
		return nil
	}
}

// Services provides the server with the list of services which should be served.
func Services(services ...Service) Option {
	return func(s *Server) error {
		for _, service := range services {
			path, handler := service()
			s.services[path] = handler
			s.logger.Info("registered service", zap.String("path", path))
		}

		return nil
	}
}

// TLSConfig allows to set the TLS configuration for the server.
func TLSConfig(cfg *tls.Config) Option {
	return func(s *Server) error {
		s.server.TLSConfig = cfg
		return nil
	}
}

// Middlewares allows to specify which middleware to use. This is also useful
// if you want to override the default Cors middleware to limit which domains
// CORS should be applicable to.
func Middlewares(middlewares ...middleware.Middleware) Option {
	return func(s *Server) error {
		s.middleware = middlewares
		return nil
	}
}

// ReadTimeout overrides the default read timeout for the HTTP Server
func ReadTimeout(d time.Duration) Option {
	return func(s *Server) error {
		s.server.ReadTimeout = d
		return nil
	}
}

// WriteTimeout overrides the default write timeout for the HTTP Server
func WriteTimeout(d time.Duration) Option {
	return func(s *Server) error {
		s.server.WriteTimeout = d
		return nil
	}
}
