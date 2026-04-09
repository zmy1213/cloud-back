package clustermonitor

import (
	"net/http"
	"strings"

	clustermonitorlogic "github.com/yanshicheng/cloud-back/application/console-api/internal/logic/clustermonitor"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/types"
)

func GetClusterOverviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := &types.GetClusterOverviewRequest{ClusterUuid: strings.TrimSpace(r.URL.Query().Get("clusterUuid"))}
		l := clustermonitorlogic.NewGetClusterOverviewLogic(r.Context(), svcCtx)
		resp, err := l.GetClusterOverview(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJSON(w, resp)
	}
}
