package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/geodask/clipboard-manager/internal/config"
)

type Server struct {
	socketPath string
	httpServer *http.Server
	listener   net.Listener
}

func NewServer(service Service, cfg config.APIConfig) *Server {
	handler := NewHandler(service)

	return &Server{
		socketPath: cfg.SocketPath,
		httpServer: &http.Server{
			Handler:      handler.Routes(),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	os.Remove(s.socketPath)

	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("failed to create unix socket: %w", err)
	}
	s.listener = listener

	if err := os.Chmod(s.socketPath, 0600); err != nil {
		listener.Close()
		return fmt.Errorf("failed to set socket permissions: %w", err)
	}

	fmt.Printf("API server listening on %s\n", s.socketPath)

	if err := s.httpServer.Serve(listener); err != http.ErrServerClosed {
		return err
	}

	return nil

}

func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down API server...")

	err := s.httpServer.Shutdown(ctx)

	os.Remove(s.socketPath)

	return err
}
