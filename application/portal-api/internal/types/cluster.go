package types

type SearchClusterRequest struct {
	Name        string `form:"name"`
	Environment string `form:"environment"`
}

type GetClusterDetailRequest struct {
	ID uint64 `path:"id"`
}

type Cluster struct {
	ID           uint64  `json:"id"`
	Name         string  `json:"name"`
	Avatar       string  `json:"avatar"`
	Environment  string  `json:"environment"`
	ClusterType  string  `json:"clusterType"`
	Version      string  `json:"version"`
	Status       int64   `json:"status"`
	HealthStatus int64   `json:"healthStatus"`
	UUID         string  `json:"uuid"`
	CpuUsage     float64 `json:"cpuUsage"`
	MemoryUsage  float64 `json:"memoryUsage"`
	PodUsage     float64 `json:"podUsage"`
	StorageUsage float64 `json:"storageUsage"`
	CreatedAt    int64   `json:"createdAt"`
}

type SearchClusterResponse struct {
	Items []Cluster `json:"items"`
	Total int       `json:"total"`
}
