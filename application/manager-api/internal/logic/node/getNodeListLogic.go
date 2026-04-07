package node

import (
	"context"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/node"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type GetNodeListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNodeListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNodeListLogic {
	return &GetNodeListLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetNodeListLogic) GetNodeList(req *types.SearchClusterNodeRequest) *types.SearchClusterNodeResponse {
	items, total := l.svcCtx.Node.Search(noderepo.SearchParams{
		Page:        req.Page,
		PageSize:    req.PageSize,
		OrderField:  req.OrderField,
		IsAsc:       req.IsAsc,
		ClusterUuid: req.ClusterUuid,
	})

	out := make([]types.ClusterNodeInfo, 0, len(items))
	for _, n := range items {
		out = append(out, types.ClusterNodeInfo{
			ID:            n.ID,
			ClusterUuid:   n.ClusterUuid,
			NodeName:      n.Name,
			NodeIp:        n.NodeIp,
			NodeStatus:    n.Status,
			CpuUsge:       usagePercent(n.Cpu),
			MemoryUsge:    usagePercent(n.Memory),
			PodTotal:      n.Pods,
			PodUsge:       podUsage(n.Pods),
			CreatedAt:     n.CreatedAt,
			UpdatedAt:     n.UpdatedAt,
			NodeRole:      n.Roles,
			Architecture:  n.Architecture,
			Unschedulable: n.Unschedulable,
		})
	}

	return &types.SearchClusterNodeResponse{
		Items: out,
		Total: total,
	}
}

func usagePercent(raw float64) float64 {
	if raw < 0 {
		return 0
	}
	if raw > 100 {
		return 100
	}
	return raw
}

func podUsage(total int64) int64 {
	if total <= 0 {
		return 0
	}
	return total / 2
}
