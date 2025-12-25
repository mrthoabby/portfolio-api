package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"

	"github.com/mrthoabby/portfolio-api/internal/common/logger"
	"github.com/mrthoabby/portfolio-api/internal/common/scope"
	"github.com/mrthoabby/portfolio-api/internal/config"
)

// Server wraps HTTP server configuration and lifecycle
type Server struct {
	httpServer *http.Server
	logger     logger.Logger
}

// NewServer creates a new server instance
func NewServer(cfg *config.Config, router *chi.Mux, appLogger logger.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Server.Port,
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		logger: appLogger,
	}
}

// Start starts the HTTP server in a goroutine
func (instance *Server) Start() error {
	go func() {
		instance.logger.Info("Server starting",
			logger.String("port", instance.httpServer.Addr),
			logger.String("scope", scope.String()),
		)
		if err := instance.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			instance.logger.Error("Server failed to start", logger.Error(err))
			os.Exit(1)
		}
	}()
	return nil
}

// Shutdown gracefully shuts down the server
func (instance *Server) Shutdown(ctx context.Context) error {
	return instance.httpServer.Shutdown(ctx)
}

// WaitForShutdown waits for interrupt signal and performs graceful shutdown
func (instance *Server) WaitForShutdown() {
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	instance.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := instance.Shutdown(ctx); err != nil {
		instance.logger.Error("Server forced to shutdown", logger.Error(err))
		os.Exit(1)
	}

	instance.logger.Info("Server exited gracefully")
}
