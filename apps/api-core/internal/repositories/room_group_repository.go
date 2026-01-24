package repositories

import (
	"context"
	"strings"
	"time"

	appContext "gitlab.bellsoft.net/rms/api-core/internal/context"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gorm.io/gorm"
)

type RoomGroupRoomFilter struct {
	RoomStatus           *models.RoomStatus
	StayStartAt          *time.Time
	StayEndAt            *time.Time
	ExcludeReservationID *uint
}

type RoomGroupRepository interface {
	Create(ctx context.Context, roomGroup *models.RoomGroup) (*models.RoomGroup, error)
	Update(ctx context.Context, roomGroup *models.RoomGroup) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.RoomGroup, error)
	FindByIDWithRooms(ctx context.Context, id uint, roomStatus *models.RoomStatus) (*models.RoomGroup, error)
	FindByIDWithFilteredRooms(ctx context.Context, id uint, filter RoomGroupRoomFilter) (*models.RoomGroup, error)
	FindAll(ctx context.Context, offset, limit int) ([]models.RoomGroup, int64, error)
	FindByIDWithUsers(ctx context.Context, id uint) (*models.RoomGroup, error)
	FindAllWithUsers(ctx context.Context, offset, limit int, sort string) ([]models.RoomGroup, int64, error)
	ExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error)
	FindByName(ctx context.Context, name string) (*models.RoomGroup, error)
}

type roomGroupRepository struct {
	db *gorm.DB
}

func NewRoomGroupRepository(db *gorm.DB) RoomGroupRepository {
	return &roomGroupRepository{db: db}
}

func (r *roomGroupRepository) Create(ctx context.Context, roomGroup *models.RoomGroup) (*models.RoomGroup, error) {
	err := r.db.WithContext(ctx).Create(roomGroup).Error
	return roomGroup, err
}

func (r *roomGroupRepository) Update(ctx context.Context, roomGroup *models.RoomGroup) error {
	return r.db.WithContext(ctx).Save(roomGroup).Error
}

func (r *roomGroupRepository) Delete(ctx context.Context, id uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"deleted_at": now,
	}

	// Set updated by field for soft delete
	if userID, ok := appContext.GetUserID(ctx); ok {
		updates["updated_by"] = userID
	}

	return r.db.WithContext(ctx).Model(&models.RoomGroup{}).Where("id = ?", id).Updates(updates).Error
}

func (r *roomGroupRepository) FindByID(ctx context.Context, id uint) (*models.RoomGroup, error) {
	var roomGroup models.RoomGroup
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, defaultDeletedAt).First(&roomGroup).Error
	if err != nil {
		return nil, err
	}
	return &roomGroup, nil
}

func (r *roomGroupRepository) FindByIDWithRooms(ctx context.Context, id uint, roomStatus *models.RoomStatus) (*models.RoomGroup, error) {
	var roomGroup models.RoomGroup
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	query := r.db.WithContext(ctx)

	if roomStatus != nil {
		query = query.Preload("Rooms", "status = ? AND deleted_at = ?", *roomStatus, defaultDeletedAt)
	} else {
		query = query.Preload("Rooms", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at = ?", defaultDeletedAt).Order("id")
		})
	}

	err := query.Where("id = ? AND deleted_at = ?", id, defaultDeletedAt).First(&roomGroup).Error
	if err != nil {
		return nil, err
	}
	return &roomGroup, nil
}

