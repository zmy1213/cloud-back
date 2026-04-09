package projectrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type ProjectCluster struct {
	ID                        uint64
	ClusterUUID               string
	ClusterName               string
	ProjectID                 uint64
	CPULimit                  float64
	CPUOvercommitRatio        float64
	CPUCapacity               float64
	CPUAllocated              float64
	MemLimit                  float64
	MemOvercommitRatio        float64
	MemCapacity               float64
	MemAllocated              float64
	StorageLimit              float64
	StorageAllocated          float64
	GPULimit                  float64
	GPUOvercommitRatio        float64
	GPUCapacity               float64
	GPUAllocated              float64
	PodsLimit                 int64
	PodsAllocated             int64
	ConfigmapLimit            int64
	ConfigmapAllocated        int64
	SecretLimit               int64
	SecretAllocated           int64
	PVCLimit                  int64
	PVCAllocated              int64
	EphemeralStorageLimit     float64
	EphemeralStorageAllocated float64
	ServiceLimit              int64
	ServiceAllocated          int64
	LoadbalancersLimit        int64
	LoadbalancersAllocated    int64
	NodeportsLimit            int64
	NodeportsAllocated        int64
	DeploymentsLimit          int64
	DeploymentsAllocated      int64
	JobsLimit                 int64
	JobsAllocated             int64
	CronjobsLimit             int64
	CronjobsAllocated         int64
	DaemonsetsLimit           int64
	DaemonsetsAllocated       int64
	StatefulsetsLimit         int64
	StatefulsetsAllocated     int64
	IngressesLimit            int64
	IngressesAllocated        int64
	CreatedBy                 string
	UpdatedBy                 string
	CreatedAt                 int64
	UpdatedAt                 int64
}

type AddProjectClusterParams struct {
	ClusterUUID           string
	ProjectID             uint64
	CPULimit              float64
	CPUOvercommitRatio    float64
	CPUCapacity           float64
	MemLimit              float64
	MemOvercommitRatio    float64
	MemCapacity           float64
	StorageLimit          float64
	GPULimit              float64
	GPUOvercommitRatio    float64
	GPUCapacity           float64
	PodsLimit             int64
	ConfigmapLimit        int64
	SecretLimit           int64
	PVCLimit              int64
	EphemeralStorageLimit float64
	ServiceLimit          int64
	LoadbalancersLimit    int64
	NodeportsLimit        int64
	DeploymentsLimit      int64
	JobsLimit             int64
	CronjobsLimit         int64
	DaemonsetsLimit       int64
	StatefulsetsLimit     int64
	IngressesLimit        int64
	Operator              string
}

type UpdateProjectClusterParams struct {
	ID                    uint64
	CPULimit              float64
	CPUOvercommitRatio    float64
	CPUCapacity           float64
	MemLimit              float64
	MemOvercommitRatio    float64
	MemCapacity           float64
	StorageLimit          float64
	GPULimit              float64
	GPUOvercommitRatio    float64
	GPUCapacity           float64
	PodsLimit             int64
	ConfigmapLimit        int64
	SecretLimit           int64
	PVCLimit              int64
	EphemeralStorageLimit float64
	ServiceLimit          int64
	LoadbalancersLimit    int64
	NodeportsLimit        int64
	DeploymentsLimit      int64
	JobsLimit             int64
	CronjobsLimit         int64
	DaemonsetsLimit       int64
	StatefulsetsLimit     int64
	IngressesLimit        int64
	Operator              string
}

type SearchProjectClusterParams struct {
	ProjectID   uint64
	ClusterUUID string
}

