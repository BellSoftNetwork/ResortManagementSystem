package repositories

import (
	"context"
	"time"

	appContext "gitlab.bellsoft.net/rms/api-core/internal/context"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gorm.io/gorm"
)

type DateBlockRepository interface {
	Create(ctx context.Context, dateBlock *models.DateBlock) (*models.DateBlock, error)
	Update(ctx context.Context, dateBlock *models.DateBlock) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.DateBlock, error)
	FindAll(ctx context.Context, filter dto.DateBlockFilter, offset, limit int) ([]models.DateBlock, int64, error)
	IsDateRangeBlocked(ctx context.Context, startDate, endDate time.Time) (bool, error)
}

type dateBlockRepository struct {
	db *gorm.DB
}

func NewDateBlockRepository(db *gorm.DB) DateBlockRepository {
	return &dateBlockRepository{db: db}
}

func (r *dateBlockRepository) Create(ctx context.Context, dateBlock *models.DateBlock) (*models.DateBlock, error) {
	err := r.db.WithContext(ctx).Create(dateBlock).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).
		Preload("CreatedByUser").
		Where("id = ?", dateBlock.ID).
		First(dateBlock).Error
	if err != nil {
		return nil, err
	}

	return dateBlock, nil
}

func (r *dateBlockRepository) Update(ctx context.Context, dateBlock *models.DateBlock) error {
	return r.db.WithContext(ctx).Save(dateBlock).Error
}

func (r *dateBlockRepository) Delete(ctx context.Context, id uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"deleted_at": now,
	}

	if userID, ok := appContext.GetUserID(ctx); ok {
		updates["updated_by"] = userID
	}

	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	return r.db.WithContext(ctx).
		Model(&models.DateBlock{}).
		Where("id = ? AND deleted_at = ?", id, defaultDeletedAt).
		Updates(updates).Error
}

func (r *dateBlockRepository) FindByID(ctx context.Context, id uint) (*models.DateBlock, error) {
	var dateBlock models.DateBlock
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).
		Preload("CreatedByUser").
		Where("id = ? AND deleted_at = ?", id, defaultDeletedAt).
		First(&dateBlock).Error
	if err != nil {
		return nil, err
	}

	return &dateBlock, nil
}

func (r *dateBlockRepository) FindAll(ctx context.Context, filter dto.DateBlockFilter, offset, limit int) ([]models.DateBlock, int64, error) {
	var dateBlocks []models.DateBlock
	var total int64

	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	query := r.db.WithContext(ctx).
		Model(&models.DateBlock{}).
		Where("deleted_at = ?", defaultDeletedAt).
		Preload("CreatedByUser")

	if filter.StartDate != "" && filter.EndDate != "" {
		query = query.Where("NOT (end_date < ? OR start_date >= ?)", filter.StartDate, filter.EndDate)
	} else if filter.StartDate != "" {
		query = query.Where("end_date >= ?", filter.StartDate)
	} else if filter.EndDate != "" {
		query = query.Where("start_date < ?", filter.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("start_date ASC, id ASC").Offset(offset).Limit(limit).Find(&dateBlocks).Error
	if err != nil {
		return nil, 0, err
	}

	return dateBlocks, total, nil
}

func (r *dateBlockRepository) IsDateRangeBlocked(ctx context.Context, startDate, endDate time.Time) (bool, error) {
	var count int64
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	err := r.db.WithContext(ctx).
		Model(&models.DateBlock{}).
		Where("NOT (end_date < ? OR start_date >= ?) AND deleted_at = ?", startDate, endDate, defaultDeletedAt).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
