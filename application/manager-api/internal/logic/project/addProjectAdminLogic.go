package project

import (
	"context"
	"errors"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type AddProjectAdminLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddProjectAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProjectAdminLogic {
	return &AddProjectAdminLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddProjectAdminLogic) AddProjectAdmin(req *types.AddProjectAdminRequest) (string, error) {
	if req.ProjectID == 0 {
		return "", errors.New("projectId is required")
	}
	if len(req.UserIDs) > 100 {
		return "", errors.New("userIds length cannot exceed 100")
	}
	if err := l.svcCtx.Project.AddAdmins(req.ProjectID, req.UserIDs); err != nil {
		return "", err
	}
	return "项目成员分配成功", nil
}
