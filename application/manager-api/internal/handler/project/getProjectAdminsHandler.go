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

func GetProjectAdminsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		rawProjectID := strings.TrimSpace(r.URL.Query().Get("projectId"))
		projectID, err := strconv.ParseUint(rawProjectID, 10, 64)
		if err != nil || projectID == 0 {
			http.Error(w, "projectId is required", http.StatusBadRequest)
			return
		}

		req := types.GetProjectAdminsRequest{ProjectID: projectID}
		l := projectlogic.NewGetProjectAdminsLogic(r.Context(), svcCtx)
		resp, err := l.GetProjectAdmins(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
