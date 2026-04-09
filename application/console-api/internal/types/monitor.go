package types

type GetPodCPUUsageRequest struct {
	ClusterUuid string `form:"clusterUuid"`
	Namespace   string `form:"namespace"`
	PodName     string `form:"podName"`
	Start       string `form:"start"`
	End         string `form:"end"`
	Step        string `form:"step"`
}

type PodCPUCurrent struct {
	Timestamp     int64   `json:"timestamp"`
	UsageCores    float64 `json:"usageCores"`
	UsagePercent  float64 `json:"usagePercent"`
	RequestCores  float64 `json:"requestCores"`
	LimitCores    float64 `json:"limitCores"`
	ThrottledTime float64 `json:"throttledTime"`
}

type PodCPUTrendPoint struct {
	Timestamp    int64   `json:"timestamp"`
	UsageCores   float64 `json:"usageCores"`
	UsagePercent float64 `json:"usagePercent"`
}

type GetPodCPUUsageResponse struct {
	Namespace string             `json:"namespace"`
	PodName   string             `json:"podName"`
	Current   PodCPUCurrent      `json:"current"`
	Trend     []PodCPUTrendPoint `json:"trend"`
}

type GetPodMemoryUsageRequest struct {
	ClusterUuid string `form:"clusterUuid"`
	Namespace   string `form:"namespace"`
	PodName     string `form:"podName"`
	Start       string `form:"start"`
	End         string `form:"end"`
	Step        string `form:"step"`
}

type PodMemoryCurrent struct {
	Timestamp    int64   `json:"timestamp"`
	UsageBytes   int64   `json:"usageBytes"`
	UsagePercent float64 `json:"usagePercent"`
	RequestBytes int64   `json:"requestBytes"`
	LimitBytes   int64   `json:"limitBytes"`
}

type PodMemoryTrendPoint struct {
	Timestamp    int64   `json:"timestamp"`
	UsageBytes   int64   `json:"usageBytes"`
	UsagePercent float64 `json:"usagePercent"`
}

type GetPodMemoryUsageResponse struct {
	Namespace string                `json:"namespace"`
	PodName   string                `json:"podName"`
	Current   PodMemoryCurrent      `json:"current"`
	Trend     []PodMemoryTrendPoint `json:"trend"`
}

type GetNodeListMetricsRequest struct {
	ClusterUuid string `form:"clusterUuid"`
}

type NodeMetricItem struct {
	NodeName    string  `json:"nodeName"`
	InternalIP  string  `json:"internalIp"`
	Instance    string  `json:"instance"`
	Ready       bool    `json:"ready"`
	CPUUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
}

type GetNodeListMetricsResponse struct {
	Items []NodeMetricItem `json:"items"`
	Total int              `json:"total"`
}

type GetNodeCPUUsageRequest struct {
	ClusterUuid string `form:"clusterUuid"`
	NodeName    string `form:"nodeName"`
	Start       string `form:"start"`
	End         string `form:"end"`
	Step        string `form:"step"`
}

type NodeCPUCurrent struct {
	Timestamp     int64   `json:"timestamp"`
	TotalCores    int64   `json:"totalCores"`
	UsagePercent  float64 `json:"usagePercent"`
	UserPercent   float64 `json:"userPercent"`
	SystemPercent float64 `json:"systemPercent"`
}

type NodeCPUTrendPoint struct {
	Timestamp    int64   `json:"timestamp"`
	UsagePercent float64 `json:"usagePercent"`
}

type GetNodeCPUUsageResponse struct {
	NodeName string              `json:"nodeName"`
	Current  NodeCPUCurrent      `json:"current"`
	Trend    []NodeCPUTrendPoint `json:"trend"`
}

type GetClusterOverviewRequest struct {
	ClusterUuid string `form:"clusterUuid"`
}

type GetClusterOverviewResponse struct {
	ClusterUuid        string  `json:"clusterUuid"`
	Timestamp          int64   `json:"timestamp"`
	NodeTotal          int64   `json:"nodeTotal"`
	NodeReady          int64   `json:"nodeReady"`
	PodRunning         int64   `json:"podRunning"`
	PodCapacity        int64   `json:"podCapacity"`
	CPUUsagePercent    float64 `json:"cpuUsagePercent"`
	MemoryUsagePercent float64 `json:"memoryUsagePercent"`
}

type GetClusterResourcesRequest struct {
	ClusterUuid string `form:"clusterUuid"`
}

type ClusterCPUResource struct {
	Capacity          float64 `json:"capacity"`
	Allocatable       float64 `json:"allocatable"`
	RequestsAllocated float64 `json:"requestsAllocated"`
	LimitsAllocated   float64 `json:"limitsAllocated"`
	Usage             float64 `json:"usage"`
	RequestsPercent   float64 `json:"requestsPercent"`
	UsagePercent      float64 `json:"usagePercent"`
}

type ClusterMemoryResource struct {
	CapacityBytes          int64   `json:"capacityBytes"`
	AllocatableBytes       int64   `json:"allocatableBytes"`
	RequestsAllocatedBytes int64   `json:"requestsAllocatedBytes"`
	LimitsAllocatedBytes   int64   `json:"limitsAllocatedBytes"`
	UsageBytes             int64   `json:"usageBytes"`
	RequestsPercent        float64 `json:"requestsPercent"`
	UsagePercent           float64 `json:"usagePercent"`
}

type ClusterPodsResource struct {
	Running      int64   `json:"running"`
	Capacity     int64   `json:"capacity"`
	UsagePercent float64 `json:"usagePercent"`
}

type ClusterStorageResource struct {
	CapacityBytes          int64   `json:"capacityBytes"`
	AllocatableBytes       int64   `json:"allocatableBytes"`
	RequestsAllocatedBytes int64   `json:"requestsAllocatedBytes"`
	LimitsAllocatedBytes   int64   `json:"limitsAllocatedBytes"`
	UsageBytes             int64   `json:"usageBytes"`
	RequestsPercent        float64 `json:"requestsPercent"`
	UsagePercent           float64 `json:"usagePercent"`
}

type ClusterGPUResource struct {
	Capacity          float64 `json:"capacity"`
	Allocatable       float64 `json:"allocatable"`
	RequestsAllocated float64 `json:"requestsAllocated"`
	LimitsAllocated   float64 `json:"limitsAllocated"`
	Usage             float64 `json:"usage"`
	RequestsPercent   float64 `json:"requestsPercent"`
	UsagePercent      float64 `json:"usagePercent"`
}

type GetClusterResourcesResponse struct {
	ClusterUuid string                 `json:"clusterUuid"`
	Timestamp   int64                  `json:"timestamp"`
	CPU         ClusterCPUResource     `json:"cpu"`
	Memory      ClusterMemoryResource  `json:"memory"`
	Storage     ClusterStorageResource `json:"storage"`
	GPU         ClusterGPUResource     `json:"gpu"`
	Pods        ClusterPodsResource    `json:"pods"`
}

// kube-nova-aligned aliases used by monitor module split packages.
type GetCPUUsageRequest = GetPodCPUUsageRequest
type GetCPUUsageResponse = GetPodCPUUsageResponse
type GetMemoryUsageRequest = GetPodMemoryUsageRequest
type GetMemoryUsageResponse = GetPodMemoryUsageResponse
type GetNodeCPURequest = GetNodeCPUUsageRequest
type GetNodeCPUResponse = GetNodeCPUUsageResponse
type ListNodesMetricsRequest = GetNodeListMetricsRequest
type ListNodesMetricsResponse = GetNodeListMetricsResponse
