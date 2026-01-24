package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
)

type DevelopmentHandler struct {
	developmentService services.DevelopmentService
}

func NewDevelopmentHandler(developmentService services.DevelopmentService) *DevelopmentHandler {
	return &DevelopmentHandler{
		developmentService: developmentService,
	}
}

// GenerateTestData godoc
// @Summary Generate test data for development
// @Description Generates test data based on the specified type (essential, reservation, or all)
// @Description Generates payment methods, room groups, and rooms for development/testing
// @Tags Development
// @Accept json
// @Produce json
// @Param request body dto.GenerateTestDataRequest true "Generate test data request"
// @Success 200 {object} dto.GenerateTestDataResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Security Bearer
// @Router /api/v1/dev/generate-test-data [post]
func (h *DevelopmentHandler) GenerateTestData(c *gin.Context) {
	var req dto.GenerateTestDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"errors":  []string{err.Error()},
		})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation failed",
			"errors":  []string{err.Error()},
		})
		return
	}

	// Convert DTO options to service options
	var serviceOptions *services.ReservationGenerationOptions
	if req.ReservationOptions != nil {
		serviceOptions = &services.ReservationGenerationOptions{
			StartDate:           req.ReservationOptions.StartDate,
			EndDate:             req.ReservationOptions.EndDate,
			RegularReservations: req.ReservationOptions.RegularReservations,
			MonthlyReservations: req.ReservationOptions.MonthlyReservations,
		}
	}

	// Generate test data
	result, err := h.developmentService.GenerateTestData(string(req.Type), serviceOptions)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate test data")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate test data",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.GenerateTestDataResponse{
		Message: "Test data generated successfully",
		Data:    result,
	})
}

// GenerateEssentialData godoc
// @Summary Generate essential data for development
// @Description Generates payment methods, room groups, and rooms for development/testing
// @Tags Development
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} middleware.ErrorResponse
// @Security Bearer
// @Router /api/v1/dev/generate-essential-data [post]
func (h *DevelopmentHandler) GenerateEssentialData(c *gin.Context) {
	result, err := h.developmentService.GenerateTestData("essential", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate essential data",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Essential data generated successfully",
		"data":    result,
	})
}

// GenerateReservationData godoc
// @Summary Generate reservation data for development
// @Description Generates sample reservations for development/testing
// @Tags Development
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} middleware.ErrorResponse
// @Security Bearer
// @Router /api/v1/dev/generate-reservation-data [post]
func (h *DevelopmentHandler) GenerateReservationData(c *gin.Context) {
	result, err := h.developmentService.GenerateTestData("reservation", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate reservation data",
			"errors":  []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reservation data generated successfully",
		"data":    result,
	})
}
