package services

import (
	"context"
	"errors"
	"time"

	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"gitlab.bellsoft.net/rms/api-core/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("잘못된 인증 정보")
	ErrUserNotActive      = errors.New("비활성 상태의 사용자")
	ErrTooManyAttempts    = errors.New("로그인 시도 횟수 초과")
	ErrUserAlreadyExists  = errors.New("이미 존재하는 사용자")
)

type AuthService interface {
	Register(ctx context.Context, userID, email, name, password string) (*models.User, error)
	Login(ctx context.Context, username, password, ipAddress string, deviceInfo *DeviceInfo) (*LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	Logout(ctx context.Context, userID uint) error
	IsDeviceChanged(ctx context.Context, username string, deviceInfo *DeviceInfo) (bool, error)
}

type DeviceInfo struct {
	OSInfo            string
	LanguageInfo      string
	UserAgent         string
	DeviceFingerprint string
}

type LoginResponse struct {
	User                 *models.User
	AccessToken          string
	RefreshToken         string
	AccessTokenExpiresIn int64
}

type TokenResponse struct {
	AccessToken          string
	RefreshToken         string
	AccessTokenExpiresIn int64
}

type authService struct {
	userRepo         repositories.UserRepository
	loginAttemptRepo repositories.LoginAttemptRepository
	jwtService       *auth.JWTService
}

func NewAuthService(userRepo repositories.UserRepository, loginAttemptRepo repositories.LoginAttemptRepository,
	jwtService *auth.JWTService) AuthService {
	return &authService{
		userRepo:         userRepo,
		loginAttemptRepo: loginAttemptRepo,
		jwtService:       jwtService,
	}
}

func (s *authService) Register(ctx context.Context, userID, email, name, password string) (*models.User, error) {
	// Check if userID already exists
	exists, err := s.userRepo.ExistsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// Check if email already exists (if provided)
	if email != "" {
		exists, err := s.userRepo.ExistsByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrUserAlreadyExists
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var emailPtr *string
	if email != "" {
		emailPtr = &email
	}

	user := &models.User{
		UserID:   userID,
		Email:    emailPtr,
		Name:     name,
		Password: "{bcrypt}" + string(hashedPassword),
		Status:   models.UserStatusActive,
		Role:     models.UserRoleNormal,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, username, password, ipAddress string, deviceInfo *DeviceInfo) (*LoginResponse, error) {
	if err := s.checkLoginAttempts(ctx, username, ipAddress); err != nil {
		return nil, err
	}

	var user *models.User
	var err error

	user, err = s.userRepo.FindByUserID(ctx, username)
	if err != nil {
		user, err = s.userRepo.FindByEmail(ctx, username)
		if err != nil {
			s.recordLoginAttempt(ctx, username, ipAddress, false, deviceInfo)
			return nil, ErrInvalidCredentials
		}
	}

	// Handle Spring Security password format
	passwordHash := user.Password
	if len(passwordHash) > 8 && passwordHash[:8] == "{bcrypt}" {
		passwordHash = passwordHash[8:]
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		s.recordLoginAttempt(ctx, username, ipAddress, false, deviceInfo)
		return nil, ErrInvalidCredentials
	}

	if !user.IsActive() {
		s.recordLoginAttempt(ctx, username, ipAddress, false, deviceInfo)
		return nil, ErrUserNotActive
	}

	s.recordLoginAttempt(ctx, username, ipAddress, true, deviceInfo)

	var deviceFingerprint string
	if deviceInfo != nil {
		deviceFingerprint = deviceInfo.DeviceFingerprint
	}

	accessToken, refreshToken, err := s.jwtService.GenerateTokenPairWithDevice(user.ID, user.UserID, user.Role.String(), deviceFingerprint)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:                 user,
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		AccessTokenExpiresIn: s.jwtService.GetAccessTokenExpiryMillis(),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID, err := s.jwtService.GetUserIDFromRefreshClaims(claims)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive() {
		return nil, ErrUserNotActive
	}

	newAccessToken, newRefreshToken, err := s.jwtService.GenerateTokenPairWithDevice(user.ID, user.UserID, user.Role.String(), claims.DeviceFingerprint)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:          newAccessToken,
		RefreshToken:         newRefreshToken,
		AccessTokenExpiresIn: s.jwtService.GetAccessTokenExpiryMillis(),
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uint) error {
	return nil
}

func (s *authService) checkLoginAttempts(ctx context.Context, username, ipAddress string) error {
	since := time.Now().Add(-15 * time.Minute)
	count, err := s.loginAttemptRepo.CountRecentFailedAttempts(ctx, username, ipAddress, since)
	if err != nil {
		return err
	}

	if count >= 5 {
		return ErrTooManyAttempts
	}

	return nil
}

func (s *authService) recordLoginAttempt(ctx context.Context, username, ipAddress string, successful bool, deviceInfo *DeviceInfo) {
	attempt := &models.LoginAttempt{
		Username:   username,
		IPAddress:  ipAddress,
		Successful: successful,
		AttemptAt:  time.Now(),
	}

	if deviceInfo != nil {
		if deviceInfo.OSInfo != "" {
			attempt.OSInfo = &deviceInfo.OSInfo
		}
		if deviceInfo.LanguageInfo != "" {
			attempt.LanguageInfo = &deviceInfo.LanguageInfo
		}
		if deviceInfo.UserAgent != "" {
			attempt.UserAgent = &deviceInfo.UserAgent
		}
		if deviceInfo.DeviceFingerprint != "" {
			attempt.DeviceFingerprint = &deviceInfo.DeviceFingerprint
		}
	}

	s.loginAttemptRepo.Create(ctx, attempt)
}

func (s *authService) IsDeviceChanged(ctx context.Context, username string, deviceInfo *DeviceInfo) (bool, error) {
	lastSuccessfulLogin, err := s.loginAttemptRepo.GetLastSuccessfulAttempt(ctx, username)
	if err != nil {
		return false, nil
	}

	if lastSuccessfulLogin.DeviceFingerprint == nil {
		return false, nil
	}

	if deviceInfo == nil || deviceInfo.DeviceFingerprint == "" {
		return true, nil
	}

	return *lastSuccessfulLogin.DeviceFingerprint != deviceInfo.DeviceFingerprint, nil
}
