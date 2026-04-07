package dashboard

import (
	"context"
	"strings"

	"github.com/yanshicheng/cloud-back/application/portal-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/portal-api/internal/types"
)

type DashboardOverviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDashboardOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DashboardOverviewLogic {
	return &DashboardOverviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DashboardOverviewLogic) DashboardOverview(req *types.DashboardOverviewRequest) (*types.DashboardOverviewResponse, bool) {
	scopeName := ""
	clusterCount := l.svcCtx.Cluster.Count()
	clusterUUID := strings.TrimSpace(req.ClusterUUID)
	if clusterUUID != "" {
		c, ok := l.svcCtx.Cluster.GetByUUID(clusterUUID)
		if !ok {
			return nil, false
		}
		scopeName = c.Name
		clusterCount = 1
	}

	resp := l.svcCtx.Dashboard.Overview(req.Username, scopeName, clusterCount)
	stats := make([]types.DashboardStat, 0, len(resp.Stats))
	for _, s := range resp.Stats {
		stats = append(stats, types.DashboardStat{
			ID:      s.ID,
			Title:   s.Title,
			Value:   s.Value,
			Trend:   s.Trend,
			Up:      s.Up,
			Icon:    s.Icon,
			Caption: s.Caption,
		})
	}
	actions := make([]types.DashboardAction, 0, len(resp.Actions))
	for _, a := range resp.Actions {
		actions = append(actions, types.DashboardAction{
			ID:          a.ID,
			Title:       a.Title,
			Description: a.Description,
			Route:       a.Route,
		})
	}
	activities := make([]types.DashboardActivity, 0, len(resp.Activities))
	for _, a := range resp.Activities {
		activities = append(activities, types.DashboardActivity{
			ID:        a.ID,
			Title:     a.Title,
			Detail:    a.Detail,
			Timestamp: a.Timestamp,
		})
	}

	return &types.DashboardOverviewResponse{
		WelcomeName: resp.WelcomeName,
		Stats:       stats,
		Actions:     actions,
		Activities:  activities,
	}, true
}
