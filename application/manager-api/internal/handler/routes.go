package handler

import (
	"net/http"

	apphandler "github.com/yanshicheng/cloud-back/application/manager-api/internal/handler/app"
	clusterhandler "github.com/yanshicheng/cloud-back/application/manager-api/internal/handler/cluster"
	healthzhandler "github.com/yanshicheng/cloud-back/application/manager-api/internal/handler/healthz"
	nodehandler "github.com/yanshicheng/cloud-back/application/manager-api/internal/handler/node"
	projecthandler "github.com/yanshicheng/cloud-back/application/manager-api/internal/handler/project"
	synchandler "github.com/yanshicheng/cloud-back/application/manager-api/internal/handler/sync"
)

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthzhandler.HealthzHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/cluster", clusterhandler.SearchClusterHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/cluster/", clusterhandler.GetClusterDetailHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/node", nodehandler.GetNodeListHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/node/", nodehandler.GetNodeDetailHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/app", apphandler.AddClusterAppHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/app/list", apphandler.GetClusterAppListHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/app/", apphandler.AppDetailOrValidateHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project", projecthandler.AddProjectHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project/search", projecthandler.SearchProjectHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project/user", projecthandler.GetProjectsByUserIdHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project/admin", projecthandler.AddProjectAdminHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project/admin/list", projecthandler.GetProjectAdminsHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project/cluster/search", projecthandler.SearchProjectClusterHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project/cluster", projecthandler.ProjectClusterEntryHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project/cluster/", projecthandler.ProjectClusterDetailHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project/workspace/search", projecthandler.SearchProjectWorkspaceHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project/workspace", projecthandler.ProjectWorkspaceEntryHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project/workspace/", projecthandler.ProjectWorkspaceDetailHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/project/", projecthandler.ProjectDetailHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/sync/cluster/all", synchandler.ClusterAllSyncHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/sync/cluster/", synchandler.ClusterOneSyncHandler(h.svcCtx))
}
