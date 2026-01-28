package handlers

import (
	"context"

	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/utils"
	pkgutils "gitlab.bellsoft.net/rms/api-core/pkg/utils"
)

// toReservationResponse converts a Reservation model to ReservationResponse DTO
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

// toReservationResponseWithDetails converts a Reservation model with details to ReservationResponse DTO
func (h *ReservationHandler) toReservationResponseWithDetails(ctx context.Context, reservation *models.Reservation) dto.ReservationResponse {
	// 기본 응답 생성은 toReservationResponse와 동일
	return h.toReservationResponse(ctx, reservation)
}

// getUserSummary retrieves user summary information
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
