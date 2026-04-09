package project

import (
	"context"
	"strings"

	projectrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type AddProjectClusterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddProjectClusterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddProjectClusterLogic {
	return &AddProjectClusterLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddProjectClusterLogic) AddProjectCluster(req *types.AddProjectClusterRequest) (string, error) {
	req.ClusterUUID = strings.TrimSpace(req.ClusterUUID)
	_, err := l.svcCtx.Project.AddCluster(projectrepo.AddProjectClusterParams{
		ClusterUUID:           req.ClusterUUID,
		ProjectID:             req.ProjectID,
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
	return "项目资源池分配成功", nil
}
