package handlers

import (
	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	appContext "gitlab.bellsoft.net/rms/api-core/internal/context"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/internal/utils"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
	pkgutils "gitlab.bellsoft.net/rms/api-core/pkg/utils"
)

type ReservationHandler struct {
	reservationService services.ReservationService
	userService        services.UserService
	historyService     services.HistoryService
}

func NewReservationHandler(reservationService services.ReservationService, userService services.UserService, historyService services.HistoryService) *ReservationHandler {
	return &ReservationHandler{
		reservationService: reservationService,
		userService:        userService,
		historyService:     historyService,
	}
}

func (h *ReservationHandler) ListReservations(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "잘못된 쿼리 파라미터", err.Error())
		return
	}

	var filterQuery dto.ReservationFilter
	if err := c.ShouldBindQuery(&filterQuery); err != nil {
		response.BadRequest(c, "잘못된 필터 파라미터", err.Error())
		return
	}

	filter := dto.ReservationRepositoryFilter{
		RoomID:    filterQuery.RoomID,
		StartDate: filterQuery.StayStartAt,
		EndDate:   filterQuery.StayEndAt,
		Search:    filterQuery.Search,
	}

	if filterQuery.Status != nil {
		switch *filterQuery.Status {
		case "REFUND":
			s := models.ReservationStatusRefund
			filter.Status = &s
		case "CANCEL":
			s := models.ReservationStatusCancel
			filter.Status = &s
		case "PENDING":
			s := models.ReservationStatusPending
			filter.Status = &s
		case "NORMAL":
			s := models.ReservationStatusNormal
			filter.Status = &s
		}
	}

	if filterQuery.Type != nil {
		switch *filterQuery.Type {
		case "STAY":
			t := models.ReservationTypeStay
			filter.Type = &t
		case "MONTHLY_RENT":
			t := models.ReservationTypeMonthlyRent
			filter.Type = &t
		}
	}

	reservations, total, err := h.reservationService.GetAll(c.Request.Context(), filter, query.Page, query.Size, query.Sort)
	if err != nil {
		response.InternalServerError(c, "예약 목록 조회 실패")
		return
	}

	reservationResponses := make([]dto.ReservationResponse, len(reservations))
	for i, reservation := range reservations {
		reservationResponses[i] = h.toReservationResponse(c.Request.Context(), &reservation)
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

	// Filter response 생성
	filterResponse := dto.ReservationFilterResponse{
		Status: filterQuery.Status,
		Type:   filterQuery.Type,
		RoomID: filterQuery.RoomID,
		Search: filterQuery.Search,
	}

	// 날짜 필터 변환
	if filterQuery.StayStartAt != nil {
		filterResponse.StayStartAt = &dto.JSONDate{Time: *filterQuery.StayStartAt}
	}
	if filterQuery.StayEndAt != nil {
		filterResponse.StayEndAt = &dto.JSONDate{Time: *filterQuery.StayEndAt}
	}

	response.SuccessListWithFilter(c, reservationResponses, pagination, filterResponse)
}

func (h *ReservationHandler) GetReservation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 예약 ID")
		return
	}

	reservation, err := h.reservationService.GetByIDWithDetails(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, services.ErrReservationNotFound) {
			response.NotFound(c, "존재하지 않는 예약")
			return
		}
		response.InternalServerError(c, "예약 조회 실패")
		return
	}

	reservationResponse := h.toReservationResponseWithDetails(c.Request.Context(), reservation)
	response.Success(c, reservationResponse)
}

