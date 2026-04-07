package node

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	nodelogic "github.com/yanshicheng/cloud-back/application/manager-api/internal/logic/node"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

func GetNodeListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := types.SearchClusterNodeRequest{
			Page:        parseUintQuery(r, "page", 1),
			PageSize:    parseUintQuery(r, "pageSize", 10),
			OrderField:  strings.TrimSpace(r.URL.Query().Get("orderField")),
			IsAsc:       parseBoolQuery(r, "isAsc", false),
			ClusterUuid: strings.TrimSpace(r.URL.Query().Get("clusterUuid")),
		}
		if req.ClusterUuid == "" {
			http.Error(w, "clusterUuid is required", http.StatusBadRequest)
			return
		}

		l := nodelogic.NewGetNodeListLogic(r.Context(), svcCtx)
		resp := l.GetNodeList(&req)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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

func parseBoolQuery(r *http.Request, key string, defaultValue bool) bool {
	raw := strings.TrimSpace(strings.ToLower(r.URL.Query().Get(key)))
	if raw == "" {
		return defaultValue
	}
	switch raw {
	case "1", "true", "yes":
		return true
	case "0", "false", "no":
		return false
	default:
		return defaultValue
	}
}
