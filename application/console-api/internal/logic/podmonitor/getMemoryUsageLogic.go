package podmonitor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/yanshicheng/cloud-back/application/console-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/types"
	promutils "github.com/yanshicheng/cloud-back/common/prometheusmanager/utils"
)

type GetMemoryUsageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMemoryUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMemoryUsageLogic {
	return &GetMemoryUsageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMemoryUsageLogic) GetMemoryUsage(req *types.GetMemoryUsageRequest) (*types.GetMemoryUsageResponse, error) {
	if strings.TrimSpace(req.ClusterUuid) == "" {
		return nil, errors.New("clusterUuid is required")
	}
	if strings.TrimSpace(req.Namespace) == "" {
		return nil, errors.New("namespace is required")
	}
	if strings.TrimSpace(req.PodName) == "" {
		return nil, errors.New("podName is required")
	}

	client, err := l.svcCtx.PrometheusManager.Get(req.ClusterUuid)
	if err != nil {
		return nil, err
	}

	timeRange := promutils.ParseTimeRange(req.Start, req.End, req.Step)
	ns := strings.TrimSpace(req.Namespace)
	pod := strings.TrimSpace(req.PodName)

	usageQuery := fmt.Sprintf(`sum(container_memory_working_set_bytes{namespace="%s",pod="%s",container!="",container!="POD"})`, ns, pod)
	requestQuery := fmt.Sprintf(`sum(kube_pod_container_resource_requests{namespace="%s",pod="%s",resource="memory"})`, ns, pod)
	limitQuery := fmt.Sprintf(`sum(kube_pod_container_resource_limits{namespace="%s",pod="%s",resource="memory"})`, ns, pod)

	var usage, requestBytes, limitBytes float64
	var timestamp time.Time

	if results, qErr := client.InstantQuery(l.ctx, usageQuery, nil); qErr == nil && len(results) > 0 {
		usage = results[0].Value
		timestamp = results[0].Time
	}
	if results, qErr := client.InstantQuery(l.ctx, requestQuery, nil); qErr == nil && len(results) > 0 {
		requestBytes = results[0].Value
	}
	if results, qErr := client.InstantQuery(l.ctx, limitQuery, nil); qErr == nil && len(results) > 0 {
		limitBytes = results[0].Value
	}
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	usagePercent := 0.0
	if limitBytes > 0 {
		usagePercent = (usage / limitBytes) * 100
	} else if requestBytes > 0 {
		usagePercent = (usage / requestBytes) * 100
	}

	resp := &types.GetMemoryUsageResponse{
		Namespace: ns,
		PodName:   pod,
		Current: types.PodMemoryCurrent{
			Timestamp:    timestamp.Unix(),
			UsageBytes:   int64(usage),
			UsagePercent: usagePercent,
			RequestBytes: int64(requestBytes),
			LimitBytes:   int64(limitBytes),
		},
		Trend: make([]types.PodMemoryTrendPoint, 0),
	}

	step := strings.TrimSpace(timeRange.Step)
	if step == "" {
		step = promutils.CalculateStep(timeRange)
	}
	if series, qErr := client.RangeQuery(l.ctx, usageQuery, timeRange.Start, timeRange.End, step); qErr == nil && len(series) > 0 {
		denominator := limitBytes
		if denominator == 0 {
			denominator = requestBytes
		}
		for _, point := range series[0].Values {
			pointPercent := 0.0
			if denominator > 0 {
				pointPercent = (point.Value / denominator) * 100
			}
			resp.Trend = append(resp.Trend, types.PodMemoryTrendPoint{
				Timestamp:    point.Time.Unix(),
				UsageBytes:   int64(point.Value),
				UsagePercent: pointPercent,
			})
		}
	}

	return resp, nil
}
