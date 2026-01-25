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
	"gitlab.bellsoft.net/rms/api-core/internal/utils"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
	pkgutils "gitlab.bellsoft.net/rms/api-core/pkg/utils"
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

func (h *RoomGroupHandler) toRoomGroupResponse(roomGroup *models.RoomGroup) dto.RoomGroupResponse {
	return dto.RoomGroupResponse{
		ID:           roomGroup.ID,
		Name:         roomGroup.Name,
		PeekPrice:    roomGroup.PeekPrice,
		OffPeekPrice: roomGroup.OffPeekPrice,
		Description:  roomGroup.Description,
		Rooms:        make([]dto.RoomLastStayDetailResponse, 0), // Initialize empty array
		CreatedAt:    dto.CustomTime{Time: roomGroup.CreatedAt},
		UpdatedAt:    dto.CustomTime{Time: roomGroup.UpdatedAt},
	}
}

func (h *RoomGroupHandler) toRoomGroupResponseWithRooms(ctx context.Context, roomGroup *models.RoomGroup) dto.RoomGroupResponse {
	resp := h.toRoomGroupResponse(roomGroup)

	// Initialize empty array
	resp.Rooms = make([]dto.RoomLastStayDetailResponse, 0)

	if len(roomGroup.Rooms) > 0 {
		resp.Rooms = make([]dto.RoomLastStayDetailResponse, len(roomGroup.Rooms))
		for i, room := range roomGroup.Rooms {
			// Create room response
			roomResponse := dto.RoomResponse{
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

			// Get last reservation for this room
			var lastReservation *dto.ReservationResponse
			if reservation, err := h.reservationService.GetLastReservationForRoom(ctx, room.ID); err == nil && reservation != nil {
				lastReservationVal := h.toReservationResponse(ctx, reservation)
				lastReservation = &lastReservationVal
			}

			resp.Rooms[i] = dto.RoomLastStayDetailResponse{
				Room:            roomResponse,
				LastReservation: lastReservation,
			}
		}
	}

	return resp
}

func (h *RoomGroupHandler) toRoomGroupResponseWithUsers(roomGroup *models.RoomGroup) dto.RoomGroupResponse {
	resp := dto.RoomGroupResponse{
		ID:           roomGroup.ID,
		Name:         roomGroup.Name,
		PeekPrice:    roomGroup.PeekPrice,
		OffPeekPrice: roomGroup.OffPeekPrice,
		Description:  roomGroup.Description,
		Rooms:        make([]dto.RoomLastStayDetailResponse, 0), // Initialize empty array
		CreatedAt:    dto.CustomTime{Time: roomGroup.CreatedAt},
		UpdatedAt:    dto.CustomTime{Time: roomGroup.UpdatedAt},
	}

	if roomGroup.CreatedByUser != nil {
		resp.CreatedBy = &dto.UserSummaryResponse{
			ID:              roomGroup.CreatedByUser.ID,
			Email:           utils.StringPtrToString(roomGroup.CreatedByUser.Email),
			Name:            roomGroup.CreatedByUser.Name,
			ProfileImageURL: pkgutils.GenerateGravatarURL(utils.StringPtrToString(roomGroup.CreatedByUser.Email)),
		}
	}

	if roomGroup.UpdatedByUser != nil {
		resp.UpdatedBy = &dto.UserSummaryResponse{
			ID:              roomGroup.UpdatedByUser.ID,
			Email:           utils.StringPtrToString(roomGroup.UpdatedByUser.Email),
			Name:            roomGroup.UpdatedByUser.Name,
			ProfileImageURL: pkgutils.GenerateGravatarURL(utils.StringPtrToString(roomGroup.UpdatedByUser.Email)),
		}
	}

	return resp
}

// getUserSummary is a helper function to get user summary
func (h *RoomGroupHandler) getUserSummary(ctx context.Context, userID uint) *dto.UserSummaryResponse {
	if userID == 0 {
		return nil
	}
	if user, err := h.userService.GetByID(ctx, userID); err == nil {
		profileImageURL := pkgutils.GenerateGravatarURL(utils.StringPtrToString(user.Email))

		return &dto.UserSummaryResponse{
			ID:              user.ID,
			UserID:          user.UserID,
			Email:           utils.StringPtrToString(user.Email),
			Name:            user.Name,
			ProfileImageURL: profileImageURL,
		}
	}
	// 사용자를 찾을 수 없는 경우 ID만 반환
	return &dto.UserSummaryResponse{
		ID: userID,
	}
}

// toReservationResponse converts a reservation model to a response DTO
func (h *RoomGroupHandler) toReservationResponse(ctx context.Context, reservation *models.Reservation) dto.ReservationResponse {
	resp := dto.ReservationResponse{
		ID:              reservation.ID,
		PaymentMethodID: reservation.PaymentMethodID,
		Name:            reservation.Name,
		Phone:           reservation.Phone,
		PeopleCount:     reservation.PeopleCount,
		StayStartAt:     dto.JSONDate{Time: reservation.StayStartAt},
		StayEndAt:       dto.JSONDate{Time: reservation.StayEndAt},
		Price:           reservation.Price,
		Deposit:         reservation.Deposit,
		PaymentAmount:   reservation.PaymentAmount,
		RefundAmount:    reservation.RefundAmount,
		BrokerFee:       reservation.BrokerFee,
		Note:            reservation.Note,
		Status:          reservation.Status.String(),
		Type:            reservation.Type.String(),
		CreatedAt:       dto.CustomTime{Time: reservation.CreatedAt},
		UpdatedAt:       dto.CustomTime{Time: reservation.UpdatedAt},
		Rooms:           []dto.RoomResponse{}, // 빈 배열로 초기화
	}

	// CheckInAt, CheckOutAt, CanceledAt 설정
	if reservation.CheckInAt != nil {
		resp.CheckInAt = &dto.CustomTime{Time: *reservation.CheckInAt}
	}
	if reservation.CheckOutAt != nil {
		resp.CheckOutAt = &dto.CustomTime{Time: *reservation.CheckOutAt}
	}
	if reservation.CanceledAt != nil {
		resp.CanceledAt = &dto.CustomTime{Time: *reservation.CanceledAt}
	}

	// CreatedBy, UpdatedBy 설정
	resp.CreatedBy = h.getUserSummary(ctx, reservation.CreatedBy)
	resp.UpdatedBy = h.getUserSummary(ctx, reservation.UpdatedBy)

	if reservation.PaymentMethod != nil {
		resp.PaymentMethod = &dto.PaymentMethodResponse{
			ID:                       reservation.PaymentMethod.ID,
			Name:                     reservation.PaymentMethod.Name,
			CommissionRate:           reservation.PaymentMethod.CommissionRate,
			RequireUnpaidAmountCheck: bool(reservation.PaymentMethod.RequireUnpaidAmountCheck),
			IsDefaultSelect:          bool(reservation.PaymentMethod.IsDefaultSelect),
			Status:                   reservation.PaymentMethod.Status.String(),
			CreatedAt:                dto.CustomTime{Time: reservation.PaymentMethod.CreatedAt},
			UpdatedAt:                dto.CustomTime{Time: reservation.PaymentMethod.UpdatedAt},
		}
	}

	// rooms 데이터가 있는 경우 설정 - Spring Boot 호환성을 위해 직접 RoomResponse 배열로 변환
	if len(reservation.Rooms) > 0 {
		resp.Rooms = make([]dto.RoomResponse, len(reservation.Rooms))
		for i, rr := range reservation.Rooms {
			if rr.Room != nil {
				resp.Rooms[i] = dto.RoomResponse{
					ID:          rr.Room.ID,
					Number:      rr.Room.Number,
					RoomGroupID: rr.Room.RoomGroupID,
					Note:        rr.Room.Note,
					Status:      rr.Room.Status.String(),
					CreatedAt:   dto.CustomTime{Time: rr.Room.CreatedAt},
					UpdatedAt:   dto.CustomTime{Time: rr.Room.UpdatedAt},
					CreatedBy:   h.getUserSummary(ctx, rr.Room.CreatedBy),
					UpdatedBy:   h.getUserSummary(ctx, rr.Room.UpdatedBy),
				}

				if rr.Room.RoomGroup != nil {
					resp.Rooms[i].RoomGroup = &dto.RoomGroupResponse{
						ID:           rr.Room.RoomGroup.ID,
						Name:         rr.Room.RoomGroup.Name,
						PeekPrice:    rr.Room.RoomGroup.PeekPrice,
						OffPeekPrice: rr.Room.RoomGroup.OffPeekPrice,
						Description:  rr.Room.RoomGroup.Description,
						CreatedAt:    dto.CustomTime{Time: rr.Room.RoomGroup.CreatedAt},
						UpdatedAt:    dto.CustomTime{Time: rr.Room.RoomGroup.UpdatedAt},
					}
				}
			}
		}
	}

	return resp
}
