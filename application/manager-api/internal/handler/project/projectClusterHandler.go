package project

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	projectlogic "github.com/yanshicheng/cloud-back/application/manager-api/internal/logic/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

func ProjectClusterEntryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req types.AddProjectClusterRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		if err := decoder.Decode(&struct{}{}); err != nil && !errors.Is(err, io.EOF) {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		l := projectlogic.NewAddProjectClusterLogic(r.Context(), svcCtx)
		resp, err := l.AddProjectCluster(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func SearchProjectClusterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req := types.SearchProjectClusterRequest{
			ProjectID:   parseUintQuery(r, "projectId", 0),
			ClusterUUID: strings.TrimSpace(r.URL.Query().Get("clusterUuid")),
		}
		l := projectlogic.NewSearchProjectClusterLogic(r.Context(), svcCtx)
		resp, err := l.SearchProjectCluster(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func ProjectClusterDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseProjectClusterIDFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid project cluster id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			l := projectlogic.NewGetProjectClusterLogic(r.Context(), svcCtx)
			resp, err := l.GetProjectCluster(id)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "not found") {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		case http.MethodPut:
			var req types.UpdateProjectClusterRequest
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&req); err != nil {
				http.Error(w, "invalid request body", http.StatusBadRequest)
				return
			}
			if err := decoder.Decode(&struct{}{}); err != nil && !errors.Is(err, io.EOF) {
				http.Error(w, "invalid request body", http.StatusBadRequest)
				return
			}
			req.ID = id

			l := projectlogic.NewUpdateProjectClusterLogic(r.Context(), svcCtx)
			resp, err := l.UpdateProjectCluster(&req)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		case http.MethodDelete:
			l := projectlogic.NewDeleteProjectClusterLogic(r.Context(), svcCtx)
			resp, err := l.DeleteProjectCluster(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func parseProjectClusterIDFromPath(path string) (uint64, error) {
	raw := strings.TrimPrefix(path, "/manager/v1/project/cluster/")
	raw = strings.Trim(raw, "/")
	if raw == "" {
		return 0, errors.New("id is required")
	}
	parts := strings.Split(raw, "/")
	if len(parts) != 1 {
		return 0, errors.New("invalid path")
	}
	id, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
