package project

import (
	"context"

	projectrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type SearchProjectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchProjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProjectLogic {
	return &SearchProjectLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SearchProjectLogic) SearchProject(req *types.SearchProjectRequest) (*types.SearchProjectResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	items, total, err := l.svcCtx.Project.Search(projectrepo.SearchParams{
		Page:     req.Page,
		PageSize: req.PageSize,
		Name:     req.Name,
		UUID:     req.Uuid,
	})
	if err != nil {
		return nil, err
	}

	resp := &types.SearchProjectResponse{
		Items: make([]types.Project, 0, len(items)),
		Total: total,
	}
	for _, item := range items {
		resp.Items = append(resp.Items, types.Project{
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
	return resp, nil
}
