package main

import (
	"context"
	"eth_parser/internal/delivery/httpserver"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	errChan := make(chan error, 1)
	// Initialize server
	server := httpserver.NewServer("8080")

	// Create signal channel for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server
	server.Start(errChan)

	// Wait for error or shutdown signal
	select {
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	case sig := <-sigChan:
		log.Printf("Received signal for graceful shutdown: %v", sig)
	}

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Stop(ctx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
