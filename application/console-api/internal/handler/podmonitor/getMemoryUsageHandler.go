package podmonitor

import (
	"net/http"
	"strings"

	podmonitorlogic "github.com/yanshicheng/cloud-back/application/console-api/internal/logic/podmonitor"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/types"
)

func GetMemoryUsageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := &types.GetMemoryUsageRequest{
			ClusterUuid: strings.TrimSpace(r.URL.Query().Get("clusterUuid")),
			Namespace:   strings.TrimSpace(r.URL.Query().Get("namespace")),
			PodName:     strings.TrimSpace(r.URL.Query().Get("podName")),
			Start:       strings.TrimSpace(r.URL.Query().Get("start")),
			End:         strings.TrimSpace(r.URL.Query().Get("end")),
			Step:        strings.TrimSpace(r.URL.Query().Get("step")),
		}

		l := podmonitorlogic.NewGetMemoryUsageLogic(r.Context(), svcCtx)
		resp, err := l.GetMemoryUsage(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJSON(w, resp)
	}
}
