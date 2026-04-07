package auth

import (
	"context"

	"github.com/yanshicheng/cloud-back/application/portal-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/types"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	resp, err := l.svcCtx.Auth.Login(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		UserID:   resp.UserID,
		Username: resp.Username,
		NickName: resp.NickName,
		UUID:     resp.UUID,
		Roles:    resp.Roles,
		Token: types.TokenResponse{
			AccessToken:      resp.Token.AccessToken,
			AccessExpiresIn:  resp.Token.AccessExpiresIn,
			RefreshToken:     resp.Token.RefreshToken,
			RefreshExpiresIn: resp.Token.RefreshExpiresIn,
		},
	}, nil
}
