package handlers

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	appContext "gitlab.bellsoft.net/rms/api-core/internal/context"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
)

type RoomGroupHandler struct {
	roomGroupService   services.RoomGroupService
	reservationService services.ReservationService
	userService        services.UserService
}

func NewRoomGroupHandler(roomGroupService services.RoomGroupService, reservationService services.ReservationService, userService services.UserService) *RoomGroupHandler {
	return &RoomGroupHandler{
		roomGroupService:   roomGroupService,
		reservationService: reservationService,
		userService:        userService,
	}
}

func (h *RoomGroupHandler) ListRoomGroups(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "잘못된 쿼리 파라미터", err.Error())
		return
	}

	roomGroups, total, err := h.roomGroupService.GetAllWithUsers(c.Request.Context(), query.Page, query.Size, query.Sort)
	if err != nil {
		response.InternalServerError(c, "객실 그룹 목록 조회 실패")
		return
	}

	roomGroupResponses := make([]dto.RoomGroupResponse, len(roomGroups))
	for i, roomGroup := range roomGroups {
		roomGroupResponses[i] = h.toRoomGroupResponseWithUsers(&roomGroup)
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

	// Room Group 리스트는 현재 필터를 지원하지 않으므로 빈 필터 객체 반환
	filterResponse := map[string]interface{}{}

	response.SuccessListWithFilter(c, roomGroupResponses, pagination, filterResponse)
}

func (h *RoomGroupHandler) GetRoomGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 객실 그룹 ID")
		return
	}

	var filter dto.RoomGroupRoomFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, "잘못된 쿼리 파라미터", err.Error())
		return
	}

	var roomStatusFilter *models.RoomStatus
	if filter.Status != nil && *filter.Status != "" {
		switch *filter.Status {
		case "DAMAGED":
			s := models.RoomStatusDamaged
			roomStatusFilter = &s
		case "CONSTRUCTION":
			s := models.RoomStatusConstruction
			roomStatusFilter = &s
		case "INACTIVE":
			s := models.RoomStatusInactive
			roomStatusFilter = &s
		case "NORMAL":
			s := models.RoomStatusNormal
			roomStatusFilter = &s
		}
	}

	var repoFilter repositories.RoomGroupRoomFilter
	repoFilter.RoomStatus = roomStatusFilter
	repoFilter.ExcludeReservationID = filter.ExcludeReservationID

	if filter.StayStartAt != nil && *filter.StayStartAt != "" {
		if t, err := time.Parse("2006-01-02", *filter.StayStartAt); err == nil {
			repoFilter.StayStartAt = &t
		}
	}
	if filter.StayEndAt != nil && *filter.StayEndAt != "" {
		if t, err := time.Parse("2006-01-02", *filter.StayEndAt); err == nil {
			repoFilter.StayEndAt = &t
		}
	}

	roomGroup, err := h.roomGroupService.GetByIDWithFilteredRooms(c.Request.Context(), uint(id), repoFilter)
	if err != nil {
		if errors.Is(err, services.ErrRoomGroupNotFound) {
			response.NotFound(c, "존재하지 않는 객실 그룹")
			return
		}
		response.InternalServerError(c, "객실 그룹 조회 실패")
		return
	}

	roomGroupResponse := h.toRoomGroupResponseWithRooms(c.Request.Context(), roomGroup)
	response.Success(c, roomGroupResponse)
}

func (h *RoomGroupHandler) CreateRoomGroup(c *gin.Context) {
	var req dto.CreateRoomGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	roomGroup := &models.RoomGroup{
		Name:         req.Name,
		PeekPrice:    req.PeekPrice,
		OffPeekPrice: req.OffPeekPrice,
		Description:  req.Description,
	}

	// Pass user ID in context
	ctx := appContext.WithUserID(c.Request.Context(), userID)
	if err := h.roomGroupService.Create(ctx, roomGroup); err != nil {
		if errors.Is(err, services.ErrRoomGroupNameExists) {
			response.Conflict(c, "이미 존재하는 객실 그룹")
			return
		}
		response.InternalServerError(c, "객실 그룹 등록 실패")
		return
	}

	// Reload with user information
	roomGroupWithUsers, err := h.roomGroupService.GetByIDWithUsers(ctx, roomGroup.ID)
	if err != nil {
		response.InternalServerError(c, "생성된 객실 그룹 조회 실패")
		return
	}

	roomGroupResponse := h.toRoomGroupResponseWithUsers(roomGroupWithUsers)
	response.Created(c, roomGroupResponse)
}

func (h *RoomGroupHandler) UpdateRoomGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 객실 그룹 ID")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	var req dto.UpdateRoomGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.PeekPrice != nil {
		updates["peekPrice"] = *req.PeekPrice
	}
	if req.OffPeekPrice != nil {
		updates["offPeekPrice"] = *req.OffPeekPrice
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	// Pass user ID in context
	ctx := appContext.WithUserID(c.Request.Context(), userID)
	roomGroup, err := h.roomGroupService.Update(ctx, uint(id), updates)
	if err != nil {
		if errors.Is(err, services.ErrRoomGroupNotFound) {
			response.NotFound(c, "존재하지 않는 객실 그룹")
			return
		}
		if errors.Is(err, services.ErrRoomGroupNameExists) {
			response.Conflict(c, "이미 존재하는 객실 그룹")
			return
		}
		response.InternalServerError(c, "객실 그룹 수정 실패")
		return
	}

	// Reload with user information
	roomGroupWithUsers, err := h.roomGroupService.GetByIDWithUsers(ctx, roomGroup.ID)
	if err != nil {
		response.InternalServerError(c, "수정된 객실 그룹 조회 실패")
		return
	}

	roomGroupResponse := h.toRoomGroupResponseWithUsers(roomGroupWithUsers)
	response.Success(c, roomGroupResponse)
}

func (h *RoomGroupHandler) DeleteRoomGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 객실 그룹 ID")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	// Pass user ID in context
	ctx := appContext.WithUserID(c.Request.Context(), userID)
	if err := h.roomGroupService.Delete(ctx, uint(id)); err != nil {
		if errors.Is(err, services.ErrRoomGroupNotFound) {
			response.NotFound(c, "존재하지 않는 객실 그룹")
			return
		}
		if errors.Is(err, services.ErrRoomGroupHasRooms) {
			response.Conflict(c, "객실이 존재하는 객실 그룹은 삭제할 수 없습니다")
			return
		}
		response.InternalServerError(c, "객실 그룹 삭제 실패")
		return
	}

	response.NoContent(c)
}
