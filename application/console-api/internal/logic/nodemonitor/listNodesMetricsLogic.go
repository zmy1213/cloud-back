package nodemonitor

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/yanshicheng/cloud-back/application/console-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/types"
)

type ListNodesMetricsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListNodesMetricsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNodesMetricsLogic {
	return &ListNodesMetricsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListNodesMetricsLogic) ListNodesMetrics(req *types.ListNodesMetricsRequest) (*types.ListNodesMetricsResponse, error) {
	if strings.TrimSpace(req.ClusterUuid) == "" {
		return nil, errors.New("clusterUuid is required")
	}

	client, err := l.svcCtx.PrometheusManager.Get(req.ClusterUuid)
	if err != nil {
		return nil, err
	}

	nodesQuery := `kube_node_info`
	readyQuery := `kube_node_status_condition{condition="Ready",status="true"}`
	cpuQuery := `100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`
	memQuery := `(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100`

	nodeMap := map[string]*types.NodeMetricItem{}
	ipToNode := map[string]string{}

	if results, qErr := client.InstantQuery(l.ctx, nodesQuery, nil); qErr == nil {
		for _, sample := range results {
			nodeName := strings.TrimSpace(sample.Metric["node"])
			if nodeName == "" {
				continue
			}
			internalIP := strings.TrimSpace(sample.Metric["internal_ip"])
			nodeMap[nodeName] = &types.NodeMetricItem{
				NodeName:   nodeName,
				InternalIP: internalIP,
			}
			if internalIP != "" {
				ipToNode[internalIP] = nodeName
			}
		}
	}

	if results, qErr := client.InstantQuery(l.ctx, readyQuery, nil); qErr == nil {
		for _, sample := range results {
			nodeName := strings.TrimSpace(sample.Metric["node"])
			if nodeName == "" {
				continue
			}
			item, ok := nodeMap[nodeName]
			if !ok {
				item = &types.NodeMetricItem{NodeName: nodeName}
				nodeMap[nodeName] = item
			}
			item.Ready = sample.Value > 0
		}
	}

	assignByInstance := func(item *types.NodeMetricItem, instance string) {
		item.Instance = instance
	}

	resolveNodeBySample := func(metric map[string]string) string {
		if node := strings.TrimSpace(metric["node"]); node != "" {
			return node
		}
		instance := strings.TrimSpace(metric["instance"])
		if instance == "" {
			return ""
		}
		if host, _, found := strings.Cut(instance, ":"); found {
			if node, ok := ipToNode[host]; ok {
				return node
			}
		}
		if node, ok := ipToNode[instance]; ok {
			return node
		}
		return instance
	}
	if results, qErr := client.InstantQuery(l.ctx, cpuQuery, nil); qErr == nil {
		for _, sample := range results {
			nodeName := resolveNodeBySample(sample.Metric)
			if nodeName == "" {
				continue
			}
			item, ok := nodeMap[nodeName]
			if !ok {
				item = &types.NodeMetricItem{NodeName: nodeName}
				nodeMap[nodeName] = item
			}
			item.CPUUsage = sample.Value
			if instance := strings.TrimSpace(sample.Metric["instance"]); instance != "" {
				assignByInstance(item, instance)
			}
		}
	}

	if results, qErr := client.InstantQuery(l.ctx, memQuery, nil); qErr == nil {
		for _, sample := range results {
			nodeName := resolveNodeBySample(sample.Metric)
			if nodeName == "" {
				continue
			}
			item, ok := nodeMap[nodeName]
			if !ok {
				item = &types.NodeMetricItem{NodeName: nodeName}
				nodeMap[nodeName] = item
			}
			item.MemoryUsage = sample.Value
			if instance := strings.TrimSpace(sample.Metric["instance"]); instance != "" && item.Instance == "" {
				assignByInstance(item, instance)
			}
		}
	}
	items := make([]types.NodeMetricItem, 0, len(nodeMap))
	for _, item := range nodeMap {
		items = append(items, *item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].NodeName < items[j].NodeName
	})

	return &types.ListNodesMetricsResponse{
		Items: items,
		Total: len(items),
	}, nil
}
