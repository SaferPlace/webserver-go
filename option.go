// Copyright 2022 SaferPlace

package webserver

import "go.uber.org/zap"

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
			name, path, handler := service()
			s.services[path] = handler
			s.logger.Info("registered service", zap.String("service", name))
		}

		return nil
	}
}
