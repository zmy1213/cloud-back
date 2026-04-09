package projectrepo

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	appcfg "github.com/yanshicheng/cloud-back/common/config"
)

type Project struct {
	ID            uint64
	Name          string
	UUID          string
	Description   string
	IsSystem      int64
	CreatedBy     string
	UpdatedBy     string
	CreatedAt     int64
	UpdatedAt     int64
	AdminCount    int64
	ResourceCount int64
}

type ProjectAdmin struct {
	ID        uint64
	ProjectID uint64
	UserID    uint64
	CreatedAt int64
}

type AddProjectParams struct {
	Name        string
	Description string
	IsSystem    int64
	Operator    string
}

type UpdateProjectParams struct {
	ID          uint64
	Name        string
	Description string
	Operator    string
}

type SearchParams struct {
	Page     uint64
	PageSize uint64
	Name     string
	UUID     string
}

type Service struct {
	db    *sql.DB
	useDB bool

	mu                sync.RWMutex
	projects          []Project
	admins            []ProjectAdmin
	projectClusters   []ProjectCluster
	projectWorkspaces []ProjectWorkspace
	nextID            uint64
	nextAdmin         uint64
	nextProjectQuota  uint64
	nextWorkspace     uint64
}

func NewService(mysqlCfg appcfg.MysqlConfig) *Service {
	s := &Service{
		projects:          defaultProjects(),
		admins:            defaultAdmins(),
		projectClusters:   defaultProjectClusters(),
		projectWorkspaces: defaultProjectWorkspaces(),
		nextID:            3,
		nextAdmin:         3,
		nextProjectQuota:  2,
		nextWorkspace:     2,
	}

	if !mysqlCfg.Enabled || strings.TrimSpace(mysqlCfg.DataSource) == "" {
		return s
	}

	db, err := sql.Open("mysql", mysqlCfg.DataSource)
	if err != nil {
		return s
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return s
	}
	s.db = db
	s.useDB = true
	return s
}

