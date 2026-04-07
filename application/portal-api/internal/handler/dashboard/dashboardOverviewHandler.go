package dashboard

import (
	"encoding/json"
	"net/http"

	dashboardlogic "github.com/yanshicheng/cloud-back/application/portal-api/internal/logic/dashboard"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/types"
)

func DashboardOverviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := types.DashboardOverviewRequest{
			Username:    r.URL.Query().Get("username"),
			ClusterUUID: r.URL.Query().Get("clusterUuid"),
		}
		l := dashboardlogic.NewDashboardOverviewLogic(r.Context(), svcCtx)
		resp, ok := l.DashboardOverview(&req)
		if !ok {
			http.Error(w, "cluster not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
