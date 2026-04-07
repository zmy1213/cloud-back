package healthz

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yanshicheng/cloud-back/application/portal-api/internal/svc"
	minioPkg "github.com/yanshicheng/cloud-back/pkg/minio"
	mysqlPkg "github.com/yanshicheng/cloud-back/pkg/mysql"
	redisPkg "github.com/yanshicheng/cloud-back/pkg/redis"
)

func HealthzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		type depStatus struct {
			Enabled bool   `json:"enabled"`
			OK      bool   `json:"ok"`
			Error   string `json:"error,omitempty"`
		}

		timeout := 2 * time.Second
		deps := map[string]depStatus{
			"mysql": {Enabled: svcCtx.Config.Mysql.Enabled, OK: true},
			"redis": {Enabled: svcCtx.Config.Redis.Enabled, OK: true},
			"minio": {Enabled: svcCtx.Config.Minio.Enabled, OK: true},
		}

		if svcCtx.Config.Mysql.Enabled {
			if err := mysqlPkg.Ping(svcCtx.Config.Mysql.DataSource, timeout); err != nil {
				v := deps["mysql"]
				v.OK = false
				v.Error = err.Error()
				deps["mysql"] = v
			}
		}

		if svcCtx.Config.Redis.Enabled {
			if err := redisPkg.Ping(
				svcCtx.Config.Redis.Addr,
				svcCtx.Config.Redis.Password,
				svcCtx.Config.Redis.DB,
				timeout,
			); err != nil {
				v := deps["redis"]
				v.OK = false
				v.Error = err.Error()
				deps["redis"] = v
			}
		}

		if svcCtx.Config.Minio.Enabled {
			if err := minioPkg.Check(
				svcCtx.Config.Minio.Endpoint,
				svcCtx.Config.Minio.AccessKey,
				svcCtx.Config.Minio.SecretKey,
				svcCtx.Config.Minio.BucketName,
				svcCtx.Config.Minio.UseSSL,
				timeout,
			); err != nil {
				v := deps["minio"]
				v.OK = false
				v.Error = err.Error()
				deps["minio"] = v
			}
		}

		status := http.StatusOK
		for _, d := range deps {
			if d.Enabled && !d.OK {
				status = http.StatusServiceUnavailable
				break
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service": svcCtx.Config.Name,
			"status":  map[bool]string{true: "ok", false: "degraded"}[status == http.StatusOK],
			"deps":    deps,
		})
	}
}
