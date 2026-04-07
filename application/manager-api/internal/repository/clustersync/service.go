package clustersyncrepo

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	appcfg "github.com/yanshicheng/cloud-back/common/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Service struct {
	db      *sql.DB
	useDB   bool
	k8sConf appcfg.K8sConfig
}

type SyncResult struct {
	ClusterID   uint64 `json:"clusterId"`
	ClusterUUID string `json:"clusterUuid"`
	ClusterName string `json:"clusterName"`
	NodeCount   int    `json:"nodeCount"`
	Source      string `json:"source"`
}

type clusterSnapshot struct {
	UUID         string
	Name         string
	Environment  string
	ClusterType  string
	Version      string
	Status       int64
	HealthStatus int64
	CpuUsage     float64
	MemoryUsage  float64
	PodUsage     float64
	StorageUsage float64
	Nodes        []nodeSnapshot
}

type nodeSnapshot struct {
	ClusterUUID     string
	NodeUUID        string
	Name            string
	Hostname        string
	Roles           string
	OSImage         string
	NodeIP          string
	KernelVersion   string
	OperatingSystem string
	Architecture    string
	CPU             float64
	Memory          float64
	Pods            int64
	IsGPU           int64
	Runtime         string
	JoinAt          time.Time
	Unschedulable   int64
	KubeletVersion  string
	Status          string
	PodCIDR         string
	PodCIDRs        string
	CreatedBy       string
	UpdatedBy       string
}

func NewService(mysqlCfg appcfg.MysqlConfig, k8sCfg appcfg.K8sConfig) *Service {
	s := &Service{k8sConf: k8sCfg}
	if !mysqlCfg.Enabled || strings.TrimSpace(mysqlCfg.DataSource) == "" {
		return s
	}

	db, err := sql.Open("mysql", mysqlCfg.DataSource)
	if err != nil {
		log.Printf("[clustersync] init_db_error=%v", err)
		return s
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Printf("[clustersync] ping_db_error=%v", err)
		_ = db.Close()
		return s
	}

	s.db = db
	s.useDB = true
	return s
}

func (s *Service) SyncAll(ctx context.Context, operator string) (*SyncResult, error) {
	return s.syncLocalCluster(ctx, operator)
}

func (s *Service) SyncByID(ctx context.Context, id uint64, operator string) (*SyncResult, error) {
	if !s.useDB {
		return nil, errors.New("sync requires mysql enabled")
	}
	var dbUUID string
	if err := s.db.QueryRowContext(ctx, "SELECT uuid FROM onec_cluster WHERE id = ? AND is_deleted = 0 LIMIT 1", id).Scan(&dbUUID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("cluster not found: id=%d", id)
		}
		return nil, fmt.Errorf("query cluster by id failed: %w", err)
	}

	result, err := s.syncLocalCluster(ctx, operator)
	if err != nil {
		return nil, err
	}
	if dbUUID != "" && result.ClusterUUID != dbUUID {
		return nil, fmt.Errorf("cluster id=%d uuid mismatch: db=%s synced=%s", id, dbUUID, result.ClusterUUID)
	}
	return result, nil
}

func (s *Service) syncLocalCluster(ctx context.Context, operator string) (*SyncResult, error) {
	if !s.useDB {
		return nil, errors.New("sync requires mysql enabled")
	}
	if !s.k8sConf.Enabled {
		return nil, errors.New("sync requires k8s.enabled=true")
	}
	if strings.TrimSpace(operator) == "" {
		operator = "system"
	}

	snapshot, err := s.loadSnapshotFromK8s(ctx, operator)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin tx failed: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	clusterID, err := upsertCluster(ctx, tx, snapshot)
	if err != nil {
		return nil, err
	}
	if err = replaceClusterNodes(ctx, tx, snapshot.UUID, snapshot.Nodes); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx failed: %w", err)
	}

	log.Printf(
		"[clustersync] op=Sync source=k8s_sync cluster_id=%d cluster_uuid=%s cluster_name=%q nodes=%d",
		clusterID, snapshot.UUID, snapshot.Name, len(snapshot.Nodes),
	)
	return &SyncResult{
		ClusterID:   clusterID,
		ClusterUUID: snapshot.UUID,
		ClusterName: snapshot.Name,
		NodeCount:   len(snapshot.Nodes),
		Source:      "k8s_sync",
	}, nil
}

