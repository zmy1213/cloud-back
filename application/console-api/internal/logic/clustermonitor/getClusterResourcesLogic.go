package clustermonitor

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/yanshicheng/cloud-back/application/console-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/types"
)

type GetClusterResourcesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClusterResourcesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClusterResourcesLogic {
	return &GetClusterResourcesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetClusterResourcesLogic) GetClusterResources(req *types.GetClusterResourcesRequest) (*types.GetClusterResourcesResponse, error) {
	if strings.TrimSpace(req.ClusterUuid) == "" {
		return nil, errors.New("clusterUuid is required")
	}

	client, err := l.svcCtx.PrometheusManager.Get(req.ClusterUuid)
	if err != nil {
		return nil, err
	}

	querySet := map[string]string{
		"cpuCapacity":        `sum(kube_node_status_capacity{resource="cpu"})`,
		"cpuAllocatable":     `sum(kube_node_status_allocatable{resource="cpu"})`,
		"cpuRequests":        `sum(kube_pod_container_resource_requests{resource="cpu"})`,
		"cpuLimits":          `sum(kube_pod_container_resource_limits{resource="cpu"})`,
		"cpuUsage":           `sum(rate(container_cpu_usage_seconds_total{container!=""}[5m]))`,
		"memoryCapacity":     `sum(kube_node_status_capacity{resource="memory"})`,
		"memoryAllocatable":  `sum(kube_node_status_allocatable{resource="memory"})`,
		"memoryRequests":     `sum(kube_pod_container_resource_requests{resource="memory"})`,
		"memoryLimits":       `sum(kube_pod_container_resource_limits{resource="memory"})`,
		"memoryUsage":        `sum(container_memory_working_set_bytes{container!=""})`,
		"storageCapacity":    `sum(kube_node_status_capacity{resource=~"ephemeral[-_]storage"})`,
		"storageAllocatable": `sum(kube_node_status_allocatable{resource=~"ephemeral[-_]storage"})`,
		"storageRequests":    `sum(kube_pod_container_resource_requests{resource=~"ephemeral[-_]storage"})`,
		"storageLimits":      `sum(kube_pod_container_resource_limits{resource=~"ephemeral[-_]storage"})`,
		"storageUsage":       `sum(container_fs_usage_bytes{container!="",pod!=""})`,
		"gpuCapacity":        `sum(kube_node_status_capacity{resource=~".*gpu.*"})`,
		"gpuAllocatable":     `sum(kube_node_status_allocatable{resource=~".*gpu.*"})`,
		"gpuRequests":        `sum(kube_pod_container_resource_requests{resource=~".*gpu.*"})`,
		"gpuLimits":          `sum(kube_pod_container_resource_limits{resource=~".*gpu.*"})`,
		"podsRunning":        `sum(kube_pod_status_phase{phase="Running"})`,
		"podsCapacity":       `sum(kube_node_status_capacity{resource="pods"})`,
	}

	values := make(map[string]float64, len(querySet))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for key, promQL := range querySet {
		wg.Add(1)
		go func(k string, q string) {
			defer wg.Done()
			if results, qErr := client.InstantQuery(l.ctx, q, nil); qErr == nil && len(results) > 0 {
				mu.Lock()
				values[k] = results[0].Value
				mu.Unlock()
			}
		}(key, promQL)
	}
	wg.Wait()

	resp := &types.GetClusterResourcesResponse{
		ClusterUuid: strings.TrimSpace(req.ClusterUuid),
		Timestamp:   time.Now().Unix(),
		CPU: types.ClusterCPUResource{
			Capacity:          values["cpuCapacity"],
			Allocatable:       values["cpuAllocatable"],
			RequestsAllocated: values["cpuRequests"],
			LimitsAllocated:   values["cpuLimits"],
			Usage:             values["cpuUsage"],
		},
		Memory: types.ClusterMemoryResource{
			CapacityBytes:          int64(values["memoryCapacity"]),
			AllocatableBytes:       int64(values["memoryAllocatable"]),
			RequestsAllocatedBytes: int64(values["memoryRequests"]),
			LimitsAllocatedBytes:   int64(values["memoryLimits"]),
			UsageBytes:             int64(values["memoryUsage"]),
		},
		Storage: types.ClusterStorageResource{
			CapacityBytes:          int64(values["storageCapacity"]),
			AllocatableBytes:       int64(values["storageAllocatable"]),
			RequestsAllocatedBytes: int64(values["storageRequests"]),
			LimitsAllocatedBytes:   int64(values["storageLimits"]),
			UsageBytes:             int64(values["storageUsage"]),
		},
		GPU: types.ClusterGPUResource{
			Capacity:          values["gpuCapacity"],
			Allocatable:       values["gpuAllocatable"],
			RequestsAllocated: values["gpuRequests"],
			LimitsAllocated:   values["gpuLimits"],
			Usage:             0,
		},
		Pods: types.ClusterPodsResource{
			Running:  int64(values["podsRunning"]),
			Capacity: int64(values["podsCapacity"]),
		},
	}

	if resp.CPU.Allocatable > 0 {
		resp.CPU.RequestsPercent = (resp.CPU.RequestsAllocated / resp.CPU.Allocatable) * 100
		resp.CPU.UsagePercent = (resp.CPU.Usage / resp.CPU.Allocatable) * 100
	}
	if resp.Memory.AllocatableBytes > 0 {
		resp.Memory.RequestsPercent = (float64(resp.Memory.RequestsAllocatedBytes) / float64(resp.Memory.AllocatableBytes)) * 100
		resp.Memory.UsagePercent = (float64(resp.Memory.UsageBytes) / float64(resp.Memory.AllocatableBytes)) * 100
	}
	if resp.Storage.AllocatableBytes > 0 {
		resp.Storage.RequestsPercent = (float64(resp.Storage.RequestsAllocatedBytes) / float64(resp.Storage.AllocatableBytes)) * 100
		resp.Storage.UsagePercent = (float64(resp.Storage.UsageBytes) / float64(resp.Storage.AllocatableBytes)) * 100
	}
	if resp.GPU.Allocatable > 0 {
		resp.GPU.RequestsPercent = (resp.GPU.RequestsAllocated / resp.GPU.Allocatable) * 100
		resp.GPU.UsagePercent = (resp.GPU.Usage / resp.GPU.Allocatable) * 100
	}
	if resp.Pods.Capacity > 0 {
		resp.Pods.UsagePercent = (float64(resp.Pods.Running) / float64(resp.Pods.Capacity)) * 100
	}

	return resp, nil
}
