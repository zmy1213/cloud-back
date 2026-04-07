package types

type TokenResponse struct {
	AccessToken      string `json:"accessToken"`
	AccessExpiresIn  int64  `json:"accessExpiresIn"`
	RefreshToken     string `json:"refreshToken"`
	RefreshExpiresIn int64  `json:"refreshExpiresIn"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID   uint64        `json:"userId"`
	Username string        `json:"username"`
	NickName string        `json:"nickName"`
	UUID     string        `json:"uuid"`
	Roles    []string      `json:"roles"`
	Token    TokenResponse `json:"token"`
}