func (s *Service) Add(params AddProjectParams) (Project, error) {
	params.Name = strings.TrimSpace(params.Name)
	params.Description = strings.TrimSpace(params.Description)
	params.Operator = strings.TrimSpace(params.Operator)
	if params.Operator == "" {
		params.Operator = "system"
	}
	if params.Name == "" {
		return Project{}, errors.New("project name is required")
	}
	if params.IsSystem != 0 && params.IsSystem != 1 {
		return Project{}, errors.New("isSystem must be 0 or 1")
	}

	if s.useDB {
		return s.addToDB(params)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().Unix()
	p := Project{
		ID:          s.nextID,
		Name:        params.Name,
		UUID:        uuid.NewString(),
		Description: params.Description,
		IsSystem:    params.IsSystem,
		CreatedBy:   params.Operator,
		UpdatedBy:   params.Operator,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.projects = append(s.projects, p)
	s.nextID++
	return p, nil
}

func (s *Service) Update(params UpdateProjectParams) error {
	params.Name = strings.TrimSpace(params.Name)
	params.Description = strings.TrimSpace(params.Description)
	params.Operator = strings.TrimSpace(params.Operator)
	if params.Operator == "" {
		params.Operator = "system"
	}
	if params.ID == 0 {
		return errors.New("project id is required")
	}
	if params.Name == "" {
		return errors.New("project name is required")
	}

	if s.useDB {
		return s.updateToDB(params)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.projects {
		if s.projects[i].ID != params.ID {
			continue
		}
		s.projects[i].Name = params.Name
		s.projects[i].Description = params.Description
		s.projects[i].UpdatedBy = params.Operator
		s.projects[i].UpdatedAt = time.Now().Unix()
		return nil
	}
	return errors.New("project not found")
}

func (s *Service) Delete(id uint64) error {
	if id == 0 {
		return errors.New("project id is required")
	}

	if s.useDB {
		_, err := s.db.Exec(`UPDATE onec_project SET is_deleted = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id)
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	filtered := make([]Project, 0, len(s.projects))
	for _, p := range s.projects {
		if p.ID == id {
			continue
		}
		filtered = append(filtered, p)
	}
	s.projects = filtered
	return nil
}

func (s *Service) GetByID(id uint64) (Project, bool, error) {
	if id == 0 {
		return Project{}, false, nil
	}

	if s.useDB {
		return s.getByIDFromDB(id)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.projects {
		if p.ID == id {
			p.AdminCount = s.memoryAdminCount(p.ID)
			p.ResourceCount = s.memoryResourceCount(p.ID)
			return p, true, nil
		}
	}
	return Project{}, false, nil
}

func (s *Service) Search(params SearchParams) ([]Project, uint64, error) {
	params.Name = strings.TrimSpace(params.Name)
	params.UUID = strings.TrimSpace(params.UUID)
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 10
	}
	if params.PageSize > 200 {
		params.PageSize = 200
	}

	if s.useDB {
		return s.searchFromDB(params)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	filtered := make([]Project, 0, len(s.projects))
	nameLower := strings.ToLower(params.Name)
	uuidLower := strings.ToLower(params.UUID)
	for _, p := range s.projects {
		if nameLower != "" && !strings.Contains(strings.ToLower(p.Name), nameLower) {
			continue
		}
		if uuidLower != "" && !strings.Contains(strings.ToLower(p.UUID), uuidLower) {
			continue
		}
		p.AdminCount = s.memoryAdminCount(p.ID)
		p.ResourceCount = s.memoryResourceCount(p.ID)
		filtered = append(filtered, p)
	}

	total := uint64(len(filtered))
	start := int((params.Page - 1) * params.PageSize)
	if start >= len(filtered) {
		return []Project{}, total, nil
	}
	end := start + int(params.PageSize)
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[start:end], total, nil
}

func (s *Service) GetByUserID(userID uint64, name string) ([]Project, error) {
	name = strings.TrimSpace(name)
	if s.useDB {
		return s.getByUserIDFromDB(userID, name)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	nameLower := strings.ToLower(name)
	out := make([]Project, 0, len(s.projects))
	for _, p := range s.projects {
		if nameLower != "" && !strings.Contains(strings.ToLower(p.Name), nameLower) {
			continue
		}
		if userID > 0 && p.IsSystem != 1 && !s.memoryIsProjectAdmin(p.ID, userID) {
			continue
		}
		p.AdminCount = s.memoryAdminCount(p.ID)
		p.ResourceCount = s.memoryResourceCount(p.ID)
		out = append(out, p)
	}
	return out, nil
}

func (s *Service) AddAdmins(projectID uint64, userIDs []uint64) error {
	if projectID == 0 {
		return errors.New("projectId is required")
	}
	if len(userIDs) == 0 {
		return nil
	}

	if s.useDB {
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		defer func() { _ = tx.Rollback() }()

		stmt, err := tx.Prepare(`
INSERT INTO onec_project_admin (project_id, user_id, is_deleted, created_at)
VALUES (?, ?, 0, CURRENT_TIMESTAMP)
ON DUPLICATE KEY UPDATE is_deleted = 0
`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, userID := range userIDs {
			if userID == 0 {
				continue
			}
			if _, err := stmt.Exec(projectID, userID); err != nil {
				return err
			}
		}
		return tx.Commit()
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, uid := range userIDs {
		if uid == 0 {
			continue
		}
		exist := false
		for _, admin := range s.admins {
			if admin.ProjectID == projectID && admin.UserID == uid {
				exist = true
				break
			}
		}
		if exist {
			continue
		}
		s.admins = append(s.admins, ProjectAdmin{
			ID:        s.nextAdmin,
			ProjectID: projectID,
			UserID:    uid,
			CreatedAt: time.Now().Unix(),
		})
		s.nextAdmin++
	}
	return nil
}

func (s *Service) GetAdmins(projectID uint64) ([]ProjectAdmin, error) {
	if projectID == 0 {
		return []ProjectAdmin{}, nil
	}
	if s.useDB {
		rows, err := s.db.Query(`
SELECT id, project_id, user_id, created_at
FROM onec_project_admin
WHERE project_id = ? AND is_deleted = 0
ORDER BY id ASC
`, projectID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		out := make([]ProjectAdmin, 0)
		for rows.Next() {
			var item ProjectAdmin
			var createdAt time.Time
			if err := rows.Scan(&item.ID, &item.ProjectID, &item.UserID, &createdAt); err != nil {
				return nil, err
			}
			item.CreatedAt = createdAt.Unix()
			out = append(out, item)
		}
		return out, rows.Err()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]ProjectAdmin, 0)
	for _, a := range s.admins {
		if a.ProjectID == projectID {
			out = append(out, a)
		}
	}
	return out, nil
}

func (s *Service) addToDB(params AddProjectParams) (Project, error) {
	projectUUID := uuid.NewString()
	_, err := s.db.Exec(`
INSERT INTO onec_project (name, uuid, description, is_system, created_by, updated_by, is_deleted)
VALUES (?, ?, ?, ?, ?, ?, 0)
`, params.Name, projectUUID, params.Description, params.IsSystem, params.Operator, params.Operator)
	if err != nil {
		return Project{}, err
	}

	row := s.db.QueryRow(`
SELECT id, name, uuid, description, is_system, created_by, updated_by, created_at, updated_at
FROM onec_project
WHERE uuid = ? AND is_deleted = 0
LIMIT 1
`, projectUUID)

	var p Project
	var createdAt time.Time
	var updatedAt time.Time
	if err := row.Scan(&p.ID, &p.Name, &p.UUID, &p.Description, &p.IsSystem, &p.CreatedBy, &p.UpdatedBy, &createdAt, &updatedAt); err != nil {
		return Project{}, err
	}
	p.CreatedAt = createdAt.Unix()
	p.UpdatedAt = updatedAt.Unix()
	return p, nil
}

func (s *Service) updateToDB(params UpdateProjectParams) error {
	result, err := s.db.Exec(`
UPDATE onec_project
SET name = ?, description = ?, updated_by = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND is_deleted = 0
`, params.Name, params.Description, params.Operator, params.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("project not found")
	}
	return nil
}

func (s *Service) getByIDFromDB(id uint64) (Project, bool, error) {
	row := s.db.QueryRow(`
SELECT p.id, p.name, p.uuid, p.description, p.is_system, p.created_by, p.updated_by, p.created_at, p.updated_at,
       (SELECT COUNT(1) FROM onec_project_admin a WHERE a.project_id = p.id AND a.is_deleted = 0) AS admin_count
FROM onec_project p
WHERE p.id = ? AND p.is_deleted = 0
LIMIT 1
`, id)

	var p Project
	var createdAt time.Time
	var updatedAt time.Time
	if err := row.Scan(&p.ID, &p.Name, &p.UUID, &p.Description, &p.IsSystem, &p.CreatedBy, &p.UpdatedBy, &createdAt, &updatedAt, &p.AdminCount); err != nil {
		if err == sql.ErrNoRows {
			return Project{}, false, nil
		}
		return Project{}, false, err
	}
	p.CreatedAt = createdAt.Unix()
	p.UpdatedAt = updatedAt.Unix()
	p.ResourceCount = 0
	return p, true, nil
}

func (s *Service) searchFromDB(params SearchParams) ([]Project, uint64, error) {
	baseWhere := " FROM onec_project p WHERE p.is_deleted = 0"
	args := make([]any, 0, 4)
	if params.Name != "" {
		baseWhere += " AND p.name LIKE ?"
		args = append(args, "%"+params.Name+"%")
	}
	if params.UUID != "" {
		baseWhere += " AND p.uuid LIKE ?"
		args = append(args, "%"+params.UUID+"%")
	}

	countQuery := "SELECT COUNT(1)" + baseWhere
	var total uint64
	if err := s.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := `
SELECT p.id, p.name, p.uuid, p.description, p.is_system, p.created_by, p.updated_by, p.created_at, p.updated_at,
       (SELECT COUNT(1) FROM onec_project_admin a WHERE a.project_id = p.id AND a.is_deleted = 0) AS admin_count
` + baseWhere + `
 ORDER BY p.id DESC
 LIMIT ? OFFSET ?
`
	offset := (params.Page - 1) * params.PageSize
	listArgs := append(args, params.PageSize, offset)
	rows, err := s.db.Query(listQuery, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]Project, 0)
	for rows.Next() {
		var p Project
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(&p.ID, &p.Name, &p.UUID, &p.Description, &p.IsSystem, &p.CreatedBy, &p.UpdatedBy, &createdAt, &updatedAt, &p.AdminCount); err != nil {
			return nil, 0, err
		}
		p.CreatedAt = createdAt.Unix()
		p.UpdatedAt = updatedAt.Unix()
		p.ResourceCount = 0
		items = append(items, p)
	}
	return items, total, rows.Err()
}

func (s *Service) getByUserIDFromDB(userID uint64, name string) ([]Project, error) {
	query := `
SELECT p.id, p.name, p.uuid, p.description, p.is_system, p.created_by, p.updated_by, p.created_at, p.updated_at,
       (SELECT COUNT(1) FROM onec_project_admin a WHERE a.project_id = p.id AND a.is_deleted = 0) AS admin_count
FROM onec_project p
WHERE p.is_deleted = 0
`
	args := make([]any, 0, 2)
	if userID > 0 {
		query += `
 AND (
      p.is_system = 1 OR EXISTS (
        SELECT 1 FROM onec_project_admin a
        WHERE a.project_id = p.id AND a.user_id = ? AND a.is_deleted = 0
      )
 )
`
		args = append(args, userID)
	}
	if name != "" {
		query += " AND p.name LIKE ?"
		args = append(args, "%"+name+"%")
	}
	query += " ORDER BY p.id DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Project, 0)
	for rows.Next() {
		var p Project
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(&p.ID, &p.Name, &p.UUID, &p.Description, &p.IsSystem, &p.CreatedBy, &p.UpdatedBy, &createdAt, &updatedAt, &p.AdminCount); err != nil {
			return nil, err
		}
		p.CreatedAt = createdAt.Unix()
		p.UpdatedAt = updatedAt.Unix()
		p.ResourceCount = 0
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Service) memoryAdminCount(projectID uint64) int64 {
	var count int64
	for _, a := range s.admins {
		if a.ProjectID == projectID {
			count++
		}
	}
	return count
}

func (s *Service) memoryIsProjectAdmin(projectID, userID uint64) bool {
	for _, a := range s.admins {
		if a.ProjectID == projectID && a.UserID == userID {
			return true
		}
	}
	return false
}

func (s *Service) memoryResourceCount(projectID uint64) int64 {
	var count int64
	for _, item := range s.projectClusters {
		if item.ProjectID == projectID {
			count++
		}
	}
	return count
}

func defaultProjects() []Project {
	return []Project{
		{
			ID:            1,
			Name:          "系统项目",
			UUID:          "00000000-0000-0000-0000-000000000001",
			Description:   "平台系统默认项目",
			IsSystem:      1,
			CreatedBy:     "system",
			UpdatedBy:     "system",
			CreatedAt:     1735689600,
			UpdatedAt:     1735689600,
			AdminCount:    1,
			ResourceCount: 0,
		},
		{
			ID:            2,
			Name:          "研发项目",
			UUID:          "00000000-0000-0000-0000-000000000002",
			Description:   "默认研发项目",
			IsSystem:      0,
			CreatedBy:     "system",
			UpdatedBy:     "system",
			CreatedAt:     1735776000,
			UpdatedAt:     1735776000,
			AdminCount:    1,
			ResourceCount: 0,
		},
	}
}

func defaultAdmins() []ProjectAdmin {
	return []ProjectAdmin{
		{ID: 1, ProjectID: 1, UserID: 1, CreatedAt: 1735689600},
		{ID: 2, ProjectID: 2, UserID: 1, CreatedAt: 1735776000},
	}
}