func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	var req dto.CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	// 프론트엔드 호환성: paymentMethod 객체 또는 paymentMethodId 모두 지원
	paymentMethodID := req.GetPaymentMethodID()
	if paymentMethodID == 0 {
		response.BadRequest(c, "결제 수단 ID 필수")
		return
	}

	reservation := &models.Reservation{
		PaymentMethodID: paymentMethodID,
		Name:            req.Name,
		Phone:           req.Phone,
		PeopleCount:     req.PeopleCount,
		StayStartAt:     req.StayStartAt.Time,
		StayEndAt:       req.StayEndAt.Time,
		Price:           req.Price,
		Deposit:         req.Deposit,
		PaymentAmount:   req.PaymentAmount,
		BrokerFee:       req.BrokerFee,
		Note:            req.Note,
	}

	if req.Type != "" {
		switch req.Type {
		case "STAY":
			reservation.Type = models.ReservationTypeStay
		case "MONTHLY_RENT":
			reservation.Type = models.ReservationTypeMonthlyRent
		}
	} else {
		reservation.Type = models.ReservationTypeStay
	}

	// Status 처리 (프론트엔드가 보낸 값 사용)
	if req.Status != "" {
		switch req.Status {
		case "NORMAL":
			reservation.Status = models.ReservationStatusNormal
		case "PENDING":
			reservation.Status = models.ReservationStatusPending
		case "CANCEL":
			reservation.Status = models.ReservationStatusCancel
		case "REFUND":
			reservation.Status = models.ReservationStatusRefund
		default:
			reservation.Status = models.ReservationStatusPending
		}
	} else {
		reservation.Status = models.ReservationStatusPending
	}

	ctx := appContext.WithUserID(c.Request.Context(), userID)
	if err := h.reservationService.Create(ctx, reservation, req.GetRoomIDs()); err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidDateRange):
			response.BadRequest(c, "잘못된 날짜 범위")
		case errors.Is(err, services.ErrPaymentMethodNotFound):
			response.BadRequest(c, "존재하지 않는 결제 수단")
		case errors.Is(err, services.ErrPaymentMethodInactive):
			response.BadRequest(c, "비활성화된 결제 수단")
		case errors.Is(err, services.ErrRoomNotAvailable):
			response.BadRequest(c, "선택한 날짜에 사용할 수 없는 객실이 있습니다")
		default:
			response.InternalServerError(c, "예약 등록 실패")
		}
		return
	}

	// Get the created reservation with all details for the response
	createdReservation, err := h.reservationService.GetByIDWithDetails(ctx, reservation.ID)
	if err != nil {
		response.InternalServerError(c, "생성된 예약 상세 조회 실패")
		return
	}

	reservationResponse := h.toReservationResponse(ctx, createdReservation)
	response.Created(c, reservationResponse)
}

func (h *ReservationHandler) UpdateReservation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 예약 ID")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	var req dto.UpdateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if paymentMethodID := req.GetPaymentMethodID(); paymentMethodID != nil {
		updates["paymentMethodId"] = *paymentMethodID
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.PeopleCount != nil {
		updates["peopleCount"] = *req.PeopleCount
	}
	if req.StayStartAt != nil {
		updates["stayStartAt"] = *req.StayStartAt
	}
	if req.StayEndAt != nil {
		updates["stayEndAt"] = *req.StayEndAt
	}
	if req.CheckInAt != nil {
		updates["checkInAt"] = req.CheckInAt
	}
	if req.CheckOutAt != nil {
		updates["checkOutAt"] = req.CheckOutAt
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Deposit != nil {
		updates["deposit"] = *req.Deposit
	}
	if req.PaymentAmount != nil {
		updates["paymentAmount"] = *req.PaymentAmount
	}
	if req.RefundAmount != nil {
		updates["refundAmount"] = *req.RefundAmount
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}
	if req.Status != nil {
		switch *req.Status {
		case "REFUND":
			updates["status"] = models.ReservationStatusRefund
		case "CANCEL":
			updates["status"] = models.ReservationStatusCancel
		case "PENDING":
			updates["status"] = models.ReservationStatusPending
		case "NORMAL":
			updates["status"] = models.ReservationStatusNormal
		}
	}
	if req.Type != nil {
		switch *req.Type {
		case "STAY":
			updates["type"] = models.ReservationTypeStay
		case "MONTHLY_RENT":
			updates["type"] = models.ReservationTypeMonthlyRent
		}
	}

	// Pass user ID in context
	ctx := appContext.WithUserID(c.Request.Context(), userID)
	roomIDs := req.GetRoomIDs()
	hasRoomsUpdate := req.HasRoomsUpdate()
	reservation, err := h.reservationService.Update(ctx, uint(id), updates, roomIDs, hasRoomsUpdate)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrReservationNotFound):
			response.NotFound(c, "존재하지 않는 예약")
		case errors.Is(err, services.ErrInvalidDateRange):
			response.BadRequest(c, "잘못된 날짜 범위")
		case errors.Is(err, services.ErrPaymentMethodNotFound):
			response.BadRequest(c, "존재하지 않는 결제 수단")
		case errors.Is(err, services.ErrPaymentMethodInactive):
			response.BadRequest(c, "비활성화된 결제 수단")
		case errors.Is(err, services.ErrRoomNotAvailable):
			response.BadRequest(c, "선택한 날짜에 사용할 수 없는 객실이 있습니다")
		default:
			response.InternalServerError(c, "예약 수정 실패")
		}
		return
	}

	reservationResponse := h.toReservationResponseWithDetails(c.Request.Context(), reservation)
	response.Success(c, reservationResponse)
}

