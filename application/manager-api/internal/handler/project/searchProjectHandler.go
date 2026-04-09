package project

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	projectlogic "github.com/yanshicheng/cloud-back/application/manager-api/internal/logic/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

func SearchProjectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := types.SearchProjectRequest{
			Page:     parseUintQuery(r, "page", 1),
			PageSize: parseUintQuery(r, "pageSize", 10),
			Name:     strings.TrimSpace(r.URL.Query().Get("name")),
			Uuid:     strings.TrimSpace(r.URL.Query().Get("uuid")),
		}

		l := projectlogic.NewSearchProjectLogic(r.Context(), svcCtx)
		resp, err := l.SearchProject(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func parseUintQuery(r *http.Request, key string, defaultValue uint64) uint64 {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return defaultValue
	}
	v, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || v == 0 {
		return defaultValue
	}
	return v
}
