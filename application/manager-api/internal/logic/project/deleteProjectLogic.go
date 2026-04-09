package project

import (
	"context"
	"errors"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
)

type DeleteProjectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteProjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProjectLogic {
	return &DeleteProjectLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteProjectLogic) DeleteProject(id uint64) (string, error) {
	if id == 0 {
		return "", errors.New("project id is required")
	}
	if err := l.svcCtx.Project.Delete(id); err != nil {
		return "", err
	}
	return "项目删除成功", nil
}
