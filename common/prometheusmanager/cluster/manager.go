package cluster

import (
	"errors"
	"strings"

	"github.com/yanshicheng/cloud-back/common/prometheusmanager/operator"
	promtypes "github.com/yanshicheng/cloud-back/common/prometheusmanager/types"
)

type ConfigResolver interface {
	Resolve(clusterUUID string) (promtypes.PrometheusConfig, error)
}

type ConfigResolverFunc func(clusterUUID string) (promtypes.PrometheusConfig, error)

func (f ConfigResolverFunc) Resolve(clusterUUID string) (promtypes.PrometheusConfig, error) {
	return f(clusterUUID)
}

type PrometheusManager struct {
	resolver ConfigResolver
}

func NewPrometheusManager(resolver ConfigResolver) *PrometheusManager {
	return &PrometheusManager{
		resolver: resolver,
	}
}

func (m *PrometheusManager) Get(clusterUUID string) (*operator.PrometheusClient, error) {
	clusterUUID = strings.TrimSpace(clusterUUID)
	if clusterUUID == "" {
		return nil, errors.New("clusterUuid is required")
	}
	if m == nil || m.resolver == nil {
		return nil, errors.New("prometheus config resolver is not initialized")
	}

	cfg, err := m.resolver.Resolve(clusterUUID)
	if err != nil {
		return nil, err
	}
	return operator.NewPrometheusClient(cfg)
}