func (s *Service) loadSnapshotFromK8s(ctx context.Context, operator string) (*clusterSnapshot, error) {
	restConf, contextName, apiServerHost, err := buildRESTConfigAndMeta(s.k8sConf)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		return nil, fmt.Errorf("create k8s client failed: %w", err)
	}

	versionInfo, err := clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("get k8s version failed: %w", err)
	}

	nodeList, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list k8s nodes failed: %w", err)
	}

	clusterName := strings.TrimSpace(s.k8sConf.ClusterName)
	if clusterName == "" {
		clusterName = contextName
	}
	if clusterName == "" {
		clusterName = "local-cluster"
	}
	environment := strings.TrimSpace(s.k8sConf.Environment)
	if environment == "" {
		environment = "prod"
	}
	clusterType := strings.TrimSpace(s.k8sConf.ClusterType)
	if clusterType == "" {
		clusterType = "standard"
	}

	clusterUUID := stableClusterUUID(contextName + "|" + apiServerHost)
	nodeSnapshots := make([]nodeSnapshot, 0, len(nodeList.Items))

	var (
		totalCPUCap      float64
		totalCPUAlloc    float64
		totalMemCapGi    float64
		totalMemAllocGi  float64
		totalPodsCap     int64
		totalPodsAlloc   int64
		readyNodeCount   int
		notReadyNodeSeen bool
	)

	for _, node := range nodeList.Items {
		roles := rolesFromNode(node)
		nodeIP := nodeIPFromNode(node)
		nodeStatus, isReady := readinessFromNode(node)
		if isReady {
			readyNodeCount++
		} else {
			notReadyNodeSeen = true
		}

		cpuCap := cpuCores(node.Status.Capacity[corev1.ResourceCPU])
		cpuAlloc := cpuCores(node.Status.Allocatable[corev1.ResourceCPU])
		memCapGi := memoryGi(node.Status.Capacity[corev1.ResourceMemory])
		memAllocGi := memoryGi(node.Status.Allocatable[corev1.ResourceMemory])
		podsCap := quantityInt64(node.Status.Capacity[corev1.ResourcePods])
		podsAlloc := quantityInt64(node.Status.Allocatable[corev1.ResourcePods])

		totalCPUCap += cpuCap
		totalCPUAlloc += cpuAlloc
		totalMemCapGi += memCapGi
		totalMemAllocGi += memAllocGi
		totalPodsCap += podsCap
		totalPodsAlloc += podsAlloc

		nodeSnapshots = append(nodeSnapshots, nodeSnapshot{
			ClusterUUID:     clusterUUID,
			NodeUUID:        string(node.UID),
			Name:            node.Name,
			Hostname:        node.Labels["kubernetes.io/hostname"],
			Roles:           roles,
			OSImage:         node.Status.NodeInfo.OSImage,
			NodeIP:          nodeIP,
			KernelVersion:   node.Status.NodeInfo.KernelVersion,
			OperatingSystem: node.Status.NodeInfo.OperatingSystem,
			Architecture:    node.Status.NodeInfo.Architecture,
			CPU:             cpuCap,
			Memory:          memCapGi,
			Pods:            podsCap,
			IsGPU:           gpuFlag(node),
			Runtime:         node.Status.NodeInfo.ContainerRuntimeVersion,
			JoinAt:          node.CreationTimestamp.Time,
			Unschedulable:   boolToInt64(node.Spec.Unschedulable),
			KubeletVersion:  node.Status.NodeInfo.KubeletVersion,
			Status:          nodeStatus,
			PodCIDR:         node.Spec.PodCIDR,
			PodCIDRs:        strings.Join(node.Spec.PodCIDRs, ","),
			CreatedBy:       operator,
			UpdatedBy:       operator,
		})
	}

	healthStatus := int64(1)
	if len(nodeList.Items) == 0 || notReadyNodeSeen {
		healthStatus = 2
	}
	status := int64(3)
	if len(nodeList.Items) == 0 || readyNodeCount == 0 {
		status = 2
	}

	return &clusterSnapshot{
		UUID:         clusterUUID,
		Name:         clusterName,
		Environment:  environment,
		ClusterType:  clusterType,
		Version:      strings.TrimSpace(versionInfo.GitVersion),
		Status:       status,
		HealthStatus: healthStatus,
		CpuUsage:     usagePercent(totalCPUCap-totalCPUAlloc, totalCPUCap),
		MemoryUsage:  usagePercent(totalMemCapGi-totalMemAllocGi, totalMemCapGi),
		PodUsage:     usagePercent(float64(totalPodsCap-totalPodsAlloc), float64(totalPodsCap)),
		StorageUsage: 0,
		Nodes:        nodeSnapshots,
	}, nil
}

