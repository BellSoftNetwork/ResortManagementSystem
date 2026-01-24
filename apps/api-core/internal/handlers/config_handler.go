package handlers

import (
	"github.com/gin-gonic/gin"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
)

type ConfigHandler struct {
	configService services.ConfigService
}

func NewConfigHandler(configService services.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configService: configService,
	}
}

// GetConfig returns the server configuration
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	config := h.configService.GetConfig()
	response.Success(c, config)
}

// GetEnvironment returns the server environment information
func (h *ConfigHandler) GetEnvironment(c *gin.Context) {
	env := h.configService.GetEnvironmentResponse()
	response.Success(c, env)
}
