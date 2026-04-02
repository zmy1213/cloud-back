package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func frontendEncode(password string) string {
	// Mirrors kube-nova-web behavior: btoa(encodeURIComponent(password))
	return base64.StdEncoding.EncodeToString([]byte(password))
}

func TestLoginSuccess(t *testing.T) {
	s := newAppServer(apiConfig{Name: "cloud-api"})

	body, _ := json.Marshal(LoginRequest{
		Username: "super_admin",
		Password: frontendEncode("admin123"),
	})

	req := httptest.NewRequest(http.MethodPost, "/portal/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body=%s", rr.Code, rr.Body.String())
	}

	var resp LoginResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Username != "super_admin" {
		t.Fatalf("expected username super_admin, got %s", resp.Username)
	}
	if resp.Token.AccessToken == "" || resp.Token.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	s := newAppServer(apiConfig{Name: "cloud-api"})

	body, _ := json.Marshal(LoginRequest{
		Username: "super_admin",
		Password: frontendEncode("wrong"),
	})

	req := httptest.NewRequest(http.MethodPost, "/portal/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestLoginMethodNotAllowed(t *testing.T) {
	s := newAppServer(apiConfig{Name: "cloud-api"})

	req := httptest.NewRequest(http.MethodGet, "/portal/v1/auth/login", nil)
	rr := httptest.NewRecorder()

	s.routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}

func TestLoginBadEncodedPassword(t *testing.T) {
	s := newAppServer(apiConfig{Name: "cloud-api"})

	body, _ := json.Marshal(LoginRequest{
		Username: "super_admin",
		Password: "not-base64",
	})

	req := httptest.NewRequest(http.MethodPost, "/portal/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	s.routes().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}
}

func TestHealthzAndPing(t *testing.T) {
	s := newAppServer(apiConfig{Name: "cloud-api"})

	for _, tc := range []struct {
		path   string
		status int
	}{
		{path: "/healthz", status: http.StatusOK},
		{path: "/cloud/v1/ping", status: http.StatusOK},
	} {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		rr := httptest.NewRecorder()
		s.routes().ServeHTTP(rr, req)
		if rr.Code != tc.status {
			t.Fatalf("%s expected %d got %d", tc.path, tc.status, rr.Code)
		}
	}
}
