package svc

import (
	"errors"
	"fmt"
	"strings"

	apprepo "github.com/yanshicheng/cloud-back/application/console-api/internal/repository/app"
	appcfg "github.com/yanshicheng/cloud-back/common/config"
	promcluster "github.com/yanshicheng/cloud-back/common/prometheusmanager/cluster"
	promtypes "github.com/yanshicheng/cloud-back/common/prometheusmanager/types"
)

type ServiceContext struct {
	Config            appcfg.AppConfig
	App               *apprepo.Service
	PrometheusManager *promcluster.PrometheusManager
}

func NewServiceContext(cfg appcfg.AppConfig) *ServiceContext {
	appService := apprepo.NewService(cfg.Mysql)
	return &ServiceContext{
		Config: cfg,
		App:    appService,
		PrometheusManager: promcluster.NewPrometheusManager(promcluster.ConfigResolverFunc(func(clusterUUID string) (promtypes.PrometheusConfig, error) {
			return resolvePrometheusConfig(appService, clusterUUID)
		})),
	}
}

func resolvePrometheusConfig(appService *apprepo.Service, clusterUUID string) (promtypes.PrometheusConfig, error) {
	clusterUUID = strings.TrimSpace(clusterUUID)
	if clusterUUID == "" {
		return promtypes.PrometheusConfig{}, errors.New("clusterUuid is required")
	}
	if appService == nil {
		return promtypes.PrometheusConfig{}, errors.New("app repository is not initialized")
	}

	apps, err := appService.ListByClusterUUID(clusterUUID)
	if err != nil {
		return promtypes.PrometheusConfig{}, fmt.Errorf("query cluster apps failed: %w", err)
	}

	app, ok := selectPrometheusApp(apps)
	if !ok {
		return promtypes.PrometheusConfig{}, errors.New("prometheus app not configured for cluster")
	}

	return buildPrometheusConfig(app)
}

func selectPrometheusApp(apps []apprepo.ClusterApp) (apprepo.ClusterApp, bool) {
	var selected apprepo.ClusterApp
	found := false
	for i := range apps {
		item := apps[i]
		if strings.ToLower(strings.TrimSpace(item.AppCode)) != "prometheus" {
			continue
		}
		if !found {
			selected = item
			found = true
		}
		if item.IsDefault == 1 {
			selected = item
			break
		}
	}
	return selected, found
}

func buildPrometheusConfig(app apprepo.ClusterApp) (promtypes.PrometheusConfig, error) {
	if strings.TrimSpace(app.AppUrl) == "" {
		return promtypes.PrometheusConfig{}, errors.New("prometheus appUrl is empty")
	}
	if app.Port < 1 || app.Port > 65535 {
		return promtypes.PrometheusConfig{}, errors.New("prometheus port is invalid")
	}

	protocol := strings.ToLower(strings.TrimSpace(app.Protocol))
	if protocol == "" {
		protocol = "http"
	}
	host := strings.TrimSpace(app.AppUrl)
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimSuffix(host, "/")
	endpoint := fmt.Sprintf("%s://%s:%d", protocol, host, app.Port)

	return promtypes.PrometheusConfig{
		Endpoint:           endpoint,
		AuthEnabled:        app.AuthEnabled == 1,
		AuthType:           app.AuthType,
		Username:           app.Username,
		Password:           app.Password,
		Token:              app.Token,
		AccessKey:          app.AccessKey,
		AccessSecret:       app.AccessSecret,
		ClientCert:         app.ClientCert,
		ClientKey:          app.ClientKey,
		CaCert:             app.CaCert,
		CaFile:             app.CaFile,
		InsecureSkipVerify: app.InsecureSkipVerify == 1,
		TimeoutSeconds:     15,
	}, nil
}
