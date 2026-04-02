# cloud-back

`cloud-back` is a kube-nova-style backend scaffold so files can map 1:1 when you rewrite features.

## Project Layout

```text
cloud-back/
├── application/           # microservice apps (API & RPC style)
│   ├── portal-api/        # auth and base capability (implemented)
│   ├── manager-api/       # cluster/project management (placeholder)
│   ├── console-api/       # console/monitoring (placeholder)
│   └── workload-api/      # workload management (placeholder)
├── common/                # shared configuration/util helpers
├── pkg/                   # wrappers: config/mysql/redis/minio/jwt
├── manifests/             # Kubernetes manifests
│   ├── basic/             # namespace + mysql/redis/minio
│   └── base/              # portal-api deployment/service/config
├── dockerfile/            # Docker build files
└── sql/                   # database initialization scripts
```

## Implemented Service

- `application/portal-api`
- login endpoint: `POST /portal/v1/auth/login`
- health endpoint: `GET /healthz`
- auth password contract matches kube-nova-web style (frontend Base64 encoding with `encodeURIComponent` flow)

## Dependency Connectivity

Connectivity config is in:

- [application/portal-api/etc/portal-api.yaml](/Users/zhuzhumingyang/githubProjects/kube-nova/cloud-back/application/portal-api/etc/portal-api.yaml)
- [manifests/base/portal-api-configmap.yaml](/Users/zhuzhumingyang/githubProjects/kube-nova/cloud-back/manifests/base/portal-api-configmap.yaml)

Switches:

- `Mysql.Enabled`
- `Redis.Enabled`
- `Minio.Enabled`

When enabled:

- MySQL uses `Mysql.DataSource`
- Redis uses `Redis.Addr/Password/DB`
- MinIO uses `Minio.Endpoint/AccessKey/SecretKey/BucketName/UseSSL`

## Local Run

```bash
cd cloud-back
go mod tidy
make run-portal-api
```

health check:

```bash
curl http://127.0.0.1:8810/healthz
```

login test:

```bash
curl -X POST http://127.0.0.1:8810/portal/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"super_admin","password":"YWRtaW4xMjM="}'
```

## Kubernetes Run

```bash
cd cloud-back
kubectl apply -k manifests/basic
kubectl apply -k manifests/base
kubectl get pods -n cloud-back
```

## SQL Bootstrap

Initialize DB script:

- [sql/init.sql](/Users/zhuzhumingyang/githubProjects/kube-nova/cloud-back/sql/init.sql)
