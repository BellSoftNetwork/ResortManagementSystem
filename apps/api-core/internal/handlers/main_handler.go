package handlers

import (
	"github.com/gin-gonic/gin"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
)

type MainHandler struct {
	configService services.ConfigService
	userRepo      repositories.UserRepository
}

func NewMainHandler(configService services.ConfigService, userRepo repositories.UserRepository) *MainHandler {
	return &MainHandler{
		configService: configService,
		userRepo:      userRepo,
	}
}

func (h *MainHandler) GetEnvironment(c *gin.Context) {
	envInfo := h.configService.GetEnvInfo()
	response.Success(c, envInfo)
}

func (h *MainHandler) GetConfig(c *gin.Context) {
	hasUsers, err := h.userRepo.HasAnyUsers(c.Request.Context())
	if err != nil {
		// If error, assume users exist (safer default)
		hasUsers = true
	}

	config := h.configService.GetAppConfig(hasUsers)
	response.Success(c, config)
}
