package project

import (
	"context"
	"errors"
	"strings"

	projectrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type UpdateProjectWorkspaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProjectWorkspaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProjectWorkspaceLogic {
	return &UpdateProjectWorkspaceLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateProjectWorkspaceLogic) UpdateProjectWorkspace(req *types.UpdateProjectWorkspaceRequest) (string, error) {
	if req.ID == 0 {
		return "", errors.New("id is required")
	}
	if err := l.svcCtx.Project.UpdateWorkspace(projectrepo.UpdateProjectWorkspaceParams{
		ID:               req.ID,
		Name:             strings.TrimSpace(req.Name),
		Description:      strings.TrimSpace(req.Description),
		CPUAllocated:     req.CPUAllocated,
		MemAllocated:     req.MemAllocated,
		StorageAllocated: req.StorageAllocated,
		GPUAllocated:     req.GPUAllocated,
		PodsAllocated:    req.PodsAllocated,
		Operator:         "system",
	}); err != nil {
		return "", err
	}
	return "工作空间更新成功", nil
}