func (s *Service) AddCluster(params AddProjectClusterParams) (ProjectCluster, error) {
	params.ClusterUUID = strings.TrimSpace(params.ClusterUUID)
	params.Operator = strings.TrimSpace(params.Operator)
	if params.Operator == "" {
		params.Operator = "system"
	}
	if params.ProjectID == 0 {
		return ProjectCluster{}, errors.New("projectId is required")
	}
	if params.ClusterUUID == "" {
		return ProjectCluster{}, errors.New("clusterUuid is required")
	}
	if params.CPUOvercommitRatio <= 0 {
		params.CPUOvercommitRatio = 1
	}
	if params.MemOvercommitRatio <= 0 {
		params.MemOvercommitRatio = 1
	}
	if params.GPUOvercommitRatio <= 0 {
		params.GPUOvercommitRatio = 1
	}
	if params.CPUCapacity <= 0 {
		params.CPUCapacity = params.CPULimit * params.CPUOvercommitRatio
	}
	if params.MemCapacity <= 0 {
		params.MemCapacity = params.MemLimit * params.MemOvercommitRatio
	}
	if params.GPUCapacity <= 0 {
		params.GPUCapacity = params.GPULimit * params.GPUOvercommitRatio
	}

	if s.useDB {
		return s.addClusterToDB(params)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range s.projectClusters {
		if item.ProjectID == params.ProjectID && item.ClusterUUID == params.ClusterUUID {
			return ProjectCluster{}, errors.New("project cluster quota already exists")
		}
	}

	now := time.Now().Unix()
	item := ProjectCluster{
		ID:                        s.nextProjectQuota,
		ClusterUUID:               params.ClusterUUID,
		ClusterName:               memoryClusterName(params.ClusterUUID),
		ProjectID:                 params.ProjectID,
		CPULimit:                  params.CPULimit,
		CPUOvercommitRatio:        params.CPUOvercommitRatio,
		CPUCapacity:               params.CPUCapacity,
		CPUAllocated:              0,
		MemLimit:                  params.MemLimit,
		MemOvercommitRatio:        params.MemOvercommitRatio,
		MemCapacity:               params.MemCapacity,
		MemAllocated:              0,
		StorageLimit:              params.StorageLimit,
		StorageAllocated:          0,
		GPULimit:                  params.GPULimit,
		GPUOvercommitRatio:        params.GPUOvercommitRatio,
		GPUCapacity:               params.GPUCapacity,
		GPUAllocated:              0,
		PodsLimit:                 params.PodsLimit,
		PodsAllocated:             0,
		ConfigmapLimit:            params.ConfigmapLimit,
		ConfigmapAllocated:        0,
		SecretLimit:               params.SecretLimit,
		SecretAllocated:           0,
		PVCLimit:                  params.PVCLimit,
		PVCAllocated:              0,
		EphemeralStorageLimit:     params.EphemeralStorageLimit,
		EphemeralStorageAllocated: 0,
		ServiceLimit:              params.ServiceLimit,
		ServiceAllocated:          0,
		LoadbalancersLimit:        params.LoadbalancersLimit,
		LoadbalancersAllocated:    0,
		NodeportsLimit:            params.NodeportsLimit,
		NodeportsAllocated:        0,
		DeploymentsLimit:          params.DeploymentsLimit,
		DeploymentsAllocated:      0,
		JobsLimit:                 params.JobsLimit,
		JobsAllocated:             0,
		CronjobsLimit:             params.CronjobsLimit,
		CronjobsAllocated:         0,
		DaemonsetsLimit:           params.DaemonsetsLimit,
		DaemonsetsAllocated:       0,
		StatefulsetsLimit:         params.StatefulsetsLimit,
		StatefulsetsAllocated:     0,
		IngressesLimit:            params.IngressesLimit,
		IngressesAllocated:        0,
		CreatedBy:                 params.Operator,
		UpdatedBy:                 params.Operator,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}
	s.projectClusters = append(s.projectClusters, item)
	s.nextProjectQuota++
	return item, nil
}

func (s *Service) UpdateCluster(params UpdateProjectClusterParams) error {
	if params.ID == 0 {
		return errors.New("id is required")
	}
	params.Operator = strings.TrimSpace(params.Operator)
	if params.Operator == "" {
		params.Operator = "system"
	}
	if params.CPUOvercommitRatio <= 0 {
		params.CPUOvercommitRatio = 1
	}
	if params.MemOvercommitRatio <= 0 {
		params.MemOvercommitRatio = 1
	}
	if params.GPUOvercommitRatio <= 0 {
		params.GPUOvercommitRatio = 1
	}
	if params.CPUCapacity <= 0 {
		params.CPUCapacity = params.CPULimit * params.CPUOvercommitRatio
	}
	if params.MemCapacity <= 0 {
		params.MemCapacity = params.MemLimit * params.MemOvercommitRatio
	}
	if params.GPUCapacity <= 0 {
		params.GPUCapacity = params.GPULimit * params.GPUOvercommitRatio
	}

	if s.useDB {
		return s.updateClusterToDB(params)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.projectClusters {
		if s.projectClusters[i].ID != params.ID {
			continue
		}
		s.projectClusters[i].CPULimit = params.CPULimit
		s.projectClusters[i].CPUOvercommitRatio = params.CPUOvercommitRatio
		s.projectClusters[i].CPUCapacity = params.CPUCapacity
		s.projectClusters[i].MemLimit = params.MemLimit
		s.projectClusters[i].MemOvercommitRatio = params.MemOvercommitRatio
		s.projectClusters[i].MemCapacity = params.MemCapacity
		s.projectClusters[i].StorageLimit = params.StorageLimit
		s.projectClusters[i].GPULimit = params.GPULimit
		s.projectClusters[i].GPUOvercommitRatio = params.GPUOvercommitRatio
		s.projectClusters[i].GPUCapacity = params.GPUCapacity
		s.projectClusters[i].PodsLimit = params.PodsLimit
		s.projectClusters[i].ConfigmapLimit = params.ConfigmapLimit
		s.projectClusters[i].SecretLimit = params.SecretLimit
		s.projectClusters[i].PVCLimit = params.PVCLimit
		s.projectClusters[i].EphemeralStorageLimit = params.EphemeralStorageLimit
		s.projectClusters[i].ServiceLimit = params.ServiceLimit
		s.projectClusters[i].LoadbalancersLimit = params.LoadbalancersLimit
		s.projectClusters[i].NodeportsLimit = params.NodeportsLimit
		s.projectClusters[i].DeploymentsLimit = params.DeploymentsLimit
		s.projectClusters[i].JobsLimit = params.JobsLimit
		s.projectClusters[i].CronjobsLimit = params.CronjobsLimit
		s.projectClusters[i].DaemonsetsLimit = params.DaemonsetsLimit
		s.projectClusters[i].StatefulsetsLimit = params.StatefulsetsLimit
		s.projectClusters[i].IngressesLimit = params.IngressesLimit
		s.projectClusters[i].UpdatedBy = params.Operator
		s.projectClusters[i].UpdatedAt = time.Now().Unix()
		return nil
	}
	return errors.New("project cluster quota not found")
}

func (s *Service) DeleteCluster(id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	if s.useDB {
		result, err := s.db.Exec(`
UPDATE onec_project_cluster
SET is_deleted = 1, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND is_deleted = 0
`, id)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return errors.New("project cluster quota not found")
		}
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	filtered := make([]ProjectCluster, 0, len(s.projectClusters))
	removed := false
	for _, item := range s.projectClusters {
		if item.ID == id {
			removed = true
			continue
		}
		filtered = append(filtered, item)
	}
	if !removed {
		return errors.New("project cluster quota not found")
	}
	s.projectClusters = filtered
	return nil
}

func (s *Service) GetClusterByID(id uint64) (ProjectCluster, bool, error) {
	if id == 0 {
		return ProjectCluster{}, false, nil
	}
	if s.useDB {
		return s.getClusterByIDFromDB(id)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.projectClusters {
		if item.ID == id {
			return item, true, nil
		}
	}
	return ProjectCluster{}, false, nil
}

func (s *Service) SearchClusters(params SearchProjectClusterParams) ([]ProjectCluster, error) {
	// projectId=0 means no project filter.
	params.ClusterUUID = strings.TrimSpace(params.ClusterUUID)
	if s.useDB {
		return s.searchClustersFromDB(params)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]ProjectCluster, 0)
	for _, item := range s.projectClusters {
		if params.ProjectID > 0 && item.ProjectID != params.ProjectID {
			continue
		}
		if params.ClusterUUID != "" && item.ClusterUUID != params.ClusterUUID {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *Service) addClusterToDB(params AddProjectClusterParams) (ProjectCluster, error) {
	if err := s.ensureProjectExists(params.ProjectID); err != nil {
		return ProjectCluster{}, err
	}
	clusterName, err := s.getClusterName(params.ClusterUUID)
	if err != nil {
		return ProjectCluster{}, err
	}

	var exists int
	if err := s.db.QueryRow(`
SELECT COUNT(1) FROM onec_project_cluster
WHERE project_id = ? AND cluster_uuid = ? AND is_deleted = 0
`, params.ProjectID, params.ClusterUUID).Scan(&exists); err != nil {
		return ProjectCluster{}, err
	}
	if exists > 0 {
		return ProjectCluster{}, errors.New("project cluster quota already exists")
	}

	result, err := s.db.Exec(`
INSERT INTO onec_project_cluster (
  project_id, cluster_uuid, cluster_name,
  cpu_limit, cpu_overcommit_ratio, cpu_capacity, cpu_allocated,
  mem_limit, mem_overcommit_ratio, mem_capacity, mem_allocated,
  storage_limit, storage_allocated,
  gpu_limit, gpu_overcommit_ratio, gpu_capacity, gpu_allocated,
  pods_limit, pods_allocated,
  configmap_limit, configmap_allocated,
  secret_limit, secret_allocated,
  pvc_limit, pvc_allocated,
  ephemeral_storage_limit, ephemeral_storage_allocated,
  service_limit, service_allocated,
  loadbalancers_limit, loadbalancers_allocated,
  nodeports_limit, nodeports_allocated,
  deployments_limit, deployments_allocated,
  jobs_limit, jobs_allocated,
  cronjobs_limit, cronjobs_allocated,
  daemonsets_limit, daemonsets_allocated,
  statefulsets_limit, statefulsets_allocated,
  ingresses_limit, ingresses_allocated,
  created_by, updated_by, is_deleted
) VALUES (
  ?, ?, ?,
  ?, ?, ?, 0,
  ?, ?, ?, 0,
  ?, 0,
  ?, ?, ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, 0,
  ?, ?,
  0
)
`,
		params.ProjectID, params.ClusterUUID, clusterName,
		params.CPULimit, params.CPUOvercommitRatio, params.CPUCapacity,
		params.MemLimit, params.MemOvercommitRatio, params.MemCapacity,
		params.StorageLimit,
		params.GPULimit, params.GPUOvercommitRatio, params.GPUCapacity,
		params.PodsLimit,
		params.ConfigmapLimit,
		params.SecretLimit,
		params.PVCLimit,
		params.EphemeralStorageLimit,
		params.ServiceLimit,
		params.LoadbalancersLimit,
		params.NodeportsLimit,
		params.DeploymentsLimit,
		params.JobsLimit,
		params.CronjobsLimit,
		params.DaemonsetsLimit,
		params.StatefulsetsLimit,
		params.IngressesLimit,
		params.Operator, params.Operator,
	)
	if err != nil {
		return ProjectCluster{}, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return ProjectCluster{}, err
	}
	item, ok, err := s.getClusterByIDFromDB(uint64(insertID))
	if err != nil {
		return ProjectCluster{}, err
	}
	if !ok {
		return ProjectCluster{}, errors.New("project cluster quota not found")
	}
	return item, nil
}

func (s *Service) updateClusterToDB(params UpdateProjectClusterParams) error {
	result, err := s.db.Exec(`
UPDATE onec_project_cluster SET
  cpu_limit = ?, cpu_overcommit_ratio = ?, cpu_capacity = ?,
  mem_limit = ?, mem_overcommit_ratio = ?, mem_capacity = ?,
  storage_limit = ?,
  gpu_limit = ?, gpu_overcommit_ratio = ?, gpu_capacity = ?,
  pods_limit = ?, configmap_limit = ?, secret_limit = ?, pvc_limit = ?,
  ephemeral_storage_limit = ?, service_limit = ?, loadbalancers_limit = ?,
  nodeports_limit = ?, deployments_limit = ?, jobs_limit = ?, cronjobs_limit = ?,
  daemonsets_limit = ?, statefulsets_limit = ?, ingresses_limit = ?,
  updated_by = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND is_deleted = 0
`,
		params.CPULimit, params.CPUOvercommitRatio, params.CPUCapacity,
		params.MemLimit, params.MemOvercommitRatio, params.MemCapacity,
		params.StorageLimit,
		params.GPULimit, params.GPUOvercommitRatio, params.GPUCapacity,
		params.PodsLimit, params.ConfigmapLimit, params.SecretLimit, params.PVCLimit,
		params.EphemeralStorageLimit, params.ServiceLimit, params.LoadbalancersLimit,
		params.NodeportsLimit, params.DeploymentsLimit, params.JobsLimit, params.CronjobsLimit,
		params.DaemonsetsLimit, params.StatefulsetsLimit, params.IngressesLimit,
		params.Operator, params.ID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("project cluster quota not found")
	}
	return nil
}

func (s *Service) getClusterByIDFromDB(id uint64) (ProjectCluster, bool, error) {
	row := s.db.QueryRow(`
SELECT id, cluster_uuid, cluster_name, project_id,
       cpu_limit, cpu_overcommit_ratio, cpu_capacity, cpu_allocated,
       mem_limit, mem_overcommit_ratio, mem_capacity, mem_allocated,
       storage_limit, storage_allocated,
       gpu_limit, gpu_overcommit_ratio, gpu_capacity, gpu_allocated,
       pods_limit, pods_allocated,
       configmap_limit, configmap_allocated,
       secret_limit, secret_allocated,
       pvc_limit, pvc_allocated,
       ephemeral_storage_limit, ephemeral_storage_allocated,
       service_limit, service_allocated,
       loadbalancers_limit, loadbalancers_allocated,
       nodeports_limit, nodeports_allocated,
       deployments_limit, deployments_allocated,
       jobs_limit, jobs_allocated,
       cronjobs_limit, cronjobs_allocated,
       daemonsets_limit, daemonsets_allocated,
       statefulsets_limit, statefulsets_allocated,
       ingresses_limit, ingresses_allocated,
       created_by, updated_by, created_at, updated_at
FROM onec_project_cluster
WHERE id = ? AND is_deleted = 0
LIMIT 1
`, id)

	var item ProjectCluster
	var createdAt time.Time
	var updatedAt time.Time
	err := row.Scan(
		&item.ID, &item.ClusterUUID, &item.ClusterName, &item.ProjectID,
		&item.CPULimit, &item.CPUOvercommitRatio, &item.CPUCapacity, &item.CPUAllocated,
		&item.MemLimit, &item.MemOvercommitRatio, &item.MemCapacity, &item.MemAllocated,
		&item.StorageLimit, &item.StorageAllocated,
		&item.GPULimit, &item.GPUOvercommitRatio, &item.GPUCapacity, &item.GPUAllocated,
		&item.PodsLimit, &item.PodsAllocated,
		&item.ConfigmapLimit, &item.ConfigmapAllocated,
		&item.SecretLimit, &item.SecretAllocated,
		&item.PVCLimit, &item.PVCAllocated,
		&item.EphemeralStorageLimit, &item.EphemeralStorageAllocated,
		&item.ServiceLimit, &item.ServiceAllocated,
		&item.LoadbalancersLimit, &item.LoadbalancersAllocated,
		&item.NodeportsLimit, &item.NodeportsAllocated,
		&item.DeploymentsLimit, &item.DeploymentsAllocated,
		&item.JobsLimit, &item.JobsAllocated,
		&item.CronjobsLimit, &item.CronjobsAllocated,
		&item.DaemonsetsLimit, &item.DaemonsetsAllocated,
		&item.StatefulsetsLimit, &item.StatefulsetsAllocated,
		&item.IngressesLimit, &item.IngressesAllocated,
		&item.CreatedBy, &item.UpdatedBy, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProjectCluster{}, false, nil
		}
		return ProjectCluster{}, false, err
	}
	item.CreatedAt = createdAt.Unix()
	item.UpdatedAt = updatedAt.Unix()
	return item, true, nil
}

func (s *Service) searchClustersFromDB(params SearchProjectClusterParams) ([]ProjectCluster, error) {
	query := `
SELECT id, cluster_uuid, cluster_name, project_id,
       cpu_limit, cpu_overcommit_ratio, cpu_capacity, cpu_allocated,
       mem_limit, mem_overcommit_ratio, mem_capacity, mem_allocated,
       storage_limit, storage_allocated,
       gpu_limit, gpu_overcommit_ratio, gpu_capacity, gpu_allocated,
       pods_limit, pods_allocated,
       configmap_limit, configmap_allocated,
       secret_limit, secret_allocated,
       pvc_limit, pvc_allocated,
       ephemeral_storage_limit, ephemeral_storage_allocated,
       service_limit, service_allocated,
       loadbalancers_limit, loadbalancers_allocated,
       nodeports_limit, nodeports_allocated,
       deployments_limit, deployments_allocated,
       jobs_limit, jobs_allocated,
       cronjobs_limit, cronjobs_allocated,
       daemonsets_limit, daemonsets_allocated,
       statefulsets_limit, statefulsets_allocated,
       ingresses_limit, ingresses_allocated,
       created_by, updated_by, created_at, updated_at
FROM onec_project_cluster
WHERE is_deleted = 0
`
	args := make([]any, 0, 2)
	if params.ProjectID > 0 {
		query += " AND project_id = ?"
		args = append(args, params.ProjectID)
	}
	if params.ClusterUUID != "" {
		query += " AND cluster_uuid = ?"
		args = append(args, params.ClusterUUID)
	}
	query += " ORDER BY id DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ProjectCluster, 0)
	for rows.Next() {
		var item ProjectCluster
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(
			&item.ID, &item.ClusterUUID, &item.ClusterName, &item.ProjectID,
			&item.CPULimit, &item.CPUOvercommitRatio, &item.CPUCapacity, &item.CPUAllocated,
			&item.MemLimit, &item.MemOvercommitRatio, &item.MemCapacity, &item.MemAllocated,
			&item.StorageLimit, &item.StorageAllocated,
			&item.GPULimit, &item.GPUOvercommitRatio, &item.GPUCapacity, &item.GPUAllocated,
			&item.PodsLimit, &item.PodsAllocated,
			&item.ConfigmapLimit, &item.ConfigmapAllocated,
			&item.SecretLimit, &item.SecretAllocated,
			&item.PVCLimit, &item.PVCAllocated,
			&item.EphemeralStorageLimit, &item.EphemeralStorageAllocated,
			&item.ServiceLimit, &item.ServiceAllocated,
			&item.LoadbalancersLimit, &item.LoadbalancersAllocated,
			&item.NodeportsLimit, &item.NodeportsAllocated,
			&item.DeploymentsLimit, &item.DeploymentsAllocated,
			&item.JobsLimit, &item.JobsAllocated,
			&item.CronjobsLimit, &item.CronjobsAllocated,
			&item.DaemonsetsLimit, &item.DaemonsetsAllocated,
			&item.StatefulsetsLimit, &item.StatefulsetsAllocated,
			&item.IngressesLimit, &item.IngressesAllocated,
			&item.CreatedBy, &item.UpdatedBy, &createdAt, &updatedAt,
		); err != nil {
			return nil, err
		}
		item.CreatedAt = createdAt.Unix()
		item.UpdatedAt = updatedAt.Unix()
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Service) ensureProjectExists(projectID uint64) error {
	var exists int
	if err := s.db.QueryRow(`
SELECT COUNT(1) FROM onec_project WHERE id = ? AND is_deleted = 0
`, projectID).Scan(&exists); err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("project not found")
	}
	return nil
}

func (s *Service) getClusterName(clusterUUID string) (string, error) {
	var name string
	err := s.db.QueryRow(`
SELECT name FROM onec_cluster WHERE uuid = ? AND is_deleted = 0 LIMIT 1
`, clusterUUID).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("cluster not found")
		}
		return "", err
	}
	return strings.TrimSpace(name), nil
}

func memoryClusterName(clusterUUID string) string {
	switch clusterUUID {
	case "11111111-1111-1111-1111-111111111111":
		return "prod-hz"
	case "22222222-2222-2222-2222-222222222222":
		return "staging-sh"
	case "33333333-3333-3333-3333-333333333333":
		return "edge-gz"
	default:
		return fmt.Sprintf("cluster-%s", clusterUUID)
	}
}

func defaultProjectClusters() []ProjectCluster {
	return []ProjectCluster{
		{
			ID:                        1,
			ClusterUUID:               "11111111-1111-1111-1111-111111111111",
			ClusterName:               "prod-hz",
			ProjectID:                 2,
			CPULimit:                  12,
			CPUOvercommitRatio:        1.5,
			CPUCapacity:               18,
			CPUAllocated:              4.5,
			MemLimit:                  32,
			MemOvercommitRatio:        1.2,
			MemCapacity:               38.4,
			MemAllocated:              12.8,
			StorageLimit:              500,
			StorageAllocated:          120,
			GPULimit:                  2,
			GPUOvercommitRatio:        1,
			GPUCapacity:               2,
			GPUAllocated:              1,
			PodsLimit:                 100,
			PodsAllocated:             36,
			ConfigmapLimit:            100,
			ConfigmapAllocated:        10,
			SecretLimit:               100,
			SecretAllocated:           12,
			PVCLimit:                  100,
			PVCAllocated:              8,
			EphemeralStorageLimit:     120,
			EphemeralStorageAllocated: 14,
			ServiceLimit:              80,
			ServiceAllocated:          10,
			LoadbalancersLimit:        5,
			LoadbalancersAllocated:    1,
			NodeportsLimit:            20,
			NodeportsAllocated:        2,
			DeploymentsLimit:          80,
			DeploymentsAllocated:      16,
			JobsLimit:                 20,
			JobsAllocated:             2,
			CronjobsLimit:             20,
			CronjobsAllocated:         1,
			DaemonsetsLimit:           10,
			DaemonsetsAllocated:       1,
			StatefulsetsLimit:         30,
			StatefulsetsAllocated:     3,
			IngressesLimit:            50,
			IngressesAllocated:        6,
			CreatedBy:                 "system",
			UpdatedBy:                 "system",
			CreatedAt:                 1735776000,
			UpdatedAt:                 1735776000,
		},
	}
}
