package project

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	projectlogic "github.com/yanshicheng/cloud-back/application/manager-api/internal/logic/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

func GetProjectsByUserIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var userID uint64
		if raw := strings.TrimSpace(r.URL.Query().Get("userId")); raw != "" {
			v, err := strconv.ParseUint(raw, 10, 64)
			if err != nil {
				http.Error(w, "invalid userId", http.StatusBadRequest)
				return
			}
			userID = v
		}
		req := types.GetProjectsByUserIdRequest{
			UserID: userID,
			Name:   strings.TrimSpace(r.URL.Query().Get("name")),
		}

		l := projectlogic.NewGetProjectsByUserIdLogic(r.Context(), svcCtx)
		resp, err := l.GetProjectsByUserId(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
