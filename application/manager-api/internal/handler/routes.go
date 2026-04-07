package handler

import (
	"net/http"

	clusterhandler "github.com/yanshicheng/cloud-back/application/manager-api/internal/handler/cluster"
	healthzhandler "github.com/yanshicheng/cloud-back/application/manager-api/internal/handler/healthz"
	nodehandler "github.com/yanshicheng/cloud-back/application/manager-api/internal/handler/node"
)

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthzhandler.HealthzHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/cluster", clusterhandler.SearchClusterHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/cluster/", clusterhandler.GetClusterDetailHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/node", nodehandler.GetNodeListHandler(h.svcCtx))
	mux.HandleFunc("/manager/v1/node/", nodehandler.GetNodeDetailHandler(h.svcCtx))
}
