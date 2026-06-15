// Daedalus is the entrypoint for the Daedalus agent orchestration backend.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cjlapao/daedalus/backend/internal/server"
)

func main() {
	// Read port from env, default to 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	srv := server.New(addr)

	// Start the server in a goroutine so we can listen for shutdown signals.
	go server.StartBlocks(srv, mustParsePort(port))

	// Wait for interrupt or SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down daedalus server...")

	// Graceful shutdown with a 5-second deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx, srv); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("daedalus server stopped")
}

func mustParsePort(s string) int {
	p, err := strconv.Atoi(s)
	if err != nil {
		// PORT might be provided as ":8080" — strip leading colon.
		p, err = strconv.Atoi(s[1:])
		if err != nil {
			panic("invalid PORT value: " + s)
		}
	}
	return p
}
