package auth

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func encodeFrontend(pwd string) string {
	// Keep compatibility with the frontend contract:
	// btoa(encodeURIComponent(password).replace(/%([0-9A-F]{2})/g, ...))
	input := []byte(pwd)
	escaped := make([]byte, 0, len(input)*3)
	for i := 0; i < len(input); i++ {
		ch := input[i]
		if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '_' || ch == '.' || ch == '~' {
			escaped = append(escaped, ch)
			continue
		}
		escaped = append(escaped, fmt.Sprintf("%%%02X", ch)...)
	}
	return base64.StdEncoding.EncodeToString(escaped)
}

func TestLoginSuccess(t *testing.T) {
	s := NewService(3600, 604800)
	resp, err := s.Login("super_admin", encodeFrontend("admin123"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Username != "super_admin" {
		t.Fatalf("expected super_admin, got %s", resp.Username)
	}
	if resp.Token.AccessToken == "" || resp.Token.RefreshToken == "" {
		t.Fatalf("expected token values")
	}
}

func TestLoginInvalid(t *testing.T) {
	s := NewService(3600, 604800)
	_, err := s.Login("super_admin", encodeFrontend("wrong"))
	if err == nil {
		t.Fatal("expected invalid password error")
	}
}

func TestDecodeFrontendPassword(t *testing.T) {
	raw := "Admin@123!"
	got, err := DecodeFrontendPassword(encodeFrontend(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != raw {
		t.Fatalf("expected %s, got %s", raw, got)
	}
}
