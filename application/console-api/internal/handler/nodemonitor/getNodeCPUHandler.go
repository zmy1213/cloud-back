package nodemonitor

import (
	"net/http"
	"strings"

	nodemonitorlogic "github.com/yanshicheng/cloud-back/application/console-api/internal/logic/nodemonitor"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/types"
)

func GetNodeCPUHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := &types.GetNodeCPURequest{
			ClusterUuid: strings.TrimSpace(r.URL.Query().Get("clusterUuid")),
			NodeName:    strings.TrimSpace(r.URL.Query().Get("nodeName")),
			Start:       strings.TrimSpace(r.URL.Query().Get("start")),
			End:         strings.TrimSpace(r.URL.Query().Get("end")),
			Step:        strings.TrimSpace(r.URL.Query().Get("step")),
		}

		l := nodemonitorlogic.NewGetNodeCPULogic(r.Context(), svcCtx)
		resp, err := l.GetNodeCPU(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJSON(w, resp)
	}
}
