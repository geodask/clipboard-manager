package api

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/geodask/clipboard-manager/internal/config"
)

type Server struct {
	socketPath string
	httpServer *http.Server
	listener   net.Listener
	logger     *slog.Logger
}

func NewServer(service Service, cfg config.APIConfig, logger *slog.Logger) *Server {
	handler := NewHandler(service)

	return &Server{
		socketPath: cfg.SocketPath,
		httpServer: &http.Server{
			Handler:      handler.Routes(),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	os.Remove(s.socketPath)

	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		s.logger.Error("failed to create unix socket", "socket_path", s.socketPath, "error", err)
		return fmt.Errorf("failed to create unix socket: %w", err)
	}
	s.listener = listener

	if err := os.Chmod(s.socketPath, 0600); err != nil {
		listener.Close()
		s.logger.Error("failed to set socket permissions", "socket_path", s.socketPath, "error", err)
		return fmt.Errorf("failed to set socket permissions: %w", err)
	}

	s.logger.Info("API server listening", "socket_path", s.socketPath)

	if err := s.httpServer.Serve(listener); err != http.ErrServerClosed {
		s.logger.Error("API server error", "error", err)
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down API server")

	err := s.httpServer.Shutdown(ctx)

	os.Remove(s.socketPath)

	if err != nil {
		s.logger.Error("API server shutdown error", "error", err)
	} else {
		s.logger.Info("API server shutdown complete")
	}

	return err
}
