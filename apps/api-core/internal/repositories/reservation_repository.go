package repositories

import (
	"context"
	"strings"
	"time"

	appContext "gitlab.bellsoft.net/rms/api-core/internal/context"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation *models.Reservation) (*models.Reservation, error)
	Update(ctx context.Context, reservation *models.Reservation) error
	Delete(ctx context.Context, id uint) error
	DeleteRooms(ctx context.Context, reservationID uint) error
	FindByID(ctx context.Context, id uint) (*models.Reservation, error)
	FindByIDWithDetails(ctx context.Context, id uint) (*models.Reservation, error)
	FindAll(ctx context.Context, filter ReservationFilter, offset, limit int, sort string) ([]models.Reservation, int64, error)
	GetStatistics(ctx context.Context, startDate, endDate time.Time, periodType string) ([]ReservationStatistics, error)
}

type ReservationFilter struct {
	Status    *models.ReservationStatus
	Type      *models.ReservationType
	RoomID    *uint
	StartDate *time.Time
	EndDate   *time.Time
	Search    string
}

type ReservationStatistics struct {
	Period           string  `json:"period"`
	ReservationCount int64   `json:"reservationCount"`
	TotalRevenue     float64 `json:"totalRevenue"`
	TotalGuests      int64   `json:"totalGuests"`
	AverageStayDays  float64 `json:"averageStayDays"`
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(ctx context.Context, reservation *models.Reservation) (*models.Reservation, error) {
	err := r.db.WithContext(ctx).Create(reservation).Error
	return reservation, err
}

func (r *reservationRepository) Update(ctx context.Context, reservation *models.Reservation) error {
	return r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Save(reservation).Error
}

func (r *reservationRepository) Delete(ctx context.Context, id uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"deleted_at": now,
	}

	if userID, ok := appContext.GetUserID(ctx); ok {
		updates["updated_by"] = userID
	}

	return r.db.WithContext(ctx).Model(&models.Reservation{}).Where("id = ?", id).Updates(updates).Error
}

func (r *reservationRepository) DeleteRooms(ctx context.Context, reservationID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&models.ReservationRoom{}).
		Where("reservation_id = ? AND deleted_at = ?", reservationID, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
		Update("deleted_at", now).Error
}

func (r *reservationRepository) FindByID(ctx context.Context, id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, defaultDeletedAt).First(&reservation).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) FindByIDWithDetails(ctx context.Context, id uint) (*models.Reservation, error) {
	var reservation models.Reservation
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).
		Preload("PaymentMethod", "deleted_at = ?", defaultDeletedAt).
		Preload("Rooms", "deleted_at = ?", defaultDeletedAt).
		Preload("Rooms.Room", "deleted_at = ?", defaultDeletedAt).
		Preload("Rooms.Room.RoomGroup", "deleted_at = ?", defaultDeletedAt).
		Where("id = ? AND deleted_at = ?", id, defaultDeletedAt).
		First(&reservation).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) FindAll(ctx context.Context, filter ReservationFilter, offset, limit int, sort string) ([]models.Reservation, int64, error) {
	var reservations []models.Reservation
	var total int64

	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	query := r.db.WithContext(ctx).Model(&models.Reservation{}).
		Where("deleted_at = ?", defaultDeletedAt).
		Preload("PaymentMethod", "deleted_at = ?", defaultDeletedAt).
		Preload("Rooms", "deleted_at = ?", defaultDeletedAt).
		Preload("Rooms.Room", "deleted_at = ?", defaultDeletedAt).
		Preload("Rooms.Room.RoomGroup", "deleted_at = ?", defaultDeletedAt)

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}

	if filter.RoomID != nil {
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
		query = query.Joins("JOIN reservation_room ON reservation_room.reservation_id = reservation.id").
			Where("reservation_room.room_id = ? AND reservation_room.deleted_at = ?", *filter.RoomID, defaultDeletedAt)
	}

	if filter.StartDate != nil {
		query = query.Where("(stay_start_at >= ? OR stay_end_at >= ?)", *filter.StartDate, *filter.StartDate)
	}
	if filter.EndDate != nil {
		query = query.Where("(stay_start_at <= ? OR stay_end_at <= ?)", *filter.EndDate, *filter.EndDate)
	}

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name LIKE ? OR phone LIKE ?", searchPattern, searchPattern)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 정렬 처리
	orderClause := r.parseSort(sort)
	if orderClause == "" {
		orderClause = "id DESC" // 기본 정렬
	}

	err = query.Offset(offset).Limit(limit).Order(orderClause).Find(&reservations).Error
	if err != nil {
		return nil, 0, err
	}

	return reservations, total, nil
}

// parseSort는 Spring Boot 형식의 정렬 파라미터를 GORM 형식으로 변환합니다.
// 예: "price,desc" -> "price DESC"
// 예: "price,desc,stayEndAt,asc" -> "price DESC, stay_end_at ASC"
func (r *reservationRepository) parseSort(sort string) string {
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
func (r *reservationRepository) mapSortField(field string) string {
	fieldMap := map[string]string{
		"id":            "id",
		"name":          "name",
		"phone":         "phone",
		"price":         "price",
		"deposit":       "deposit",
		"paymentAmount": "payment_amount",
		"peopleCount":   "people_count",
		"stayStartAt":   "stay_start_at",
		"stayEndAt":     "stay_end_at",
		"checkInAt":     "check_in_at",
		"checkOutAt":    "check_out_at",
		"status":        "status",
		"type":          "type",
		"createdAt":     "created_at",
		"updatedAt":     "updated_at",
	}

	if dbField, ok := fieldMap[field]; ok {
		return dbField
	}
	return ""
}

func (r *reservationRepository) GetStatistics(ctx context.Context, startDate, endDate time.Time, periodType string) ([]ReservationStatistics, error) {
	var stats []ReservationStatistics

	var dateFormat string
	switch periodType {
	case "DAILY":
		dateFormat = "%Y-%m-%d"
	case "MONTHLY":
		dateFormat = "%Y-%m"
	case "YEARLY":
		dateFormat = "%Y"
	default:
		dateFormat = "%Y-%m"
	}

	err := r.db.WithContext(ctx).
		Model(&models.Reservation{}).
		Select(`
			DATE_FORMAT(stay_start_at, ?) as period,
			COUNT(*) as reservation_count,
			SUM(price) as total_revenue,
			SUM(people_count) as total_guests,
			AVG(DATEDIFF(stay_end_at, stay_start_at)) as average_stay_days
		`, dateFormat).
		Where("stay_start_at >= ? AND stay_end_at <= ?", startDate, endDate).
		Where("status IN ?", []models.ReservationStatus{models.ReservationStatusNormal, models.ReservationStatusPending}).
		Where("deleted_at = ?", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
		Group("period").
		Order("period").
		Scan(&stats).Error

	return stats, err
}
