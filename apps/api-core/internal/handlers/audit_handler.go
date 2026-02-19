package handlers

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.bellsoft.net/rms/api-core/internal/audit"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
)

type AuditHandler struct {
	auditService audit.AuditService
}

func NewAuditHandler(auditService audit.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

func (h *AuditHandler) ListAuditLogs(c *gin.Context) {
	var query dto.AuditLogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "Invalid query parameters", err.Error())
		return
	}

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		response.BadRequest(c, "Invalid pagination", err.Error())
		return
	}

	filter := audit.AuditLogFilter{
		EntityType: query.EntityType,
		Action:     query.Action,
		UserID:     query.UserID,
		EntityID:   query.EntityID,
	}

	if query.StartDate != "" {
		t, err := time.Parse("2006-01-02", query.StartDate)
		if err == nil {
			filter.StartDate = &t
		}
	}
	if query.EndDate != "" {
		t, err := time.Parse("2006-01-02", query.EndDate)
		if err == nil {
			t = t.Add(24 * time.Hour)
			filter.EndDate = &t
		}
	}

	if filter.StartDate == nil && filter.EndDate == nil {
		sevenDaysAgo := time.Now().AddDate(0, 0, -7)
		filter.StartDate = &sevenDaysAgo
	}

	logs, total, err := h.auditService.GetAllHistory(c.Request.Context(), filter, pagination.Page, pagination.Size)
	if err != nil {
		response.InternalServerError(c, "Failed to get audit logs")
		return
	}

	items := make([]dto.AuditLogListResponse, 0)
	for _, log := range logs {
		changedFields := make([]string, 0)
		if log.ChangedFields != nil {
			json.Unmarshal(log.ChangedFields, &changedFields)
		}
		items = append(items, dto.AuditLogListResponse{
			ID:            log.ID,
			EntityType:    log.EntityType,
			EntityID:      log.EntityID,
			Action:        string(log.Action),
			ChangedFields: changedFields,
			UserID:        log.UserID,
			Username:      log.Username,
			CreatedAt:     log.CreatedAt,
		})
	}

	totalPages := int(total) / pagination.Size
	if int(total)%pagination.Size != 0 {
		totalPages++
	}

	paginationResp := &response.Pagination{
		Page:          pagination.Page,
		Size:          pagination.Size,
		TotalPages:    totalPages,
		TotalElements: total,
	}

	response.SuccessList(c, items, paginationResp)
}

func (h *AuditHandler) GetAuditLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid ID", err.Error())
		return
	}

	log, err := h.auditService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Audit log not found")
		return
	}

	response.Success(c, dto.AuditLogResponse{
		ID:            log.ID,
		EntityType:    log.EntityType,
		EntityID:      log.EntityID,
		Action:        string(log.Action),
		OldValues:     log.OldValues,
		NewValues:     log.NewValues,
		ChangedFields: log.ChangedFields,
		UserID:        log.UserID,
		Username:      log.Username,
		CreatedAt:     log.CreatedAt,
	})
}
