package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/yanshicheng/cloud-back/pkg/config"
)

// apiConfig maps application/cloud-api/etc/cloud-api.yaml.
type apiConfig struct {
	Name    string `yaml:"Name"`
	Host    string `yaml:"Host"`
	Port    int    `yaml:"Port"`
	Timeout int    `yaml:"Timeout"`
}

func main() {
	// -f supports custom config path, same style as kube-nova services.
	cfgPath := flag.String("f", "./application/cloud-api/etc/cloud-api.yaml", "config file path")
	flag.Parse()

	var cfg apiConfig
	config.MustLoad(*cfgPath, &cfg)

	mux := http.NewServeMux()
	// Health endpoint for liveness/readiness probes.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	// Minimal API endpoint for quick connectivity verification.
	mux.HandleFunc("/cloud/v1/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service": cfg.Name,
			"time":    time.Now().Format(time.RFC3339),
			"status":  "pong",
		})
	})

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf("Starting %s at %s...\n", cfg.Name, addr)
	// Use standard net/http server for a lightweight scaffold.
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
