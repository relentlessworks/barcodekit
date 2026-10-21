package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/relentlessworks/barcodekit/internal/api"
	"github.com/relentlessworks/barcodekit/internal/config"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	handler := api.NewHandler(cfg.NoAuth)
	handler.Register(mux)

	log.Printf("barcodekit starting on %s (no-auth=%v)", cfg.Addr, cfg.NoAuth)
	log.Printf("  help: http://localhost%s/help", cfg.Addr)
	log.Printf("  types: http://localhost%s/types", cfg.Addr)

	server := &http.Server{
		Addr:    cfg.Addr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}

	_ = fmt.Sprintf // keep import if needed
}
