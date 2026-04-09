package project

import (
	"context"
	"errors"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type GetProjectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProjectLogic {
	return &GetProjectLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *GetProjectLogic) GetProject(id uint64) (types.Project, error) {
	if id == 0 {
		return types.Project{}, errors.New("project id is required")
	}
	item, ok, err := l.svcCtx.Project.GetByID(id)
	if err != nil {
		return types.Project{}, err
	}
	if !ok {
		return types.Project{}, errors.New("project not found")
	}
	return types.Project{
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
	}, nil
}
