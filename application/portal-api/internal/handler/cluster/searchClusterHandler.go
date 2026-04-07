package cluster

import (
	"encoding/json"
	"net/http"

	clusterlogic "github.com/yanshicheng/cloud-back/application/portal-api/internal/logic/cluster"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/types"
)

func SearchClusterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := types.SearchClusterRequest{
			Name:        r.URL.Query().Get("name"),
			Environment: r.URL.Query().Get("environment"),
		}

		l := clusterlogic.NewSearchClusterLogic(r.Context(), svcCtx)
		items, total := l.SearchCluster(&req)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(types.SearchClusterResponse{
			Items: items,
			Total: total,
		})
	}
}
