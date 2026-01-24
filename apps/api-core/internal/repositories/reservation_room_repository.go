package repositories

import (
	"context"
	"time"

	appContext "gitlab.bellsoft.net/rms/api-core/internal/context"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gorm.io/gorm"
)

type ReservationRoomRepository interface {
	Create(ctx context.Context, reservationRoom *models.ReservationRoom) (*models.ReservationRoom, error)
	Delete(ctx context.Context, id uint) error
	FindByReservationID(ctx context.Context, reservationID uint) ([]models.ReservationRoom, error)
}

type reservationRoomRepository struct {
	db *gorm.DB
}

func NewReservationRoomRepository(db *gorm.DB) ReservationRoomRepository {
	return &reservationRoomRepository{db: db}
}

func (r *reservationRoomRepository) Create(ctx context.Context, reservationRoom *models.ReservationRoom) (*models.ReservationRoom, error) {
	err := r.db.WithContext(ctx).Create(reservationRoom).Error
	return reservationRoom, err
}

func (r *reservationRoomRepository) Delete(ctx context.Context, id uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"deleted_at": now,
	}

	// Set updated by field for soft delete
	if userID, ok := appContext.GetUserID(ctx); ok {
		updates["updated_by"] = userID
	}

	return r.db.WithContext(ctx).Model(&models.ReservationRoom{}).Where("id = ?", id).Updates(updates).Error
}

func (r *reservationRoomRepository) FindByReservationID(ctx context.Context, reservationID uint) ([]models.ReservationRoom, error) {
	var reservationRooms []models.ReservationRoom
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).
		Where("reservation_id = ? AND deleted_at = ?", reservationID, defaultDeletedAt).
		Find(&reservationRooms).Error
	return reservationRooms, err
}
