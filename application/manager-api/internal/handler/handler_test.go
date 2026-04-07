package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	appcfg "github.com/yanshicheng/cloud-back/common/config"
)

func TestHealthzHandler(t *testing.T) {
	h := New(appcfg.AppConfig{Name: "manager-api"})
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestSearchClusterHandler(t *testing.T) {
	h := New(appcfg.AppConfig{Name: "manager-api"})
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/manager/v1/cluster?environment=prod", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}

	var body struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if body.Total < 1 || len(body.Items) < 1 {
		t.Fatalf("expected at least one cluster, got total=%d", body.Total)
	}
}

func TestGetClusterDetailHandler(t *testing.T) {
	h := New(appcfg.AppConfig{Name: "manager-api"})
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/manager/v1/cluster/1", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestGetNodeListHandler(t *testing.T) {
	h := New(appcfg.AppConfig{Name: "manager-api"})
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(
		http.MethodGet,
		"/manager/v1/node?clusterUuid=11111111-1111-1111-1111-111111111111&page=1&pageSize=10",
		nil,
	)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}

	var body struct {
		Items []map[string]any `json:"items"`
		Total uint64           `json:"total"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	if body.Total < 1 || len(body.Items) < 1 {
		t.Fatalf("expected at least one node, got total=%d", body.Total)
	}
}

func TestGetNodeDetailHandler(t *testing.T) {
	h := New(appcfg.AppConfig{Name: "manager-api"})
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/manager/v1/node/1", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
}
