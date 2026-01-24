package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/auth"
	"gorm.io/gorm"
)

type LoginAttemptServiceTestSuite struct {
	suite.Suite
	ctx                  context.Context
	authService          services.AuthService
	mockUserRepo         *MockUserRepository
	mockLoginAttemptRepo *MockLoginAttemptRepository
	jwtService           *auth.JWTService
	testUser             *models.User
}

func (suite *LoginAttemptServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()

	// JWT 서비스 설정
	suite.jwtService = auth.NewJWTService("test-secret", 15*time.Minute, 7*24*time.Hour)

	// Mock 설정
	suite.mockUserRepo = new(MockUserRepository)
	suite.mockLoginAttemptRepo = new(MockLoginAttemptRepository)

	// AuthService 생성
	suite.authService = services.NewAuthService(
		suite.mockUserRepo,
		suite.mockLoginAttemptRepo,
		suite.jwtService,
	)

	// 테스트 사용자 설정
	suite.testUser = &models.User{
		UserID:   "testuser",
		Email:    "testuser@example.com",
		Name:     "Test User",
		Password: "{bcrypt}$2a$10$yS/Y3Y0OcBZ9VFaNeTmpEuI6Vk1jbl5dke9prZNYZOduhmy2xu7T2", // password: password123
		Status:   models.UserStatusActive,
		Role:     models.UserRoleNormal,
	}
	suite.testUser.ID = 1
}

func (suite *LoginAttemptServiceTestSuite) TearDownTest() {
}

func (suite *LoginAttemptServiceTestSuite) TestRecordLoginAttempt() {
	// Given - 로그인 시도 서비스
	username := "testuser"
	ipAddress := "192.168.1.100"
	deviceInfo := &services.DeviceInfo{
		OSInfo:            "Windows",
		LanguageInfo:      "ko-KR",
		UserAgent:         "Mozilla/5.0",
		DeviceFingerprint: "test-fingerprint",
	}

	// 로그인 성공 케이스 설정
	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, username, ipAddress, mock.Anything).Return(int64(0), nil)
	suite.mockUserRepo.On("FindByUserID", suite.ctx, username).Return(suite.testUser, nil)
	suite.mockLoginAttemptRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.LoginAttempt")).Return(nil)

	// When - 로그인 시도를 기록하면
	loginResp, err := suite.authService.Login(suite.ctx, username, "password123", ipAddress, deviceInfo)

	// Then - 데이터베이스에 기록된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), loginResp)

	// 로그인 시도가 기록되었는지 확인
	suite.mockLoginAttemptRepo.AssertCalled(suite.T(), "Create", suite.ctx, mock.MatchedBy(func(attempt *models.LoginAttempt) bool {
		return attempt.Username == username &&
			attempt.IPAddress == ipAddress &&
			attempt.Successful == true &&
			attempt.DeviceFingerprint != nil && *attempt.DeviceFingerprint == deviceInfo.DeviceFingerprint
	}))
}

func (suite *LoginAttemptServiceTestSuite) TestLoginAttemptWindowExpiry() {
	// Given - 시간이 지나 이전 실패 시도가 윈도우를 벗어나면
	username := "testuser"
	ipAddress := "192.168.1.100"
	deviceInfo := &services.DeviceInfo{
		OSInfo: "Windows",
	}

	// 15분 이전의 실패 시도는 카운트에 포함되지 않음
	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, username, ipAddress, mock.Anything).Return(int64(0), nil)
	suite.mockUserRepo.On("FindByUserID", suite.ctx, username).Return(suite.testUser, nil)
	suite.mockLoginAttemptRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.LoginAttempt")).Return(nil)

	// When/Then - 로그인 시도가 허용된다
	loginResp, err := suite.authService.Login(suite.ctx, username, "password123", ipAddress, deviceInfo)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), loginResp)
}

