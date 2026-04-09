package operator

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	promtypes "github.com/yanshicheng/cloud-back/common/prometheusmanager/types"
)

type PrometheusClient struct {
	base *BaseOperator
}

func NewPrometheusClient(config promtypes.PrometheusConfig) (*PrometheusClient, error) {
	base, err := NewBaseOperator(config)
	if err != nil {
		return nil, err
	}
	return &PrometheusClient{base: base}, nil
}

func (c *PrometheusClient) InstantQuery(ctx context.Context, query string, ts *time.Time) ([]promtypes.InstantSample, error) {
	params := url.Values{}
	params.Set("query", query)
	if ts != nil {
		params.Set("time", fmt.Sprintf("%.3f", float64(ts.Unix())+float64(ts.Nanosecond())/1e9))
	}

	var respBody struct {
		Status string `json:"status"`
		Data   struct {
			Result []struct {
				Metric map[string]string `json:"metric"`
				Value  []interface{}     `json:"value"`
			} `json:"result"`
		} `json:"data"`
		Error string `json:"error"`
	}

	if err := c.base.DoRequest(ctx, "/api/v1/query", params, &respBody); err != nil {
		return nil, err
	}
	if respBody.Status != "success" {
		return nil, fmt.Errorf("prometheus query failed: %s", respBody.Error)
	}

	out := make([]promtypes.InstantSample, 0, len(respBody.Data.Result))
	for _, item := range respBody.Data.Result {
		if len(item.Value) != 2 {
			continue
		}
		tsFloat, ok := item.Value[0].(float64)
		if !ok {
			continue
		}
		value, ok := parsePromValue(item.Value[1])
		if !ok {
			continue
		}
		out = append(out, promtypes.InstantSample{
			Metric: item.Metric,
			Value:  value,
			Time:   parsePromTime(tsFloat),
		})
	}
	return out, nil
}

func (c *PrometheusClient) RangeQuery(ctx context.Context, query string, start, end time.Time, step string) ([]promtypes.RangeSeries, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("start", fmt.Sprintf("%.3f", float64(start.Unix())+float64(start.Nanosecond())/1e9))
	params.Set("end", fmt.Sprintf("%.3f", float64(end.Unix())+float64(end.Nanosecond())/1e9))
	params.Set("step", step)

	var respBody struct {
		Status string `json:"status"`
		Data   struct {
			Result []struct {
				Metric map[string]string `json:"metric"`
				Values [][]interface{}   `json:"values"`
			} `json:"result"`
		} `json:"data"`
		Error string `json:"error"`
	}

	if err := c.base.DoRequest(ctx, "/api/v1/query_range", params, &respBody); err != nil {
		return nil, err
	}
	if respBody.Status != "success" {
		return nil, fmt.Errorf("prometheus range query failed: %s", respBody.Error)
	}

	series := make([]promtypes.RangeSeries, 0, len(respBody.Data.Result))
	for _, item := range respBody.Data.Result {
		points := make([]promtypes.RangePoint, 0, len(item.Values))
		for _, point := range item.Values {
			if len(point) != 2 {
				continue
			}
			tsFloat, ok := point[0].(float64)
			if !ok {
				continue
			}
			value, ok := parsePromValue(point[1])
			if !ok {
				continue
			}
			points = append(points, promtypes.RangePoint{
				Time:  parsePromTime(tsFloat),
				Value: value,
			})
		}
		series = append(series, promtypes.RangeSeries{
			Metric: item.Metric,
			Values: points,
		})
	}
	return series, nil
}

func parsePromValue(raw interface{}) (float64, bool) {
	switch v := raw.(type) {
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, false
		}
		return f, true
	case float64:
		return v, true
	case int64:
		return float64(v), true
	case int:
		return float64(v), true
	default:
		return 0, false
	}
}

func parsePromTime(ts float64) time.Time {
	sec := int64(ts)
	nsec := int64((ts - float64(sec)) * 1e9)
	return time.Unix(sec, nsec)
}
