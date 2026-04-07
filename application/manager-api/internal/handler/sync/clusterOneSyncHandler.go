package sync

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	synclogic "github.com/yanshicheng/cloud-back/application/manager-api/internal/logic/sync"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

func ClusterOneSyncHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rawID := strings.TrimPrefix(r.URL.Path, "/manager/v1/sync/cluster/")
		if rawID == "" || strings.Contains(rawID, "/") {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.ParseUint(rawID, 10, 64)
		if err != nil || id == 0 {
			http.Error(w, "invalid cluster id", http.StatusBadRequest)
			return
		}

		l := synclogic.NewClusterOneSyncLogic(r.Context(), svcCtx)
		resp, err := l.ClusterOneSync(&types.SyncClusterRequest{ID: id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
