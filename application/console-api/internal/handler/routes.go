package handler

import (
	"net/http"

	clustermonitorhandler "github.com/yanshicheng/cloud-back/application/console-api/internal/handler/clustermonitor"
	healthzhandler "github.com/yanshicheng/cloud-back/application/console-api/internal/handler/healthz"
	nodemonitorhandler "github.com/yanshicheng/cloud-back/application/console-api/internal/handler/nodemonitor"
	podmonitorhandler "github.com/yanshicheng/cloud-back/application/console-api/internal/handler/podmonitor"
)

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthzhandler.HealthzHandler(h.svcCtx))

	// kube-nova style monitor routes.
	mux.HandleFunc("/console/v1/pod-monitor/cpu/usage", podmonitorhandler.GetCPUUsageHandler(h.svcCtx))
	mux.HandleFunc("/console/v1/pod-monitor/memory/usage", podmonitorhandler.GetMemoryUsageHandler(h.svcCtx))
	mux.HandleFunc("/console/v1/node-monitor/list", nodemonitorhandler.ListNodesMetricsHandler(h.svcCtx))
	mux.HandleFunc("/console/v1/node-monitor/cpu", nodemonitorhandler.GetNodeCPUHandler(h.svcCtx))
	mux.HandleFunc("/console/v1/cluster-monitor/overview", clustermonitorhandler.GetClusterOverviewHandler(h.svcCtx))
	mux.HandleFunc("/console/v1/cluster-monitor/resources", clustermonitorhandler.GetClusterResourcesHandler(h.svcCtx))
}
