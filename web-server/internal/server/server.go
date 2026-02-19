package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/nkrypt-xyz/nkrypt-xyz-web-server/internal/config"
)

// Server wraps the HTTP server and graceful shutdown handling.
type Server struct {
	httpServer *http.Server
}

// New creates a new Server instance.
func New(cfg *config.Config, handler http.Handler) *Server {
	addr := net.JoinHostPort("", itoa(cfg.Server.HTTP.Port))

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 15 * time.Second,
	}

	return &Server{httpServer: srv}
}

// ListenAndServe starts the HTTP server and blocks until it is shut down.
func (s *Server) ListenAndServe(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		log.Info().Str("addr", s.httpServer.Addr).Msg("HTTP server starting")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Wait for signal or context cancellation.
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Info().Msg("context cancelled, shutting down HTTP server")
	case sig := <-stopCh:
		log.Info().Str("signal", sig.String()).Msg("received signal, shutting down HTTP server")
	case err := <-errCh:
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("HTTP server shutdown error")
		return err
	}

	log.Info().Msg("HTTP server stopped gracefully")
	return nil
}

// small helper to avoid importing strconv just for this.
func itoa(v int) string {
	return fmt.Sprintf("%d", v)
}

