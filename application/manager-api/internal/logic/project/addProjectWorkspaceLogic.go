package project

import (
	"context"
	"strings"

	projectrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type AddProjectWorkspaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddProjectWorkspaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProjectWorkspaceLogic {
	return &AddProjectWorkspaceLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddProjectWorkspaceLogic) AddProjectWorkspace(req *types.AddProjectWorkspaceRequest) (string, error) {
	_, err := l.svcCtx.Project.AddWorkspace(projectrepo.AddProjectWorkspaceParams{
		ProjectClusterID: req.ProjectClusterID,
		Name:             strings.TrimSpace(req.Name),
		Namespace:        strings.TrimSpace(req.Namespace),
		Description:      strings.TrimSpace(req.Description),
		CPUAllocated:     req.CPUAllocated,
		MemAllocated:     req.MemAllocated,
		StorageAllocated: req.StorageAllocated,
		GPUAllocated:     req.GPUAllocated,
		PodsAllocated:    req.PodsAllocated,
		Operator:         "system",
	})
	if err != nil {
		return "", err
	}
	return "工作空间创建成功", nil
}
