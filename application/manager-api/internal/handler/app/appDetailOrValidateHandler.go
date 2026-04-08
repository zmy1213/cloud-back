package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	applogic "github.com/yanshicheng/cloud-back/application/manager-api/internal/logic/app"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

func AppDetailOrValidateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimPrefix(r.URL.Path, "/manager/v1/app/")
		raw = strings.Trim(raw, "/")
		if raw == "" {
			http.NotFound(w, r)
			return
		}

		parts := strings.Split(raw, "/")
		switch {
		case len(parts) == 1:
			handleGetAppDetail(svcCtx, w, r, parts[0])
		case len(parts) == 2 && parts[1] == "validate":
			handleValidateApp(svcCtx, w, r, parts[0])
		default:
			http.NotFound(w, r)
		}
	}
}

func handleGetAppDetail(svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request, rawID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "invalid app id", http.StatusBadRequest)
		return
	}

	l := applogic.NewGetClusterAppDetailLogic(r.Context(), svcCtx)
	resp, err := l.GetClusterAppDetail(&types.ClusterAppDetailRequest{ID: id})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func handleValidateApp(svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request, rawID string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		http.Error(w, "invalid app id", http.StatusBadRequest)
		return
	}

	l := applogic.NewValidateClusterAppLogic(r.Context(), svcCtx)
	resp, err := l.ValidateClusterApp(&types.ClusterAppValidateRequest{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
