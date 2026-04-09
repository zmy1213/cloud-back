package project

import (
	"context"
	"errors"

	projectrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type UpdateProjectClusterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProjectClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProjectClusterLogic {
	return &UpdateProjectClusterLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateProjectClusterLogic) UpdateProjectCluster(req *types.UpdateProjectClusterRequest) (string, error) {
	if req.ID == 0 {
		return "", errors.New("id is required")
	}
	err := l.svcCtx.Project.UpdateCluster(projectrepo.UpdateProjectClusterParams{
		ID:                    req.ID,
		CPULimit:              req.CPULimit,
		CPUOvercommitRatio:    req.CPUOvercommitRatio,
		CPUCapacity:           req.CPUCapacity,
		MemLimit:              req.MemLimit,
		MemOvercommitRatio:    req.MemOvercommitRatio,
		MemCapacity:           req.MemCapacity,
		StorageLimit:          req.StorageLimit,
		GPULimit:              req.GPULimit,
		GPUOvercommitRatio:    req.GPUOvercommitRatio,
		GPUCapacity:           req.GPUCapacity,
		PodsLimit:             req.PodsLimit,
		ConfigmapLimit:        req.ConfigmapLimit,
		SecretLimit:           req.SecretLimit,
		PVCLimit:              req.PVCLimit,
		EphemeralStorageLimit: req.EphemeralStorageLimit,
		ServiceLimit:          req.ServiceLimit,
		LoadbalancersLimit:    req.LoadbalancersLimit,
		NodeportsLimit:        req.NodeportsLimit,
		DeploymentsLimit:      req.DeploymentsLimit,
		JobsLimit:             req.JobsLimit,
		CronjobsLimit:         req.CronjobsLimit,
		DaemonsetsLimit:       req.DaemonsetsLimit,
		StatefulsetsLimit:     req.StatefulsetsLimit,
		IngressesLimit:        req.IngressesLimit,
		Operator:              "system",
	})
	if err != nil {
		return "", err
	}
	return "项目资源池更新成功", nil
}
