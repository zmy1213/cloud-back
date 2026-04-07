package clusterrepo

import (
	"context"
	"database/sql"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	appcfg "github.com/yanshicheng/cloud-back/common/config"
)

// Cluster models a registered Kubernetes cluster endpoint.
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

// Service provides cross-cluster registry queries.
// For scaffold stage this is in-memory and can be replaced by DB later.
type Service struct {
	clusters []Cluster
	db       *sql.DB
	useDB    bool
}

func NewService(mysqlCfg appcfg.MysqlConfig) *Service {
	s := &Service{
		clusters: defaultClusters(),
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

func (s *Service) Search(name, environment string) []Cluster {
	if s.useDB {
		if items, err := s.searchFromDB(name, environment); err == nil {
			return items
		}
	}

	name = strings.ToLower(strings.TrimSpace(name))
	environment = strings.ToLower(strings.TrimSpace(environment))

	out := make([]Cluster, 0, len(s.clusters))
	for _, c := range s.clusters {
		if name != "" && !strings.Contains(strings.ToLower(c.Name), name) {
			continue
		}
		if environment != "" && strings.ToLower(c.Environment) != environment {
			continue
		}
		out = append(out, c)
	}
	return out
}

func (s *Service) GetByID(id uint64) (Cluster, bool) {
	if s.useDB {
		if item, ok := s.getByIDFromDB(id); ok {
			return item, true
		}
	}

	for _, c := range s.clusters {
		if c.ID == id {
			return c, true
		}
	}
	return Cluster{}, false
}

func (s *Service) GetByUUID(uuid string) (Cluster, bool) {
	if s.useDB {
		if item, ok := s.getByUUIDFromDB(uuid); ok {
			return item, true
		}
	}

	for _, c := range s.clusters {
		if c.UUID == uuid {
			return c, true
		}
	}
	return Cluster{}, false
}

func (s *Service) Count() int {
	if s.useDB {
		if n, ok := s.countFromDB(); ok {
			return n
		}
	}
	return len(s.clusters)
}

func (s *Service) searchFromDB(name, environment string) ([]Cluster, error) {
	name = strings.TrimSpace(name)
	environment = strings.TrimSpace(environment)

	query := `
SELECT id, name, avatar, environment, cluster_type, version, status, health_status, uuid,
       cpu_usage, memory_usage, pod_usage, storage_usage, created_at
FROM onec_cluster
WHERE is_deleted = 0
`
	args := make([]any, 0, 2)
	if name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+name+"%")
	}
	if environment != "" {
		query += " AND environment = ?"
		args = append(args, environment)
	}
	query += " ORDER BY id ASC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Cluster, 0)
	for rows.Next() {
		c, err := scanCluster(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	return items, rows.Err()
}

func (s *Service) getByIDFromDB(id uint64) (Cluster, bool) {
	query := `
SELECT id, name, avatar, environment, cluster_type, version, status, health_status, uuid,
       cpu_usage, memory_usage, pod_usage, storage_usage, created_at
FROM onec_cluster
WHERE id = ? AND is_deleted = 0
LIMIT 1
`
	row := s.db.QueryRow(query, id)
	c, err := scanCluster(row)
	if err != nil {
		return Cluster{}, false
	}
	return c, true
}

func (s *Service) getByUUIDFromDB(uuid string) (Cluster, bool) {
	query := `
SELECT id, name, avatar, environment, cluster_type, version, status, health_status, uuid,
       cpu_usage, memory_usage, pod_usage, storage_usage, created_at
FROM onec_cluster
WHERE uuid = ? AND is_deleted = 0
LIMIT 1
`
	row := s.db.QueryRow(query, uuid)
	c, err := scanCluster(row)
	if err != nil {
		return Cluster{}, false
	}
	return c, true
}

func (s *Service) countFromDB() (int, bool) {
	row := s.db.QueryRow("SELECT COUNT(1) FROM onec_cluster WHERE is_deleted = 0")
	var n int
	if err := row.Scan(&n); err != nil {
		return 0, false
	}
	return n, true
}

type clusterScanner interface {
	Scan(dest ...any) error
}

func scanCluster(s clusterScanner) (Cluster, error) {
	var c Cluster
	var createdAt time.Time
	err := s.Scan(
		&c.ID, &c.Name, &c.Avatar, &c.Environment, &c.ClusterType, &c.Version,
		&c.Status, &c.HealthStatus, &c.UUID,
		&c.CpuUsage, &c.MemoryUsage, &c.PodUsage, &c.StorageUsage,
		&createdAt,
	)
	if err != nil {
		return Cluster{}, err
	}
	c.CreatedAt = createdAt.Unix()
	return c, nil
}

func defaultClusters() []Cluster {
	return []Cluster{
		{
			ID: 1, Name: "prod-hz", UUID: "11111111-1111-1111-1111-111111111111",
			Environment: "prod", ClusterType: "standard", Version: "v1.29.4",
			Status: 3, HealthStatus: 1, CpuUsage: 62.1, MemoryUsage: 57.4, PodUsage: 48.8, StorageUsage: 39.2, CreatedAt: 1735689600,
		},
		{
			ID: 2, Name: "staging-sh", UUID: "22222222-2222-2222-2222-222222222222",
			Environment: "staging", ClusterType: "standard", Version: "v1.28.7",
			Status: 3, HealthStatus: 1, CpuUsage: 38.7, MemoryUsage: 41.2, PodUsage: 28.5, StorageUsage: 33.9, CreatedAt: 1735776000,
		},
		{
			ID: 3, Name: "edge-gz", UUID: "33333333-3333-3333-3333-333333333333",
			Environment: "edge", ClusterType: "edge", Version: "v1.27.12",
			Status: 1, HealthStatus: 2, CpuUsage: 71.5, MemoryUsage: 66.8, PodUsage: 59.9, StorageUsage: 44.7, CreatedAt: 1735862400,
		},
	}
}
