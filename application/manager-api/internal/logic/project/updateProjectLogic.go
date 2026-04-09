package project

import (
	"context"
	"errors"
	"strings"

	projectrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type UpdateProjectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProjectLogic {
	return &UpdateProjectLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateProjectLogic) UpdateProject(req *types.UpdateProjectRequest) (string, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	if req.ID == 0 {
		return "", errors.New("project id is required")
	}
	if req.Name == "" {
		return "", errors.New("project name is required")
	}

	err := l.svcCtx.Project.Update(projectrepo.UpdateProjectParams{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Operator:    "system",
	})
	if err != nil {
		return "", err
	}
	return "项目更新成功", nil
}
