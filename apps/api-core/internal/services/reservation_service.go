package services

import (
	"context"
	"errors"
	"time"

	"gitlab.bellsoft.net/rms/api-core/internal/audit"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
)

var (
	ErrReservationNotFound   = errors.New("존재하지 않는 예약")
	ErrRoomNotAvailable      = errors.New("해당 기간에 예약이 불가능한 객실")
	ErrInvalidDateRange      = errors.New("잘못된 날짜 범위")
	ErrPaymentMethodInactive = errors.New("비활성 상태의 결제 수단")
)

type ReservationService interface {
	GetByID(ctx context.Context, id uint) (*models.Reservation, error)
	GetByIDWithDetails(ctx context.Context, id uint) (*models.Reservation, error)
	GetAll(ctx context.Context, filter dto.ReservationRepositoryFilter, page, size int, sort string) ([]models.Reservation, int64, error)
	GetStatistics(ctx context.Context, startDate, endDate time.Time, periodType string) ([]repositories.ReservationStatistics, error)
	Create(ctx context.Context, reservation *models.Reservation, roomIDs []uint) error
	Update(ctx context.Context, id uint, updates map[string]interface{}, roomIDs []uint, hasRoomsUpdate bool) (*models.Reservation, error)
	Delete(ctx context.Context, id uint) error
	GetAvailableRooms(ctx context.Context, startDate, endDate time.Time, excludeReservationID *uint) ([]models.Room, error)
	GetLastReservationForRoom(ctx context.Context, roomID uint) (*models.Reservation, error)
}

type reservationService struct {
	reservationRepo   repositories.ReservationRepository
	roomRepo          repositories.RoomRepository
	paymentMethodRepo repositories.PaymentMethodRepository
	auditService      audit.AuditService
}

func NewReservationService(reservationRepo repositories.ReservationRepository, roomRepo repositories.RoomRepository,
	paymentMethodRepo repositories.PaymentMethodRepository, auditService audit.AuditService) ReservationService {
	return &reservationService{
		reservationRepo:   reservationRepo,
		roomRepo:          roomRepo,
		paymentMethodRepo: paymentMethodRepo,
		auditService:      auditService,
	}
}

func (s *reservationService) GetByID(ctx context.Context, id uint) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrReservationNotFound
	}
	return reservation, nil
}

func (s *reservationService) GetByIDWithDetails(ctx context.Context, id uint) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.FindByIDWithDetails(ctx, id)
	if err != nil {
		return nil, ErrReservationNotFound
	}
	return reservation, nil
}

func (s *reservationService) GetAll(ctx context.Context, filter dto.ReservationRepositoryFilter, page, size int, sort string) ([]models.Reservation, int64, error) {
	offset := page * size
	return s.reservationRepo.FindAll(ctx, filter, offset, size, sort)
}

func (s *reservationService) GetStatistics(ctx context.Context, startDate, endDate time.Time, periodType string) ([]repositories.ReservationStatistics, error) {
	if periodType == "" {
		periodType = "MONTHLY"
	}
	return s.reservationRepo.GetStatistics(ctx, startDate, endDate, periodType)
}

func (s *reservationService) Create(ctx context.Context, reservation *models.Reservation, roomIDs []uint) error {
	if reservation.StayStartAt.After(reservation.StayEndAt) || reservation.StayStartAt.Equal(reservation.StayEndAt) {
		return ErrInvalidDateRange
	}

	paymentMethod, err := s.paymentMethodRepo.FindByID(ctx, reservation.PaymentMethodID)
	if err != nil {
		return ErrPaymentMethodNotFound
	}
	reservation.PaymentMethod = paymentMethod // 추가: audit 로깅용
	if !paymentMethod.IsActive() {
		return ErrPaymentMethodInactive
	}

	for _, roomID := range roomIDs {
		available, err := s.roomRepo.IsRoomAvailable(ctx, roomID, reservation.StayStartAt, reservation.StayEndAt, nil)
		if err != nil {
			return err
		}
		if !available {
			return ErrRoomNotAvailable
		}
	}

	reservation.Rooms = make([]models.ReservationRoom, len(roomIDs))
	for i, roomID := range roomIDs {
		room, err := s.roomRepo.FindByID(ctx, roomID)
		if err != nil {
			return ErrRoomNotFound
		}
		reservation.Rooms[i] = models.ReservationRoom{
			RoomID: roomID,
			Room:   room, // 추가: audit 로깅용
		}
	}

	reservation.BrokerFee = int(float64(reservation.Price) * paymentMethod.CommissionRate)

	_, err = s.reservationRepo.Create(ctx, reservation)
	return err
}

