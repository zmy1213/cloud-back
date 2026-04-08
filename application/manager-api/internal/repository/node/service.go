package noderepo

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	appcfg "github.com/yanshicheng/cloud-back/common/config"
)

type Node struct {
	ID              uint64
	ClusterUuid     string
	NodeUuid        string
	Name            string
	Hostname        string
	Roles           string
	OsImage         string
	NodeIp          string
	KernelVersion   string
	OperatingSystem string
	Architecture    string
	Cpu             float64
	Memory          float64
	Pods            int64
	IsGpu           int64
	Runtime         string
	JoinAt          int64
	Unschedulable   int64
	KubeletVersion  string
	Status          string
	PodCidr         string
	PodCidrs        string
	CreatedBy       string
	UpdatedBy       string
	CreatedAt       int64
	UpdatedAt       int64
}

type SearchParams struct {
	Page        uint64
	PageSize    uint64
	OrderField  string
	IsAsc       bool
	ClusterUuid string
}

type Service struct {
	nodes []Node
	db    *sql.DB
	useDB bool
}

func NewService(mysqlCfg appcfg.MysqlConfig) *Service {
	s := &Service{nodes: defaultNodes()}
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

func (s *Service) Search(params SearchParams) ([]Node, uint64) {
	params = normalizeSearchParams(params)
	if s.useDB {
		items, total, err := s.searchFromDB(params)
		if err != nil {
			log.Printf(
				"[noderepo] op=Search source=db cluster_uuid=%q page=%d page_size=%d query_error=true err=%v",
				params.ClusterUuid, params.Page, params.PageSize, err,
			)
			return []Node{}, 0
		}
		log.Printf(
			"[noderepo] op=Search source=db cluster_uuid=%q page=%d page_size=%d total=%d count=%d",
			params.ClusterUuid, params.Page, params.PageSize, total, len(items),
		)
		return items, total
	}

	clusterUUID := strings.TrimSpace(params.ClusterUuid)
	filtered := make([]Node, 0, len(s.nodes))
	for _, node := range s.nodes {
		if clusterUUID != "" && node.ClusterUuid != clusterUUID {
			continue
		}
		filtered = append(filtered, node)
	}

	sortNodes(filtered, params.OrderField, params.IsAsc)
	total := uint64(len(filtered))
	start := int((params.Page - 1) * params.PageSize)
	if start >= len(filtered) {
		log.Printf(
			"[noderepo] op=Search source=default fallback_reason=db_disabled cluster_uuid=%q page=%d page_size=%d total=%d count=0",
			params.ClusterUuid, params.Page, params.PageSize, total,
		)
		return []Node{}, total
	}
	end := start + int(params.PageSize)
	if end > len(filtered) {
		end = len(filtered)
	}
	items := filtered[start:end]
	log.Printf(
		"[noderepo] op=Search source=default fallback_reason=db_disabled cluster_uuid=%q page=%d page_size=%d total=%d count=%d",
		params.ClusterUuid, params.Page, params.PageSize, total, len(items),
	)
	return items, total
}

func (s *Service) GetByID(id uint64) (Node, bool) {
	if s.useDB {
		item, ok, err := s.getByIDFromDB(id)
		if err != nil {
			log.Printf("[noderepo] op=GetByID source=db id=%d query_error=true err=%v", id, err)
			return Node{}, false
		}
		log.Printf("[noderepo] op=GetByID source=db id=%d hit=%t", id, ok)
		return item, ok
	}
	for _, n := range s.nodes {
		if n.ID == id {
			log.Printf("[noderepo] op=GetByID source=default fallback_reason=db_disabled id=%d hit=true", id)
			return n, true
		}
	}
	log.Printf("[noderepo] op=GetByID source=default fallback_reason=db_disabled id=%d hit=false", id)
	return Node{}, false
}

func normalizeSearchParams(params SearchParams) SearchParams {
	if params.Page == 0 {
		params.Page = 1
	}
	if params.PageSize == 0 {
		params.PageSize = 10
	}
	if params.PageSize > 200 {
		params.PageSize = 200
	}
	if params.OrderField == "" {
		params.OrderField = "id"
	}
	return params
}

func sortNodes(nodes []Node, orderField string, isAsc bool) {
	less := func(a, b Node) bool {
		switch orderField {
		case "nodeName":
			return a.Name < b.Name
		case "nodeIp":
			return a.NodeIp < b.NodeIp
		case "nodeStatus":
			return a.Status < b.Status
		case "createdAt":
			return a.CreatedAt < b.CreatedAt
		case "updatedAt":
			return a.UpdatedAt < b.UpdatedAt
		default:
			return a.ID < b.ID
		}
	}

	// Stable insertion sort is enough for this small in-memory fallback.
	for i := 1; i < len(nodes); i++ {
		j := i
		for j > 0 {
			shouldSwap := less(nodes[j], nodes[j-1])
			if !isAsc {
				shouldSwap = less(nodes[j-1], nodes[j])
			}
			if !shouldSwap {
				break
			}
			nodes[j], nodes[j-1] = nodes[j-1], nodes[j]
			j--
		}
	}
}

func (s *Service) searchFromDB(params SearchParams) ([]Node, uint64, error) {
	where := " WHERE is_deleted = 0"
	args := make([]any, 0, 3)
	if strings.TrimSpace(params.ClusterUuid) != "" {
		where += " AND cluster_uuid = ?"
		args = append(args, params.ClusterUuid)
	}

	var total uint64
	countQuery := "SELECT COUNT(1) FROM onec_cluster_node" + where
	if err := s.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []Node{}, 0, nil
	}

	orderField := map[string]string{
		"id":         "id",
		"nodeName":   "name",
		"nodeIp":     "node_ip",
		"nodeStatus": "status",
		"createdAt":  "created_at",
		"updatedAt":  "updated_at",
	}[params.OrderField]
	if orderField == "" {
		orderField = "id"
	}
	orderDirection := "DESC"
	if params.IsAsc {
		orderDirection = "ASC"
	}

	offset := (params.Page - 1) * params.PageSize
	query := `
SELECT id, cluster_uuid, node_uuid, name, hostname, roles, os_image, node_ip, kernel_version,
       operating_system, architecture, cpu, memory, pods, is_gpu, runtime, join_at,
       unschedulable, kubelet_version, status, pod_cidr, pod_cidrs, created_by, updated_by,
       created_at, updated_at
FROM onec_cluster_node` + where + `
ORDER BY ` + orderField + ` ` + orderDirection + `
LIMIT ? OFFSET ?`

	queryArgs := append(append([]any{}, args...), params.PageSize, offset)
	rows, err := s.db.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]Node, 0)
	for rows.Next() {
		node, err := scanNode(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, node)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *Service) getByIDFromDB(id uint64) (Node, bool, error) {
	query := `
SELECT id, cluster_uuid, node_uuid, name, hostname, roles, os_image, node_ip, kernel_version,
       operating_system, architecture, cpu, memory, pods, is_gpu, runtime, join_at,
       unschedulable, kubelet_version, status, pod_cidr, pod_cidrs, created_by, updated_by,
       created_at, updated_at
FROM onec_cluster_node
WHERE id = ? AND is_deleted = 0
LIMIT 1`

	node, err := scanNode(s.db.QueryRow(query, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return Node{}, false, nil
		}
		return Node{}, false, err
	}
	return node, true, nil
}