func (suite *LoginAttemptServiceTestSuite) TestDeviceFingerprint() {
	// Given - 디바이스 핑거프린트 관련 테스트
	username := "testuser"
	ipAddress := "192.168.1.100"

	deviceOneInfo := &services.DeviceInfo{
		DeviceFingerprint: "windows-fingerprint",
	}
	deviceTwoInfo := &services.DeviceInfo{
		DeviceFingerprint: "android-fingerprint",
	}
	deviceThreeInfo := &services.DeviceInfo{
		DeviceFingerprint: "",
	}

	// 첫 번째 디바이스로 성공한 로그인 시도 기록
	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, username, ipAddress, mock.Anything).Return(int64(0), nil).Once()
	suite.mockUserRepo.On("FindByUserID", suite.ctx, username).Return(suite.testUser, nil).Once()
	suite.mockLoginAttemptRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.LoginAttempt")).Return(nil).Once()

	loginResp, err := suite.authService.Login(suite.ctx, username, "password123", ipAddress, deviceOneInfo)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), loginResp)

	// 로그인 시도가 올바른 디바이스 정보로 기록되었는지 확인
	suite.mockLoginAttemptRepo.AssertCalled(suite.T(), "Create", suite.ctx, mock.MatchedBy(func(attempt *models.LoginAttempt) bool {
		return attempt.DeviceFingerprint != nil && *attempt.DeviceFingerprint == deviceOneInfo.DeviceFingerprint
	}))

	// 다른 디바이스로 로그인 시도
	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, username, ipAddress, mock.Anything).Return(int64(0), nil).Once()
	suite.mockUserRepo.On("FindByUserID", suite.ctx, username).Return(suite.testUser, nil).Once()
	suite.mockLoginAttemptRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.LoginAttempt")).Return(nil).Once()

	loginResp2, err := suite.authService.Login(suite.ctx, username, "password123", ipAddress, deviceTwoInfo)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), loginResp2)

	// 빈 디바이스 핑거프린트로 로그인 시도
	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, username, ipAddress, mock.Anything).Return(int64(0), nil).Once()
	suite.mockUserRepo.On("FindByUserID", suite.ctx, username).Return(suite.testUser, nil).Once()
	suite.mockLoginAttemptRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.LoginAttempt")).Return(nil).Once()

	loginResp3, err := suite.authService.Login(suite.ctx, username, "password123", ipAddress, deviceThreeInfo)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), loginResp3)
}

func (suite *LoginAttemptServiceTestSuite) TestLoginAttemptLimit_4Failures() {
	// Given - 동일 IP와 ID로 로그인을 4번 실패하고 5번째 시도 시
	username := "testuser"
	ipAddress := "192.168.1.100"
	deviceInfo := &services.DeviceInfo{}

	// 4번의 실패한 로그인 시도 기록
	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, username, ipAddress, mock.Anything).Return(int64(4), nil)
	suite.mockUserRepo.On("FindByUserID", suite.ctx, username).Return(suite.testUser, nil)
	suite.mockLoginAttemptRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.LoginAttempt")).Return(nil)

	// When/Then - 5번째 시도에서는 로그인이 허용되어야 한다
	loginResp, err := suite.authService.Login(suite.ctx, username, "password123", ipAddress, deviceInfo)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), loginResp)
}

func (suite *LoginAttemptServiceTestSuite) TestLoginAttemptLimit_5Failures() {
	// Given - 동일 IP와 ID로 로그인을 5번 실패하고 6번째 시도 시
	username := "testuser"
	ipAddress := "192.168.1.100"
	deviceInfo := &services.DeviceInfo{}

	// 5번의 실패한 로그인 시도 기록
	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, username, ipAddress, mock.Anything).Return(int64(5), nil)

	// When - 6번째 시도에서는 TooManyAttempts 에러가 발생해야 한다
	loginResp, err := suite.authService.Login(suite.ctx, username, "password123", ipAddress, deviceInfo)

	// Then
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrTooManyAttempts, err)
	assert.Nil(suite.T(), loginResp)
}

func (suite *LoginAttemptServiceTestSuite) TestLoginAttemptLimit_DifferentIPs() {
	// Given - 1번 IP에서 ID를 연속 로그인 실패하여 429 상태가 되고 2번 IP로 같은 ID 로그인 시도 시
	adminUsername := "admin"
	ipAddress1 := "192.168.1.100"
	ipAddress2 := "192.168.1.200"
	deviceInfo1 := &services.DeviceInfo{}
	deviceInfo2 := &services.DeviceInfo{}

	// 1번 IP에서 관리자 ID로 5번 실패하여 429 상태
	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, adminUsername, ipAddress1, mock.Anything).Return(int64(5), nil)

	// 1번 IP에서는 로그인 시도가 거부되는지 확인
	loginResp1, err := suite.authService.Login(suite.ctx, adminUsername, "wrongpassword", ipAddress1, deviceInfo1)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrTooManyAttempts, err)
	assert.Nil(suite.T(), loginResp1)

	// 2번 IP에서는 동일 ID로 정상 로그인이 가능해야 한다
	adminUser := &models.User{
		UserID:   "admin",
		Email:    "admin@example.com",
		Name:     "Admin User",
		Password: "{bcrypt}$2a$10$yS/Y3Y0OcBZ9VFaNeTmpEuI6Vk1jbl5dke9prZNYZOduhmy2xu7T2",
		Status:   models.UserStatusActive,
		Role:     models.UserRoleAdmin,
	}
	adminUser.ID = 2

	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, adminUsername, ipAddress2, mock.Anything).Return(int64(0), nil)
	suite.mockUserRepo.On("FindByUserID", suite.ctx, adminUsername).Return(adminUser, nil)
	suite.mockLoginAttemptRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.LoginAttempt")).Return(nil)

	loginResp2, err := suite.authService.Login(suite.ctx, adminUsername, "password123", ipAddress2, deviceInfo2)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), loginResp2)
}

