package nodemonitor

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/yanshicheng/cloud-back/application/console-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/console-api/internal/types"
	promutils "github.com/yanshicheng/cloud-back/common/prometheusmanager/utils"
)

type GetNodeCPULogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNodeCPULogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNodeCPULogic {
	return &GetNodeCPULogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNodeCPULogic) GetNodeCPU(req *types.GetNodeCPURequest) (*types.GetNodeCPUResponse, error) {
	if strings.TrimSpace(req.ClusterUuid) == "" {
		return nil, errors.New("clusterUuid is required")
	}
	if strings.TrimSpace(req.NodeName) == "" {
		return nil, errors.New("nodeName is required")
	}

	client, err := l.svcCtx.PrometheusManager.Get(req.ClusterUuid)
	if err != nil {
		return nil, err
	}

	timeRange := promutils.ParseTimeRange(req.Start, req.End, req.Step)
	window := promutils.CalculateRateWindow(timeRange)
	nodePattern := regexp.QuoteMeta(strings.TrimSpace(req.NodeName))

	coresQuery := fmt.Sprintf(`count(node_cpu_seconds_total{instance=~".*%s.*",mode="idle"})`, nodePattern)
	usageQuery := fmt.Sprintf(`100 - (avg(irate(node_cpu_seconds_total{instance=~".*%s.*",mode="idle"}[%s])) * 100)`, nodePattern, window)
	userQuery := fmt.Sprintf(`avg(irate(node_cpu_seconds_total{instance=~".*%s.*",mode="user"}[%s])) * 100`, nodePattern, window)
	systemQuery := fmt.Sprintf(`avg(irate(node_cpu_seconds_total{instance=~".*%s.*",mode="system"}[%s])) * 100`, nodePattern, window)

	var current types.NodeCPUCurrent
	current.Timestamp = time.Now().Unix()

	if results, qErr := client.InstantQuery(l.ctx, coresQuery, nil); qErr == nil && len(results) > 0 {
		current.TotalCores = int64(results[0].Value)
	}
	if results, qErr := client.InstantQuery(l.ctx, usageQuery, nil); qErr == nil && len(results) > 0 {
		current.UsagePercent = results[0].Value
		current.Timestamp = results[0].Time.Unix()
	}
	if results, qErr := client.InstantQuery(l.ctx, userQuery, nil); qErr == nil && len(results) > 0 {
		current.UserPercent = results[0].Value
	}
	if results, qErr := client.InstantQuery(l.ctx, systemQuery, nil); qErr == nil && len(results) > 0 {
		current.SystemPercent = results[0].Value
	}

	resp := &types.GetNodeCPUResponse{
		NodeName: strings.TrimSpace(req.NodeName),
		Current:  current,
		Trend:    make([]types.NodeCPUTrendPoint, 0),
	}

	step := strings.TrimSpace(timeRange.Step)
	if step == "" {
		step = promutils.CalculateStep(timeRange)
	}
	if series, qErr := client.RangeQuery(l.ctx, usageQuery, timeRange.Start, timeRange.End, step); qErr == nil && len(series) > 0 {
		for _, point := range series[0].Values {
			resp.Trend = append(resp.Trend, types.NodeCPUTrendPoint{
				Timestamp:    point.Time.Unix(),
				UsagePercent: point.Value,
			})
		}
	}

	return resp, nil
}
