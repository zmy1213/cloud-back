package types

type DashboardOverviewRequest struct {
	Username    string `form:"username"`
	ClusterUUID string `form:"clusterUuid"`
}

type DashboardStat struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Value   string `json:"value"`
	Trend   string `json:"trend"`
	Up      bool   `json:"up"`
	Icon    string `json:"icon"`
	Caption string `json:"caption"`
}

type DashboardAction struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Route       string `json:"route"`
}

type DashboardActivity struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	Timestamp int64  `json:"timestamp"`
}

type DashboardOverviewResponse struct {
	WelcomeName string              `json:"welcomeName"`
	Stats       []DashboardStat     `json:"stats"`
	Actions     []DashboardAction   `json:"actions"`
	Activities  []DashboardActivity `json:"activities"`
}
