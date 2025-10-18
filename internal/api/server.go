package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

type Server struct {
	socketPath string
	httpServer *http.Server
	listener   net.Listener
}

func NewServer(service Service, socketPath string) *Server {
	handler := NewHandler(service)

	return &Server{
		socketPath: socketPath,
		httpServer: &http.Server{
			Handler:      handler.Routes(),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
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
