package app

import (
	"context"
	"errors"
	"strings"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type GetClusterAppListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClusterAppListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClusterAppListLogic {
	return &GetClusterAppListLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetClusterAppListLogic) GetClusterAppList(req *types.ClusterAppListRequest) ([]types.ClusterAppDetail, error) {
	req.ClusterUuid = strings.TrimSpace(req.ClusterUuid)
	if req.ClusterUuid == "" {
		return nil, errors.New("clusterUuid is required")
	}

	items, err := l.svcCtx.App.ListByClusterUUID(req.ClusterUuid)
	if err != nil {
		return nil, err
	}

	out := make([]types.ClusterAppDetail, 0, len(items))
	for _, item := range items {
		out = append(out, types.ClusterAppDetail{
			ID:                 item.ID,
			ClusterUuid:        item.ClusterUuid,
			AppName:            item.AppName,
			AppCode:            item.AppCode,
			AppType:            item.AppType,
			IsDefault:          item.IsDefault,
			AppUrl:             item.AppUrl,
			Port:               item.Port,
			Protocol:           item.Protocol,
			AuthEnabled:        item.AuthEnabled,
			AuthType:           item.AuthType,
			Username:           item.Username,
			Password:           item.Password,
			Token:              item.Token,
			AccessKey:          item.AccessKey,
			AccessSecret:       item.AccessSecret,
			TlsEnabled:         item.TlsEnabled,
			CaFile:             item.CaFile,
			CaKey:              item.CaKey,
			CaCert:             item.CaCert,
			ClientCert:         item.ClientCert,
			ClientKey:          item.ClientKey,
			InsecureSkipVerify: item.InsecureSkipVerify,
			Status:             item.Status,
			CreatedBy:          item.CreatedBy,
			UpdatedBy:          item.UpdatedBy,
			CreatedAt:          item.CreatedAt,
			UpdatedAt:          item.UpdatedAt,
		})
	}

	return out, nil
}