func (suite *LoginAttemptServiceTestSuite) TestLoginAttemptLimit_BruteForceProtection() {
	// Given - 악의적인 사용자가 관리자 계정에 대해 지속적인 인증 실패를 발생시킬 때
	adminUsername := "admin"
	hackerIP := "1.2.3.4"       // 공격자 IP
	normalIP := "192.168.1.100" // 정상 사용자 IP
	hackerDeviceInfo := &services.DeviceInfo{}
	normalDeviceInfo := &services.DeviceInfo{}

	adminUser := &models.User{
		UserID:   "admin",
		Email:    "admin@example.com",
		Name:     "Admin User",
		Password: "{bcrypt}$2a$10$yS/Y3Y0OcBZ9VFaNeTmpEuI6Vk1jbl5dke9prZNYZOduhmy2xu7T2",
		Status:   models.UserStatusActive,
		Role:     models.UserRoleAdmin,
	}
	adminUser.ID = 2

	// 공격자 IP에서 관리자 ID로 5번 실패하여 429 상태
	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, adminUsername, hackerIP, mock.Anything).Return(int64(5), nil)

	// 공격자 IP에서는 로그인 시도가 거부되는지 확인
	loginResp1, err := suite.authService.Login(suite.ctx, adminUsername, "wrongpassword", hackerIP, hackerDeviceInfo)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrTooManyAttempts, err)
	assert.Nil(suite.T(), loginResp1)

	// 실제 관리자는 본인 IP에서 정상적으로 로그인이 가능해야 한다
	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, adminUsername, normalIP, mock.Anything).Return(int64(0), nil)
	suite.mockUserRepo.On("FindByUserID", suite.ctx, adminUsername).Return(adminUser, nil)
	suite.mockLoginAttemptRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.LoginAttempt")).Return(nil)

	loginResp2, err := suite.authService.Login(suite.ctx, adminUsername, "password123", normalIP, normalDeviceInfo)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), loginResp2)
}

func (suite *LoginAttemptServiceTestSuite) TestLoginAttemptReset_AfterSuccess() {
	// Given - 동일 IP, 동일 ID로 4번 로그인 실패 후 성공하고 다시 1번 실패 후 시도할 때
	username := "resetuser"
	ipAddress := "192.168.1.100"
	deviceInfo := &services.DeviceInfo{}

	resetUser := &models.User{
		UserID:   "resetuser",
		Email:    "resetuser@example.com",
		Name:     "Reset User",
		Password: "{bcrypt}$2a$10$yS/Y3Y0OcBZ9VFaNeTmpEuI6Vk1jbl5dke9prZNYZOduhmy2xu7T2",
		Status:   models.UserStatusActive,
		Role:     models.UserRoleNormal,
	}
	resetUser.ID = 3

	// 4번의 실패한 로그인 시도 후 성공
	// 성공 후에는 실패 카운트가 리셋되지 않으므로 이전 실패 이력이 남아있음
	// 하지만 시간 창(15분) 내에서만 카운트됨

	// When - 성공 후 1번 실패한 로그인 기록이 있을 때
	// 실패 횟수가 5번 미만이므로 로그인이 허용되어야 한다
	suite.mockLoginAttemptRepo.On("CountRecentFailedAttempts", suite.ctx, username, ipAddress, mock.Anything).Return(int64(1), nil)
	suite.mockUserRepo.On("FindByUserID", suite.ctx, username).Return(resetUser, nil)
	suite.mockLoginAttemptRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.LoginAttempt")).Return(nil)

	// Then - 로그인이 허용되어야 한다
	loginResp, err := suite.authService.Login(suite.ctx, username, "password123", ipAddress, deviceInfo)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), loginResp)
}

