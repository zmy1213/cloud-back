package dashboard

import (
	"fmt"
	"time"
)

// Service provides dashboard data for the portal workbench page.
type Service struct{}

func NewService() *Service {
	return &Service{}
}

type StatCard struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Value   string `json:"value"`
	Trend   string `json:"trend"`
	Up      bool   `json:"up"`
	Icon    string `json:"icon"`
	Caption string `json:"caption"`
}

type QuickAction struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Route       string `json:"route"`
}

type ActivityItem struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	Timestamp int64  `json:"timestamp"`
}

type OverviewResponse struct {
	WelcomeName string         `json:"welcomeName"`
	Stats       []StatCard     `json:"stats"`
	Actions     []QuickAction  `json:"actions"`
	Activities  []ActivityItem `json:"activities"`
}

// Overview returns mock data for the first dashboard version.
// The structure is designed to keep frontend modules stable for future real data wiring.
func (s *Service) Overview(username, scopeName string, clusterCount int) *OverviewResponse {
	now := time.Now()
	name := username
	if name == "" {
		name = "super_admin"
	}
	scopeLabel := "All clusters"
	if scopeName != "" {
		scopeLabel = scopeName
	}
	if clusterCount <= 0 {
		clusterCount = 1
	}

	return &OverviewResponse{
		WelcomeName: name,
		Stats: []StatCard{
			{
				ID: "cluster", Title: "Clusters", Value: fmt.Sprintf("%d", clusterCount), Trend: "+0", Up: true, Icon: "cloud",
				Caption: fmt.Sprintf("Scope: %s", scopeLabel),
			},
			{
				ID: "workload", Title: "Workloads", Value: fmt.Sprintf("%d", 14*clusterCount), Trend: "+8%", Up: true, Icon: "apps",
				Caption: "Compared to last 7 days",
			},
			{
				ID: "alerts", Title: "Alerts", Value: fmt.Sprintf("%d", 2*clusterCount), Trend: "-21%", Up: true, Icon: "notifications",
				Caption: "Critical alerts reduced",
			},
			{
				ID: "members", Title: "Members", Value: fmt.Sprintf("%d", 6*clusterCount), Trend: "+1", Up: true, Icon: "group",
				Caption: "New collaborators",
			},
		},
		Actions: []QuickAction{
			{ID: "create-project", Title: "Create Project", Description: "Bootstrap a new project with baseline resources.", Route: "/console/project/create"},
			{ID: "add-cluster", Title: "Add Cluster", Description: "Register a new Kubernetes cluster endpoint.", Route: "/manager/cluster/create"},
			{ID: "deploy-app", Title: "Deploy App", Description: "Ship an application workload from template.", Route: "/workload/deploy"},
			{ID: "view-alerts", Title: "Open Alerts", Description: "Inspect unresolved incidents and notification status.", Route: "/portal/alert"},
		},
		Activities: []ActivityItem{
			{ID: "a1", Title: "Cluster sync completed", Detail: "production-cn synced in 37s.", Timestamp: now.Add(-12 * time.Minute).Unix()},
			{ID: "a2", Title: "Workload deployed", Detail: "payment-api v1.9.2 rolled out to staging.", Timestamp: now.Add(-46 * time.Minute).Unix()},
			{ID: "a3", Title: "New member invited", Detail: "alice joined workspace cloud-core.", Timestamp: now.Add(-2 * time.Hour).Unix()},
			{ID: "a4", Title: "Policy updated", Detail: "RBAC rule set updated by super_admin.", Timestamp: now.Add(-5 * time.Hour).Unix()},
		},
	}
}
