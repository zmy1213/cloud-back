package app

import (
	"context"
	"errors"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type GetClusterAppDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetClusterAppDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetClusterAppDetailLogic {
	return &GetClusterAppDetailLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetClusterAppDetailLogic) GetClusterAppDetail(req *types.ClusterAppDetailRequest) (*types.ClusterAppDetail, error) {
	if req.ID == 0 {
		return nil, errors.New("id is required")
	}

	item, ok, err := l.svcCtx.App.GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("app not found")
	}

	return &types.ClusterAppDetail{
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
	}, nil
}
