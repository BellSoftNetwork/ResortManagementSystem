package handlers

import (
	"context"
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
	"gitlab.bellsoft.net/rms/api-core/pkg/utils"
)

type RoomHandler struct {
	roomService    services.RoomService
	userService    services.UserService
	historyService services.HistoryService
}

func NewRoomHandler(roomService services.RoomService, userService services.UserService, historyService services.HistoryService) *RoomHandler {
	return &RoomHandler{
		roomService:    roomService,
		userService:    userService,
		historyService: historyService,
	}
}

func (h *RoomHandler) ListRooms(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "잘못된 요청 파라미터", err.Error())
		return
	}

	var filterQuery dto.RoomFilter
	if err := c.ShouldBindQuery(&filterQuery); err != nil {
		response.BadRequest(c, "잘못된 필터 파라미터", err.Error())
		return
	}

	filter := repositories.RoomFilter{
		RoomGroupID: filterQuery.RoomGroupID,
		Search:      filterQuery.Search,
	}

	if filterQuery.Status != nil {
		switch *filterQuery.Status {
		case "DAMAGED":
			s := models.RoomStatusDamaged
			filter.Status = &s
		case "CONSTRUCTION":
			s := models.RoomStatusConstruction
			filter.Status = &s
		case "INACTIVE":
			s := models.RoomStatusInactive
			filter.Status = &s
		case "NORMAL":
			s := models.RoomStatusNormal
			filter.Status = &s
		}
	}

	var rooms []models.Room
	var total int64
	var err error

	if filterQuery.StayStartAt != nil && filterQuery.StayEndAt != nil {
		startDate, parseErr := time.Parse("2006-01-02", *filterQuery.StayStartAt)
		if parseErr != nil {
			response.BadRequest(c, "잘못된 시작 날짜 형식")
			return
		}
		endDate, parseErr := time.Parse("2006-01-02", *filterQuery.StayEndAt)
		if parseErr != nil {
			response.BadRequest(c, "잘못된 종료 날짜 형식")
			return
		}

		rooms, err = h.roomService.GetAvailableRooms(c.Request.Context(), startDate, endDate, filterQuery.ExcludeReservationID)
		if err != nil {
			response.InternalServerError(c, "객실 목록 조회 실패")
			return
		}
		total = int64(len(rooms))

		offset := query.Page * query.Size
		end := offset + query.Size
		if offset > len(rooms) {
			rooms = []models.Room{}
		} else if end > len(rooms) {
			rooms = rooms[offset:]
		} else {
			rooms = rooms[offset:end]
		}
	} else {
		rooms, total, err = h.roomService.GetAll(c.Request.Context(), filter, query.Page, query.Size, query.Sort)
		if err != nil {
			response.InternalServerError(c, "객실 목록 조회 실패")
			return
		}
	}

	roomResponses := make([]dto.RoomResponse, len(rooms))
	for i, room := range rooms {
		roomResponses[i] = h.toRoomResponse(c.Request.Context(), &room)
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

	filterResponse := dto.RoomFilterResponse{
		RoomGroupID:          filterQuery.RoomGroupID,
		Status:               filterQuery.Status,
		Search:               filterQuery.Search,
		StayStartAt:          filterQuery.StayStartAt,
		StayEndAt:            filterQuery.StayEndAt,
		ExcludeReservationID: filterQuery.ExcludeReservationID,
	}

	response.SuccessListWithFilter(c, roomResponses, pagination, filterResponse)
}

