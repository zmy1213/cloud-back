package project

import (
	"context"
	"errors"
	"strings"

	projectrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type AddProjectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddProjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProjectLogic {
	return &AddProjectLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddProjectLogic) AddProject(req *types.AddProjectRequest) (string, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	if req.Name == "" {
		return "", errors.New("project name is required")
	}
	if req.IsSystem != 0 && req.IsSystem != 1 {
		return "", errors.New("isSystem must be 0 or 1")
	}

	_, err := l.svcCtx.Project.Add(projectrepo.AddProjectParams{
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    req.IsSystem,
		Operator:    "system",
	})
	if err != nil {
		return "", err
	}
	return "项目创建成功", nil
}
