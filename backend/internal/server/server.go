// Package server provides HTTP server construction and lifecycle management.
package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cjlapao/daedalus/backend/internal/handlers"
)

// New creates an [http.Server] configured with the given address and default
// routes. It returns immediately — call [Server.ListenAndServe] or
// [Server.Shutdown] to manage the lifecycle.
func New(addr string) *http.Server {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/health", handlers.HealthHandler())

	// Production static file serving: serve frontend/dist for non-API routes
	// when the directory exists; silently no-op otherwise.
	distPath := filepath.Join("..", "..", "frontend", "dist")
	if _, err := os.Stat(distPath); err == nil {
		fileServer := http.FileServer(http.Dir(distPath))
		mux.Handle("/", fileServer)
	}

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

// StartBlocks runs [srv.ListenAndServe] and blocks until the server exits.
// It logs the port the server bound to before blocking.
func StartBlocks(srv *http.Server, port int) {
	log.Printf("daedalus server listening on port %d", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("daedalus server terminated: %v", err)
	}
}

// Shutdown gracefully shuts down the server, allowing in-flight requests to
// complete within the given timeout.
func Shutdown(ctx context.Context, srv *http.Server) error {
	return srv.Shutdown(ctx)
}

// startsWith reports whether the string s begins with prefix.
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
