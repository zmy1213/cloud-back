package project

import (
	"context"
	"errors"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type GetProjectClusterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProjectClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProjectClusterLogic {
	return &GetProjectClusterLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetProjectClusterLogic) GetProjectCluster(id uint64) (types.ProjectCluster, error) {
	if id == 0 {
		return types.ProjectCluster{}, errors.New("id is required")
	}
	item, ok, err := l.svcCtx.Project.GetClusterByID(id)
	if err != nil {
		return types.ProjectCluster{}, err
	}
	if !ok {
		return types.ProjectCluster{}, errors.New("project cluster quota not found")
	}
	return types.ProjectCluster{
		ID:                        item.ID,
		ClusterUUID:               item.ClusterUUID,
		ClusterName:               item.ClusterName,
		ProjectID:                 item.ProjectID,
		CPULimit:                  item.CPULimit,
		CPUOvercommitRatio:        item.CPUOvercommitRatio,
		CPUCapacity:               item.CPUCapacity,
		CPUAllocated:              item.CPUAllocated,
		MemLimit:                  item.MemLimit,
		MemOvercommitRatio:        item.MemOvercommitRatio,
		MemCapacity:               item.MemCapacity,
		MemAllocated:              item.MemAllocated,
		StorageLimit:              item.StorageLimit,
		StorageAllocated:          item.StorageAllocated,
		GPULimit:                  item.GPULimit,
		GPUOvercommitRatio:        item.GPUOvercommitRatio,
		GPUCapacity:               item.GPUCapacity,
		GPUAllocated:              item.GPUAllocated,
		PodsLimit:                 item.PodsLimit,
		PodsAllocated:             item.PodsAllocated,
		ConfigmapLimit:            item.ConfigmapLimit,
		ConfigmapAllocated:        item.ConfigmapAllocated,
		SecretLimit:               item.SecretLimit,
		SecretAllocated:           item.SecretAllocated,
		PVCLimit:                  item.PVCLimit,
		PVCAllocated:              item.PVCAllocated,
		EphemeralStorageLimit:     item.EphemeralStorageLimit,
		EphemeralStorageAllocated: item.EphemeralStorageAllocated,
		ServiceLimit:              item.ServiceLimit,
		ServiceAllocated:          item.ServiceAllocated,
		LoadbalancersLimit:        item.LoadbalancersLimit,
		LoadbalancersAllocated:    item.LoadbalancersAllocated,
		NodeportsLimit:            item.NodeportsLimit,
		NodeportsAllocated:        item.NodeportsAllocated,
		DeploymentsLimit:          item.DeploymentsLimit,
		DeploymentsAllocated:      item.DeploymentsAllocated,
		JobsLimit:                 item.JobsLimit,
		JobsAllocated:             item.JobsAllocated,
		CronjobsLimit:             item.CronjobsLimit,
		CronjobsAllocated:         item.CronjobsAllocated,
		DaemonsetsLimit:           item.DaemonsetsLimit,
		DaemonsetsAllocated:       item.DaemonsetsAllocated,
		StatefulsetsLimit:         item.StatefulsetsLimit,
		StatefulsetsAllocated:     item.StatefulsetsAllocated,
		IngressesLimit:            item.IngressesLimit,
		IngressesAllocated:        item.IngressesAllocated,
		CreatedBy:                 item.CreatedBy,
		UpdatedBy:                 item.UpdatedBy,
		CreatedAt:                 item.CreatedAt,
		UpdatedAt:                 item.UpdatedAt,
	}, nil
}
