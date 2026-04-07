package sync

import (
	"context"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type ClusterAllSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClusterAllSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClusterAllSyncLogic {
	return &ClusterAllSyncLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ClusterAllSyncLogic) ClusterAllSync() (*types.SyncClusterResponse, error) {
	result, err := l.svcCtx.Sync.SyncAll(l.ctx, "system")
	if err != nil {
		return nil, err
	}
	return &types.SyncClusterResponse{
		Message:     "cluster sync completed",
		ClusterID:   result.ClusterID,
		ClusterUUID: result.ClusterUUID,
		ClusterName: result.ClusterName,
		NodeCount:   result.NodeCount,
		Source:      result.Source,
	}, nil
}
