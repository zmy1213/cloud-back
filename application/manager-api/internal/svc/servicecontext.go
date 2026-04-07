package svc

import (
	clusterrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/cluster"
	clustersyncrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/clustersync"
	noderepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/node"
	appcfg "github.com/yanshicheng/cloud-back/common/config"
)

type ServiceContext struct {
	Config  appcfg.AppConfig
	Cluster *clusterrepo.Service
	Node    *noderepo.Service
	Sync    *clustersyncrepo.Service
}

func NewServiceContext(cfg appcfg.AppConfig) *ServiceContext {
	return &ServiceContext{
		Config:  cfg,
		Cluster: clusterrepo.NewService(cfg.Mysql),
		Node:    noderepo.NewService(cfg.Mysql),
		Sync:    clustersyncrepo.NewService(cfg.Mysql, cfg.K8s),
	}
}
