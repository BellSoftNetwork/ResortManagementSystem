package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
	"gitlab.bellsoft.net/rms/api-core/pkg/utils"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.UserID, req.Email, req.Name, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			response.Conflict(c, "이미 존재하는 사용자")
			return
		}
		response.InternalServerError(c, "사용자 등록 실패")
		return
	}

	userResponse := dto.UserResponse{
		ID:              user.ID,
		UserID:          user.UserID,
		Email:           user.Email,
		Name:            user.Name,
		Status:          user.Status.String(),
		Role:            user.Role.String(),
		ProfileImageURL: utils.GenerateGravatarURL(user.Email),
		CreatedAt:       dto.CustomTime{Time: user.CreatedAt},
		UpdatedAt:       dto.CustomTime{Time: user.UpdatedAt},
	}

	response.Created(c, userResponse)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	ipAddress := c.ClientIP()
	deviceInfo := h.extractDeviceInfo(&req)

	loginResp, err := h.authService.Login(c.Request.Context(), req.Username, req.Password, ipAddress, deviceInfo)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			response.Unauthorized(c, "잘못된 사용자명 또는 비밀번호")
		case errors.Is(err, services.ErrUserNotActive):
			response.Forbidden(c, "비활성화된 계정")
		case errors.Is(err, services.ErrTooManyAttempts):
			response.TooManyRequests(c, "로그인 시도 횟수 초과. 잠시 후 다시 시도하세요")
		default:
			response.InternalServerError(c, "로그인 실패")
		}
		return
	}

	resp := dto.LoginResponse{
		AccessToken:          loginResp.AccessToken,
		RefreshToken:         loginResp.RefreshToken,
		AccessTokenExpiresIn: loginResp.AccessTokenExpiresIn,
	}

	response.Success(c, resp)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	tokenResp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "유효하지 않은 리프레시 토큰입니다.")
		return
	}

	resp := dto.TokenResponse{
		AccessToken:          tokenResp.AccessToken,
		RefreshToken:         tokenResp.RefreshToken,
		AccessTokenExpiresIn: tokenResp.AccessTokenExpiresIn,
	}

	response.Success(c, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		response.Unauthorized(c, "인증되지 않음")
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID.(uint)); err != nil {
		response.InternalServerError(c, "로그아웃 실패")
		return
	}

	response.Success(c, gin.H{"message": "Successfully logged out"})
}

func (h *AuthHandler) extractDeviceInfo(req *dto.LoginRequest) *services.DeviceInfo {
	if req.Device == nil {
		return nil
	}

	return &services.DeviceInfo{
		OSInfo:            req.Device.OSInfo,
		LanguageInfo:      req.Device.LanguageInfo,
		UserAgent:         req.Device.UserAgent,
		DeviceFingerprint: req.Device.DeviceFingerprint,
	}
}
