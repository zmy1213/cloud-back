package cluster

import (
	"context"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type SearchClusterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchClusterLogic {
	return &SearchClusterLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SearchClusterLogic) SearchCluster(req *types.SearchClusterRequest) ([]types.Cluster, int) {
	items := l.svcCtx.Cluster.Search(req.Name, req.Environment)
	out := make([]types.Cluster, 0, len(items))
	for _, c := range items {
		out = append(out, types.Cluster{
			ID: c.ID, Name: c.Name, Avatar: c.Avatar, Environment: c.Environment,
			ClusterType: c.ClusterType, Version: c.Version, Status: c.Status,
			HealthStatus: c.HealthStatus, UUID: c.UUID, CpuUsage: c.CpuUsage,
			MemoryUsage: c.MemoryUsage, PodUsage: c.PodUsage, StorageUsage: c.StorageUsage,
			CreatedAt: c.CreatedAt,
		})
	}
	return out, len(out)
}