func buildRESTConfigAndMeta(cfg appcfg.K8sConfig) (*rest.Config, string, string, error) {
	mode := strings.ToLower(strings.TrimSpace(cfg.Mode))
	if mode == "" {
		mode = "kubeconfig"
	}
	switch mode {
	case "incluster":
		restConf, err := rest.InClusterConfig()
		if err != nil {
			return nil, "", "", fmt.Errorf("load incluster config failed: %w", err)
		}
		return restConf, "incluster", restConf.Host, nil
	case "kubeconfig":
		kubeconfig := strings.TrimSpace(cfg.Kubeconfig)
		if kubeconfig == "" {
			home := homedir.HomeDir()
			if home == "" {
				return nil, "", "", errors.New("home directory not found for default kubeconfig")
			}
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
		if _, err := os.Stat(kubeconfig); err != nil {
			return nil, "", "", fmt.Errorf("kubeconfig not found: %s", kubeconfig)
		}

		loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
		overrides := &clientcmd.ConfigOverrides{}
		if strings.TrimSpace(cfg.Context) != "" {
			overrides.CurrentContext = strings.TrimSpace(cfg.Context)
		}

		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
		restConf, err := clientConfig.ClientConfig()
		if err != nil {
			return nil, "", "", fmt.Errorf("build rest config from kubeconfig failed: %w", err)
		}
		rawCfg, err := clientConfig.RawConfig()
		if err != nil {
			return nil, "", "", fmt.Errorf("read kubeconfig metadata failed: %w", err)
		}

		contextName := strings.TrimSpace(cfg.Context)
		if contextName == "" {
			contextName = rawCfg.CurrentContext
		}
		return restConf, contextName, restConf.Host, nil
	default:
		return nil, "", "", fmt.Errorf("unsupported k8s mode: %s", cfg.Mode)
	}
}

func upsertCluster(ctx context.Context, tx *sql.Tx, snap *clusterSnapshot) (uint64, error) {
	query := `
INSERT INTO onec_cluster
  (name, avatar, environment, cluster_type, version, status, health_status, uuid,
   cpu_usage, memory_usage, pod_usage, storage_usage, is_deleted)
VALUES
  (?, '', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  environment = VALUES(environment),
  cluster_type = VALUES(cluster_type),
  version = VALUES(version),
  status = VALUES(status),
  health_status = VALUES(health_status),
  cpu_usage = VALUES(cpu_usage),
  memory_usage = VALUES(memory_usage),
  pod_usage = VALUES(pod_usage),
  storage_usage = VALUES(storage_usage),
  is_deleted = 0,
  updated_at = CURRENT_TIMESTAMP
`
	if _, err := tx.ExecContext(
		ctx,
		query,
		snap.Name,
		snap.Environment,
		snap.ClusterType,
		snap.Version,
		snap.Status,
		snap.HealthStatus,
		snap.UUID,
		snap.CpuUsage,
		snap.MemoryUsage,
		snap.PodUsage,
		snap.StorageUsage,
	); err != nil {
		return 0, fmt.Errorf("upsert cluster failed: %w", err)
	}

	var clusterID uint64
	if err := tx.QueryRowContext(ctx, "SELECT id FROM onec_cluster WHERE uuid = ? AND is_deleted = 0 LIMIT 1", snap.UUID).Scan(&clusterID); err != nil {
		return 0, fmt.Errorf("query upserted cluster id failed: %w", err)
	}
	return clusterID, nil
}

func replaceClusterNodes(ctx context.Context, tx *sql.Tx, clusterUUID string, nodes []nodeSnapshot) error {
	if _, err := tx.ExecContext(ctx, "DELETE FROM onec_cluster_node WHERE cluster_uuid = ?", clusterUUID); err != nil {
		return fmt.Errorf("delete old cluster nodes failed: %w", err)
	}

	if len(nodes) == 0 {
		return nil
	}

	query := `
INSERT INTO onec_cluster_node
  (cluster_uuid, node_uuid, name, hostname, roles, os_image, node_ip, kernel_version,
   operating_system, architecture, cpu, memory, pods, is_gpu, runtime, join_at,
   unschedulable, kubelet_version, status, pod_cidr, pod_cidrs, created_by, updated_by, is_deleted)
VALUES
  (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
ON DUPLICATE KEY UPDATE
  cluster_uuid = VALUES(cluster_uuid),
  name = VALUES(name),
  hostname = VALUES(hostname),
  roles = VALUES(roles),
  os_image = VALUES(os_image),
  node_ip = VALUES(node_ip),
  kernel_version = VALUES(kernel_version),
  operating_system = VALUES(operating_system),
  architecture = VALUES(architecture),
  cpu = VALUES(cpu),
  memory = VALUES(memory),
  pods = VALUES(pods),
  is_gpu = VALUES(is_gpu),
  runtime = VALUES(runtime),
  join_at = VALUES(join_at),
  unschedulable = VALUES(unschedulable),
  kubelet_version = VALUES(kubelet_version),
  status = VALUES(status),
  pod_cidr = VALUES(pod_cidr),
  pod_cidrs = VALUES(pod_cidrs),
  updated_by = VALUES(updated_by),
  is_deleted = 0,
  updated_at = CURRENT_TIMESTAMP
`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare node upsert failed: %w", err)
	}
	defer stmt.Close()

	for _, node := range nodes {
		joinAt := node.JoinAt
		if joinAt.IsZero() {
			joinAt = time.Now()
		}

		if _, err := stmt.ExecContext(
			ctx,
			node.ClusterUUID,
			node.NodeUUID,
			node.Name,
			node.Hostname,
			node.Roles,
			node.OSImage,
			node.NodeIP,
			node.KernelVersion,
			node.OperatingSystem,
			node.Architecture,
			node.CPU,
			node.Memory,
			node.Pods,
			node.IsGPU,
			node.Runtime,
			joinAt,
			node.Unschedulable,
			node.KubeletVersion,
			node.Status,
			node.PodCIDR,
			node.PodCIDRs,
			node.CreatedBy,
			node.UpdatedBy,
		); err != nil {
			return fmt.Errorf("upsert node %s failed: %w", node.Name, err)
		}
	}
	return nil
}

