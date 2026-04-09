package nodemonitor

import (
	"net/http"
	"strings"

	nodemonitorlogic "github.com/yanshicheng/cloud-back/application/console-api/internal/logic/nodemonitor"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/types"
)

func ListNodesMetricsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := &types.ListNodesMetricsRequest{
			ClusterUuid: strings.TrimSpace(r.URL.Query().Get("clusterUuid")),
		}

		l := nodemonitorlogic.NewListNodesMetricsLogic(r.Context(), svcCtx)
		resp, err := l.ListNodesMetrics(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJSON(w, resp)
	}
}
