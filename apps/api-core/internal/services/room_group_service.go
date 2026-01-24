package services

import (
	"context"
	"errors"

	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
)

var (
	ErrRoomGroupNotFound   = errors.New("존재하지 않는 객실 그룹")
	ErrRoomGroupNameExists = errors.New("이미 존재하는 객실 그룹 이름")
	ErrRoomGroupHasRooms   = errors.New("객실이 존재하는 객실 그룹")
)

type RoomGroupService interface {
	GetByID(ctx context.Context, id uint) (*models.RoomGroup, error)
	GetByIDWithRooms(ctx context.Context, id uint, roomStatus *models.RoomStatus) (*models.RoomGroup, error)
	GetByIDWithFilteredRooms(ctx context.Context, id uint, filter repositories.RoomGroupRoomFilter) (*models.RoomGroup, error)
	GetAll(ctx context.Context, page, size int) ([]models.RoomGroup, int64, error)
	GetByIDWithUsers(ctx context.Context, id uint) (*models.RoomGroup, error)
	GetAllWithUsers(ctx context.Context, page, size int, sort string) ([]models.RoomGroup, int64, error)
	Create(ctx context.Context, roomGroup *models.RoomGroup) error
	Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.RoomGroup, error)
	Delete(ctx context.Context, id uint) error
}

type roomGroupService struct {
	roomGroupRepo repositories.RoomGroupRepository
}

func NewRoomGroupService(roomGroupRepo repositories.RoomGroupRepository) RoomGroupService {
	return &roomGroupService{
		roomGroupRepo: roomGroupRepo,
	}
}

func (s *roomGroupService) GetByID(ctx context.Context, id uint) (*models.RoomGroup, error) {
	roomGroup, err := s.roomGroupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrRoomGroupNotFound
	}
	return roomGroup, nil
}

func (s *roomGroupService) GetByIDWithRooms(ctx context.Context, id uint, roomStatus *models.RoomStatus) (*models.RoomGroup, error) {
	roomGroup, err := s.roomGroupRepo.FindByIDWithRooms(ctx, id, roomStatus)
	if err != nil {
		return nil, ErrRoomGroupNotFound
	}
	return roomGroup, nil
}

func (s *roomGroupService) GetByIDWithFilteredRooms(ctx context.Context, id uint, filter repositories.RoomGroupRoomFilter) (*models.RoomGroup, error) {
	roomGroup, err := s.roomGroupRepo.FindByIDWithFilteredRooms(ctx, id, filter)
	if err != nil {
		return nil, ErrRoomGroupNotFound
	}
	return roomGroup, nil
}

func (s *roomGroupService) GetAll(ctx context.Context, page, size int) ([]models.RoomGroup, int64, error) {
	offset := page * size
	return s.roomGroupRepo.FindAll(ctx, offset, size)
}

func (s *roomGroupService) Create(ctx context.Context, roomGroup *models.RoomGroup) error {
	exists, err := s.roomGroupRepo.ExistsByName(ctx, roomGroup.Name, nil)
	if err != nil {
		return err
	}
	if exists {
		return ErrRoomGroupNameExists
	}

	_, err = s.roomGroupRepo.Create(ctx, roomGroup)
	return err
}

func (s *roomGroupService) Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.RoomGroup, error) {
	roomGroup, err := s.roomGroupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrRoomGroupNotFound
	}

	if name, ok := updates["name"].(string); ok {
		if name != roomGroup.Name {
			exists, err := s.roomGroupRepo.ExistsByName(ctx, name, &id)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrRoomGroupNameExists
			}
			roomGroup.Name = name
		}
	}

	if peekPrice, ok := updates["peekPrice"].(int); ok {
		roomGroup.PeekPrice = peekPrice
	}

	if offPeekPrice, ok := updates["offPeekPrice"].(int); ok {
		roomGroup.OffPeekPrice = offPeekPrice
	}

	if description, ok := updates["description"].(string); ok {
		roomGroup.Description = description
	}

	if err := s.roomGroupRepo.Update(ctx, roomGroup); err != nil {
		return nil, err
	}

	return roomGroup, nil
}

func (s *roomGroupService) Delete(ctx context.Context, id uint) error {
	roomGroup, err := s.roomGroupRepo.FindByIDWithRooms(ctx, id, nil)
	if err != nil {
		return ErrRoomGroupNotFound
	}

	if len(roomGroup.Rooms) > 0 {
		return ErrRoomGroupHasRooms
	}

	return s.roomGroupRepo.Delete(ctx, id)
}

func (s *roomGroupService) GetByIDWithUsers(ctx context.Context, id uint) (*models.RoomGroup, error) {
	roomGroup, err := s.roomGroupRepo.FindByIDWithUsers(ctx, id)
	if err != nil {
		return nil, ErrRoomGroupNotFound
	}
	return roomGroup, nil
}

func (s *roomGroupService) GetAllWithUsers(ctx context.Context, page, size int, sort string) ([]models.RoomGroup, int64, error) {
	offset := page * size
	return s.roomGroupRepo.FindAllWithUsers(ctx, offset, size, sort)
}