func usagePercent(used, total float64) float64 {
	if total <= 0 {
		return 0
	}
	value := (used / total) * 100
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func cpuCores(qty resource.Quantity) float64 {
	return float64(qty.MilliValue()) / 1000.0
}

func memoryGi(qty resource.Quantity) float64 {
	const gib = 1024.0 * 1024.0 * 1024.0
	return float64(qty.Value()) / gib
}

func quantityInt64(qty resource.Quantity) int64 {
	return qty.Value()
}

func boolToInt64(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

func rolesFromNode(node corev1.Node) string {
	roles := make([]string, 0, 2)
	for key := range node.Labels {
		if strings.HasPrefix(key, "node-role.kubernetes.io/") {
			role := strings.TrimPrefix(key, "node-role.kubernetes.io/")
			if role == "" {
				role = "worker"
			}
			roles = append(roles, role)
		}
	}
	if len(roles) == 0 {
		if role := strings.TrimSpace(node.Labels["kubernetes.io/role"]); role != "" {
			roles = append(roles, role)
		}
	}
	if len(roles) == 0 {
		roles = append(roles, "worker")
	}

	sort.Strings(roles)
	return strings.Join(roles, ",")
}

func nodeIPFromNode(node corev1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			return addr.Address
		}
	}
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeExternalIP {
			return addr.Address
		}
	}
	return ""
}

func readinessFromNode(node corev1.Node) (string, bool) {
	for _, cond := range node.Status.Conditions {
		if cond.Type == corev1.NodeReady {
			if cond.Status == corev1.ConditionTrue {
				return "Ready", true
			}
			return "NotReady", false
		}
	}
	return "Unknown", false
}

func gpuFlag(node corev1.Node) int64 {
	if q, ok := node.Status.Capacity["nvidia.com/gpu"]; ok && q.Value() > 0 {
		return 1
	}
	return 0
}

func stableClusterUUID(seed string) string {
	base := strings.TrimSpace(seed)
	if base == "" {
		base = "local-cluster"
	}
	sum := sha1.Sum([]byte(base))
	hex := fmt.Sprintf("%x", sum)
	return fmt.Sprintf("%s-%s-%s-%s-%s", hex[0:8], hex[8:12], hex[12:16], hex[16:20], hex[20:32])
}
