package app

import (
	"encoding/json"
	"net/http"
	"strings"

	applogic "github.com/yanshicheng/cloud-back/application/manager-api/internal/logic/app"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

func GetClusterAppListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := types.ClusterAppListRequest{
			ClusterUuid: strings.TrimSpace(r.URL.Query().Get("clusterUuid")),
		}

		l := applogic.NewGetClusterAppListLogic(r.Context(), svcCtx)
		resp, err := l.GetClusterAppList(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
