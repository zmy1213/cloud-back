package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yanshicheng/cloud-back/pkg/config"
)

// apiConfig maps application/cloud-api/etc/cloud-api.yaml.
type apiConfig struct {
	Name    string `yaml:"Name"`
	Host    string `yaml:"Host"`
	Port    int    `yaml:"Port"`
	Timeout int    `yaml:"Timeout"`
}

// LoginRequest keeps the same field naming as kube-nova-web login payload.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"` // base64 encoded by frontend
}

// TokenResponse mirrors kube-nova auth token response shape.
type TokenResponse struct {
	AccessToken      string `json:"accessToken"`
	AccessExpiresIn  int64  `json:"accessExpiresIn"`
	RefreshToken     string `json:"refreshToken"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
}

// LoginResponse mirrors kube-nova auth login response shape.
type LoginResponse struct {
	UserID   uint64        `json:"userId"`
	Username string        `json:"username"`
	NickName string        `json:"nickName"`
	UUID     string        `json:"uuid"`
	Roles    []string      `json:"roles"`
	Token    TokenResponse `json:"token"`
}

type appServer struct {
	cfg         apiConfig
	validUsers  map[string]string
	tokenExpiry int64
}

func newAppServer(cfg apiConfig) *appServer {
	// Empty scaffold user. Replace with DB/RPC auth in real implementation.
	return &appServer{
		cfg: cfg,
		validUsers: map[string]string{
			"super_admin": "admin123",
		},
		tokenExpiry: 3600,
	}
}

func (s *appServer) routes() http.Handler {
	mux := http.NewServeMux()

	// Health endpoint for liveness/readiness probes.
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Minimal API endpoint for quick connectivity verification.
	mux.HandleFunc("/cloud/v1/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"service": s.cfg.Name,
			"time":    time.Now().Format(time.RFC3339),
			"status":  "pong",
		})
	})

	// Keep login endpoint contract consistent with kube-nova-web.
	mux.HandleFunc("/portal/v1/auth/login", s.handleLogin)

	return mux
}

func (s *appServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	decodedPassword, err := decodeFrontendPassword(req.Password)
	if err != nil {
		http.Error(w, "password decode failed", http.StatusBadRequest)
		return
	}

	expectedPassword, ok := s.validUsers[req.Username]
	if !ok || expectedPassword != decodedPassword {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	accessToken, err := randomToken(32)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}
	refreshToken, err := randomToken(40)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	resp := LoginResponse{
		UserID:   1,
		Username: req.Username,
		NickName: "Cloud Admin",
		UUID:     "cloud-user-0001",
		Roles:    []string{"super_admin"},
		Token: TokenResponse{
			AccessToken:      accessToken,
			AccessExpiresIn:  s.tokenExpiry,
			RefreshToken:     refreshToken,
			RefreshExpiresIn: 7 * 24 * 3600,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func decodeFrontendPassword(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	password, err := decodeURIComponentBytes(decoded)
	if err != nil {
		return "", err
	}
	return password, nil
}

// decodeURIComponentBytes mirrors: decodeURIComponent(str) in browser.
func decodeURIComponentBytes(b []byte) (string, error) {
	input := string(b)
	if !strings.Contains(input, "%") {
		return input, nil
	}

	var out []byte
	for i := 0; i < len(input); i++ {
		if input[i] != '%' {
			out = append(out, input[i])
			continue
		}

		if i+2 >= len(input) {
			return "", errors.New("invalid percent encoding")
		}
		hexPair := input[i+1 : i+3]
		decodedByte, err := decodeHexByte(hexPair)
		if err != nil {
			return "", err
		}
		out = append(out, decodedByte)
		i += 2
	}

	return string(out), nil
}

func decodeHexByte(s string) (byte, error) {
	const hex = "0123456789ABCDEF"
	if len(s) != 2 {
		return 0, errors.New("hex pair length must be 2")
	}
	high := strings.IndexByte(hex, byte(strings.ToUpper(s)[0]))
	low := strings.IndexByte(hex, byte(strings.ToUpper(s)[1]))
	if high < 0 || low < 0 {
		return 0, errors.New("invalid hex pair")
	}
	return byte((high << 4) | low), nil
}

func randomToken(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func main() {
	// -f supports custom config path, same style as kube-nova services.
	cfgPath := flag.String("f", "./application/cloud-api/etc/cloud-api.yaml", "config file path")
	flag.Parse()

	var cfg apiConfig
	config.MustLoad(*cfgPath, &cfg)

	server := newAppServer(cfg)
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf("Starting %s at %s...\n", cfg.Name, addr)

	// Use standard net/http server for a lightweight scaffold.
	if err := http.ListenAndServe(addr, server.routes()); err != nil {
		panic(err)
	}
}
