package clustermonitor

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/yanshicheng/cloud-back/application/console-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/types"
)

type GetClusterOverviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClusterOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClusterOverviewLogic {
	return &GetClusterOverviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetClusterOverviewLogic) GetClusterOverview(req *types.GetClusterOverviewRequest) (*types.GetClusterOverviewResponse, error) {
	if strings.TrimSpace(req.ClusterUuid) == "" {
		return nil, errors.New("clusterUuid is required")
	}

	client, err := l.svcCtx.PrometheusManager.Get(req.ClusterUuid)
	if err != nil {
		return nil, err
	}

	queries := map[string]string{
		"nodeTotal":         `count(kube_node_info)`,
		"nodeReady":         `sum(kube_node_status_condition{condition="Ready",status="true"})`,
		"podRunning":        `sum(kube_pod_status_phase{phase="Running"})`,
		"podCapacity":       `sum(kube_node_status_capacity{resource="pods"})`,
		"cpuUsage":          `sum(rate(container_cpu_usage_seconds_total{container!=""}[5m]))`,
		"cpuAllocatable":    `sum(kube_node_status_allocatable{resource="cpu"})`,
		"memoryUsage":       `sum(container_memory_working_set_bytes{container!=""})`,
		"memoryAllocatable": `sum(kube_node_status_allocatable{resource="memory"})`,
	}

	values := make(map[string]float64, len(queries))
	for key, query := range queries {
		if results, qErr := client.InstantQuery(l.ctx, query, nil); qErr == nil && len(results) > 0 {
			values[key] = results[0].Value
		}
	}

	resp := &types.GetClusterOverviewResponse{
		ClusterUuid: strings.TrimSpace(req.ClusterUuid),
		Timestamp:   time.Now().Unix(),
		NodeTotal:   int64(values["nodeTotal"]),
		NodeReady:   int64(values["nodeReady"]),
		PodRunning:  int64(values["podRunning"]),
		PodCapacity: int64(values["podCapacity"]),
	}
	if values["cpuAllocatable"] > 0 {
		resp.CPUUsagePercent = (values["cpuUsage"] / values["cpuAllocatable"]) * 100
	}
	if values["memoryAllocatable"] > 0 {
		resp.MemoryUsagePercent = (values["memoryUsage"] / values["memoryAllocatable"]) * 100
	}

	return resp, nil
}
