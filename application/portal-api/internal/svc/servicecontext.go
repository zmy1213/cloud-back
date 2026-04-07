package svc

import (
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/auth"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/dashboard"
	clusterrepo "github.com/yanshicheng/cloud-back/application/portal-api/internal/repository/cluster"
	appcfg "github.com/yanshicheng/cloud-back/common/config"
)

// ServiceContext keeps shared dependencies for handlers/logic.
// The shape follows kube-nova's internal/svc pattern.
type ServiceContext struct {
	Config    appcfg.AppConfig
	Auth      *auth.Service
	Dashboard *dashboard.Service
	Cluster   *clusterrepo.Service
}

func NewServiceContext(cfg appcfg.AppConfig) *ServiceContext {
	return &ServiceContext{
		Config:    cfg,
		Auth:      auth.NewService(cfg.Auth.AccessExpiresIn, cfg.Auth.RefreshExpiresIn),
		Dashboard: dashboard.NewService(),
		Cluster:   clusterrepo.NewService(cfg.Mysql),
	}
}
