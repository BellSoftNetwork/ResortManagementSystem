package services

import (
	"os"
	"time"

	"gitlab.bellsoft.net/rms/api-core/internal/config"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
)

type ConfigService interface {
	GetEnvironment() string
	GetEnvInfo() EnvInfo
	GetAppConfig(hasUsers bool) AppConfig
	GetConfig() *dto.ConfigResponse
	GetEnvironmentResponse() *dto.EnvironmentResponse
}

type AppConfig struct {
	IsAvailableRegistration bool `json:"isAvailableRegistration"`
}

type EnvInfo struct {
	ApplicationFullName  string `json:"applicationFullName"`
	ApplicationShortName string `json:"applicationShortName"`
	CommitSHA            string `json:"commitSha"`
	CommitShortSHA       string `json:"commitShortSha"`
	CommitTitle          string `json:"commitTitle"`
	CommitTimestamp      string `json:"commitTimestamp"`
}

type configService struct {
	config *config.Config
}

func NewConfigService(cfg *config.Config) ConfigService {
	return &configService{
		config: cfg,
	}
}

func (s *configService) GetEnvironment() string {
	return s.config.Environment
}

func (s *configService) GetEnvInfo() EnvInfo {
	return EnvInfo{
		ApplicationFullName:  "Resort Management System",
		ApplicationShortName: "RMS",
		CommitSHA:            s.config.App.Deploy.CommitSHA,
		CommitShortSHA:       s.config.App.Deploy.CommitShortSHA,
		CommitTitle:          s.config.App.Deploy.CommitTitle,
		CommitTimestamp:      s.config.App.Deploy.CommitTimestamp,
	}
}

func (s *configService) GetAppConfig(hasUsers bool) AppConfig {
	return AppConfig{
		IsAvailableRegistration: !hasUsers, // Registration is available when no users exist
	}
}

var startTime = time.Now()

func (s *configService) GetConfig() *dto.ConfigResponse {
	hostname, _ := os.Hostname()

	return &dto.ConfigResponse{
		API: dto.APIConfig{
			Host:    hostname,
			Port:    s.config.Server.Port,
			Profile: s.config.Environment,
		},
		Database: dto.DatabaseConfig{
			Host:     s.config.Database.Host,
			Port:     s.config.Database.Port,
			Database: s.config.Database.Database,
		},
		Redis: dto.RedisConfig{
			Host: s.config.Redis.Host,
			Port: s.config.Redis.Port,
		},
	}
}

func (s *configService) GetEnvironmentResponse() *dto.EnvironmentResponse {
	hostname, _ := os.Hostname()
	uptime := time.Since(startTime).String()

	return &dto.EnvironmentResponse{
		Profile:  s.config.Environment,
		Hostname: hostname,
		Version:  s.config.App.Version,
		Uptime:   uptime,
	}
}
