package clustermonitor

import (
	"net/http"
	"strings"

	clustermonitorlogic "github.com/yanshicheng/cloud-back/application/console-api/internal/logic/clustermonitor"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/types"
)

func GetClusterResourcesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := &types.GetClusterResourcesRequest{ClusterUuid: strings.TrimSpace(r.URL.Query().Get("clusterUuid"))}
		l := clustermonitorlogic.NewGetClusterResourcesLogic(r.Context(), svcCtx)
		resp, err := l.GetClusterResources(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJSON(w, resp)
	}
}