func (r *roomGroupRepository) FindByIDWithFilteredRooms(ctx context.Context, id uint, filter RoomGroupRoomFilter) (*models.RoomGroup, error) {
	var roomGroup models.RoomGroup
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at = ?", id, defaultDeletedAt).First(&roomGroup).Error
	if err != nil {
		return nil, err
	}

	var rooms []models.Room
	roomQuery := r.db.WithContext(ctx).Where("room_group_id = ? AND deleted_at = ?", id, defaultDeletedAt)

	if filter.RoomStatus != nil {
		roomQuery = roomQuery.Where("status = ?", *filter.RoomStatus)
	}

	if filter.StayStartAt != nil && filter.StayEndAt != nil {
		reservedRoomIDs := r.db.Model(&models.ReservationRoom{}).
			Select("room_id").
			Joins("JOIN reservation ON reservation.id = reservation_room.reservation_id").
			Where("reservation_room.deleted_at = ?", defaultDeletedAt).
			Where("reservation.deleted_at = ?", defaultDeletedAt).
			Where("reservation.status IN ?", []models.ReservationStatus{models.ReservationStatusNormal, models.ReservationStatusPending}).
			Where("NOT (reservation.stay_end_at <= ? OR reservation.stay_start_at >= ?)", *filter.StayStartAt, *filter.StayEndAt)

		if filter.ExcludeReservationID != nil {
			reservedRoomIDs = reservedRoomIDs.Where("reservation.id != ?", *filter.ExcludeReservationID)
		}

		roomQuery = roomQuery.Where("id NOT IN (?)", reservedRoomIDs)
	}

	err = roomQuery.Order("id").Find(&rooms).Error
	if err != nil {
		return nil, err
	}

	roomGroup.Rooms = rooms
	return &roomGroup, nil
}

func (r *roomGroupRepository) FindAll(ctx context.Context, offset, limit int) ([]models.RoomGroup, int64, error) {
	var roomGroups []models.RoomGroup
	var total int64

	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	query := r.db.WithContext(ctx).Model(&models.RoomGroup{}).Where("deleted_at = ?", defaultDeletedAt)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(limit).Order("id").Find(&roomGroups).Error
	if err != nil {
		return nil, 0, err
	}

	return roomGroups, total, nil
}

func (r *roomGroupRepository) ExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error) {
	var count int64
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	query := r.db.WithContext(ctx).Model(&models.RoomGroup{}).Where("name = ? AND deleted_at = ?", name, defaultDeletedAt)

	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *roomGroupRepository) FindByIDWithUsers(ctx context.Context, id uint) (*models.RoomGroup, error) {
	var roomGroup models.RoomGroup
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("id = ? AND deleted_at = ?", id, defaultDeletedAt).
		First(&roomGroup).Error
	if err != nil {
		return nil, err
	}
	return &roomGroup, nil
}

func (r *roomGroupRepository) FindAllWithUsers(ctx context.Context, offset, limit int, sort string) ([]models.RoomGroup, int64, error) {
	var roomGroups []models.RoomGroup
	var total int64

	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	query := r.db.WithContext(ctx).Model(&models.RoomGroup{}).Where("deleted_at = ?", defaultDeletedAt)

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

	err = query.
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Offset(offset).Limit(limit).Find(&roomGroups).Error
	if err != nil {
		return nil, 0, err
	}

	return roomGroups, total, nil
}

func (r *roomGroupRepository) FindByName(ctx context.Context, name string) (*models.RoomGroup, error) {
	var roomGroup models.RoomGroup
	defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	err := r.db.WithContext(ctx).Where("name = ? AND deleted_at = ?", name, defaultDeletedAt).First(&roomGroup).Error
	if err != nil {
		return nil, err
	}
	return &roomGroup, nil
}

// parseSort는 Spring Boot 형식의 정렬 파라미터를 GORM 형식으로 변환합니다.
// 예: "name,desc" -> "name DESC"
// 예: "name,desc,peekPrice,asc" -> "name DESC, peek_price ASC"
func (r *roomGroupRepository) parseSort(sort string) string {
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
func (r *roomGroupRepository) mapSortField(field string) string {
	fieldMap := map[string]string{
		"id":           "id",
		"name":         "name",
		"peekPrice":    "peek_price",
		"offPeekPrice": "off_peek_price",
		"description":  "description",
		"createdAt":    "created_at",
		"updatedAt":    "updated_at",
	}

	if dbField, ok := fieldMap[field]; ok {
		return dbField
	}
	return ""
}
