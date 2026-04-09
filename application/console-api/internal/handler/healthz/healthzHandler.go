package healthz

import (
	"encoding/json"
	"net/http"

	"github.com/yanshicheng/cloud-back/application/console-api/internal/svc"
)

func HealthzHandler(_ *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service": "console-api",
			"status":  "ok",
		})
	}
}
