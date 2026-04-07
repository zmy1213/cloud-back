package types

type SyncClusterRequest struct {
	ID uint64 `path:"id"`
}

type SyncClusterResponse struct {
	Message     string `json:"message"`
	ClusterID   uint64 `json:"clusterId"`
	ClusterUUID string `json:"clusterUuid"`
	ClusterName string `json:"clusterName"`
	NodeCount   int    `json:"nodeCount"`
	Source      string `json:"source"`
}
