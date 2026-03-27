package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/messages", func(w http.ResponseWriter, r *http.Request) {
		handleMessages(w, r, cfg)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "claude-nvidia-proxy",
			"health":  "ok",
		})
	})

	srv := &http.Server{
		Addr:              cfg.addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      0, // allow streaming
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("listening on %s", cfg.addr)
	log.Printf("upstream: %s", cfg.upstreamURL)
	if cfg.serverAPIKey != "" {
		log.Printf("inbound auth: enabled")
	} else {
		log.Printf("inbound auth: disabled (SERVER_API_KEY not set)")
	}
	log.Fatal(srv.ListenAndServe())
}
