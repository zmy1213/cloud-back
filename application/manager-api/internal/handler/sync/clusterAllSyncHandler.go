package sync

import (
	"encoding/json"
	"net/http"

	synclogic "github.com/yanshicheng/cloud-back/application/manager-api/internal/logic/sync"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
)

func ClusterAllSyncHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		l := synclogic.NewClusterAllSyncLogic(r.Context(), svcCtx)
		resp, err := l.ClusterAllSync()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
