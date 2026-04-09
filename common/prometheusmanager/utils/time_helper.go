package utils

import (
	"errors"
	"strings"
	"time"

	promtypes "github.com/yanshicheng/cloud-back/common/prometheusmanager/types"
)

func ParseTimeRange(startRaw, endRaw, stepRaw string) promtypes.TimeRange {
	end := time.Now()
	if v, err := parseRFC3339(strings.TrimSpace(endRaw)); err == nil {
		end = v
	}

	start := end.Add(-1 * time.Hour)
	if v, err := parseRFC3339(strings.TrimSpace(startRaw)); err == nil {
		start = v
	}
	if !start.Before(end) {
		end = time.Now()
		start = end.Add(-1 * time.Hour)
	}

	step := strings.TrimSpace(stepRaw)
	if step == "" {
		step = "1m"
	}

	return promtypes.TimeRange{
		Start: start,
		End:   end,
		Step:  step,
	}
}

func parseRFC3339(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, errors.New("empty")
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.999Z07:00",
		"2006-01-02T15:04:05.999Z",
		"2006-01-02T15:04:05Z",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("invalid time format")
}

func CalculateRateWindow(timeRange promtypes.TimeRange) string {
	duration := timeRange.End.Sub(timeRange.Start)
	switch {
	case duration <= 5*time.Minute:
		return "1m"
	case duration <= 30*time.Minute:
		return "5m"
	case duration <= 2*time.Hour:
		return "10m"
	case duration <= 6*time.Hour:
		return "15m"
	case duration <= 24*time.Hour:
		return "30m"
	default:
		return "1h"
	}
}

func CalculateStep(timeRange promtypes.TimeRange) string {
	duration := timeRange.End.Sub(timeRange.Start)
	switch {
	case duration <= 1*time.Hour:
		return "15s"
	case duration <= 6*time.Hour:
		return "30s"
	case duration <= 24*time.Hour:
		return "1m"
	case duration <= 7*24*time.Hour:
		return "5m"
	default:
		return "15m"
	}
}
