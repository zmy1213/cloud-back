package node

import (
	"context"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type GetNodeDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNodeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNodeDetailLogic {
	return &GetNodeDetailLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetNodeDetailLogic) GetNodeDetail(req *types.NodeIdRequest) (types.ClusterNodeDetail, bool) {
	n, ok := l.svcCtx.Node.GetByID(req.ID)
	if !ok {
		return types.ClusterNodeDetail{}, false
	}

	return types.ClusterNodeDetail{
		ID:              n.ID,
		ClusterUuid:     n.ClusterUuid,
		NodeUuid:        n.NodeUuid,
		Name:            n.Name,
		Hostname:        n.Hostname,
		Roles:           n.Roles,
		OsImage:         n.OsImage,
		NodeIp:          n.NodeIp,
		KernelVersion:   n.KernelVersion,
		OperatingSystem: n.OperatingSystem,
		Architecture:    n.Architecture,
		Cpu:             n.Cpu,
		Memory:          n.Memory,
		Pods:            n.Pods,
		IsGpu:           n.IsGpu,
		Runtime:         n.Runtime,
		JoinAt:          n.JoinAt,
		Unschedulable:   n.Unschedulable,
		KubeletVersion:  n.KubeletVersion,
		Status:          n.Status,
		PodCidr:         n.PodCidr,
		PodCidrs:        n.PodCidrs,
		CreatedBy:       n.CreatedBy,
		UpdatedBy:       n.UpdatedBy,
		CreatedAt:       n.CreatedAt,
		UpdatedAt:       n.UpdatedAt,
	}, true
}
