CREATE DATABASE IF NOT EXISTS cloud_back DEFAULT CHARACTER SET utf8mb4;
USE cloud_back;

CREATE TABLE IF NOT EXISTS sys_user (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(64) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  nick_name VARCHAR(128) NOT NULL,
  role_code VARCHAR(64) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO sys_user (username, password, nick_name, role_code)
VALUES ('super_admin', 'admin123', 'Cloud Admin', 'super_admin')
ON DUPLICATE KEY UPDATE updated_at = CURRENT_TIMESTAMP;

CREATE TABLE IF NOT EXISTS onec_cluster (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL,
  avatar VARCHAR(255) NOT NULL DEFAULT '',
  environment VARCHAR(32) NOT NULL DEFAULT 'prod',
  cluster_type VARCHAR(32) NOT NULL DEFAULT 'standard',
  version VARCHAR(64) NOT NULL DEFAULT '',
  status TINYINT NOT NULL DEFAULT 3,
  health_status TINYINT NOT NULL DEFAULT 1,
  uuid VARCHAR(64) NOT NULL UNIQUE,
  cpu_usage DOUBLE NOT NULL DEFAULT 0,
  memory_usage DOUBLE NOT NULL DEFAULT 0,
  pod_usage DOUBLE NOT NULL DEFAULT 0,
  storage_usage DOUBLE NOT NULL DEFAULT 0,
  is_deleted TINYINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT INTO onec_cluster
  (id, name, environment, cluster_type, version, status, health_status, uuid, cpu_usage, memory_usage, pod_usage, storage_usage, is_deleted)
VALUES
  (1, 'prod-hz', 'prod', 'standard', 'v1.29.4', 3, 1, '11111111-1111-1111-1111-111111111111', 62.1, 57.4, 48.8, 39.2, 0),
  (2, 'staging-sh', 'staging', 'standard', 'v1.28.7', 3, 1, '22222222-2222-2222-2222-222222222222', 38.7, 41.2, 28.5, 33.9, 0),
  (3, 'edge-gz', 'edge', 'edge', 'v1.27.12', 1, 2, '33333333-3333-3333-3333-333333333333', 71.5, 66.8, 59.9, 44.7, 0)
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
  is_deleted = VALUES(is_deleted),
  updated_at = CURRENT_TIMESTAMP;

CREATE TABLE IF NOT EXISTS onec_cluster_node (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  cluster_uuid VARCHAR(64) NOT NULL,
  node_uuid CHAR(36) NOT NULL UNIQUE,
  name VARCHAR(64) NOT NULL,
  hostname VARCHAR(128) NOT NULL DEFAULT '',
  roles VARCHAR(255) NOT NULL DEFAULT '',
  os_image VARCHAR(128) NOT NULL DEFAULT '',
  node_ip VARCHAR(64) NOT NULL DEFAULT '',
  kernel_version VARCHAR(64) NOT NULL DEFAULT '',
  operating_system VARCHAR(64) NOT NULL DEFAULT '',
  architecture VARCHAR(32) NOT NULL DEFAULT '',
  cpu DOUBLE NOT NULL DEFAULT 0,
  memory DOUBLE NOT NULL DEFAULT 0,
  pods BIGINT NOT NULL DEFAULT 0,
  is_gpu TINYINT NOT NULL DEFAULT 0,
  runtime VARCHAR(128) NOT NULL DEFAULT '',
  join_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  unschedulable INT NOT NULL DEFAULT 1,
  kubelet_version VARCHAR(64) NOT NULL DEFAULT '',
  status VARCHAR(32) NOT NULL DEFAULT 'Unknown',
  pod_cidr VARCHAR(64) NOT NULL DEFAULT '',
  pod_cidrs VARCHAR(255) NOT NULL DEFAULT '',
  created_by VARCHAR(64) NOT NULL DEFAULT 'system',
  updated_by VARCHAR(64) NOT NULL DEFAULT 'system',
  is_deleted TINYINT NOT NULL DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_cluster_node_name (cluster_uuid, name),
  KEY idx_cluster_uuid (cluster_uuid),
  KEY idx_status (status),
  KEY idx_is_deleted (is_deleted)
);

INSERT INTO onec_cluster_node
  (id, cluster_uuid, node_uuid, name, hostname, roles, os_image, node_ip, kernel_version, operating_system, architecture, cpu, memory, pods, is_gpu, runtime, unschedulable, kubelet_version, status, pod_cidr, pod_cidrs, created_by, updated_by, is_deleted)
VALUES
  (1, '11111111-1111-1111-1111-111111111111', '8f6a1d2e-a234-47af-a0a4-d4c7ef230001', 'prod-hz-master-01', 'prod-hz-master-01', 'control-plane,master', 'Ubuntu 22.04.4 LTS', '10.0.10.11', '6.5.0-28-generic', 'linux', 'amd64', 16, 64, 110, 0, 'containerd://1.7.12', 1, 'v1.29.4', 'Ready', '10.244.0.0/24', '10.244.0.0/24', 'system', 'system', 0),
  (2, '11111111-1111-1111-1111-111111111111', 'de3e6cb5-9681-4f86-b3cd-e4d887f00002', 'prod-hz-worker-01', 'prod-hz-worker-01', 'worker', 'Ubuntu 22.04.4 LTS', '10.0.10.21', '6.5.0-28-generic', 'linux', 'amd64', 32, 128, 220, 1, 'containerd://1.7.12', 1, 'v1.29.4', 'Ready', '10.244.1.0/24', '10.244.1.0/24', 'system', 'system', 0),
  (3, '22222222-2222-2222-2222-222222222222', 'f2f49524-707e-4382-baf2-fea5e7000003', 'staging-sh-master-01', 'staging-sh-master-01', 'control-plane,master', 'Ubuntu 22.04.4 LTS', '10.1.10.11', '6.5.0-26-generic', 'linux', 'amd64', 8, 32, 110, 0, 'containerd://1.7.10', 1, 'v1.28.7', 'Ready', '10.245.0.0/24', '10.245.0.0/24', 'system', 'system', 0),
  (4, '22222222-2222-2222-2222-222222222222', 'e3ddd318-1433-4b9d-b9b5-c42659000004', 'staging-sh-worker-01', 'staging-sh-worker-01', 'worker', 'Ubuntu 20.04.6 LTS', '10.1.10.21', '5.15.0-105-generic', 'linux', 'amd64', 16, 64, 180, 0, 'containerd://1.6.24', 1, 'v1.28.7', 'NotReady', '10.245.1.0/24', '10.245.1.0/24', 'system', 'system', 0),
  (5, '33333333-3333-3333-3333-333333333333', '5d9a08df-fbd4-4f33-8a23-98f832000005', 'edge-gz-master-01', 'edge-gz-master-01', 'control-plane,master', 'CentOS Stream 9', '10.2.10.11', '5.14.0-427.el9', 'linux', 'arm64', 8, 16, 80, 0, 'containerd://1.7.6', 1, 'v1.27.12', 'Ready', '10.246.0.0/24', '10.246.0.0/24', 'system', 'system', 0)
ON DUPLICATE KEY UPDATE
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
  unschedulable = VALUES(unschedulable),
  kubelet_version = VALUES(kubelet_version),
  status = VALUES(status),
  pod_cidr = VALUES(pod_cidr),
  pod_cidrs = VALUES(pod_cidrs),
  updated_by = VALUES(updated_by),
  is_deleted = VALUES(is_deleted),
  updated_at = CURRENT_TIMESTAMP;

CREATE TABLE IF NOT EXISTS onec_cluster_app (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  cluster_uuid VARCHAR(64) NOT NULL DEFAULT '',
  app_name VARCHAR(64) NOT NULL DEFAULT '',
  app_code VARCHAR(64) NOT NULL DEFAULT '',
  app_type INT NOT NULL DEFAULT 1,
  is_default TINYINT(1) NOT NULL DEFAULT 0,
  app_url VARCHAR(500) NOT NULL DEFAULT '',
  port INT NOT NULL DEFAULT 0,
  protocol VARCHAR(16) NOT NULL DEFAULT 'http',
  auth_enabled TINYINT(1) NOT NULL DEFAULT 0,
  auth_type VARCHAR(32) NOT NULL DEFAULT 'none',
  username VARCHAR(128) NOT NULL DEFAULT '',
  password VARCHAR(500) NOT NULL DEFAULT '',
  token TEXT,
  access_key VARCHAR(128) NOT NULL DEFAULT '',
  access_secret VARCHAR(128) NOT NULL DEFAULT '',
  tls_enabled TINYINT(1) NOT NULL DEFAULT 0,
  ca_file TEXT,
  ca_key TEXT,
  ca_cert TEXT,
  client_cert TEXT,
  client_key TEXT,
  insecure_skip_verify TINYINT(1) NOT NULL DEFAULT 0,
  status TINYINT(1) NOT NULL DEFAULT 1,
  created_by VARCHAR(32) NOT NULL DEFAULT 'system',
  updated_by VARCHAR(32) NOT NULL DEFAULT 'system',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  is_deleted TINYINT(1) NOT NULL DEFAULT 0,
  UNIQUE KEY uk_cluster_app (cluster_uuid, app_code, app_type),
  KEY idx_cluster_uuid (cluster_uuid),
  KEY idx_app_type (app_type),
  KEY idx_status (status),
  KEY idx_is_deleted (is_deleted)
);

INSERT INTO onec_cluster_app
  (id, cluster_uuid, app_name, app_code, app_type, is_default, app_url, port, protocol, auth_enabled, auth_type, status, created_by, updated_by, is_deleted)
VALUES
  (1, '11111111-1111-1111-1111-111111111111', 'Prometheus-prod-hz', 'prometheus', 1, 0, 'prometheus.example.local', 9090, 'http', 0, 'none', 1, 'system', 'system', 0),
  (2, '11111111-1111-1111-1111-111111111111', 'Grafana-prod-hz', 'grafana', 1, 0, 'grafana.example.local', 3000, 'http', 0, 'none', 1, 'system', 'system', 0),
  (3, '22222222-2222-2222-2222-222222222222', 'Jaeger-staging-sh', 'jaeger', 3, 0, 'jaeger.example.local', 16686, 'http', 0, 'none', 0, 'system', 'system', 0)
ON DUPLICATE KEY UPDATE
  app_name = VALUES(app_name),
  app_url = VALUES(app_url),
  port = VALUES(port),
  protocol = VALUES(protocol),
  auth_enabled = VALUES(auth_enabled),
  auth_type = VALUES(auth_type),
  status = VALUES(status),
  updated_by = VALUES(updated_by),
  is_deleted = VALUES(is_deleted),
  updated_at = CURRENT_TIMESTAMP;
