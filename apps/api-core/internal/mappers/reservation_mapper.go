package mappers

import (
	"context"

	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
)

// ToReservationResponse converts a Reservation model to ReservationResponse DTO
func ToReservationResponse(ctx context.Context, reservation *models.Reservation, getUserSummary GetUserSummaryFunc) dto.ReservationResponse {
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
		Rooms:           []dto.RoomResponse{},
	}

	if reservation.CheckInAt != nil {
		resp.CheckInAt = &dto.CustomTime{Time: *reservation.CheckInAt}
	}
	if reservation.CheckOutAt != nil {
		resp.CheckOutAt = &dto.CustomTime{Time: *reservation.CheckOutAt}
	}
	if reservation.CanceledAt != nil {
		resp.CanceledAt = &dto.CustomTime{Time: *reservation.CanceledAt}
	}

	resp.CreatedBy = getUserSummary(ctx, reservation.CreatedBy)
	resp.UpdatedBy = getUserSummary(ctx, reservation.UpdatedBy)

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
					CreatedBy:   getUserSummary(ctx, rr.Room.CreatedBy),
					UpdatedBy:   getUserSummary(ctx, rr.Room.UpdatedBy),
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