func (suite *LoginAttemptServiceTestSuite) TestIsDeviceChanged_DifferentFingerprint() {
	// Given - 디바이스 핑거프린트가 변경된 경우
	username := "testuser"
	windowsFingerprint := "windows-fingerprint"
	lastSuccessfulAttempt := &models.LoginAttempt{
		Username:          username,
		IPAddress:         "192.168.1.100",
		Successful:        true,
		DeviceFingerprint: &windowsFingerprint,
	}

	androidDeviceInfo := &services.DeviceInfo{
		DeviceFingerprint: "android-fingerprint",
	}

	suite.mockLoginAttemptRepo.On("GetLastSuccessfulAttempt", suite.ctx, username).Return(lastSuccessfulAttempt, nil)

	// When - 디바이스 변경 여부를 확인하면
	isChanged, err := suite.authService.IsDeviceChanged(suite.ctx, username, androidDeviceInfo)

	// Then - true가 반환된다
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isChanged)
}

func (suite *LoginAttemptServiceTestSuite) TestIsDeviceChanged_SameFingerprint() {
	// Given - 디바이스 핑거프린트가 동일한 경우
	username := "testuser"
	windowsFingerprint := "windows-fingerprint"
	lastSuccessfulAttempt := &models.LoginAttempt{
		Username:          username,
		IPAddress:         "192.168.1.100",
		Successful:        true,
		DeviceFingerprint: &windowsFingerprint,
	}

	sameDeviceInfo := &services.DeviceInfo{
		DeviceFingerprint: "windows-fingerprint",
	}

	suite.mockLoginAttemptRepo.On("GetLastSuccessfulAttempt", suite.ctx, username).Return(lastSuccessfulAttempt, nil)

	// When - 디바이스 변경 여부를 확인하면
	isChanged, err := suite.authService.IsDeviceChanged(suite.ctx, username, sameDeviceInfo)

	// Then - false가 반환된다
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isChanged)
}

func (suite *LoginAttemptServiceTestSuite) TestIsDeviceChanged_EmptyFingerprint() {
	// Given - 새로운 디바이스 핑거프린트가 빈 값인 경우
	username := "testuser"
	windowsFingerprint := "windows-fingerprint"
	lastSuccessfulAttempt := &models.LoginAttempt{
		Username:          username,
		IPAddress:         "192.168.1.100",
		Successful:        true,
		DeviceFingerprint: &windowsFingerprint,
	}

	emptyDeviceInfo := &services.DeviceInfo{
		DeviceFingerprint: "",
	}

	suite.mockLoginAttemptRepo.On("GetLastSuccessfulAttempt", suite.ctx, username).Return(lastSuccessfulAttempt, nil)

	// When - 디바이스 변경 여부를 확인하면
	isChanged, err := suite.authService.IsDeviceChanged(suite.ctx, username, emptyDeviceInfo)

	// Then - true가 반환된다 (빈 값은 변경으로 간주)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isChanged)
}

func (suite *LoginAttemptServiceTestSuite) TestIsDeviceChanged_NoSuccessfulLoginHistory() {
	// Given - 이전에 성공한 로그인 기록이 없는 경우
	username := "newuser"
	deviceInfo := &services.DeviceInfo{
		DeviceFingerprint: "new-fingerprint",
	}

	suite.mockLoginAttemptRepo.On("GetLastSuccessfulAttempt", suite.ctx, username).Return(nil, gorm.ErrRecordNotFound)

	// When - 디바이스 변경 여부를 확인하면
	isChanged, err := suite.authService.IsDeviceChanged(suite.ctx, username, deviceInfo)

	// Then - false가 반환된다 (비교 대상 없음)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isChanged)
}

func (suite *LoginAttemptServiceTestSuite) TestIsDeviceChanged_NilLastFingerprint() {
	// Given - 마지막 성공한 로그인에 핑거프린트가 없는 경우
	username := "testuser"
	lastSuccessfulAttempt := &models.LoginAttempt{
		Username:          username,
		IPAddress:         "192.168.1.100",
		Successful:        true,
		DeviceFingerprint: nil,
	}

	deviceInfo := &services.DeviceInfo{
		DeviceFingerprint: "new-fingerprint",
	}

	suite.mockLoginAttemptRepo.On("GetLastSuccessfulAttempt", suite.ctx, username).Return(lastSuccessfulAttempt, nil)

	// When - 디바이스 변경 여부를 확인하면
	isChanged, err := suite.authService.IsDeviceChanged(suite.ctx, username, deviceInfo)

	// Then - false가 반환된다 (비교 대상 없음)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isChanged)
}

func TestLoginAttemptServiceTestSuite(t *testing.T) {
	suite.Run(t, new(LoginAttemptServiceTestSuite))
}
