package repositories

import (
	"context"
	"time"

	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByUserID(ctx context.Context, userID string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindAll(ctx context.Context, offset, limit int) ([]models.User, int64, error)
	ExistsByUserID(ctx context.Context, userID string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	HasAnyUsers(ctx context.Context) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	// Set deleted_at to default value for new records
	user.DeletedAt = models.DefaultDeletedAt()
	return r.db.WithContext(ctx).Unscoped().Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Unscoped().Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	// Soft delete by updating deleted_at
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("deleted_at", now).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Unscoped().Where("id = ? AND deleted_at = ?", id, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUserID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Unscoped().Where("user_id = ? AND deleted_at = ?", userID, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Unscoped().Where("email = ? AND deleted_at = ?", email, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	query := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).Where("deleted_at = ?", defaultDeletedAt)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	var count int64
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).Where("user_id = ? AND deleted_at = ?", userID, defaultDeletedAt).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).Where("email = ? AND deleted_at = ?", email, defaultDeletedAt).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) HasAnyUsers(ctx context.Context) (bool, error) {
	var count int64
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Unscoped().Model(&models.User{}).Where("deleted_at = ?", defaultDeletedAt).Count(&count).Error
	return count > 0, err
}
