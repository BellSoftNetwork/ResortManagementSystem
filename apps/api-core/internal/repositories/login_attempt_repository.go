package repositories

import (
	"context"
	"time"

	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gorm.io/gorm"
)

type LoginAttemptRepository interface {
	Create(ctx context.Context, attempt *models.LoginAttempt) error
	CountRecentFailedAttempts(ctx context.Context, username string, ipAddress string, since time.Time) (int64, error)
	GetLastSuccessfulAttempt(ctx context.Context, username string) (*models.LoginAttempt, error)
	CleanOldAttempts(ctx context.Context, before time.Time) error
}

type loginAttemptRepository struct {
	db *gorm.DB
}

func NewLoginAttemptRepository(db *gorm.DB) LoginAttemptRepository {
	return &loginAttemptRepository{db: db}
}

func (r *loginAttemptRepository) Create(ctx context.Context, attempt *models.LoginAttempt) error {
	return r.db.WithContext(ctx).Create(attempt).Error
}

func (r *loginAttemptRepository) CountRecentFailedAttempts(ctx context.Context, username string, ipAddress string, since time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.LoginAttempt{}).
		Where("(username = ? OR ip_address = ?) AND successful = ? AND attempt_at >= ?",
			username, ipAddress, false, since).
		Count(&count).Error
	return count, err
}

func (r *loginAttemptRepository) GetLastSuccessfulAttempt(ctx context.Context, username string) (*models.LoginAttempt, error) {
	var attempt models.LoginAttempt
	err := r.db.WithContext(ctx).
		Where("username = ? AND successful = ?", username, true).
		Order("attempt_at DESC").
		First(&attempt).Error
	if err != nil {
		return nil, err
	}
	return &attempt, nil
}

func (r *loginAttemptRepository) CleanOldAttempts(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).
		Where("attempt_at < ?", before).
		Delete(&models.LoginAttempt{}).Error
}
