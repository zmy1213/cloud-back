package types

type SearchClusterNodeRequest struct {
	Page        uint64 `form:"page"`
	PageSize    uint64 `form:"pageSize"`
	OrderField  string `form:"orderField"`
	IsAsc       bool   `form:"isAsc"`
	ClusterUuid string `form:"clusterUuid"`
}

type ClusterNodeInfo struct {
	ID            uint64  `json:"id"`
	ClusterUuid   string  `json:"clusterUuid"`
	NodeName      string  `json:"nodeName"`
	NodeIp        string  `json:"nodeIp"`
	NodeStatus    string  `json:"nodeStatus"`
	CpuUsge       float64 `json:"cpuUsge"`
	MemoryUsge    float64 `json:"memoryUsge"`
	PodTotal      int64   `json:"podTotal"`
	PodUsge       int64   `json:"podUsge"`
	CreatedAt     int64   `json:"createdAt"`
	UpdatedAt     int64   `json:"updatedAt"`
	NodeRole      string  `json:"nodeRole"`
	Architecture  string  `json:"architecture"`
	Unschedulable int64   `json:"unschedulable"`
}

type SearchClusterNodeResponse struct {
	Items []ClusterNodeInfo `json:"items"`
	Total uint64            `json:"total"`
}

type NodeIdRequest struct {
	ID uint64 `path:"id"`
}

type ClusterNodeDetail struct {
	ID              uint64  `json:"id"`
	ClusterUuid     string  `json:"clusterUuid"`
	NodeUuid        string  `json:"nodeUuid"`
	Name            string  `json:"name"`
	Hostname        string  `json:"hostname"`
	Roles           string  `json:"roles"`
	OsImage         string  `json:"osImage"`
	NodeIp          string  `json:"nodeIp"`
	KernelVersion   string  `json:"kernelVersion"`
	OperatingSystem string  `json:"operatingSystem"`
	Architecture    string  `json:"architecture"`
	Cpu             float64 `json:"cpu"`
	Memory          float64 `json:"memory"`
	Pods            int64   `json:"pods"`
	IsGpu           int64   `json:"isGpu"`
	Runtime         string  `json:"runtime"`
	JoinAt          int64   `json:"joinAt"`
	Unschedulable   int64   `json:"unschedulable"`
	KubeletVersion  string  `json:"kubeletVersion"`
	Status          string  `json:"status"`
	PodCidr         string  `json:"podCidr"`
	PodCidrs        string  `json:"podCidrs"`
	CreatedBy       string  `json:"createdBy"`
	UpdatedBy       string  `json:"updatedBy"`
	CreatedAt       int64   `json:"createdAt"`
	UpdatedAt       int64   `json:"updatedAt"`
}