func (s *reservationService) Update(ctx context.Context, id uint, updates map[string]interface{}, roomIDs []uint, hasRoomsUpdate bool) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.FindByIDWithDetails(ctx, id)
	if err != nil {
		return nil, ErrReservationNotFound
	}

	if name, ok := updates["name"].(string); ok {
		reservation.Name = name
	}

	if phone, ok := updates["phone"].(string); ok {
		reservation.Phone = phone
	}

	if peopleCount, ok := updates["peopleCount"].(int); ok {
		reservation.PeopleCount = peopleCount
	}

	if stayStartAt, ok := updates["stayStartAt"].(time.Time); ok {
		reservation.StayStartAt = stayStartAt
	}

	if stayEndAt, ok := updates["stayEndAt"].(time.Time); ok {
		reservation.StayEndAt = stayEndAt
	}

	if reservation.StayStartAt.After(reservation.StayEndAt) || reservation.StayStartAt.Equal(reservation.StayEndAt) {
		return nil, ErrInvalidDateRange
	}

	if checkInAt, ok := updates["checkInAt"].(*time.Time); ok {
		reservation.CheckInAt = checkInAt
	}

	if checkOutAt, ok := updates["checkOutAt"].(*time.Time); ok {
		reservation.CheckOutAt = checkOutAt
	}

	if price, ok := updates["price"].(int); ok {
		reservation.Price = price
	}

	if deposit, ok := updates["deposit"].(int); ok {
		reservation.Deposit = deposit
	}

	if paymentAmount, ok := updates["paymentAmount"].(int); ok {
		reservation.PaymentAmount = paymentAmount
	}

	if refundAmount, ok := updates["refundAmount"].(int); ok {
		reservation.RefundAmount = refundAmount
	}

	if note, ok := updates["note"].(string); ok {
		reservation.Note = note
	}

	if status, ok := updates["status"].(models.ReservationStatus); ok {
		reservation.Status = status
		if status == models.ReservationStatusCancel || status == models.ReservationStatusRefund {
			now := time.Now()
			reservation.CanceledAt = &now
		}
	}

	if type_, ok := updates["type"].(models.ReservationType); ok {
		reservation.Type = type_
	}

	if paymentMethodID, ok := updates["paymentMethodId"].(uint); ok {
		if paymentMethodID != reservation.PaymentMethodID {
			paymentMethod, err := s.paymentMethodRepo.FindByID(ctx, paymentMethodID)
			if err != nil {
				return nil, ErrPaymentMethodNotFound
			}
			if !paymentMethod.IsActive() {
				return nil, ErrPaymentMethodInactive
			}
			reservation.PaymentMethod = nil // GORM Save 충돌 방지: Preload된 association을 nil로 설정
			reservation.PaymentMethodID = paymentMethodID
			reservation.BrokerFee = int(float64(reservation.Price) * paymentMethod.CommissionRate)
		}
	}

	if hasRoomsUpdate {
		for _, roomID := range roomIDs {
			available, err := s.roomRepo.IsRoomAvailable(ctx, roomID, reservation.StayStartAt, reservation.StayEndAt, &id)
			if err != nil {
				return nil, err
			}
			if !available {
				return nil, ErrRoomNotAvailable
			}
		}

		if err := s.reservationRepo.DeleteRooms(ctx, id); err != nil {
			return nil, err
		}

		reservation.Rooms = make([]models.ReservationRoom, len(roomIDs))
		for i, roomID := range roomIDs {
			room, err := s.roomRepo.FindByID(ctx, roomID)
			if err != nil {
				return nil, ErrRoomNotFound
			}
			reservation.Rooms[i] = models.ReservationRoom{
				RoomID: roomID,
				Room:   room,
			}
		}
	}

	if err := s.reservationRepo.Update(ctx, reservation); err != nil {
		return nil, err
	}

	return s.reservationRepo.FindByIDWithDetails(ctx, id)
}

func (s *reservationService) Delete(ctx context.Context, id uint) error {
	reservation, err := s.reservationRepo.FindByID(ctx, id)
	if err != nil {
		return ErrReservationNotFound
	}

	// Perform the deletion
	if err := s.reservationRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Log the deletion manually since soft delete doesn't trigger GORM hooks
	if s.auditService != nil {
		if err := s.auditService.LogDelete(ctx, reservation); err != nil {
			// Log error but don't fail the deletion
			// The deletion has already succeeded, we just couldn't log it
		}
	}

	return nil
}

func (s *reservationService) GetAvailableRooms(ctx context.Context, startDate, endDate time.Time, excludeReservationID *uint) ([]models.Room, error) {
	if startDate.After(endDate) || startDate.Equal(endDate) {
		return nil, ErrInvalidDateRange
	}

	return s.roomRepo.FindAvailableRooms(ctx, startDate, endDate, excludeReservationID)
}

func (s *reservationService) GetLastReservationForRoom(ctx context.Context, roomID uint) (*models.Reservation, error) {
	return s.reservationRepo.FindLastReservationForRoom(ctx, roomID)
}
