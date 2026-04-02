package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	appcfg "github.com/yanshicheng/cloud-back/common/config"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/auth"
	minioPkg "github.com/yanshicheng/cloud-back/pkg/minio"
	mysqlPkg "github.com/yanshicheng/cloud-back/pkg/mysql"
	redisPkg "github.com/yanshicheng/cloud-back/pkg/redis"
)

type Handler struct {
	cfg  appcfg.AppConfig
	auth *auth.Service
}

func New(cfg appcfg.AppConfig) *Handler {
	return &Handler{
		cfg:  cfg,
		auth: auth.NewService(cfg.Auth.AccessExpiresIn, cfg.Auth.RefreshExpiresIn),
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.Healthz)
	mux.HandleFunc("/portal/v1/auth/login", h.Login)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	resp, err := h.auth.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Healthz(w http.ResponseWriter, _ *http.Request) {
	type depStatus struct {
		Enabled bool   `json:"enabled"`
		OK      bool   `json:"ok"`
		Error   string `json:"error,omitempty"`
	}

	timeout := 2 * time.Second
	deps := map[string]depStatus{
		"mysql": {Enabled: h.cfg.Mysql.Enabled, OK: true},
		"redis": {Enabled: h.cfg.Redis.Enabled, OK: true},
		"minio": {Enabled: h.cfg.Minio.Enabled, OK: true},
	}

	if h.cfg.Mysql.Enabled {
		if err := mysqlPkg.Ping(h.cfg.Mysql.DataSource, timeout); err != nil {
			v := deps["mysql"]
			v.OK = false
			v.Error = err.Error()
			deps["mysql"] = v
		}
	}

	if h.cfg.Redis.Enabled {
		if err := redisPkg.Ping(h.cfg.Redis.Addr, h.cfg.Redis.Password, h.cfg.Redis.DB, timeout); err != nil {
			v := deps["redis"]
			v.OK = false
			v.Error = err.Error()
			deps["redis"] = v
		}
	}

	if h.cfg.Minio.Enabled {
		if err := minioPkg.Check(
			h.cfg.Minio.Endpoint,
			h.cfg.Minio.AccessKey,
			h.cfg.Minio.SecretKey,
			h.cfg.Minio.BucketName,
			h.cfg.Minio.UseSSL,
			timeout,
		); err != nil {
			v := deps["minio"]
			v.OK = false
			v.Error = err.Error()
			deps["minio"] = v
		}
	}

	status := http.StatusOK
	for _, d := range deps {
		if d.Enabled && !d.OK {
			status = http.StatusServiceUnavailable
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"service": h.cfg.Name,
		"status":  map[bool]string{true: "ok", false: "degraded"}[status == http.StatusOK],
		"deps":    deps,
	})
}
