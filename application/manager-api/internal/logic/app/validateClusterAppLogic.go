package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	apprepo "github.com/yanshicheng/cloud-back/application/manager-api/internal/repository/app"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/svc"
	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

type ValidateClusterAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidateClusterAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateClusterAppLogic {
	return &ValidateClusterAppLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ValidateClusterAppLogic) ValidateClusterApp(req *types.ClusterAppValidateRequest) (string, error) {
	if req.ID == 0 {
		return "", errors.New("id is required")
	}

	item, ok, err := l.svcCtx.App.GetByID(req.ID)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("app not found")
	}

	validateErr := validateClusterAppConnection(item)
	if validateErr != nil {
		_ = l.svcCtx.App.UpdateStatus(req.ID, 0, "system")
		return "", fmt.Errorf("测试连接失败: %w", validateErr)
	}

	if err := l.svcCtx.App.UpdateStatus(req.ID, 1, "system"); err != nil {
		return "", err
	}
	return "测试连接成功", nil
}

func validateClusterAppConnection(app apprepo.ClusterApp) error {
	protocol := strings.ToLower(strings.TrimSpace(app.Protocol))
	switch protocol {
	case "http", "https":
		return validateHTTPConnection(app)
	case "grpc":
		return validateGRPCConnection(app)
	default:
		return fmt.Errorf("unsupported protocol: %s", app.Protocol)
	}
}

func validateHTTPConnection(app apprepo.ClusterApp) error {
	if strings.TrimSpace(app.AppUrl) == "" || app.Port < 1 || app.Port > 65535 {
		return errors.New("invalid app address")
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: app.InsecureSkipVerify == 1}

	if app.Protocol == "https" || app.TlsEnabled == 1 || normalizeAuthType(app.AuthType) == "certificate" {
		if rootCA := firstNonEmpty(app.CaCert, app.CaFile); rootCA != "" {
			pool := x509.NewCertPool()
			if !pool.AppendCertsFromPEM([]byte(rootCA)) {
				return errors.New("invalid CA certificate content")
			}
			tlsConfig.RootCAs = pool
		}

		if normalizeAuthType(app.AuthType) == "certificate" {
			cert, err := tls.X509KeyPair([]byte(app.ClientCert), []byte(app.ClientKey))
			if err != nil {
				return fmt.Errorf("invalid client certificate pair: %w", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	target := fmt.Sprintf("%s://%s:%d", strings.ToLower(app.Protocol), strings.TrimSpace(app.AppUrl), app.Port)
	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return err
	}

	switch normalizeAuthType(app.AuthType) {
	case "basic":
		req.SetBasicAuth(app.Username, app.Password)
	case "token":
		token := strings.TrimSpace(app.Token)
		if token != "" {
			if strings.HasPrefix(strings.ToLower(token), "bearer ") {
				req.Header.Set("Authorization", token)
			} else {
				req.Header.Set("Authorization", "Bearer "+token)
			}
		}
	case "apiKey":
		if strings.TrimSpace(app.AccessKey) != "" {
			req.Header.Set("X-Access-Key", app.AccessKey)
		}
		if strings.TrimSpace(app.AccessSecret) != "" {
			req.Header.Set("X-Access-Secret", app.AccessSecret)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("remote service returned status %d", resp.StatusCode)
	}
	return nil
}

func validateGRPCConnection(app apprepo.ClusterApp) error {
	if strings.TrimSpace(app.AppUrl) == "" || app.Port < 1 || app.Port > 65535 {
		return errors.New("invalid app address")
	}
	endpoint := net.JoinHostPort(strings.TrimSpace(app.AppUrl), strconv.Itoa(int(app.Port)))
	conn, err := net.DialTimeout("tcp", endpoint, 5*time.Second)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
