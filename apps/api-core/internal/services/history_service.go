package services

import (
	"context"
	"encoding/json"
	"time"

	"gitlab.bellsoft.net/rms/api-core/internal/audit"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
)

type HistoryService interface {
	GetRoomHistory(ctx context.Context, roomID uint, page, size int) ([]dto.RoomRevisionResponse, int64, error)
	GetReservationHistory(ctx context.Context, reservationID uint, page, size int) ([]dto.ReservationRevisionResponse, int64, error)
}

type historyService struct {
	auditService audit.AuditService
	userService  UserService
}

func NewHistoryService(auditService audit.AuditService, userService UserService) HistoryService {
	return &historyService{
		auditService: auditService,
		userService:  userService,
	}
}

func (s *historyService) GetRoomHistory(ctx context.Context, roomID uint, page, size int) ([]dto.RoomRevisionResponse, int64, error) {
	logs, total, err := s.auditService.GetHistory(ctx, "room", roomID, page, size)
	if err != nil {
		return nil, 0, err
	}

	revisions := make([]dto.RoomRevisionResponse, len(logs))
	for i, log := range logs {
		revisions[i] = s.convertToRoomRevision(ctx, &log)
	}

	return revisions, total, nil
}

func (s *historyService) GetReservationHistory(ctx context.Context, reservationID uint, page, size int) ([]dto.ReservationRevisionResponse, int64, error) {
	logs, total, err := s.auditService.GetHistory(ctx, "reservation", reservationID, page, size)
	if err != nil {
		return nil, 0, err
	}

	revisions := make([]dto.ReservationRevisionResponse, len(logs))
	for i, log := range logs {
		revisions[i] = s.convertToReservationRevision(ctx, &log)
	}

	return revisions, total, nil
}

func (s *historyService) convertToRoomRevision(ctx context.Context, log *audit.AuditLog) dto.RoomRevisionResponse {
	var roomEntity dto.RoomResponse

	valuesJSON := log.NewValues
	if log.Action == audit.ActionDelete {
		valuesJSON = log.OldValues
	}

	if valuesJSON != nil && len(valuesJSON) > 0 {
		var snapshot dto.RoomHistorySnapshot
		if err := json.Unmarshal(valuesJSON, &snapshot); err == nil {
			roomEntity = dto.RoomResponse{
				ID:          snapshot.ID,
				Number:      snapshot.Number,
				RoomGroupID: snapshot.RoomGroupID,
				Note:        snapshot.Note,
				Status:      snapshot.Status,
				CreatedBy:   s.getUserSummary(ctx, snapshot.CreatedBy),
				UpdatedBy:   s.getUserSummary(ctx, snapshot.UpdatedBy),
			}
		}
	}

	return dto.RoomRevisionResponse{
		Entity:           roomEntity,
		HistoryType:      dto.ActionToHistoryType(string(log.Action)),
		HistoryCreatedAt: dto.CustomTime{Time: log.CreatedAt},
		UpdatedFields:    dto.ParseChangedFields(log.ChangedFields),
	}
}

func (s *historyService) convertToReservationRevision(ctx context.Context, log *audit.AuditLog) dto.ReservationRevisionResponse {
	var reservationEntity dto.ReservationResponse

	valuesJSON := log.NewValues
	if log.Action == audit.ActionDelete {
		valuesJSON = log.OldValues
	}

	if valuesJSON != nil && len(valuesJSON) > 0 {
		var snapshot dto.ReservationHistorySnapshot
		if err := json.Unmarshal(valuesJSON, &snapshot); err == nil {
			reservationEntity = dto.ReservationResponse{
				ID:              snapshot.ID,
				PaymentMethodID: snapshot.PaymentMethodID,
				Name:            snapshot.Name,
				Phone:           snapshot.Phone,
				PeopleCount:     snapshot.PeopleCount,
				Price:           snapshot.Price,
				Deposit:         snapshot.Deposit,
				PaymentAmount:   snapshot.PaymentAmount,
				RefundAmount:    snapshot.RefundAmount,
				BrokerFee:       snapshot.BrokerFee,
				Note:            snapshot.Note,
				Status:          snapshot.Status,
				Type:            snapshot.Type,
				CreatedBy:       s.getUserSummary(ctx, snapshot.CreatedBy),
				UpdatedBy:       s.getUserSummary(ctx, snapshot.UpdatedBy),
				Rooms:           []dto.RoomResponse{},
			}

			if stayStartAt, err := time.Parse("2006-01-02", snapshot.StayStartAt); err == nil {
				reservationEntity.StayStartAt = dto.JSONDate{Time: stayStartAt}
			}
			if stayEndAt, err := time.Parse("2006-01-02", snapshot.StayEndAt); err == nil {
				reservationEntity.StayEndAt = dto.JSONDate{Time: stayEndAt}
			}
			if snapshot.CheckInAt != nil {
				if checkInAt, err := time.Parse(time.RFC3339, *snapshot.CheckInAt); err == nil {
					reservationEntity.CheckInAt = &dto.CustomTime{Time: checkInAt}
				}
			}
			if snapshot.CheckOutAt != nil {
				if checkOutAt, err := time.Parse(time.RFC3339, *snapshot.CheckOutAt); err == nil {
					reservationEntity.CheckOutAt = &dto.CustomTime{Time: checkOutAt}
				}
			}
			if snapshot.CanceledAt != nil {
				if canceledAt, err := time.Parse(time.RFC3339, *snapshot.CanceledAt); err == nil {
					reservationEntity.CanceledAt = &dto.CustomTime{Time: canceledAt}
				}
			}
		}
	}

	return dto.ReservationRevisionResponse{
		Entity:           reservationEntity,
		HistoryType:      dto.ActionToHistoryType(string(log.Action)),
		HistoryCreatedAt: dto.CustomTime{Time: log.CreatedAt},
		UpdatedFields:    dto.ParseChangedFields(log.ChangedFields),
	}
}

func (s *historyService) getUserSummary(ctx context.Context, userID uint) *dto.UserSummaryResponse {
	if userID == 0 {
		return nil
	}
	if user, err := s.userService.GetByID(ctx, userID); err == nil {
		email := ""
		if user.Email != nil {
			email = *user.Email
		}
		return &dto.UserSummaryResponse{
			ID:     user.ID,
			UserID: user.UserID,
			Email:  email,
			Name:   user.Name,
		}
	}
	return &dto.UserSummaryResponse{
		ID: userID,
	}
}
