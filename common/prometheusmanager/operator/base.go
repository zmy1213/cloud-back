package operator

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	promtypes "github.com/yanshicheng/cloud-back/common/prometheusmanager/types"
)

type BaseOperator struct {
	endpoint string
	http     *http.Client
	config   promtypes.PrometheusConfig
}

func NewBaseOperator(config promtypes.PrometheusConfig) (*BaseOperator, error) {
	if strings.TrimSpace(config.Endpoint) == "" {
		return nil, errors.New("prometheus endpoint is empty")
	}

	timeout := config.TimeoutSeconds
	if timeout <= 0 {
		timeout = 15
	}

	tlsConfig := &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify}
	needTLS := strings.HasPrefix(strings.ToLower(strings.TrimSpace(config.Endpoint)), "https://") || normalizeAuthType(config.AuthType) == "certificate"
	if needTLS {
		rootCA := firstNonEmpty(config.CaCert, config.CaFile)
		if rootCA != "" {
			pool := x509.NewCertPool()
			if !pool.AppendCertsFromPEM([]byte(rootCA)) {
				return nil, errors.New("invalid prometheus CA certificate content")
			}
			tlsConfig.RootCAs = pool
		}

		if normalizeAuthType(config.AuthType) == "certificate" {
			cert, err := tls.X509KeyPair([]byte(config.ClientCert), []byte(config.ClientKey))
			if err != nil {
				return nil, fmt.Errorf("invalid prometheus client certificate pair: %w", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
	}

	return &BaseOperator{
		endpoint: strings.TrimSuffix(strings.TrimSpace(config.Endpoint), "/"),
		http: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
		config: config,
	}, nil
}

func normalizeAuthType(authType string) string {
	authType = strings.TrimSpace(authType)
	if authType == "" {
		return "none"
	}
	switch strings.ToLower(authType) {
	case "none":
		return "none"
	case "basic":
		return "basic"
	case "token", "bearer":
		return "token"
	case "apikey", "api_key", "api-key":
		return "apiKey"
	case "certificate", "cert":
		return "certificate"
	default:
		return authType
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func (b *BaseOperator) DoRequest(ctx context.Context, path string, params url.Values, out interface{}) error {
	endpoint := b.endpoint + path
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	if b.config.AuthEnabled {
		switch normalizeAuthType(b.config.AuthType) {
		case "basic":
			req.SetBasicAuth(b.config.Username, b.config.Password)
		case "token":
			token := strings.TrimSpace(b.config.Token)
			if token != "" {
				if strings.HasPrefix(strings.ToLower(token), "bearer ") {
					req.Header.Set("Authorization", token)
				} else {
					req.Header.Set("Authorization", "Bearer "+token)
				}
			}
		case "apiKey":
			if strings.TrimSpace(b.config.AccessKey) != "" {
				req.Header.Set("X-Access-Key", b.config.AccessKey)
			}
			if strings.TrimSpace(b.config.AccessSecret) != "" {
				req.Header.Set("X-Access-Secret", b.config.AccessSecret)
			}
		}
	}

	resp, err := b.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("prometheus status=%d body=%s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode prometheus response failed: %w", err)
	}
	return nil
}
