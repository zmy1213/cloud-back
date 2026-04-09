package project

import (
	"context"
	"errors"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
)

type DeleteProjectWorkspaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProjectWorkspaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProjectWorkspaceLogic {
	return &DeleteProjectWorkspaceLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteProjectWorkspaceLogic) DeleteProjectWorkspace(id uint64) (string, error) {
	if id == 0 {
		return "", errors.New("id is required")
	}
	if err := l.svcCtx.Project.DeleteWorkspace(id); err != nil {
		return "", err
	}
	return "工作空间删除成功", nil
}
