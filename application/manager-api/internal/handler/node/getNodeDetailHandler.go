package node

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	nodelogic "github.com/yanshicheng/cloud-back/application/manager-api/internal/logic/node"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

func GetNodeDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rawID := strings.TrimPrefix(r.URL.Path, "/manager/v1/node/")
		if rawID == "" || strings.Contains(rawID, "/") {
			http.NotFound(w, r)
			return
		}
		id, err := strconv.ParseUint(rawID, 10, 64)
		if err != nil || id == 0 {
			http.Error(w, "invalid node id", http.StatusBadRequest)
			return
		}

		l := nodelogic.NewGetNodeDetailLogic(r.Context(), svcCtx)
		item, ok := l.GetNodeDetail(&types.NodeIdRequest{ID: id})
		if !ok {
			http.Error(w, "node not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(item)
	}
}
