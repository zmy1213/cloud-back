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

type GetCPUUsageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCPUUsageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCPUUsageLogic {
	return &GetCPUUsageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCPUUsageLogic) GetCPUUsage(req *types.GetCPUUsageRequest) (*types.GetCPUUsageResponse, error) {
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
	window := promutils.CalculateRateWindow(timeRange)

	ns := strings.TrimSpace(req.Namespace)
	pod := strings.TrimSpace(req.PodName)

	usageQuery := fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{namespace="%s",pod="%s",container!="",container!="POD"}[%s]))`, ns, pod, window)
	requestQuery := fmt.Sprintf(`sum(kube_pod_container_resource_requests{namespace="%s",pod="%s",resource="cpu"})`, ns, pod)
	limitQuery := fmt.Sprintf(`sum(kube_pod_container_resource_limits{namespace="%s",pod="%s",resource="cpu"})`, ns, pod)
	throttleQuery := fmt.Sprintf(`sum(rate(container_cpu_cfs_throttled_seconds_total{namespace="%s",pod="%s",container!="",container!="POD"}[%s]))`, ns, pod, window)

	var usage, requestCores, limitCores, throttle float64
	var timestamp time.Time

	if results, qErr := client.InstantQuery(l.ctx, usageQuery, nil); qErr == nil && len(results) > 0 {
		usage = results[0].Value
		timestamp = results[0].Time
	}
	if results, qErr := client.InstantQuery(l.ctx, requestQuery, nil); qErr == nil && len(results) > 0 {
		requestCores = results[0].Value
	}
	if results, qErr := client.InstantQuery(l.ctx, limitQuery, nil); qErr == nil && len(results) > 0 {
		limitCores = results[0].Value
	}
	if results, qErr := client.InstantQuery(l.ctx, throttleQuery, nil); qErr == nil && len(results) > 0 {
		throttle = results[0].Value
	}
	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	usagePercent := 0.0
	if limitCores > 0 {
		usagePercent = (usage / limitCores) * 100
	} else if requestCores > 0 {
		usagePercent = (usage / requestCores) * 100
	}

	resp := &types.GetCPUUsageResponse{
		Namespace: ns,
		PodName:   pod,
		Current: types.PodCPUCurrent{
			Timestamp:     timestamp.Unix(),
			UsageCores:    usage,
			UsagePercent:  usagePercent,
			RequestCores:  requestCores,
			LimitCores:    limitCores,
			ThrottledTime: throttle,
		},
		Trend: make([]types.PodCPUTrendPoint, 0),
	}

	step := strings.TrimSpace(timeRange.Step)
	if step == "" {
		step = promutils.CalculateStep(timeRange)
	}

	if series, qErr := client.RangeQuery(l.ctx, usageQuery, timeRange.Start, timeRange.End, step); qErr == nil && len(series) > 0 {
		denominator := limitCores
		if denominator == 0 {
			denominator = requestCores
		}
		for _, point := range series[0].Values {
			pointPercent := 0.0
			if denominator > 0 {
				pointPercent = (point.Value / denominator) * 100
			}
			resp.Trend = append(resp.Trend, types.PodCPUTrendPoint{
				Timestamp:    point.Time.Unix(),
				UsageCores:   point.Value,
				UsagePercent: pointPercent,
			})
		}
	}

	return resp, nil
}