type nodeScanner interface {
	Scan(dest ...any) error
}

func scanNode(s nodeScanner) (Node, error) {
	var node Node
	var joinAt time.Time
	var createdAt time.Time
	var updatedAt time.Time
	err := s.Scan(
		&node.ID,
		&node.ClusterUuid,
		&node.NodeUuid,
		&node.Name,
		&node.Hostname,
		&node.Roles,
		&node.OsImage,
		&node.NodeIp,
		&node.KernelVersion,
		&node.OperatingSystem,
		&node.Architecture,
		&node.Cpu,
		&node.Memory,
		&node.Pods,
		&node.IsGpu,
		&node.Runtime,
		&joinAt,
		&node.Unschedulable,
		&node.KubeletVersion,
		&node.Status,
		&node.PodCidr,
		&node.PodCidrs,
		&node.CreatedBy,
		&node.UpdatedBy,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return Node{}, err
	}
	node.JoinAt = joinAt.Unix()
	node.CreatedAt = createdAt.Unix()
	node.UpdatedAt = updatedAt.Unix()
	return node, nil
}

func defaultNodes() []Node {
	return []Node{
		{
			ID:              1,
			ClusterUuid:     "11111111-1111-1111-1111-111111111111",
			NodeUuid:        "8f6a1d2e-a234-47af-a0a4-d4c7ef230001",
			Name:            "prod-hz-master-01",
			Hostname:        "prod-hz-master-01",
			Roles:           "control-plane,master",
			OsImage:         "Ubuntu 22.04.4 LTS",
			NodeIp:          "10.0.10.11",
			KernelVersion:   "6.5.0-28-generic",
			OperatingSystem: "linux",
			Architecture:    "amd64",
			Cpu:             16,
			Memory:          64,
			Pods:            110,
			IsGpu:           0,
			Runtime:         "containerd://1.7.12",
			JoinAt:          1735689600,
			Unschedulable:   1,
			KubeletVersion:  "v1.29.4",
			Status:          "Ready",
			PodCidr:         "10.244.0.0/24",
			PodCidrs:        "10.244.0.0/24",
			CreatedBy:       "system",
			UpdatedBy:       "system",
			CreatedAt:       1735689600,
			UpdatedAt:       1735776000,
		},
		{
			ID:              2,
			ClusterUuid:     "11111111-1111-1111-1111-111111111111",
			NodeUuid:        "de3e6cb5-9681-4f86-b3cd-e4d887f00002",
			Name:            "prod-hz-worker-01",
			Hostname:        "prod-hz-worker-01",
			Roles:           "worker",
			OsImage:         "Ubuntu 22.04.4 LTS",
			NodeIp:          "10.0.10.21",
			KernelVersion:   "6.5.0-28-generic",
			OperatingSystem: "linux",
			Architecture:    "amd64",
			Cpu:             32,
			Memory:          128,
			Pods:            220,
			IsGpu:           1,
			Runtime:         "containerd://1.7.12",
			JoinAt:          1735689600,
			Unschedulable:   1,
			KubeletVersion:  "v1.29.4",
			Status:          "Ready",
			PodCidr:         "10.244.1.0/24",
			PodCidrs:        "10.244.1.0/24",
			CreatedBy:       "system",
			UpdatedBy:       "system",
			CreatedAt:       1735689600,
			UpdatedAt:       1735948800,
		},
		{
			ID:              3,
			ClusterUuid:     "22222222-2222-2222-2222-222222222222",
			NodeUuid:        "f2f49524-707e-4382-baf2-fea5e7000003",
			Name:            "staging-sh-master-01",
			Hostname:        "staging-sh-master-01",
			Roles:           "control-plane,master",
			OsImage:         "Ubuntu 22.04.4 LTS",
			NodeIp:          "10.1.10.11",
			KernelVersion:   "6.5.0-26-generic",
			OperatingSystem: "linux",
			Architecture:    "amd64",
			Cpu:             8,
			Memory:          32,
			Pods:            110,
			IsGpu:           0,
			Runtime:         "containerd://1.7.10",
			JoinAt:          1735776000,
			Unschedulable:   1,
			KubeletVersion:  "v1.28.7",
			Status:          "Ready",
			PodCidr:         "10.245.0.0/24",
			PodCidrs:        "10.245.0.0/24",
			CreatedBy:       "system",
			UpdatedBy:       "system",
			CreatedAt:       1735776000,
			UpdatedAt:       1735862400,
		},
		{
			ID:              4,
			ClusterUuid:     "22222222-2222-2222-2222-222222222222",
			NodeUuid:        "e3ddd318-1433-4b9d-b9b5-c42659000004",
			Name:            "staging-sh-worker-01",
			Hostname:        "staging-sh-worker-01",
			Roles:           "worker",
			OsImage:         "Ubuntu 20.04.6 LTS",
			NodeIp:          "10.1.10.21",
			KernelVersion:   "5.15.0-105-generic",
			OperatingSystem: "linux",
			Architecture:    "amd64",
			Cpu:             16,
			Memory:          64,
			Pods:            180,
			IsGpu:           0,
			Runtime:         "containerd://1.6.24",
			JoinAt:          1735776000,
			Unschedulable:   1,
			KubeletVersion:  "v1.28.7",
			Status:          "NotReady",
			PodCidr:         "10.245.1.0/24",
			PodCidrs:        "10.245.1.0/24",
			CreatedBy:       "system",
			UpdatedBy:       "system",
			CreatedAt:       1735776000,
			UpdatedAt:       1736035200,
		},
		{
			ID:              5,
			ClusterUuid:     "33333333-3333-3333-3333-333333333333",
			NodeUuid:        "5d9a08df-fbd4-4f33-8a23-98f832000005",
			Name:            "edge-gz-master-01",
			Hostname:        "edge-gz-master-01",
			Roles:           "control-plane,master",
			OsImage:         "CentOS Stream 9",
			NodeIp:          "10.2.10.11",
			KernelVersion:   "5.14.0-427.el9",
			OperatingSystem: "linux",
			Architecture:    "arm64",
			Cpu:             8,
			Memory:          16,
			Pods:            80,
			IsGpu:           0,
			Runtime:         "containerd://1.7.6",
			JoinAt:          1735862400,
			Unschedulable:   1,
			KubeletVersion:  "v1.27.12",
			Status:          "Ready",
			PodCidr:         "10.246.0.0/24",
			PodCidrs:        "10.246.0.0/24",
			CreatedBy:       "system",
			UpdatedBy:       "system",
			CreatedAt:       1735862400,
			UpdatedAt:       1736121600,
		},
	}
}
