package project

import (
	"context"
	"errors"

	projectrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type GetProjectWorkspaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProjectWorkspaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProjectWorkspaceLogic {
	return &GetProjectWorkspaceLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetProjectWorkspaceLogic) GetProjectWorkspace(id uint64) (types.ProjectWorkspace, error) {
	if id == 0 {
		return types.ProjectWorkspace{}, errors.New("id is required")
	}
	item, ok, err := l.svcCtx.Project.GetWorkspaceByID(id)
	if err != nil {
		return types.ProjectWorkspace{}, err
	}
	if !ok {
		return types.ProjectWorkspace{}, errors.New("workspace not found")
	}
	return mapProjectWorkspace(item), nil
}

func mapProjectWorkspace(item projectrepo.ProjectWorkspace) types.ProjectWorkspace {
	return types.ProjectWorkspace{
		ID:               item.ID,
		ProjectClusterID: item.ProjectClusterID,
		ProjectID:        item.ProjectID,
		ClusterUUID:      item.ClusterUUID,
		ClusterName:      item.ClusterName,
		Name:             item.Name,
		Namespace:        item.Namespace,
		Description:      item.Description,
		CPUAllocated:     item.CPUAllocated,
		MemAllocated:     item.MemAllocated,
		StorageAllocated: item.StorageAllocated,
		GPUAllocated:     item.GPUAllocated,
		PodsAllocated:    item.PodsAllocated,
		CreatedBy:        item.CreatedBy,
		UpdatedBy:        item.UpdatedBy,
		CreatedAt:        item.CreatedAt,
		UpdatedAt:        item.UpdatedAt,
	}
}
