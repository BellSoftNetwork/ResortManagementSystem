package services

import (
	"context"
	"errors"
	"time"

	"gitlab.bellsoft.net/rms/api-core/internal/audit"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"gorm.io/gorm"
)

var (
	ErrRoomNotFound        = errors.New("존재하지 않는 객실")
	ErrRoomNumberExists    = errors.New("이미 존재하는 객실 번호")
	ErrRoomHasReservations = errors.New("예약이 존재하는 객실")
)

type RoomService interface {
	GetByID(ctx context.Context, id uint) (*models.Room, error)
	GetByIDWithGroup(ctx context.Context, id uint) (*models.Room, error)
	GetAll(ctx context.Context, filter repositories.RoomFilter, page, size int, sort string) ([]models.Room, int64, error)
	GetAvailableRooms(ctx context.Context, startDate, endDate time.Time, excludeReservationID *uint) ([]models.Room, error)
	Create(ctx context.Context, room *models.Room) error
	Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.Room, error)
	Delete(ctx context.Context, id uint) error
}

type roomService struct {
	roomRepo      repositories.RoomRepository
	roomGroupRepo repositories.RoomGroupRepository
	auditService  audit.AuditService
}

func NewRoomService(roomRepo repositories.RoomRepository, roomGroupRepo repositories.RoomGroupRepository, auditService audit.AuditService) RoomService {
	return &roomService{
		roomRepo:      roomRepo,
		roomGroupRepo: roomGroupRepo,
		auditService:  auditService,
	}
}

func (s *roomService) GetByID(ctx context.Context, id uint) (*models.Room, error) {
	room, err := s.roomRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrRoomNotFound
	}
	return room, nil
}

func (s *roomService) GetByIDWithGroup(ctx context.Context, id uint) (*models.Room, error) {
	room, err := s.roomRepo.FindByIDWithGroup(ctx, id)
	if err != nil {
		return nil, ErrRoomNotFound
	}
	return room, nil
}

func (s *roomService) GetAll(ctx context.Context, filter repositories.RoomFilter, page, size int, sort string) ([]models.Room, int64, error) {
	offset := page * size
	return s.roomRepo.FindAll(ctx, filter, offset, size, sort)
}

func (s *roomService) GetAvailableRooms(ctx context.Context, startDate, endDate time.Time, excludeReservationID *uint) ([]models.Room, error) {
	return s.roomRepo.FindAvailableRooms(ctx, startDate, endDate, excludeReservationID)
}

func (s *roomService) Create(ctx context.Context, room *models.Room) error {
	exists, err := s.roomRepo.ExistsByNumber(ctx, room.Number, nil)
	if err != nil {
		return err
	}
	if exists {
		return ErrRoomNumberExists
	}

	_, err = s.roomGroupRepo.FindByID(ctx, room.RoomGroupID)
	if err != nil {
		return ErrRoomGroupNotFound
	}

	_, err = s.roomRepo.Create(ctx, room)
	return err
}

func (s *roomService) Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.Room, error) {
	room, err := s.roomRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrRoomNotFound
	}

	if number, ok := updates["number"].(string); ok {
		if number != room.Number {
			exists, err := s.roomRepo.ExistsByNumber(ctx, number, &id)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrRoomNumberExists
			}
			room.Number = number
		}
	}

	if roomGroupID, ok := updates["room_group_id"].(uint); ok {
		if roomGroupID != room.RoomGroupID {
			_, err := s.roomGroupRepo.FindByID(ctx, roomGroupID)
			if err != nil {
				return nil, ErrRoomGroupNotFound
			}
			room.RoomGroupID = roomGroupID
		}
	}

	if note, ok := updates["note"].(string); ok {
		room.Note = note
	}

	if status, ok := updates["status"].(models.RoomStatus); ok {
		room.Status = status
	}

	if err := s.roomRepo.Update(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomService) Delete(ctx context.Context, id uint) error {
	room, err := s.roomRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRoomNotFound
		}
		return err
	}

	// Perform the deletion
	if err := s.roomRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Log the deletion manually since soft delete doesn't trigger GORM hooks
	if s.auditService != nil {
		if err := s.auditService.LogDelete(ctx, room); err != nil {
			// Log error but don't fail the deletion
			// The deletion has already succeeded, we just couldn't log it
		}
	}

	return nil
}
