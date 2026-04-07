package sync

import (
	"context"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type ClusterOneSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClusterOneSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClusterOneSyncLogic {
	return &ClusterOneSyncLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ClusterOneSyncLogic) ClusterOneSync(req *types.SyncClusterRequest) (*types.SyncClusterResponse, error) {
	result, err := l.svcCtx.Sync.SyncByID(l.ctx, req.ID, "system")
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
