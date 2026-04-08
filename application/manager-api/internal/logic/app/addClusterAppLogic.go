package app

import (
	"context"
	"fmt"

	apprepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/app"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type AddClusterAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddClusterAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddClusterAppLogic {
	return &AddClusterAppLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AddClusterAppLogic) AddClusterApp(req *types.AddClusterAppRequest) (string, error) {
	sanitizeAddClusterAppRequest(req)
	if err := validateAddClusterAppRequest(req); err != nil {
		return "", err
	}

	err := l.svcCtx.App.Upsert(apprepo.UpsertParams{
		ClusterUuid:        req.ClusterUuid,
		AppName:            req.AppName,
		AppCode:            req.AppCode,
		AppType:            req.AppType,
		IsDefault:          req.IsDefault,
		AppUrl:             req.AppUrl,
		Port:               req.Port,
		Protocol:           req.Protocol,
		AuthEnabled:        req.AuthEnabled,
		AuthType:           req.AuthType,
		Username:           req.Username,
		Password:           req.Password,
		Token:              req.Token,
		AccessKey:          req.AccessKey,
		AccessSecret:       req.AccessSecret,
		TlsEnabled:         req.TlsEnabled,
		CaFile:             req.CaFile,
		CaKey:              req.CaKey,
		CaCert:             req.CaCert,
		ClientCert:         req.ClientCert,
		ClientKey:          req.ClientKey,
		InsecureSkipVerify: req.InsecureSkipVerify,
		UpdatedBy:          req.UpdatedBy,
	})
	if err != nil {
		return "", fmt.Errorf("保存应用配置失败: %w", err)
	}

	return "中间件保存成功", nil
}
