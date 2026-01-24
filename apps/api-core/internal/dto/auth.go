package dto

type RegisterRequest struct {
	UserID   string `json:"userId" binding:"required,min=3,max=30"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Name     string `json:"name" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Username string      `json:"username" binding:"required"`
	Password string      `json:"password" binding:"required"`
	Device   *DeviceInfo `json:"device,omitempty"`
}

type DeviceInfo struct {
	OSInfo            string `json:"osInfo,omitempty"`
	LanguageInfo      string `json:"languageInfo,omitempty"`
	UserAgent         string `json:"userAgent,omitempty"`
	DeviceFingerprint string `json:"deviceFingerprint,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LoginResponse struct {
	AccessToken          string `json:"accessToken"`
	RefreshToken         string `json:"refreshToken"`
	AccessTokenExpiresIn int64  `json:"accessTokenExpiresIn"`
}

type TokenResponse struct {
	AccessToken          string `json:"accessToken"`
	RefreshToken         string `json:"refreshToken"`
	AccessTokenExpiresIn int64  `json:"accessTokenExpiresIn"`
}
