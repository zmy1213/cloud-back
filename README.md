# cloud-back

`cloud-back` 是按 `kube-nova` 风格拆分的后端脚手架，便于和原项目 1:1 对照改写。

## 目录结构

```text
cloud-back/
├── application/           # 微服务应用层 (API/RPC 风格)
│   ├── portal-api/        # 认证与基础能力（已实现）
│   ├── manager-api/       # 集群配置管理（已实现）
│   ├── console-api/       # 监控查询接口（已实现）
│   └── workload-api/      # 占位
├── common/                # 公共配置/工具
├── pkg/                   # mysql/redis/minio/jwt 等封装
├── manifests/             # K8s 清单
├── dockerfile/            # 镜像构建文件
└── sql/                   # 数据库初始化脚本
```

## 本地依赖安装（MySQL / Redis / MinIO）

下面以 macOS + Homebrew 为例。

1. 安装并启动 MySQL

```bash
brew install mysql
brew services start mysql
mysql -uroot
```

进入 MySQL 后执行：

```sql
ALTER USER 'root'@'localhost' IDENTIFIED BY '515117';
CREATE DATABASE IF NOT EXISTS cloud_back DEFAULT CHARACTER SET utf8mb4;
```

2. 安装并启动 Redis

```bash
brew install redis
brew services start redis
redis-cli ping
```

返回 `PONG` 即正常。

3. 安装并启动 MinIO

```bash
brew install minio/stable/minio
mkdir -p ~/minio-data
minio server ~/minio-data --address ":9000" --console-address ":9001"
```

另开终端初始化桶：

```bash
brew install minio/stable/mc
mc alias set local http://127.0.0.1:9000 minioadmin minioadmin
mc mb local/cloud-back || true
```

## 本地安装 Prometheus（用于中间件连接测试）

下面以 macOS + 本地 Kubernetes（kind / Docker Desktop Kubernetes）为例，安装 `kube-prometheus-stack`。

1. 准备 Helm 仓库

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

2. 安装 Prometheus 栈（包含 Prometheus / Alertmanager / Grafana）

```bash
kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -
helm upgrade --install kps prometheus-community/kube-prometheus-stack \
  -n monitoring \
  --wait \
  --timeout 10m
```

3. 检查安装状态

```bash
helm status kps -n monitoring
kubectl -n monitoring get pods
```

4. 暴露 Prometheus 到本地 `9090` 端口

```bash
kubectl -n monitoring port-forward svc/kps-kube-prometheus-stack-prometheus 9090:9090
```

访问：

```text
http://127.0.0.1:9090/graph
```

5. （可选）采集 Mac 主机 CPU/内存指标

先安装并启动主机侧 `node_exporter`：

```bash
brew install node_exporter
brew services start node_exporter
curl http://127.0.0.1:9100/metrics | head
```

然后给 Prometheus 增加抓取任务：

```bash
cat > /tmp/kps-mac-host.yaml <<'EOF'
prometheus:
  prometheusSpec:
    additionalScrapeConfigs:
      - job_name: mac-host
        scrape_interval: 15s
        static_configs:
          - targets:
              - host.docker.internal:9100
EOF

helm upgrade --install kps prometheus-community/kube-prometheus-stack \
  -n monitoring \
  -f /tmp/kps-mac-host.yaml \
  --wait \
  --timeout 10m
```

验证：

```bash
kubectl -n monitoring exec -it prometheus-kps-kube-prometheus-stack-prometheus-0 -c prometheus -- \
  wget -qO- 'http://localhost:9090/api/v1/query?query=up{job="mac-host"}'
```

## 如何修改连接信息

本地运行时修改这个文件：

- [portal-api.yaml](/Users/zhuzhumingyang/githubProjects/kube-nova/cloud-back/application/portal-api/etc/portal-api.yaml)

关键配置项：

```yaml
Mysql:
  Enabled: true
  DataSource: root:515117@tcp(127.0.0.1:3306)/cloud_back?charset=utf8mb4&parseTime=True&loc=Local

Redis:
  Enabled: true
  Addr: 127.0.0.1:6379
  Password: ""
  DB: 0

Minio:
  Enabled: true
  Endpoint: 127.0.0.1:9000
  AccessKey: minioadmin
  SecretKey: minioadmin
  BucketName: cloud-back
  UseSSL: false
```

说明：

- MySQL 用户/密码/库名改 `DataSource`
- Redis 地址/密码/库改 `Addr/Password/DB`
- MinIO 改 `Endpoint/AccessKey/SecretKey/BucketName`
- 开关由 `Enabled` 控制，设为 `true` 才会纳入健康检查

## 启动服务

```bash
cd /Users/zhuzhumingyang/githubProjects/kube-nova/cloud-back
go mod tidy
make run-portal-api
make run-manager-api
make run-console-api
```

## 连接测试（推荐）

1. 健康检查（同时检查 MySQL/Redis/MinIO）

```bash
curl -sS -i http://127.0.0.1:8810/healthz
```
示至少有一个依赖未连通，查看 `deps.xxx.error` 即可定位原因

2. 登录接口检查

```bash
curl -X POST http://127.0.0.1:8810/portal/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"super_admin","password":"YWRtaW4xMjM="}'
```
## SQL 初始化脚本

- [init.sql](/Users/zhuzhumingyang/githubProjects/kube-nova/cloud-back/sql/init.sql)
