package project

import (
	"context"
	"errors"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
)

type DeleteProjectClusterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProjectClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProjectClusterLogic {
	return &DeleteProjectClusterLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteProjectClusterLogic) DeleteProjectCluster(id uint64) (string, error) {
	if id == 0 {
		return "", errors.New("id is required")
	}
	if err := l.svcCtx.Project.DeleteCluster(id); err != nil {
		return "", err
	}
	return "项目资源池删除成功", nil
}
