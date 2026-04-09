package project

import (
	"context"
	"strings"

	projectrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type SearchProjectClusterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchProjectClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProjectClusterLogic {
	return &SearchProjectClusterLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SearchProjectClusterLogic) SearchProjectCluster(req *types.SearchProjectClusterRequest) ([]types.ProjectCluster, error) {
	// projectId=0 means searching cluster quotas across all projects.
	req.ClusterUUID = strings.TrimSpace(req.ClusterUUID)

	items, err := l.svcCtx.Project.SearchClusters(projectrepo.SearchProjectClusterParams{
		ProjectID:   req.ProjectID,
		ClusterUUID: req.ClusterUUID,
	})
	if err != nil {
		return nil, err
	}

	resp := make([]types.ProjectCluster, 0, len(items))
	for _, item := range items {
		resp = append(resp, types.ProjectCluster{
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
		})
	}
	return resp, nil
}
