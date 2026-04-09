package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/yanshicheng/cloud-back/application/console-api/internal/handler"
	appcfg "github.com/yanshicheng/cloud-back/common/config"
	"github.com/yanshicheng/cloud-back/pkg/config"
)

func main() {
	cfgPath := flag.String("f", "./application/console-api/etc/console-api.yaml", "config file path")
	flag.Parse()

	var cfg appcfg.AppConfig
	config.MustLoad(*cfgPath, &cfg)

	h := handler.New(cfg)
	mux := http.NewServeMux()
	h.Register(mux)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf("Starting %s at %s...\n", cfg.Name, addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
