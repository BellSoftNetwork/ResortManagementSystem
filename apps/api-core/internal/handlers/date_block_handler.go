package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	appContext "gitlab.bellsoft.net/rms/api-core/internal/context"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
)

type DateBlockHandler struct {
	service        services.DateBlockService
	historyService services.HistoryService
}

func NewDateBlockHandler(service services.DateBlockService, historyService services.HistoryService) *DateBlockHandler {
	return &DateBlockHandler{service: service, historyService: historyService}
}

func (h *DateBlockHandler) ListDateBlocks(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "잘못된 쿼리 파라미터", err.Error())
		return
	}

	var filter dto.DateBlockFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, "잘못된 필터 파라미터", err.Error())
		return
	}

	dateBlocks, total, err := h.service.GetAll(c.Request.Context(), filter, query.Page, query.Size)
	if err != nil {
		response.InternalServerError(c, "날짜 차단 목록 조회 실패")
		return
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

	response.SuccessListWithFilter(c, dateBlocks, pagination, filter)
}

func (h *DateBlockHandler) CreateDateBlock(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	var req dto.CreateDateBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	ctx := appContext.WithUserID(c.Request.Context(), userID)
	created, err := h.service.Create(ctx, req)
	if err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	response.Created(c, created)
}

func (h *DateBlockHandler) DeleteDateBlock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 날짜 차단 ID")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	ctx := appContext.WithUserID(c.Request.Context(), userID)
	if err := h.service.Delete(ctx, uint(id)); err != nil {
		if errors.Is(err, services.ErrDateBlockNotFound) {
			response.NotFound(c, "존재하지 않는 날짜 차단")
			return
		}
		response.InternalServerError(c, "날짜 차단 삭제 실패")
		return
	}

	response.NoContent(c)
}

func (h *DateBlockHandler) GetDateBlock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 날짜 차단 ID")
		return
	}

	dateBlock, err := h.service.GetDateBlock(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, services.ErrDateBlockNotFound) {
			response.NotFound(c, "존재하지 않는 날짜 차단")
			return
		}
		response.InternalServerError(c, "날짜 차단 조회 실패")
		return
	}

	response.Success(c, dateBlock)
}

func (h *DateBlockHandler) UpdateDateBlock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 날짜 차단 ID")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	var req dto.UpdateDateBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	ctx := appContext.WithUserID(c.Request.Context(), userID)
	updated, err := h.service.UpdateDateBlock(ctx, uint(id), req)
	if err != nil {
		if errors.Is(err, services.ErrDateBlockNotFound) {
			response.NotFound(c, "존재하지 않는 날짜 차단")
			return
		}
		if errors.Is(err, services.ErrInvalidDateBlockRequest) {
			response.BadRequest(c, "잘못된 요청", err.Error())
			return
		}
		response.InternalServerError(c, "날짜 차단 수정 실패")
		return
	}

	response.Success(c, updated)
}

func (h *DateBlockHandler) GetDateBlockHistories(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 날짜 차단 ID")
		return
	}

	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "잘못된 요청 파라미터", err.Error())
		return
	}

	histories, total, err := h.historyService.GetDateBlockHistory(c.Request.Context(), uint(id), query.Page, query.Size)
	if err != nil {
		response.InternalServerError(c, "날짜 차단 이력 조회 실패")
		return
	}

	totalPages := int(total) / query.Size
	if int(total)%query.Size != 0 {
		totalPages++
	}

	pagination := &response.Pagination{
		Page:          query.Page,
		Size:          query.Size,
		TotalPages:    totalPages,
		TotalElements: total,
	}

	response.SuccessList(c, histories, pagination)
}
