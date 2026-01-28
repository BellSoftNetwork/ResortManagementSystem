package handlers

import (
	"context"

	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/utils"
	pkgutils "gitlab.bellsoft.net/rms/api-core/pkg/utils"
)

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
