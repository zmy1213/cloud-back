package types

type Project struct {
	ID            uint64 `json:"id"`
	Name          string `json:"name"`
	Uuid          string `json:"uuid"`
	Description   string `json:"description"`
	IsSystem      int64  `json:"isSystem"`
	CreatedBy     string `json:"createdBy"`
	UpdatedBy     string `json:"updatedBy"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
	AdminCount    int64  `json:"adminCount"`
	ResourceCount int64  `json:"resourceCount"`
}

type AddProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsSystem    int64  `json:"isSystem"`
}

type UpdateProjectRequest struct {
	ID          uint64 `path:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SearchProjectRequest struct {
	Page     uint64 `form:"page"`
	PageSize uint64 `form:"pageSize"`
	Name     string `form:"name"`
	Uuid     string `form:"uuid"`
}

type SearchProjectResponse struct {
	Items []Project `json:"items"`
	Total uint64    `json:"total"`
}

type GetProjectRequest struct {
	ID uint64 `path:"id"`
}

type GetProjectsByUserIdRequest struct {
	UserID uint64 `form:"userId"`
	Name   string `form:"name"`
}

type ProjectAdmin struct {
	ID        uint64 `json:"id"`
	ProjectID uint64 `json:"projectId"`
	UserID    uint64 `json:"userId"`
	CreatedAt int64  `json:"createdAt"`
}

type AddProjectAdminRequest struct {
	ProjectID uint64   `json:"projectId"`
	UserIDs   []uint64 `json:"userIds"`
}

type GetProjectAdminsRequest struct {
	ProjectID uint64 `form:"projectId"`
}

type ProjectCluster struct {
	ID                        uint64  `json:"id"`
	ClusterUUID               string  `json:"clusterUuid"`
	ClusterName               string  `json:"clusterName"`
	ProjectID                 uint64  `json:"projectId"`
	CPULimit                  float64 `json:"cpuLimit"`
	CPUOvercommitRatio        float64 `json:"cpuOvercommitRatio"`
	CPUCapacity               float64 `json:"cpuCapacity"`
	CPUAllocated              float64 `json:"cpuAllocated"`
	MemLimit                  float64 `json:"memLimit"`
	MemOvercommitRatio        float64 `json:"memOvercommitRatio"`
	MemCapacity               float64 `json:"memCapacity"`
	MemAllocated              float64 `json:"memAllocated"`
	StorageLimit              float64 `json:"storageLimit"`
	StorageAllocated          float64 `json:"storageAllocated"`
	GPULimit                  float64 `json:"gpuLimit"`
	GPUOvercommitRatio        float64 `json:"gpuOvercommitRatio"`
	GPUCapacity               float64 `json:"gpuCapacity"`
	GPUAllocated              float64 `json:"gpuAllocated"`
	PodsLimit                 int64   `json:"podsLimit"`
	PodsAllocated             int64   `json:"podsAllocated"`
	ConfigmapLimit            int64   `json:"configmapLimit"`
	ConfigmapAllocated        int64   `json:"configmapAllocated"`
	SecretLimit               int64   `json:"secretLimit"`
	SecretAllocated           int64   `json:"secretAllocated"`
	PVCLimit                  int64   `json:"pvcLimit"`
	PVCAllocated              int64   `json:"pvcAllocated"`
	EphemeralStorageLimit     float64 `json:"ephemeralStorageLimit"`
	EphemeralStorageAllocated float64 `json:"ephemeralStorageAllocated"`
	ServiceLimit              int64   `json:"serviceLimit"`
	ServiceAllocated          int64   `json:"serviceAllocated"`
	LoadbalancersLimit        int64   `json:"loadbalancersLimit"`
	LoadbalancersAllocated    int64   `json:"loadbalancersAllocated"`
	NodeportsLimit            int64   `json:"nodeportsLimit"`
	NodeportsAllocated        int64   `json:"nodeportsAllocated"`
	DeploymentsLimit          int64   `json:"deploymentsLimit"`
	DeploymentsAllocated      int64   `json:"deploymentsAllocated"`
	JobsLimit                 int64   `json:"jobsLimit"`
	JobsAllocated             int64   `json:"jobsAllocated"`
	CronjobsLimit             int64   `json:"cronjobsLimit"`
	CronjobsAllocated         int64   `json:"cronjobsAllocated"`
	DaemonsetsLimit           int64   `json:"daemonsetsLimit"`
	DaemonsetsAllocated       int64   `json:"daemonsetsAllocated"`
	StatefulsetsLimit         int64   `json:"statefulsetsLimit"`
	StatefulsetsAllocated     int64   `json:"statefulsetsAllocated"`
	IngressesLimit            int64   `json:"ingressesLimit"`
	IngressesAllocated        int64   `json:"ingressesAllocated"`
	CreatedBy                 string  `json:"createdBy"`
	UpdatedBy                 string  `json:"updatedBy"`
	CreatedAt                 int64   `json:"createdAt"`
	UpdatedAt                 int64   `json:"updatedAt"`
}

type AddProjectClusterRequest struct {
	ClusterUUID           string  `json:"clusterUuid"`
	ProjectID             uint64  `json:"projectId"`
	CPULimit              float64 `json:"cpuLimit"`
	CPUOvercommitRatio    float64 `json:"cpuOvercommitRatio"`
	CPUCapacity           float64 `json:"cpuCapacity"`
	MemLimit              float64 `json:"memLimit"`
	MemOvercommitRatio    float64 `json:"memOvercommitRatio"`
	MemCapacity           float64 `json:"memCapacity"`
	StorageLimit          float64 `json:"storageLimit"`
	GPULimit              float64 `json:"gpuLimit"`
	GPUOvercommitRatio    float64 `json:"gpuOvercommitRatio"`
	GPUCapacity           float64 `json:"gpuCapacity"`
	PodsLimit             int64   `json:"podsLimit"`
	ConfigmapLimit        int64   `json:"configmapLimit"`
	SecretLimit           int64   `json:"secretLimit"`
	PVCLimit              int64   `json:"pvcLimit"`
	EphemeralStorageLimit float64 `json:"ephemeralStorageLimit"`
	ServiceLimit          int64   `json:"serviceLimit"`
	LoadbalancersLimit    int64   `json:"loadbalancersLimit"`
	NodeportsLimit        int64   `json:"nodeportsLimit"`
	DeploymentsLimit      int64   `json:"deploymentsLimit"`
	JobsLimit             int64   `json:"jobsLimit"`
	CronjobsLimit         int64   `json:"cronjobsLimit"`
	DaemonsetsLimit       int64   `json:"daemonsetsLimit"`
	StatefulsetsLimit     int64   `json:"statefulsetsLimit"`
	IngressesLimit        int64   `json:"ingressesLimit"`
}

type UpdateProjectClusterRequest struct {
	ID                    uint64  `path:"id"`
	CPULimit              float64 `json:"cpuLimit"`
	CPUOvercommitRatio    float64 `json:"cpuOvercommitRatio"`
	CPUCapacity           float64 `json:"cpuCapacity"`
	MemLimit              float64 `json:"memLimit"`
	MemOvercommitRatio    float64 `json:"memOvercommitRatio"`
	MemCapacity           float64 `json:"memCapacity"`
	StorageLimit          float64 `json:"storageLimit"`
	GPULimit              float64 `json:"gpuLimit"`
	GPUOvercommitRatio    float64 `json:"gpuOvercommitRatio"`
	GPUCapacity           float64 `json:"gpuCapacity"`
	PodsLimit             int64   `json:"podsLimit"`
	ConfigmapLimit        int64   `json:"configmapLimit"`
	SecretLimit           int64   `json:"secretLimit"`
	PVCLimit              int64   `json:"pvcLimit"`
	EphemeralStorageLimit float64 `json:"ephemeralStorageLimit"`
	ServiceLimit          int64   `json:"serviceLimit"`
	LoadbalancersLimit    int64   `json:"loadbalancersLimit"`
	NodeportsLimit        int64   `json:"nodeportsLimit"`
	DeploymentsLimit      int64   `json:"deploymentsLimit"`
	JobsLimit             int64   `json:"jobsLimit"`
	CronjobsLimit         int64   `json:"cronjobsLimit"`
	DaemonsetsLimit       int64   `json:"daemonsetsLimit"`
	StatefulsetsLimit     int64   `json:"statefulsetsLimit"`
	IngressesLimit        int64   `json:"ingressesLimit"`
}

type SearchProjectClusterRequest struct {
	ProjectID   uint64 `form:"projectId"`
	ClusterUUID string `form:"clusterUuid"`
}

type ProjectWorkspace struct {
	ID               uint64  `json:"id"`
	ProjectClusterID uint64  `json:"projectClusterId"`
	ProjectID        uint64  `json:"projectId"`
	ClusterUUID      string  `json:"clusterUuid"`
	ClusterName      string  `json:"clusterName"`
	Name             string  `json:"name"`
	Namespace        string  `json:"namespace"`
	Description      string  `json:"description"`
	CPUAllocated     float64 `json:"cpuAllocated"`
	MemAllocated     float64 `json:"memAllocated"`
	StorageAllocated float64 `json:"storageAllocated"`
	GPUAllocated     float64 `json:"gpuAllocated"`
	PodsAllocated    int64   `json:"podsAllocated"`
	CreatedBy        string  `json:"createdBy"`
	UpdatedBy        string  `json:"updatedBy"`
	CreatedAt        int64   `json:"createdAt"`
	UpdatedAt        int64   `json:"updatedAt"`
}

type AddProjectWorkspaceRequest struct {
	ProjectClusterID uint64  `json:"projectClusterId"`
	Name             string  `json:"name"`
	Namespace        string  `json:"namespace"`
	Description      string  `json:"description"`
	CPUAllocated     float64 `json:"cpuAllocated"`
	MemAllocated     float64 `json:"memAllocated"`
	StorageAllocated float64 `json:"storageAllocated"`
	GPUAllocated     float64 `json:"gpuAllocated"`
	PodsAllocated    int64   `json:"podsAllocated"`
}

type UpdateProjectWorkspaceRequest struct {
	ID               uint64  `path:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	CPUAllocated     float64 `json:"cpuAllocated"`
	MemAllocated     float64 `json:"memAllocated"`
	StorageAllocated float64 `json:"storageAllocated"`
	GPUAllocated     float64 `json:"gpuAllocated"`
	PodsAllocated    int64   `json:"podsAllocated"`
}

type SearchProjectWorkspaceRequest struct {
	ProjectClusterID uint64 `form:"projectClusterId"`
	Name             string `form:"name"`
	Namespace        string `form:"namespace"`
}
