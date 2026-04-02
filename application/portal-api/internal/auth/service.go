package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"
)

var (
	ErrMissingCredentials = errors.New("username and password are required")
	ErrPasswordDecode     = errors.New("password decode failed")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type Service struct {
	users            map[string]string
	accessExpiresIn  int64
	refreshExpiresIn int64
}

type TokenResponse struct {
	AccessToken      string `json:"accessToken"`
	AccessExpiresIn  int64  `json:"accessExpiresIn"`
	RefreshToken     string `json:"refreshToken"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
}

type LoginResponse struct {
	UserID   uint64        `json:"userId"`
	Username string        `json:"username"`
	NickName string        `json:"nickName"`
	UUID     string        `json:"uuid"`
	Roles    []string      `json:"roles"`
	Token    TokenResponse `json:"token"`
}

func NewService(accessExpiresIn, refreshExpiresIn int64) *Service {
	if accessExpiresIn <= 0 {
		accessExpiresIn = 3600
	}
	if refreshExpiresIn <= 0 {
		refreshExpiresIn = 7 * 24 * 3600
	}
	return &Service{
		users: map[string]string{
			"super_admin": "admin123",
		},
		accessExpiresIn:  accessExpiresIn,
		refreshExpiresIn: refreshExpiresIn,
	}
}

func (s *Service) Login(username, encodedPassword string) (*LoginResponse, error) {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(encodedPassword) == "" {
		return nil, ErrMissingCredentials
	}

	decodedPassword, err := DecodeFrontendPassword(encodedPassword)
	if err != nil {
		return nil, ErrPasswordDecode
	}

	expected, ok := s.users[username]
	if !ok || expected != decodedPassword {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := randomToken(32)
	if err != nil {
		return nil, err
	}
	refreshToken, err := randomToken(40)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		UserID:   1,
		Username: username,
		NickName: "Cloud Admin",
		UUID:     "cloud-user-0001",
		Roles:    []string{"super_admin"},
		Token: TokenResponse{
			AccessToken:      accessToken,
			AccessExpiresIn:  s.accessExpiresIn,
			RefreshToken:     refreshToken,
			RefreshExpiresIn: s.refreshExpiresIn,
		},
	}, nil
}

func DecodeFrontendPassword(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return decodeURIComponentBytes(decoded)
}

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
		v, err := decodeHexByte(input[i+1 : i+3])
		if err != nil {
			return "", err
		}
		out = append(out, v)
		i += 2
	}
	return string(out), nil
}

func decodeHexByte(s string) (byte, error) {
	const hex = "0123456789ABCDEF"
	if len(s) != 2 {
		return 0, errors.New("hex pair length must be 2")
	}
	u := strings.ToUpper(s)
	h := strings.IndexByte(hex, u[0])
	l := strings.IndexByte(hex, u[1])
	if h < 0 || l < 0 {
		return 0, errors.New("invalid hex pair")
	}
	return byte((h << 4) | l), nil
}

func randomToken(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
