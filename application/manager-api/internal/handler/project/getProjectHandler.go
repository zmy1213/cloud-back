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

func ProjectDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseProjectIDFromPath(r.URL.Path)
		if err != nil {
			http.Error(w, "invalid project id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetProject(svcCtx, w, r, id)
		case http.MethodPut:
			handleUpdateProject(svcCtx, w, r, id)
		case http.MethodDelete:
			handleDeleteProject(svcCtx, w, r, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleGetProject(svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request, id uint64) {
	l := projectlogic.NewGetProjectLogic(r.Context(), svcCtx)
	resp, err := l.GetProject(id)
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
}

func handleUpdateProject(svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request, id uint64) {
	var req types.UpdateProjectRequest
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

	l := projectlogic.NewUpdateProjectLogic(r.Context(), svcCtx)
	resp, err := l.UpdateProject(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func handleDeleteProject(svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request, id uint64) {
	l := projectlogic.NewDeleteProjectLogic(r.Context(), svcCtx)
	resp, err := l.DeleteProject(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func parseProjectIDFromPath(path string) (uint64, error) {
	raw := strings.TrimPrefix(path, "/manager/v1/project/")
	raw = strings.Trim(raw, "/")
	if raw == "" {
		return 0, errors.New("project id is required")
	}
	parts := strings.Split(raw, "/")
	if len(parts) != 1 {
		return 0, errors.New("invalid path")
	}
	id, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("invalid project id")
	}
	return id, nil
}
