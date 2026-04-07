package handler

import (
	"net/http"

	authhandler "github.com/yanshicheng/cloud-back/application/portal-api/internal/handler/auth"
	dashboardhandler "github.com/yanshicheng/cloud-back/application/portal-api/internal/handler/dashboard"
	healthzhandler "github.com/yanshicheng/cloud-back/application/portal-api/internal/handler/healthz"
)

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthzhandler.HealthzHandler(h.svcCtx))
	mux.HandleFunc("/portal/v1/auth/login", authhandler.LoginHandler(h.svcCtx))
	mux.HandleFunc("/portal/v1/dashboard/overview", dashboardhandler.DashboardOverviewHandler(h.svcCtx))
}
