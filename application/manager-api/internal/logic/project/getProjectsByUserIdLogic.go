package project

import (
	"context"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type GetProjectsByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProjectsByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProjectsByUserIdLogic {
	return &GetProjectsByUserIdLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetProjectsByUserIdLogic) GetProjectsByUserId(req *types.GetProjectsByUserIdRequest) ([]types.Project, error) {
	items, err := l.svcCtx.Project.GetByUserID(req.UserID, req.Name)
	if err != nil {
		return nil, err
	}
	out := make([]types.Project, 0, len(items))
	for _, item := range items {
		out = append(out, types.Project{
			ID:            item.ID,
			Name:          item.Name,
			Uuid:          item.UUID,
			Description:   item.Description,
			IsSystem:      item.IsSystem,
			CreatedBy:     item.CreatedBy,
			UpdatedBy:     item.UpdatedBy,
			CreatedAt:     item.CreatedAt,
			UpdatedAt:     item.UpdatedAt,
			AdminCount:    item.AdminCount,
			ResourceCount: item.ResourceCount,
		})
	}
	return out, nil
}
