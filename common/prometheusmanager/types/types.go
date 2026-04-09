package types

import "time"

type TimeRange struct {
	Start time.Time
	End   time.Time
	Step  string
}

type InstantSample struct {
	Metric map[string]string
	Value  float64
	Time   time.Time
}

type RangePoint struct {
	Time  time.Time
	Value float64
}

type RangeSeries struct {
	Metric map[string]string
	Values []RangePoint
}

type PrometheusConfig struct {
	Endpoint           string
	AuthEnabled        bool
	AuthType           string
	Username           string
	Password           string
	Token              string
	AccessKey          string
	AccessSecret       string
	ClientCert         string
	ClientKey          string
	CaCert             string
	CaFile             string
	InsecureSkipVerify bool
	TimeoutSeconds     int
}
