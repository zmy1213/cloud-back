package project

import (
	"context"
	"errors"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type GetProjectAdminsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProjectAdminsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProjectAdminsLogic {
	return &GetProjectAdminsLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetProjectAdminsLogic) GetProjectAdmins(req *types.GetProjectAdminsRequest) ([]types.ProjectAdmin, error) {
	if req.ProjectID == 0 {
		return nil, errors.New("projectId is required")
	}
	items, err := l.svcCtx.Project.GetAdmins(req.ProjectID)
	if err != nil {
		return nil, err
	}
	out := make([]types.ProjectAdmin, 0, len(items))
	for _, item := range items {
		out = append(out, types.ProjectAdmin{
			ID:        item.ID,
			ProjectID: item.ProjectID,
			UserID:    item.UserID,
			CreatedAt: item.CreatedAt,
		})
	}
	return out, nil
}
