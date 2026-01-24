package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
)

type PaymentMethodHandler struct {
	paymentMethodService services.PaymentMethodService
}

func NewPaymentMethodHandler(paymentMethodService services.PaymentMethodService) *PaymentMethodHandler {
	return &PaymentMethodHandler{
		paymentMethodService: paymentMethodService,
	}
}

func (h *PaymentMethodHandler) ListPaymentMethods(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "잘못된 쿼리 파라미터", err.Error())
		return
	}

	paymentMethods, total, err := h.paymentMethodService.GetAll(c.Request.Context(), query.Page, query.Size, query.Sort)
	if err != nil {
		response.InternalServerError(c, "결제 수단 목록 조회 실패")
		return
	}

	paymentMethodResponses := make([]dto.PaymentMethodResponse, len(paymentMethods))
	for i, paymentMethod := range paymentMethods {
		paymentMethodResponses[i] = h.toPaymentMethodResponse(&paymentMethod)
	}

	totalPages := int(total) / query.Size
	if int(total)%query.Size > 0 {
		totalPages++
	}

	pagination := &response.Pagination{
		Page:          query.Page,
		Size:          query.Size,
		TotalPages:    totalPages,
		TotalElements: total,
	}

	// Payment Method 리스트는 현재 필터를 지원하지 않으므로 빈 필터 객체 반환
	filterResponse := map[string]interface{}{}

	response.SuccessListWithFilter(c, paymentMethodResponses, pagination, filterResponse)
}

func (h *PaymentMethodHandler) GetPaymentMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 결제 수단 ID")
		return
	}

	paymentMethod, err := h.paymentMethodService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, services.ErrPaymentMethodNotFound) {
			response.NotFound(c, "존재하지 않는 결제 수단")
			return
		}
		response.InternalServerError(c, "결제 수단 조회 실패")
		return
	}

	paymentMethodResponse := h.toPaymentMethodResponse(paymentMethod)
	response.Success(c, paymentMethodResponse)
}

func (h *PaymentMethodHandler) CreatePaymentMethod(c *gin.Context) {
	var req dto.CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	paymentMethod := &models.PaymentMethod{
		Name:                     req.Name,
		CommissionRate:           req.CommissionRate,
		RequireUnpaidAmountCheck: models.BitBool(req.RequireUnpaidAmountCheck),
		IsDefaultSelect:          models.BitBool(false),
		Status:                   models.PaymentMethodStatusInactive,
	}

	if err := h.paymentMethodService.Create(c.Request.Context(), paymentMethod); err != nil {
		if errors.Is(err, services.ErrPaymentMethodNameExists) {
			response.Conflict(c, "이미 존재하는 결제 수단")
			return
		}
		// Log the actual error for debugging
		gin.DefaultErrorWriter.Write([]byte("PaymentMethod creation error: " + err.Error() + "\n"))
		response.InternalServerError(c, "결제 수단 등록 실패")
		return
	}

	paymentMethodResponse := h.toPaymentMethodResponse(paymentMethod)
	response.Created(c, paymentMethodResponse)
}

func (h *PaymentMethodHandler) UpdatePaymentMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 결제 수단 ID")
		return
	}

	var req dto.UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.CommissionRate != nil {
		updates["commissionRate"] = *req.CommissionRate
	}
	if req.RequireUnpaidAmountCheck != nil {
		updates["requireUnpaidAmountCheck"] = models.BitBool(*req.RequireUnpaidAmountCheck)
	}
	if req.IsDefaultSelect != nil {
		updates["isDefaultSelect"] = models.BitBool(*req.IsDefaultSelect)
	}
	if req.Status != nil {
		switch *req.Status {
		case "ACTIVE":
			updates["status"] = models.PaymentMethodStatusActive
		case "INACTIVE":
			updates["status"] = models.PaymentMethodStatusInactive
		}
	}

	paymentMethod, err := h.paymentMethodService.Update(c.Request.Context(), uint(id), updates)
	if err != nil {
		if errors.Is(err, services.ErrPaymentMethodNotFound) {
			response.NotFound(c, "존재하지 않는 결제 수단")
			return
		}
		if errors.Is(err, services.ErrPaymentMethodNameExists) {
			response.Conflict(c, "이미 존재하는 결제 수단")
			return
		}
		response.InternalServerError(c, "결제 수단 수정 실패")
		return
	}

	paymentMethodResponse := h.toPaymentMethodResponse(paymentMethod)
	response.Success(c, paymentMethodResponse)
}

func (h *PaymentMethodHandler) DeletePaymentMethod(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 결제 수단 ID")
		return
	}

	if err := h.paymentMethodService.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, services.ErrPaymentMethodNotFound) {
			response.NotFound(c, "존재하지 않는 결제 수단")
			return
		}
		if errors.Is(err, services.ErrPaymentMethodInUse) {
			response.Conflict(c, "사용 중인 결제 수단은 삭제할 수 없습니다")
			return
		}
		response.InternalServerError(c, "결제 수단 삭제 실패")
		return
	}

	response.NoContent(c)
}

func (h *PaymentMethodHandler) toPaymentMethodResponse(paymentMethod *models.PaymentMethod) dto.PaymentMethodResponse {
	return dto.PaymentMethodResponse{
		ID:                       paymentMethod.ID,
		Name:                     paymentMethod.Name,
		CommissionRate:           paymentMethod.CommissionRate,
		RequireUnpaidAmountCheck: bool(paymentMethod.RequireUnpaidAmountCheck),
		IsDefaultSelect:          bool(paymentMethod.IsDefaultSelect),
		Status:                   paymentMethod.Status.String(),
		CreatedAt:                dto.CustomTime{Time: paymentMethod.CreatedAt},
		UpdatedAt:                dto.CustomTime{Time: paymentMethod.UpdatedAt},
	}
}
