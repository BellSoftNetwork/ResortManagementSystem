package repositories

import (
	"context"
	"strings"
	"time"

	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gorm.io/gorm"
)

type PaymentMethodRepository interface {
	Create(ctx context.Context, paymentMethod *models.PaymentMethod) (*models.PaymentMethod, error)
	Update(ctx context.Context, paymentMethod *models.PaymentMethod) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.PaymentMethod, error)
	FindAll(ctx context.Context, offset, limit int, sort string) ([]models.PaymentMethod, int64, error)
	FindActive(ctx context.Context) ([]models.PaymentMethod, error)
	ExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error)
	FindByName(ctx context.Context, name string) (*models.PaymentMethod, error)
	ResetAllDefaultSelects(ctx context.Context) error
}

type paymentMethodRepository struct {
	db *gorm.DB
}

func NewPaymentMethodRepository(db *gorm.DB) PaymentMethodRepository {
	return &paymentMethodRepository{db: db}
}

func (r *paymentMethodRepository) Create(ctx context.Context, paymentMethod *models.PaymentMethod) (*models.PaymentMethod, error) {
	err := r.db.WithContext(ctx).Create(paymentMethod).Error
	return paymentMethod, err
}

func (r *paymentMethodRepository) Update(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	return r.db.WithContext(ctx).Save(paymentMethod).Error
}

func (r *paymentMethodRepository) Delete(ctx context.Context, id uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"deleted_at": now,
	}

	return r.db.WithContext(ctx).Model(&models.PaymentMethod{}).Where("id = ?", id).Updates(updates).Error
}

func (r *paymentMethodRepository) FindByID(ctx context.Context, id uint) (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, defaultDeletedAt).First(&paymentMethod).Error
	if err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

func (r *paymentMethodRepository) FindAll(ctx context.Context, offset, limit int, sort string) ([]models.PaymentMethod, int64, error) {
	var paymentMethods []models.PaymentMethod
	var total int64

	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	query := r.db.WithContext(ctx).Model(&models.PaymentMethod{}).Where("deleted_at = ?", defaultDeletedAt)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 정렬 처리
	orderClause := r.parseSort(sort)
	if orderClause == "" {
		orderClause = "id DESC" // 기본 정렬
	}
	query = query.Order(orderClause)

	err = query.Offset(offset).Limit(limit).Find(&paymentMethods).Error
	if err != nil {
		return nil, 0, err
	}

	return paymentMethods, total, nil
}

func (r *paymentMethodRepository) FindActive(ctx context.Context) ([]models.PaymentMethod, error) {
	var paymentMethods []models.PaymentMethod
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).
		Where("status = ? AND deleted_at = ?", models.PaymentMethodStatusActive, defaultDeletedAt).
		Order("id").
		Find(&paymentMethods).Error
	return paymentMethods, err
}

func (r *paymentMethodRepository) ExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error) {
	var count int64
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	query := r.db.WithContext(ctx).Model(&models.PaymentMethod{}).Where("name = ? AND deleted_at = ?", name, defaultDeletedAt)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *paymentMethodRepository) FindByName(ctx context.Context, name string) (*models.PaymentMethod, error) {
	var paymentMethod models.PaymentMethod
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Where("name = ? AND deleted_at = ?", name, defaultDeletedAt).First(&paymentMethod).Error
	if err != nil {
		return nil, err
	}
	return &paymentMethod, nil
}

func (r *paymentMethodRepository) ResetAllDefaultSelects(ctx context.Context) error {
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	return r.db.WithContext(ctx).
		Model(&models.PaymentMethod{}).
		Where("deleted_at = ?", defaultDeletedAt).
		Update("is_default_select", models.BitBool(false)).Error
}

// parseSort는 Spring Boot 형식의 정렬 파라미터를 GORM 형식으로 변환합니다.
// 예: "name,desc" -> "name DESC"
// 예: "name,desc,commissionRate,asc" -> "name DESC, commission_rate ASC"
func (r *paymentMethodRepository) parseSort(sort string) string {
	if sort == "" {
		return ""
	}

	parts := strings.Split(sort, ",")
	if len(parts) < 2 {
		return ""
	}

	var orderClauses []string
	for i := 0; i < len(parts); i += 2 {
		if i+1 >= len(parts) {
			break
		}

		field := parts[i]
		direction := strings.ToUpper(parts[i+1])

		// 필드명을 데이터베이스 컬럼명으로 매핑
		dbField := r.mapSortField(field)
		if dbField == "" {
			continue
		}

		// 방향 검증
		if direction != "ASC" && direction != "DESC" {
			direction = "ASC"
		}

		orderClauses = append(orderClauses, dbField+" "+direction)
	}

	return strings.Join(orderClauses, ", ")
}

// mapSortField는 API 필드명을 데이터베이스 컬럼명으로 매핑합니다.
func (r *paymentMethodRepository) mapSortField(field string) string {
	fieldMap := map[string]string{
		"id":                       "id",
		"name":                     "name",
		"commissionRate":           "commission_rate",
		"requireUnpaidAmountCheck": "require_unpaid_amount_check",
		"isDefaultSelect":          "is_default_select",
		"status":                   "status",
		"createdAt":                "created_at",
		"updatedAt":                "updated_at",
	}

	if dbField, ok := fieldMap[field]; ok {
		return dbField
	}
	return ""
}