func (h *RoomHandler) GetRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 객실 ID")
		return
	}

	room, err := h.roomService.GetByIDWithGroup(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, services.ErrRoomNotFound) {
			response.NotFound(c, "존재하지 않는 객실")
			return
		}
		response.InternalServerError(c, "객실 조회 실패")
		return
	}

	roomResponse := h.toRoomResponse(c.Request.Context(), room)
	response.Success(c, roomResponse)
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req dto.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	if req.RoomGroup != nil && req.RoomGroup.ID != 0 {
		req.RoomGroupID = req.RoomGroup.ID
	}

	if req.RoomGroupID == 0 {
		response.BadRequest(c, "객실 그룹 ID 필수")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	room := &models.Room{
		Number:      req.Number,
		RoomGroupID: req.RoomGroupID,
		Note:        req.Note,
	}

	if req.Status != "" {
		switch req.Status {
		case "DAMAGED":
			room.Status = models.RoomStatusDamaged
		case "CONSTRUCTION":
			room.Status = models.RoomStatusConstruction
		case "INACTIVE":
			room.Status = models.RoomStatusInactive
		case "NORMAL":
			room.Status = models.RoomStatusNormal
		}
	} else {
		room.Status = models.RoomStatusInactive
	}

	ctx := appContext.WithUserID(c.Request.Context(), userID)
	if err := h.roomService.Create(ctx, room); err != nil {
		if errors.Is(err, services.ErrRoomNumberExists) {
			response.Conflict(c, "이미 존재하는 객실")
			return
		}
		if errors.Is(err, services.ErrRoomGroupNotFound) {
			response.BadRequest(c, "존재하지 않는 객실 그룹")
			return
		}
		response.InternalServerError(c, "객실 등록 실패")
		return
	}

	roomResponse := h.toRoomResponse(ctx, room)
	response.Created(c, roomResponse)
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 객실 ID")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	var req dto.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Number != nil {
		updates["number"] = *req.Number
	}
	if req.RoomGroupID != nil {
		updates["room_group_id"] = *req.RoomGroupID
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}
	if req.Status != nil {
		switch *req.Status {
		case "DAMAGED":
			updates["status"] = models.RoomStatusDamaged
		case "CONSTRUCTION":
			updates["status"] = models.RoomStatusConstruction
		case "INACTIVE":
			updates["status"] = models.RoomStatusInactive
		case "NORMAL":
			updates["status"] = models.RoomStatusNormal
		}
	}

	ctx := appContext.WithUserID(c.Request.Context(), userID)
	room, err := h.roomService.Update(ctx, uint(id), updates)
	if err != nil {
		if errors.Is(err, services.ErrRoomNotFound) {
			response.NotFound(c, "존재하지 않는 객실")
			return
		}
		if errors.Is(err, services.ErrRoomNumberExists) {
			response.Conflict(c, "이미 존재하는 객실")
			return
		}
		if errors.Is(err, services.ErrRoomGroupNotFound) {
			response.BadRequest(c, "존재하지 않는 객실 그룹")
			return
		}
		response.InternalServerError(c, "객실 수정 실패")
		return
	}

	roomResponse := h.toRoomResponse(ctx, room)
	response.Success(c, roomResponse)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 객실 ID")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	ctx := appContext.WithUserID(c.Request.Context(), userID)
	if err := h.roomService.Delete(ctx, uint(id)); err != nil {
		if errors.Is(err, services.ErrRoomNotFound) {
			response.NotFound(c, "존재하지 않는 객실")
			return
		}
		if errors.Is(err, services.ErrRoomHasReservations) {
			response.Conflict(c, "예약이 존재하는 객실은 삭제할 수 없습니다")
			return
		}
		response.InternalServerError(c, "객실 삭제 실패")
		return
	}

	response.NoContent(c)
}

func (h *RoomHandler) GetRoomHistories(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 객실 ID")
		return
	}

	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "잘못된 요청 파라미터", err.Error())
		return
	}

	if _, err := h.roomService.GetByIDWithGroup(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "존재하지 않는 객실")
		return
	}

	histories, total, err := h.historyService.GetRoomHistory(c.Request.Context(), uint(id), query.Page, query.Size)
	if err != nil {
		response.InternalServerError(c, "객실 이력 조회 실패")
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

func (h *RoomHandler) toRoomResponse(ctx context.Context, room *models.Room) dto.RoomResponse {
	resp := dto.RoomResponse{
		ID:          room.ID,
		Number:      room.Number,
		RoomGroupID: room.RoomGroupID,
		Note:        room.Note,
		Status:      room.Status.String(),
		CreatedAt:   dto.CustomTime{Time: room.CreatedAt},
		UpdatedAt:   dto.CustomTime{Time: room.UpdatedAt},
		CreatedBy:   h.getUserSummary(ctx, room.CreatedBy),
		UpdatedBy:   h.getUserSummary(ctx, room.UpdatedBy),
	}

	if room.RoomGroup != nil {
		resp.RoomGroup = &dto.RoomGroupResponse{
			ID:           room.RoomGroup.ID,
			Name:         room.RoomGroup.Name,
			PeekPrice:    room.RoomGroup.PeekPrice,
			OffPeekPrice: room.RoomGroup.OffPeekPrice,
			Description:  room.RoomGroup.Description,
			CreatedAt:    dto.CustomTime{Time: room.RoomGroup.CreatedAt},
			UpdatedAt:    dto.CustomTime{Time: room.RoomGroup.UpdatedAt},
		}
	}

	return resp
}

func (h *RoomHandler) getUserSummary(ctx context.Context, userID uint) *dto.UserSummaryResponse {
	if userID == 0 {
		return nil
	}
	if user, err := h.userService.GetByID(ctx, userID); err == nil {
		profileImageURL := utils.GenerateGravatarURL(user.Email)

		return &dto.UserSummaryResponse{
			ID:              user.ID,
			UserID:          user.UserID,
			Email:           user.Email,
			Name:            user.Name,
			ProfileImageURL: profileImageURL,
		}
	}
	return &dto.UserSummaryResponse{
		ID: userID,
	}
}
