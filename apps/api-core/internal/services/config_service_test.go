package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.bellsoft.net/rms/api-core/internal/config"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
)

type ConfigServiceTestSuite struct {
	suite.Suite
	service services.ConfigService
}

func (suite *ConfigServiceTestSuite) SetupTest() {
	// ConfigService를 위한 mock config 생성
	mockConfig := &config.Config{
		Environment: "test",
		App: config.AppConfig{
			Deploy: config.DeployConfig{
				CommitSHA:       "unknown",
				CommitShortSHA:  "unknown",
				CommitTitle:     "unknown",
				CommitTimestamp: "unknown",
			},
		},
	}
	suite.service = services.NewConfigService(mockConfig)
}

func (suite *ConfigServiceTestSuite) TestGetAppConfig_NoUsers() {
	// Given - 가입한 사용자가 없는 상황에서
	hasUsers := false

	// When - 앱 설정 정보 요청 시
	appConfig := suite.service.GetAppConfig(hasUsers)

	// Then - 회원 가입 가능 여부 플래그가 활성화된다
	assert.True(suite.T(), appConfig.IsAvailableRegistration)
}

func (suite *ConfigServiceTestSuite) TestGetAppConfig_WithUsers() {
	// Given - 기존에 가입한 사용자가 있는 상황에서
	hasUsers := true

	// When - 앱 설정 정보 요청 시
	appConfig := suite.service.GetAppConfig(hasUsers)

	// Then - 회원 가입 가능 여부 플래그가 비활성화된다
	assert.False(suite.T(), appConfig.IsAvailableRegistration)
}

func (suite *ConfigServiceTestSuite) TestGetEnvironment() {
	// When - 환경 정보를 요청하면
	env := suite.service.GetEnvironment()

	// Then - 환경 정보가 반환된다
	assert.NotEmpty(suite.T(), env)
	// test 환경이므로 test를 기대
	expectedEnv := "test"
	assert.Equal(suite.T(), expectedEnv, env)
}

func (suite *ConfigServiceTestSuite) TestGetEnvInfo() {
	// When - 애플리케이션 정보를 요청하면
	envInfo := suite.service.GetEnvInfo()

	// Then - 애플리케이션 정보가 반환된다
	assert.NotNil(suite.T(), envInfo)
	assert.Equal(suite.T(), "Resort Management System", envInfo.ApplicationFullName)
	assert.Equal(suite.T(), "RMS", envInfo.ApplicationShortName)
	// 빌드 시 설정되지 않은 경우 unknown 값 확인
	assert.Equal(suite.T(), "unknown", envInfo.CommitSHA)
	assert.Equal(suite.T(), "unknown", envInfo.CommitShortSHA)
	assert.Equal(suite.T(), "unknown", envInfo.CommitTitle)
	assert.Equal(suite.T(), "unknown", envInfo.CommitTimestamp)
}

func TestConfigServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigServiceTestSuite))
}
