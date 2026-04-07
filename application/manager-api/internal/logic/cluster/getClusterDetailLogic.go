package cluster

import (
	"context"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type GetClusterDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClusterDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClusterDetailLogic {
	return &GetClusterDetailLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetClusterDetailLogic) GetClusterDetail(req *types.GetClusterDetailRequest) (types.Cluster, bool) {
	c, ok := l.svcCtx.Cluster.GetByID(req.ID)
	if !ok {
		return types.Cluster{}, false
	}
	return types.Cluster{
		ID: c.ID, Name: c.Name, Avatar: c.Avatar, Environment: c.Environment,
		ClusterType: c.ClusterType, Version: c.Version, Status: c.Status,
		HealthStatus: c.HealthStatus, UUID: c.UUID, CpuUsage: c.CpuUsage,
		MemoryUsage: c.MemoryUsage, PodUsage: c.PodUsage, StorageUsage: c.StorageUsage,
		CreatedAt: c.CreatedAt,
	}, true
}
