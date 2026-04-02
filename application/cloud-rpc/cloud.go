package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/yanshicheng/cloud-back/pkg/config"
)

// rpcConfig maps application/cloud-rpc/etc/cloud-rpc.yaml.
type rpcConfig struct {
	Name     string `yaml:"Name"`
	ListenOn string `yaml:"ListenOn"`
	Timeout  int    `yaml:"Timeout"`
}

func main() {
	// -f supports custom config path, same style as kube-nova services.
	cfgPath := flag.String("f", "./application/cloud-rpc/etc/cloud-rpc.yaml", "config file path")
	flag.Parse()

	var cfg rpcConfig
	config.MustLoad(*cfgPath, &cfg)

	mux := http.NewServeMux()
	// Health endpoint for liveness/readiness probes.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	// Placeholder RPC-like endpoint for local testing.
	mux.HandleFunc("/rpc/ping", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("pong " + time.Now().Format(time.RFC3339)))
	})

	fmt.Printf("Starting %s at %s...\n", cfg.Name, cfg.ListenOn)
	// Kept as HTTP for simplicity in this empty scaffold.
	if err := http.ListenAndServe(cfg.ListenOn, mux); err != nil {
		panic(err)
	}
}
