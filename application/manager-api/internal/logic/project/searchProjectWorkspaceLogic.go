package project

import (
	"context"
	"strings"

	projectrepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/project"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type SearchProjectWorkspaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchProjectWorkspaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchProjectWorkspaceLogic {
	return &SearchProjectWorkspaceLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *SearchProjectWorkspaceLogic) SearchProjectWorkspace(req *types.SearchProjectWorkspaceRequest) ([]types.ProjectWorkspace, error) {
	items, err := l.svcCtx.Project.SearchWorkspaces(projectrepo.SearchProjectWorkspaceParams{
		ProjectClusterID: req.ProjectClusterID,
		Name:             strings.TrimSpace(req.Name),
		Namespace:        strings.TrimSpace(req.Namespace),
	})
	if err != nil {
		return nil, err
	}

	resp := make([]types.ProjectWorkspace, 0, len(items))
	for _, item := range items {
		resp = append(resp, mapProjectWorkspace(item))
	}
	return resp, nil
}
