package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	appcfg "github.com/yanshicheng/cloud-back/common/config"
)

func TestLoginHandler(t *testing.T) {
	h := New(appcfg.AppConfig{Name: "portal-api"})
	mux := http.NewServeMux()
	h.Register(mux)

	body, _ := json.Marshal(map[string]string{
		"username": "super_admin",
		"password": base64.StdEncoding.EncodeToString([]byte("admin123")),
	})

	req := httptest.NewRequest(http.MethodPost, "/portal/v1/auth/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestHealthzHandler(t *testing.T) {
	h := New(appcfg.AppConfig{Name: "portal-api"})
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestDashboardOverviewHandler(t *testing.T) {
	h := New(appcfg.AppConfig{Name: "portal-api"})
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/portal/v1/dashboard/overview?username=super_admin", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if rr.Body.Len() == 0 {
		t.Fatal("expected non-empty response body")
	}
}

func TestDashboardOverviewWithClusterScope(t *testing.T) {
	h := New(appcfg.AppConfig{Name: "portal-api"})
	mux := http.NewServeMux()
	h.Register(mux)

	clusterUUID := url.QueryEscape("11111111-1111-1111-1111-111111111111")
	req := httptest.NewRequest(http.MethodGet, "/portal/v1/dashboard/overview?clusterUuid="+clusterUUID, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
}
