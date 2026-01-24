package repositories

import (
	"context"
	"strings"
	"time"

	appContext "gitlab.bellsoft.net/rms/api-core/internal/context"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, room *models.Room) (*models.Room, error)
	Update(ctx context.Context, room *models.Room) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Room, error)
	FindByIDWithGroup(ctx context.Context, id uint) (*models.Room, error)
	FindAll(ctx context.Context, filter RoomFilter, offset, limit int, sort string) ([]models.Room, int64, error)
	FindAvailableRooms(ctx context.Context, startDate, endDate time.Time, excludeReservationID *uint) ([]models.Room, error)
	ExistsByNumber(ctx context.Context, number string, excludeID *uint) (bool, error)
	IsRoomAvailable(ctx context.Context, roomID uint, startDate, endDate time.Time, excludeReservationID *uint) (bool, error)
	FindByNumber(ctx context.Context, number string) (*models.Room, error)
	FindByStatus(ctx context.Context, status models.RoomStatus) ([]models.Room, error)
}

type RoomFilter struct {
	RoomGroupID *uint
	Status      *models.RoomStatus
	Search      string
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *models.Room) (*models.Room, error) {
	err := r.db.WithContext(ctx).Create(room).Error
	return room, err
}

func (r *roomRepository) Update(ctx context.Context, room *models.Room) error {
	return r.db.WithContext(ctx).Save(room).Error
}

func (r *roomRepository) Delete(ctx context.Context, id uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"deleted_at": now,
	}

	// Set updated by field for soft delete
	if userID, ok := appContext.GetUserID(ctx); ok {
		updates["updated_by"] = userID
	}

	return r.db.WithContext(ctx).Model(&models.Room{}).Where("id = ?", id).Updates(updates).Error
}

func (r *roomRepository) FindByID(ctx context.Context, id uint) (*models.Room, error) {
	var room models.Room
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, defaultDeletedAt).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) FindByIDWithGroup(ctx context.Context, id uint) (*models.Room, error) {
	var room models.Room
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Preload("RoomGroup", "deleted_at = ?", defaultDeletedAt).Where("id = ? AND deleted_at = ?", id, defaultDeletedAt).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) FindAll(ctx context.Context, filter RoomFilter, offset, limit int, sort string) ([]models.Room, int64, error) {
	var rooms []models.Room
	var total int64

	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	query := r.db.WithContext(ctx).Model(&models.Room{}).Where("deleted_at = ?", defaultDeletedAt).Preload("RoomGroup", "deleted_at = ?", defaultDeletedAt)

	if filter.RoomGroupID != nil {
		query = query.Where("room_group_id = ?", *filter.RoomGroupID)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Search != "" {
		query = query.Where("number LIKE ?", "%"+filter.Search+"%")
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
	query = query.Order(orderClause)

	err = query.Offset(offset).Limit(limit).Find(&rooms).Error
	if err != nil {
		return nil, 0, err
	}

	return rooms, total, nil
}

func (r *roomRepository) FindAvailableRooms(ctx context.Context, startDate, endDate time.Time, excludeReservationID *uint) ([]models.Room, error) {
	var rooms []models.Room

	subQuery := r.db.Model(&models.ReservationRoom{}).
		Select("room_id").
		Joins("JOIN reservation ON reservation.id = reservation_room.reservation_id").
		Where("reservation_room.deleted_at = ?", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
		Where("reservation.deleted_at = ?", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
		Where("reservation.status IN ?", []models.ReservationStatus{models.ReservationStatusNormal, models.ReservationStatusPending}).
		Where("NOT (reservation.stay_end_at <= ? OR reservation.stay_start_at >= ?)", startDate, endDate)

	if excludeReservationID != nil {
		subQuery = subQuery.Where("reservation.id != ?", *excludeReservationID)
	}

	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).
		Preload("RoomGroup", "deleted_at = ?", defaultDeletedAt).
		Where("status = ? AND deleted_at = ?", models.RoomStatusNormal, defaultDeletedAt).
		Where("id NOT IN (?)", subQuery).
		Order("room_group_id, number").
		Find(&rooms).Error

	return rooms, err
}

func (r *roomRepository) ExistsByNumber(ctx context.Context, number string, excludeID *uint) (bool, error) {
	var count int64
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	query := r.db.WithContext(ctx).Model(&models.Room{}).Where("number = ? AND deleted_at = ?", number, defaultDeletedAt)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *roomRepository) IsRoomAvailable(ctx context.Context, roomID uint, startDate, endDate time.Time, excludeReservationID *uint) (bool, error) {
	var count int64

	query := r.db.WithContext(ctx).
		Model(&models.ReservationRoom{}).
		Joins("JOIN reservation ON reservation.id = reservation_room.reservation_id").
		Where("reservation_room.room_id = ?", roomID).
		Where("reservation_room.deleted_at = ?", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
		Where("reservation.deleted_at = ?", time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)).
		Where("reservation.status IN ?", []models.ReservationStatus{models.ReservationStatusNormal, models.ReservationStatusPending}).
		Where("NOT (reservation.stay_end_at <= ? OR reservation.stay_start_at >= ?)", startDate, endDate)

	if excludeReservationID != nil {
		query = query.Where("reservation.id != ?", *excludeReservationID)
	}

	err := query.Count(&count).Error
	return count == 0, err
}

func (r *roomRepository) FindByNumber(ctx context.Context, number string) (*models.Room, error) {
	var room models.Room
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Where("number = ? AND deleted_at = ?", number, defaultDeletedAt).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) FindByStatus(ctx context.Context, status models.RoomStatus) ([]models.Room, error) {
	var rooms []models.Room
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).
		Preload("RoomGroup", "deleted_at = ?", defaultDeletedAt).
		Where("status = ? AND deleted_at = ?", status, defaultDeletedAt).
		Order("room_group_id, number").
		Find(&rooms).Error
	return rooms, err
}

// parseSort는 Spring Boot 형식의 정렬 파라미터를 GORM 형식으로 변환합니다.
// 예: "number,desc" -> "number DESC"
// 예: "number,desc,roomGroupId,asc" -> "number DESC, room_group_id ASC"
func (r *roomRepository) parseSort(sort string) string {
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
func (r *roomRepository) mapSortField(field string) string {
	fieldMap := map[string]string{
		"id":          "id",
		"number":      "number",
		"roomGroupId": "room_group_id",
		"note":        "note",
		"status":      "status",
		"createdAt":   "created_at",
		"updatedAt":   "updated_at",
	}

	if dbField, ok := fieldMap[field]; ok {
		return dbField
	}
	return ""
}