func (h *ReservationHandler) DeleteReservation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 예약 ID")
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	// Pass user ID in context
	ctx := appContext.WithUserID(c.Request.Context(), userID)
	if err := h.reservationService.Delete(ctx, uint(id)); err != nil {
		if errors.Is(err, services.ErrReservationNotFound) {
			response.NotFound(c, "존재하지 않는 예약")
			return
		}
		response.InternalServerError(c, "예약 삭제 실패")
		return
	}

	response.NoContent(c)
}

func (h *ReservationHandler) GetReservationStatistics(c *gin.Context) {
	var query dto.ReservationStatisticsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "잘못된 쿼리 파라미터", err.Error())
		return
	}

	// Set default periodType if not provided
	if query.PeriodType == "" {
		query.PeriodType = "MONTHLY"
	}

	// Validate date range
	if query.StartDate.After(query.EndDate) {
		response.BadRequest(c, "잘못된 날짜 범위", "시작일은 종료일보다 이전이거나 같아야 합니다")
		return
	}

	statistics, err := h.reservationService.GetStatistics(c.Request.Context(), query.StartDate, query.EndDate, query.PeriodType)
	if err != nil {
		response.InternalServerError(c, "예약 통계 조회 실패")
		return
	}

	// Convert to Spring Boot compatible format
	statsData := make([]dto.StatisticsData, len(statistics))
	monthlyStats := make([]dto.MonthlyStats, 0)

	for i, stat := range statistics {
		statsData[i] = dto.StatisticsData{
			Period:            stat.Period,
			TotalSales:        int(stat.TotalRevenue),
			TotalReservations: int(stat.ReservationCount),
			TotalGuests:       int(stat.TotalGuests),
		}

		// Add to monthlyStats for backward compatibility
		if query.PeriodType == "MONTHLY" || query.PeriodType == "" {
			monthlyStats = append(monthlyStats, dto.MonthlyStats{
				YearMonth:         stat.Period,
				TotalSales:        int(stat.TotalRevenue),
				TotalReservations: int(stat.ReservationCount),
				TotalGuests:       int(stat.TotalGuests),
			})
		}
	}

	responseData := dto.ReservationStatisticsResponse{
		PeriodType:   query.PeriodType,
		Stats:        statsData,
		MonthlyStats: monthlyStats,
	}

	if responseData.PeriodType == "" {
		responseData.PeriodType = "MONTHLY"
	}

	response.Success(c, responseData)
}

func (h *ReservationHandler) GetReservationHistories(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 예약 ID")
		return
	}

	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "잘못된 쿼리 파라미터", err.Error())
		return
	}

	if _, err := h.reservationService.GetByIDWithDetails(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "존재하지 않는 예약")
		return
	}

	histories, total, err := h.historyService.GetReservationHistory(c.Request.Context(), uint(id), query.Page, query.Size)
	if err != nil {
		response.InternalServerError(c, "예약 이력 조회 실패")
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

func (h *ReservationHandler) toReservationResponse(ctx context.Context, reservation *models.Reservation) dto.ReservationResponse {
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

func (h *ReservationHandler) toReservationResponseWithDetails(ctx context.Context, reservation *models.Reservation) dto.ReservationResponse {
	// 기본 응답 생성은 toReservationResponse와 동일
	return h.toReservationResponse(ctx, reservation)
}

func (h *ReservationHandler) getUserSummary(ctx context.Context, userID uint) *dto.UserSummaryResponse {
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
	return &dto.UserSummaryResponse{
		ID: userID,
	}
}
