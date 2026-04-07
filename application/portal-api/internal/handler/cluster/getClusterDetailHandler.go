package cluster

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	clusterlogic "github.com/yanshicheng/cloud-back/application/portal-api/internal/logic/cluster"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/types"
)

func GetClusterDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		prefix := "/manager/v1/cluster/"
		rawID := strings.TrimPrefix(r.URL.Path, prefix)
		if rawID == "" || strings.Contains(rawID, "/") {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.ParseUint(rawID, 10, 64)
		if err != nil {
			http.Error(w, "invalid cluster id", http.StatusBadRequest)
			return
		}
		req := types.GetClusterDetailRequest{ID: id}

		l := clusterlogic.NewGetClusterDetailLogic(r.Context(), svcCtx)
		item, ok := l.GetClusterDetail(&req)
		if !ok {
			http.Error(w, "cluster not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(item)
	}
}
