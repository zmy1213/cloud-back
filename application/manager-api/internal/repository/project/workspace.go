package projectrepo

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type ProjectWorkspace struct {
	ID               uint64
	ProjectClusterID uint64
	ProjectID        uint64
	ClusterUUID      string
	ClusterName      string
	Name             string
	Namespace        string
	Description      string
	CPUAllocated     float64
	MemAllocated     float64
	StorageAllocated float64
	GPUAllocated     float64
	PodsAllocated    int64
	CreatedBy        string
	UpdatedBy        string
	CreatedAt        int64
	UpdatedAt        int64
}

type AddProjectWorkspaceParams struct {
	ProjectClusterID uint64
	Name             string
	Namespace        string
	Description      string
	CPUAllocated     float64
	MemAllocated     float64
	StorageAllocated float64
	GPUAllocated     float64
	PodsAllocated    int64
	Operator         string
}

type UpdateProjectWorkspaceParams struct {
	ID               uint64
	Name             string
	Description      string
	CPUAllocated     float64
	MemAllocated     float64
	StorageAllocated float64
	GPUAllocated     float64
	PodsAllocated    int64
	Operator         string
}

type SearchProjectWorkspaceParams struct {
	ProjectClusterID uint64
	Name             string
	Namespace        string
}

func (s *Service) AddWorkspace(params AddProjectWorkspaceParams) (ProjectWorkspace, error) {
	params.Name = strings.TrimSpace(params.Name)
	params.Namespace = normalizeNamespace(params.Namespace)
	params.Description = strings.TrimSpace(params.Description)
	params.Operator = strings.TrimSpace(params.Operator)
	if params.Operator == "" {
		params.Operator = "system"
	}
	if params.ProjectClusterID == 0 {
		return ProjectWorkspace{}, errors.New("projectClusterId is required")
	}
	if params.Name == "" {
		return ProjectWorkspace{}, errors.New("workspace name is required")
	}
	if params.Namespace == "" {
		return ProjectWorkspace{}, errors.New("namespace is required")
	}
	if params.CPUAllocated < 0 || params.MemAllocated < 0 || params.StorageAllocated < 0 || params.GPUAllocated < 0 || params.PodsAllocated < 0 {
		return ProjectWorkspace{}, errors.New("workspace resources cannot be negative")
	}

	if s.useDB {
		return s.addWorkspaceToDB(params)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var cluster ProjectCluster
	found := false
	for _, item := range s.projectClusters {
		if item.ID == params.ProjectClusterID {
			cluster = item
			found = true
			break
		}
	}
	if !found {
		return ProjectWorkspace{}, errors.New("project cluster quota not found")
	}

	for _, item := range s.projectWorkspaces {
		if item.ProjectClusterID == params.ProjectClusterID && item.Namespace == params.Namespace {
			return ProjectWorkspace{}, errors.New("workspace namespace already exists in this cluster quota")
		}
	}

	now := time.Now().Unix()
	workspace := ProjectWorkspace{
		ID:               s.nextWorkspace,
		ProjectClusterID: params.ProjectClusterID,
		ProjectID:        cluster.ProjectID,
		ClusterUUID:      cluster.ClusterUUID,
		ClusterName:      cluster.ClusterName,
		Name:             params.Name,
		Namespace:        params.Namespace,
		Description:      params.Description,
		CPUAllocated:     params.CPUAllocated,
		MemAllocated:     params.MemAllocated,
		StorageAllocated: params.StorageAllocated,
		GPUAllocated:     params.GPUAllocated,
		PodsAllocated:    params.PodsAllocated,
		CreatedBy:        params.Operator,
		UpdatedBy:        params.Operator,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	s.projectWorkspaces = append(s.projectWorkspaces, workspace)
	s.nextWorkspace++
	return workspace, nil
}

func (s *Service) UpdateWorkspace(params UpdateProjectWorkspaceParams) error {
	params.Name = strings.TrimSpace(params.Name)
	params.Description = strings.TrimSpace(params.Description)
	params.Operator = strings.TrimSpace(params.Operator)
	if params.Operator == "" {
		params.Operator = "system"
	}
	if params.ID == 0 {
		return errors.New("id is required")
	}
	if params.Name == "" {
		return errors.New("workspace name is required")
	}
	if params.CPUAllocated < 0 || params.MemAllocated < 0 || params.StorageAllocated < 0 || params.GPUAllocated < 0 || params.PodsAllocated < 0 {
		return errors.New("workspace resources cannot be negative")
	}

	if s.useDB {
		return s.updateWorkspaceToDB(params)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.projectWorkspaces {
		if s.projectWorkspaces[i].ID != params.ID {
			continue
		}
		s.projectWorkspaces[i].Name = params.Name
		s.projectWorkspaces[i].Description = params.Description
		s.projectWorkspaces[i].CPUAllocated = params.CPUAllocated
		s.projectWorkspaces[i].MemAllocated = params.MemAllocated
		s.projectWorkspaces[i].StorageAllocated = params.StorageAllocated
		s.projectWorkspaces[i].GPUAllocated = params.GPUAllocated
		s.projectWorkspaces[i].PodsAllocated = params.PodsAllocated
		s.projectWorkspaces[i].UpdatedBy = params.Operator
		s.projectWorkspaces[i].UpdatedAt = time.Now().Unix()
		return nil
	}
	return errors.New("workspace not found")
}

func (s *Service) DeleteWorkspace(id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}

	if s.useDB {
		result, err := s.db.Exec(`
UPDATE onec_project_workspace
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
			return errors.New("workspace not found")
		}
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	filtered := make([]ProjectWorkspace, 0, len(s.projectWorkspaces))
	removed := false
	for _, item := range s.projectWorkspaces {
		if item.ID == id {
			removed = true
			continue
		}
		filtered = append(filtered, item)
	}
	if !removed {
		return errors.New("workspace not found")
	}
	s.projectWorkspaces = filtered
	return nil
}

func (s *Service) GetWorkspaceByID(id uint64) (ProjectWorkspace, bool, error) {
	if id == 0 {
		return ProjectWorkspace{}, false, nil
	}

	if s.useDB {
		return s.getWorkspaceByIDFromDB(id)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, item := range s.projectWorkspaces {
		if item.ID == id {
			return item, true, nil
		}
	}
	return ProjectWorkspace{}, false, nil
}

func (s *Service) SearchWorkspaces(params SearchProjectWorkspaceParams) ([]ProjectWorkspace, error) {
	if params.ProjectClusterID == 0 {
		return nil, errors.New("projectClusterId is required")
	}
	params.Name = strings.ToLower(strings.TrimSpace(params.Name))
	params.Namespace = strings.ToLower(strings.TrimSpace(params.Namespace))

	if s.useDB {
		return s.searchWorkspacesFromDB(params)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]ProjectWorkspace, 0)
	for _, item := range s.projectWorkspaces {
		if item.ProjectClusterID != params.ProjectClusterID {
			continue
		}
		if params.Name != "" && !strings.Contains(strings.ToLower(item.Name), params.Name) {
			continue
		}
		if params.Namespace != "" && !strings.Contains(strings.ToLower(item.Namespace), params.Namespace) {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}

func (s *Service) addWorkspaceToDB(params AddProjectWorkspaceParams) (ProjectWorkspace, error) {
	clusterMeta, err := s.getProjectClusterMeta(params.ProjectClusterID)
	if err != nil {
		return ProjectWorkspace{}, err
	}

	var exists int
	if err := s.db.QueryRow(`
SELECT COUNT(1) FROM onec_project_workspace
WHERE project_cluster_id = ? AND namespace = ? AND is_deleted = 0
`, params.ProjectClusterID, params.Namespace).Scan(&exists); err != nil {
		return ProjectWorkspace{}, err
	}
	if exists > 0 {
		return ProjectWorkspace{}, errors.New("workspace namespace already exists in this cluster quota")
	}

	result, err := s.db.Exec(`
INSERT INTO onec_project_workspace (
  project_cluster_id, project_id, cluster_uuid, cluster_name,
  name, namespace, description,
  cpu_allocated, mem_allocated, storage_allocated, gpu_allocated, pods_allocated,
  created_by, updated_by, is_deleted
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
`,
		params.ProjectClusterID, clusterMeta.ProjectID, clusterMeta.ClusterUUID, clusterMeta.ClusterName,
		params.Name, params.Namespace, params.Description,
		params.CPUAllocated, params.MemAllocated, params.StorageAllocated, params.GPUAllocated, params.PodsAllocated,
		params.Operator, params.Operator,
	)
	if err != nil {
		return ProjectWorkspace{}, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		return ProjectWorkspace{}, err
	}
	item, ok, err := s.getWorkspaceByIDFromDB(uint64(insertID))
	if err != nil {
		return ProjectWorkspace{}, err
	}
	if !ok {
		return ProjectWorkspace{}, errors.New("workspace not found")
	}
	return item, nil
}

func (s *Service) updateWorkspaceToDB(params UpdateProjectWorkspaceParams) error {
	result, err := s.db.Exec(`
UPDATE onec_project_workspace SET
  name = ?, description = ?,
  cpu_allocated = ?, mem_allocated = ?, storage_allocated = ?, gpu_allocated = ?, pods_allocated = ?,
  updated_by = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND is_deleted = 0
`,
		params.Name, params.Description,
		params.CPUAllocated, params.MemAllocated, params.StorageAllocated, params.GPUAllocated, params.PodsAllocated,
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
		return errors.New("workspace not found")
	}
	return nil
}

func (s *Service) getWorkspaceByIDFromDB(id uint64) (ProjectWorkspace, bool, error) {
	row := s.db.QueryRow(`
SELECT id, project_cluster_id, project_id, cluster_uuid, cluster_name,
       name, namespace, description,
       cpu_allocated, mem_allocated, storage_allocated, gpu_allocated, pods_allocated,
       created_by, updated_by, created_at, updated_at
FROM onec_project_workspace
WHERE id = ? AND is_deleted = 0
LIMIT 1
`, id)

	var item ProjectWorkspace
	var createdAt time.Time
	var updatedAt time.Time
	err := row.Scan(
		&item.ID, &item.ProjectClusterID, &item.ProjectID, &item.ClusterUUID, &item.ClusterName,
		&item.Name, &item.Namespace, &item.Description,
		&item.CPUAllocated, &item.MemAllocated, &item.StorageAllocated, &item.GPUAllocated, &item.PodsAllocated,
		&item.CreatedBy, &item.UpdatedBy, &createdAt, &updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProjectWorkspace{}, false, nil
		}
		return ProjectWorkspace{}, false, err
	}
	item.CreatedAt = createdAt.Unix()
	item.UpdatedAt = updatedAt.Unix()
	return item, true, nil
}

func (s *Service) searchWorkspacesFromDB(params SearchProjectWorkspaceParams) ([]ProjectWorkspace, error) {
	query := `
SELECT id, project_cluster_id, project_id, cluster_uuid, cluster_name,
       name, namespace, description,
       cpu_allocated, mem_allocated, storage_allocated, gpu_allocated, pods_allocated,
       created_by, updated_by, created_at, updated_at
FROM onec_project_workspace
WHERE project_cluster_id = ? AND is_deleted = 0
`
	args := []any{params.ProjectClusterID}
	if params.Name != "" {
		query += " AND LOWER(name) LIKE ?"
		args = append(args, "%"+params.Name+"%")
	}
	if params.Namespace != "" {
		query += " AND LOWER(namespace) LIKE ?"
		args = append(args, "%"+params.Namespace+"%")
	}
	query += " ORDER BY id DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ProjectWorkspace, 0)
	for rows.Next() {
		var item ProjectWorkspace
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(
			&item.ID, &item.ProjectClusterID, &item.ProjectID, &item.ClusterUUID, &item.ClusterName,
			&item.Name, &item.Namespace, &item.Description,
			&item.CPUAllocated, &item.MemAllocated, &item.StorageAllocated, &item.GPUAllocated, &item.PodsAllocated,
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

type projectClusterMeta struct {
	ProjectID   uint64
	ClusterUUID string
	ClusterName string
}

func (s *Service) getProjectClusterMeta(projectClusterID uint64) (projectClusterMeta, error) {
	row := s.db.QueryRow(`
SELECT project_id, cluster_uuid, cluster_name
FROM onec_project_cluster
WHERE id = ? AND is_deleted = 0
LIMIT 1
`, projectClusterID)

	var meta projectClusterMeta
	if err := row.Scan(&meta.ProjectID, &meta.ClusterUUID, &meta.ClusterName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return projectClusterMeta{}, errors.New("project cluster quota not found")
		}
		return projectClusterMeta{}, err
	}
	return meta, nil
}

func normalizeNamespace(namespace string) string {
	return strings.ToLower(strings.TrimSpace(namespace))
}

func defaultProjectWorkspaces() []ProjectWorkspace {
	return []ProjectWorkspace{
		{
			ID:               1,
			ProjectClusterID: 1,
			ProjectID:        2,
			ClusterUUID:      "11111111-1111-1111-1111-111111111111",
			ClusterName:      "prod-hz",
			Name:             "默认研发空间",
			Namespace:        "dev-default",
			Description:      "默认研发工作空间",
			CPUAllocated:     2,
			MemAllocated:     4,
			StorageAllocated: 20,
			GPUAllocated:     0,
			PodsAllocated:    30,
			CreatedBy:        "system",
			UpdatedBy:        "system",
			CreatedAt:        1735776000,
			UpdatedAt:        1735776000,
		},
	}
}
