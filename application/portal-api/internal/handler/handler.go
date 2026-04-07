package handler

import (
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/svc"
	appcfg "github.com/yanshicheng/cloud-back/common/config"
)

type Handler struct {
	svcCtx *svc.ServiceContext
}

func New(cfg appcfg.AppConfig) *Handler {
	return &Handler{
		svcCtx: svc.NewServiceContext(cfg),
	}
}
