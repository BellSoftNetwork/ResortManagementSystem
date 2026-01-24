package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/auth"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUserID(ctx context.Context, userID string) (*models.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindAll(ctx context.Context, offset, limit int) ([]models.User, int64, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) HasAnyUsers(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

// MockLoginAttemptRepository is a mock implementation of LoginAttemptRepository
type MockLoginAttemptRepository struct {
	mock.Mock
}

func (m *MockLoginAttemptRepository) Create(ctx context.Context, attempt *models.LoginAttempt) error {
	args := m.Called(ctx, attempt)
	return args.Error(0)
}

func (m *MockLoginAttemptRepository) CountRecentFailedAttempts(ctx context.Context, username, ipAddress string, since time.Time) (int64, error) {
	args := m.Called(ctx, username, ipAddress, since)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLoginAttemptRepository) GetLastSuccessfulAttempt(ctx context.Context, username string) (*models.LoginAttempt, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoginAttempt), args.Error(1)
}

func (m *MockLoginAttemptRepository) CleanOldAttempts(ctx context.Context, before time.Time) error {
	args := m.Called(ctx, before)
	return args.Error(0)
}

func TestAuthService_LoginAndTokenUsage(t *testing.T) {
	// Given: 테스트 환경 설정
	ctx := context.Background()

	jwtService := auth.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)
	mockUserRepo := new(MockUserRepository)
	mockLoginAttemptRepo := new(MockLoginAttemptRepository)

	authService := services.NewAuthService(mockUserRepo, mockLoginAttemptRepo, jwtService)

	// 테스트 사용자 설정
	testUser := &models.User{
		UserID:   "testadmin",
		Email:    "testadmin@example.com",
		Name:     "Test Admin",
		Password: "{bcrypt}$2a$10$yS/Y3Y0OcBZ9VFaNeTmpEuI6Vk1jbl5dke9prZNYZOduhmy2xu7T2", // password: password123
		Status:   models.UserStatusActive,
		Role:     models.UserRoleSuperAdmin,
	}
	testUser.ID = 1

	t.Run("로그인 성공 후 반환된 토큰이 즉시 사용 가능하다", func(t *testing.T) {
		// Given: Repository 모의 설정
		mockLoginAttemptRepo.On("CountRecentFailedAttempts", ctx, "testadmin", "127.0.0.1", mock.Anything).Return(int64(0), nil)
		mockUserRepo.On("FindByUserID", ctx, "testadmin").Return(testUser, nil)
		mockLoginAttemptRepo.On("Create", ctx, mock.Anything).Return(nil)

		// When: 로그인 수행
		loginResp, err := authService.Login(ctx, "testadmin", "password123", "127.0.0.1", nil)

		// Then: 로그인 성공
		assert.NoError(t, err)
		assert.NotNil(t, loginResp)
		assert.NotEmpty(t, loginResp.AccessToken)
		assert.NotEmpty(t, loginResp.RefreshToken)

		// And: 생성된 액세스 토큰이 즉시 검증 가능하다
		claims, err := jwtService.ValidateToken(loginResp.AccessToken)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "testadmin", claims.Username)
		assert.Equal(t, "SUPER_ADMIN", claims.Authorities)
		assert.Equal(t, "1", claims.Subject)

		// And: 리프레시 토큰이 유효하게 검증된다
		refreshClaims, err := jwtService.ValidateRefreshToken(loginResp.RefreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, refreshClaims)
		assert.Equal(t, "1", refreshClaims.Subject)
	})

	t.Run("리프레시 토큰으로 새 토큰을 발급받을 수 있다", func(t *testing.T) {
		// Given: 먼저 로그인하여 토큰 획득
		mockLoginAttemptRepo.On("CountRecentFailedAttempts", ctx, "testadmin", "127.0.0.1", mock.Anything).Return(int64(0), nil).Once()
		mockUserRepo.On("FindByUserID", ctx, "testadmin").Return(testUser, nil).Once()
		mockLoginAttemptRepo.On("Create", ctx, mock.Anything).Return(nil).Once()

		loginResp, err := authService.Login(ctx, "testadmin", "password123", "127.0.0.1", nil)
		assert.NoError(t, err)

		// When: 리프레시 토큰으로 새 토큰 발급
		mockUserRepo.On("FindByID", ctx, uint(1)).Return(testUser, nil)

		tokenResp, err := authService.RefreshToken(ctx, loginResp.RefreshToken)

		// Then: 새 토큰 발급 성공
		assert.NoError(t, err)
		assert.NotNil(t, tokenResp)
		assert.NotEmpty(t, tokenResp.AccessToken)
		assert.NotEmpty(t, tokenResp.RefreshToken)

		// And: 새로 발급받은 액세스 토큰도 즉시 사용 가능하다
		claims, err := jwtService.ValidateToken(tokenResp.AccessToken)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "testadmin", claims.Username)
		assert.Equal(t, "SUPER_ADMIN", claims.Authorities)
		assert.Equal(t, "1", claims.Subject)

		// And: 새 리프레시 토큰이 유효하게 검증된다
		newRefreshClaims, err := jwtService.ValidateRefreshToken(tokenResp.RefreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, newRefreshClaims)
		assert.Equal(t, "1", newRefreshClaims.Subject)
	})

	t.Run("잘못된 리프레시 토큰으로는 새 토큰을 발급받을 수 없다", func(t *testing.T) {
		// Given: 잘못된 리프레시 토큰
		invalidToken := "invalid.refresh.token"

		// When: 리프레시 시도
		tokenResp, err := authService.RefreshToken(ctx, invalidToken)

		// Then: 에러 발생
		assert.Error(t, err)
		assert.Nil(t, tokenResp)
	})

	t.Run("존재하지 않는 사용자의 리프레시 토큰으로는 새 토큰을 발급받을 수 없다", func(t *testing.T) {
		// Given: 존재하지 않는 사용자의 유효한 JWT 토큰
		orphanToken, err := jwtService.GenerateRefreshToken(999, "nonexistent")
		assert.NoError(t, err)

		// 사용자 조회 실패로 설정
		mockUserRepo.On("FindByID", ctx, uint(999)).Return(nil, assert.AnError)

		// When: 리프레시 시도
		tokenResp, err := authService.RefreshToken(ctx, orphanToken)

		// Then: 에러 발생
		assert.Error(t, err)
		assert.Nil(t, tokenResp)
	})

	t.Run("비활성화된 사용자는 리프레시 토큰 사용 불가", func(t *testing.T) {
		// Given: 유효한 리프레시 토큰이지만 사용자가 비활성화됨
		validToken, err := jwtService.GenerateRefreshToken(1, "device123")
		assert.NoError(t, err)

		inactiveUser := &models.User{
			UserID: "testadmin",
			Status: models.UserStatusInactive,
		}
		inactiveUser.ID = 1

		// Clear previous mock calls to avoid conflicts
		mockUserRepo.ExpectedCalls = nil
		mockUserRepo.Calls = nil

		mockUserRepo.On("FindByID", ctx, uint(1)).Return(inactiveUser, nil).Once()

		// When: 리프레시 시도
		tokenResp, err := authService.RefreshToken(ctx, validToken)

		// Then: 사용자 비활성화 에러 발생
		assert.Error(t, err)
		assert.Equal(t, services.ErrUserNotActive, err)
		assert.Nil(t, tokenResp)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("JWT 기반 로그아웃은 서버 사이드 액션 없이 성공", func(t *testing.T) {
		// Given: 유효한 사용자 ID
		userID := uint(1)

		// When: 로그아웃 시도
		err := authService.Logout(ctx, userID)

		// Then: 항상 성공 (JWT는 stateless)
		assert.NoError(t, err)
	})
}

func TestAuthService_LoginFailureCases(t *testing.T) {
	// Given: 테스트 환경 설정
	ctx := context.Background()
	jwtService := auth.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)

	t.Run("존재하지 않는 사용자로 로그인 시도 시 실패한다", func(t *testing.T) {
		// Given: 존재하지 않는 사용자
		mockUserRepo := new(MockUserRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)
		authService := services.NewAuthService(mockUserRepo, mockLoginAttemptRepo, jwtService)

		mockLoginAttemptRepo.On("CountRecentFailedAttempts", ctx, "nonexistent", "127.0.0.1", mock.Anything).Return(int64(0), nil)
		mockUserRepo.On("FindByUserID", ctx, "nonexistent").Return(nil, assert.AnError)
		mockUserRepo.On("FindByEmail", ctx, "nonexistent").Return(nil, assert.AnError)
		mockLoginAttemptRepo.On("Create", ctx, mock.Anything).Return(nil)

		// When: 로그인 시도
		loginResp, err := authService.Login(ctx, "nonexistent", "password123", "127.0.0.1", nil)

		// Then: ErrInvalidCredentials 에러 발생
		assert.Error(t, err)
		assert.Equal(t, services.ErrInvalidCredentials, err)
		assert.Nil(t, loginResp)
		mockLoginAttemptRepo.AssertExpectations(t)
	})

	t.Run("잘못된 비밀번호로 로그인 시도 시 실패한다", func(t *testing.T) {
		// Given: 올바른 사용자이지만 잘못된 비밀번호
		mockUserRepo := new(MockUserRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)
		authService := services.NewAuthService(mockUserRepo, mockLoginAttemptRepo, jwtService)

		testUser := &models.User{
			UserID:   "testuser",
			Email:    "testuser@example.com",
			Name:     "Test User",
			Password: "{bcrypt}$2a$10$yS/Y3Y0OcBZ9VFaNeTmpEuI6Vk1jbl5dke9prZNYZOduhmy2xu7T2", // password: password123
			Status:   models.UserStatusActive,
			Role:     models.UserRoleNormal,
		}
		testUser.ID = 1

		mockLoginAttemptRepo.On("CountRecentFailedAttempts", ctx, "testuser", "127.0.0.1", mock.Anything).Return(int64(0), nil)
		mockUserRepo.On("FindByUserID", ctx, "testuser").Return(testUser, nil)
		mockLoginAttemptRepo.On("Create", ctx, mock.Anything).Return(nil)

		// When: 잘못된 비밀번호로 로그인 시도
		loginResp, err := authService.Login(ctx, "testuser", "wrongpassword", "127.0.0.1", nil)

		// Then: ErrInvalidCredentials 에러 발생
		assert.Error(t, err)
		assert.Equal(t, services.ErrInvalidCredentials, err)
		assert.Nil(t, loginResp)
		mockLoginAttemptRepo.AssertExpectations(t)
	})

	t.Run("비활성 사용자로 로그인 시도 시 실패한다", func(t *testing.T) {
		// Given: 비활성화된 사용자
		mockUserRepo := new(MockUserRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)
		authService := services.NewAuthService(mockUserRepo, mockLoginAttemptRepo, jwtService)

		inactiveUser := &models.User{
			UserID:   "inactiveuser",
			Email:    "inactive@example.com",
			Name:     "Inactive User",
			Password: "{bcrypt}$2a$10$yS/Y3Y0OcBZ9VFaNeTmpEuI6Vk1jbl5dke9prZNYZOduhmy2xu7T2", // password: password123
			Status:   models.UserStatusInactive,
			Role:     models.UserRoleNormal,
		}
		inactiveUser.ID = 2

		mockLoginAttemptRepo.On("CountRecentFailedAttempts", ctx, "inactiveuser", "127.0.0.1", mock.Anything).Return(int64(0), nil)
		mockUserRepo.On("FindByUserID", ctx, "inactiveuser").Return(inactiveUser, nil)
		mockLoginAttemptRepo.On("Create", ctx, mock.Anything).Return(nil)

		// When: 비활성 사용자로 로그인 시도
		loginResp, err := authService.Login(ctx, "inactiveuser", "password123", "127.0.0.1", nil)

		// Then: ErrUserNotActive 에러 발생
		assert.Error(t, err)
		assert.Equal(t, services.ErrUserNotActive, err)
		assert.Nil(t, loginResp)
		mockLoginAttemptRepo.AssertExpectations(t)
	})

	t.Run("이메일로 로그인 성공한다", func(t *testing.T) {
		// Given: 이메일로 로그인 시도
		mockUserRepo := new(MockUserRepository)
		mockLoginAttemptRepo := new(MockLoginAttemptRepository)
		authService := services.NewAuthService(mockUserRepo, mockLoginAttemptRepo, jwtService)

		testUser := &models.User{
			UserID:   "testuser",
			Email:    "testuser@example.com",
			Name:     "Test User",
			Password: "{bcrypt}$2a$10$yS/Y3Y0OcBZ9VFaNeTmpEuI6Vk1jbl5dke9prZNYZOduhmy2xu7T2", // password: password123
			Status:   models.UserStatusActive,
			Role:     models.UserRoleNormal,
		}
		testUser.ID = 1

		// UserID로는 찾지 못하고, 이메일로 찾음
		mockLoginAttemptRepo.On("CountRecentFailedAttempts", ctx, "testuser@example.com", "127.0.0.1", mock.Anything).Return(int64(0), nil)
		mockUserRepo.On("FindByUserID", ctx, "testuser@example.com").Return(nil, assert.AnError)
		mockUserRepo.On("FindByEmail", ctx, "testuser@example.com").Return(testUser, nil)
		mockLoginAttemptRepo.On("Create", ctx, mock.Anything).Return(nil)

		// When: 이메일로 로그인 시도
		loginResp, err := authService.Login(ctx, "testuser@example.com", "password123", "127.0.0.1", nil)

		// Then: 로그인 성공
		assert.NoError(t, err)
		assert.NotNil(t, loginResp)
		assert.Equal(t, "testuser", loginResp.User.UserID)
		assert.NotEmpty(t, loginResp.AccessToken)
		assert.NotEmpty(t, loginResp.RefreshToken)

		// 토큰 검증
		claims, err := jwtService.ValidateToken(loginResp.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, "testuser", claims.Username)
	})
}
