# cloud-back

`cloud-back` 是按 `kube-nova` 风格拆分的后端脚手架，便于和原项目 1:1 对照改写。

## 目录结构

```text
cloud-back/
├── application/           # 微服务应用层 (API/RPC 风格)
│   ├── portal-api/        # 认证与基础能力（已实现）
│   ├── manager-api/       # 占位
│   ├── console-api/       # 占位
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
