package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.comm/kevinlucasklein/madden-discord-bot/pkg/config"
	"github.comm/kevinlucasklein/madden-discord-bot/pkg/madden"
	"github.comm/kevinlucasklein/madden-discord-bot/pkg/utils"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger, err := utils.NewLogger(cfg.LogLevel, cfg.LogToFile, cfg.LogDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	// Initialize the Madden service
	maddenService := madden.NewService(cfg.DataDir)
	maddenService.SetLogger(logger)

	// Create server mux and register routes
	mux := http.NewServeMux()
	maddenService.RegisterRoutes(mux, cfg.ExportURL)

	// Set up the server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	// Start the server in a goroutine
	go func() {
		logger.Info("Starting Madden Companion Export server on http://localhost:%d", cfg.Port)
		logger.Info("Export endpoint available at http://localhost:%d%s", cfg.Port, cfg.ExportURL)

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("Server failed: %v", err)
			os.Exit(1)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	logger.Info("Server gracefully stopped")
}
