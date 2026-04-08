package app

import (
	"errors"
	"strings"

	"github.com/yanshicheng/cloud-back/application/manager-api/internal/types"
)

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

func sanitizeAddClusterAppRequest(req *types.AddClusterAppRequest) {
	req.ClusterUuid = strings.TrimSpace(req.ClusterUuid)
	req.AppName = strings.TrimSpace(req.AppName)
	req.AppCode = strings.TrimSpace(req.AppCode)
	req.AppUrl = strings.TrimSpace(req.AppUrl)
	req.Protocol = strings.ToLower(strings.TrimSpace(req.Protocol))
	if req.Protocol == "" {
		req.Protocol = "http"
	}

	req.UpdatedBy = strings.TrimSpace(req.UpdatedBy)
	if req.UpdatedBy == "" {
		req.UpdatedBy = "system"
	}

	req.AuthType = normalizeAuthType(req.AuthType)
	if req.AuthEnabled != 1 {
		req.AuthEnabled = 0
		req.AuthType = "none"
		req.Username = ""
		req.Password = ""
		req.Token = ""
		req.AccessKey = ""
		req.AccessSecret = ""
		req.TlsEnabled = 0
		req.CaFile = ""
		req.CaKey = ""
		req.CaCert = ""
		req.ClientCert = ""
		req.ClientKey = ""
		req.InsecureSkipVerify = 0
		return
	}

	if req.AuthType != "apiKey" {
		req.TlsEnabled = 0
		req.CaFile = ""
		req.CaKey = ""
	}
	if req.AuthType != "certificate" {
		req.CaCert = ""
		req.ClientCert = ""
		req.ClientKey = ""
	}
}

func validateAddClusterAppRequest(req *types.AddClusterAppRequest) error {
	if req.ClusterUuid == "" {
		return errors.New("clusterUuid is required")
	}
	if req.AppName == "" {
		return errors.New("appName is required")
	}
	if req.AppCode == "" {
		return errors.New("appCode is required")
	}
	if req.AppType <= 0 {
		return errors.New("appType is invalid")
	}
	if req.AppUrl == "" {
		return errors.New("appUrl is required")
	}
	if strings.Contains(req.AppUrl, "://") {
		return errors.New("appUrl should not include protocol prefix")
	}
	if req.Port < 1 || req.Port > 65535 {
		return errors.New("port must be within 1-65535")
	}

	switch req.Protocol {
	case "http", "https", "grpc":
	default:
		return errors.New("protocol must be one of http/https/grpc")
	}

	if req.AuthEnabled == 0 {
		return nil
	}

	switch req.AuthType {
	case "basic":
		if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
			return errors.New("basic auth requires username and password")
		}
	case "token":
		if strings.TrimSpace(req.Token) == "" {
			return errors.New("token auth requires token")
		}
	case "apiKey":
		if strings.TrimSpace(req.AccessKey) == "" || strings.TrimSpace(req.AccessSecret) == "" {
			return errors.New("apiKey auth requires accessKey and accessSecret")
		}
		if req.TlsEnabled == 1 {
			if strings.TrimSpace(req.CaFile) == "" || strings.TrimSpace(req.CaKey) == "" {
				return errors.New("tls enabled for apiKey auth requires caFile and caKey")
			}
		}
	case "certificate":
		if strings.TrimSpace(req.CaCert) == "" || strings.TrimSpace(req.ClientCert) == "" || strings.TrimSpace(req.ClientKey) == "" {
			return errors.New("certificate auth requires caCert, clientCert and clientKey")
		}
	case "none":
	default:
		return errors.New("unsupported authType")
	}

	if req.InsecureSkipVerify != 0 && req.InsecureSkipVerify != 1 {
		return errors.New("insecureSkipVerify must be 0 or 1")
	}

	return nil
}
