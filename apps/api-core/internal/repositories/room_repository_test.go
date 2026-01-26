package repositories_test

import (
	"context"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type RoomRepositoryTestSuite struct {
	suite.Suite
	ctx  context.Context
	db   *gorm.DB
	mock sqlmock.Sqlmock
	repo repositories.RoomRepository
}

func (suite *RoomRepositoryTestSuite) SetupTest() {
	suite.ctx = context.Background()

	// Create mock DB
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)

	suite.mock = mock

	// Create GORM DB with mock
	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	suite.db, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(suite.T(), err)

	suite.repo = repositories.NewRoomRepository(suite.db)
}

func (suite *RoomRepositoryTestSuite) TearDownTest() {
	sqlDB, err := suite.db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// =====================================================
// Create Tests
// =====================================================

func (suite *RoomRepositoryTestSuite) TestCreate() {
	suite.Run("객실을 생성하면 ID가 할당된다", func() {
		// Given
		room := &models.Room{
			Number:      "101",
			RoomGroupID: 1,
			Note:        "Test room",
			Status:      models.RoomStatusNormal,
		}

		// Mock expectations
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `room`")).
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				sqlmock.AnyArg(), // created_by
				sqlmock.AnyArg(), // updated_by
				room.Number,
				room.RoomGroupID,
				room.Note,
				room.Status,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		suite.mock.ExpectCommit()

		// When
		result, err := suite.repo.Create(suite.ctx, room)

		// Then
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), result)
		assert.Equal(suite.T(), uint(1), result.ID)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("객실 생성 시 기본값이 설정된다", func() {
		// Given
		room := &models.Room{
			Number:      "102",
			RoomGroupID: 1,
			// Note와 Status는 기본값 사용
		}

		// Mock expectations
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `room`")).
			WillReturnResult(sqlmock.NewResult(2, 1))
		suite.mock.ExpectCommit()

		// When
		result, err := suite.repo.Create(suite.ctx, room)

		// Then
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), result)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})
}

// =====================================================
// FindByID Tests
// =====================================================

func (suite *RoomRepositoryTestSuite) TestFindByID() {
	suite.Run("ID로 객실을 조회하면 객실 정보가 반환된다", func() {
		// Given
		roomID := uint(1)
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		rows := sqlmock.NewRows([]string{
			"id", "number", "room_group_id", "note", "status",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(
			1, "101", 1, "Test room", models.RoomStatusNormal,
			time.Now(), time.Now(), defaultDeletedAt, 1, 1,
		)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room` WHERE id = ? AND deleted_at = ? ORDER BY `room`.`id` LIMIT ?")).
			WithArgs(roomID, defaultDeletedAt, 1).
			WillReturnRows(rows)

		// When
		result, err := suite.repo.FindByID(suite.ctx, roomID)

		// Then
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), result)
		assert.Equal(suite.T(), roomID, result.ID)
		assert.Equal(suite.T(), "101", result.Number)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("존재하지 않는 ID로 조회하면 에러가 발생한다", func() {
		// Given
		roomID := uint(999)
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room` WHERE id = ? AND deleted_at = ? ORDER BY `room`.`id` LIMIT ?")).
			WithArgs(roomID, defaultDeletedAt, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// When
		result, err := suite.repo.FindByID(suite.ctx, roomID)

		// Then
		assert.Error(suite.T(), err)
		assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
		assert.Nil(suite.T(), result)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("삭제된 객실은 조회되지 않는다", func() {
		// Given
		roomID := uint(1)
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		// 삭제된 객실은 deleted_at이 기본값이 아니므로 조회되지 않음
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room` WHERE id = ? AND deleted_at = ? ORDER BY `room`.`id` LIMIT ?")).
			WithArgs(roomID, defaultDeletedAt, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// When
		result, err := suite.repo.FindByID(suite.ctx, roomID)

		// Then
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), result)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})
}

// =====================================================
// FindAll Tests
// =====================================================

func (suite *RoomRepositoryTestSuite) TestFindAll() {
	suite.Run("전체 객실 목록을 조회하면 모든 객실이 반환된다", func() {
		// Given
		filter := dto.RoomRepositoryFilter{}
		offset := 0
		limit := 10
		sort := ""
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		// Count query
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `room` WHERE deleted_at = ?")).
			WithArgs(defaultDeletedAt).
			WillReturnRows(countRows)

		// Data query with preload
		dataRows := sqlmock.NewRows([]string{
			"id", "number", "room_group_id", "note", "status",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).
			AddRow(1, "101", 1, "Room 1", models.RoomStatusNormal, time.Now(), time.Now(), defaultDeletedAt, 1, 1).
			AddRow(2, "102", 1, "Room 2", models.RoomStatusNormal, time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room` WHERE deleted_at = ? ORDER BY id DESC LIMIT ?")).
			WithArgs(defaultDeletedAt, limit).
			WillReturnRows(dataRows)

		// Preload RoomGroup query
		roomGroupRows := sqlmock.NewRows([]string{
			"id", "name", "peek_price", "off_peek_price", "description",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(1, "Standard", 100000, 80000, "Standard room", time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room_group` WHERE `room_group`.`id` = ? AND deleted_at = ?")).
			WithArgs(1, defaultDeletedAt).
			WillReturnRows(roomGroupRows)

		// When
		rooms, total, err := suite.repo.FindAll(suite.ctx, filter, offset, limit, sort)

		// Then
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), rooms, 2)
		assert.Equal(suite.T(), int64(2), total)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("페이지네이션이 적용된다", func() {
		// Given
		filter := dto.RoomRepositoryFilter{}
		offset := 10
		limit := 5
		sort := ""
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		// Count query
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(15)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `room` WHERE deleted_at = ?")).
			WithArgs(defaultDeletedAt).
			WillReturnRows(countRows)

		// Data query with offset and limit
		dataRows := sqlmock.NewRows([]string{
			"id", "number", "room_group_id", "note", "status",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).
			AddRow(11, "111", 1, "Room 11", models.RoomStatusNormal, time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room` WHERE deleted_at = ? ORDER BY id DESC LIMIT ? OFFSET ?")).
			WithArgs(defaultDeletedAt, limit, offset).
			WillReturnRows(dataRows)

		// Preload RoomGroup query
		roomGroupRows := sqlmock.NewRows([]string{
			"id", "name", "peek_price", "off_peek_price", "description",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(1, "Standard", 100000, 80000, "Standard room", time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room_group` WHERE `room_group`.`id` = ? AND deleted_at = ?")).
			WithArgs(1, defaultDeletedAt).
			WillReturnRows(roomGroupRows)

		// When
		rooms, total, err := suite.repo.FindAll(suite.ctx, filter, offset, limit, sort)

		// Then
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), rooms, 1)
		assert.Equal(suite.T(), int64(15), total)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("필터링이 적용된다 - RoomGroupID", func() {
		// Given
		roomGroupID := uint(2)
		filter := dto.RoomRepositoryFilter{
			RoomGroupID: &roomGroupID,
		}
		offset := 0
		limit := 10
		sort := ""
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		// Count query with filter
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `room` WHERE deleted_at = ? AND room_group_id = ?")).
			WithArgs(defaultDeletedAt, roomGroupID).
			WillReturnRows(countRows)

		// Data query with filter
		dataRows := sqlmock.NewRows([]string{
			"id", "number", "room_group_id", "note", "status",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(1, "201", 2, "Deluxe room", models.RoomStatusNormal, time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room` WHERE deleted_at = ? AND room_group_id = ? ORDER BY id DESC LIMIT ?")).
			WithArgs(defaultDeletedAt, roomGroupID, limit).
			WillReturnRows(dataRows)

		// Preload RoomGroup query
		roomGroupRows := sqlmock.NewRows([]string{
			"id", "name", "peek_price", "off_peek_price", "description",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(2, "Deluxe", 150000, 120000, "Deluxe room", time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room_group` WHERE `room_group`.`id` = ? AND deleted_at = ?")).
			WithArgs(2, defaultDeletedAt).
			WillReturnRows(roomGroupRows)

		// When
		rooms, total, err := suite.repo.FindAll(suite.ctx, filter, offset, limit, sort)

		// Then
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), rooms, 1)
		assert.Equal(suite.T(), int64(1), total)
		assert.Equal(suite.T(), uint(2), rooms[0].RoomGroupID)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("필터링이 적용된다 - Status", func() {
		// Given
		status := models.RoomStatusNormal
		filter := dto.RoomRepositoryFilter{
			Status: &status,
		}
		offset := 0
		limit := 10
		sort := ""
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		// Count query with filter
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `room` WHERE deleted_at = ? AND status = ?")).
			WithArgs(defaultDeletedAt, status).
			WillReturnRows(countRows)

		// Data query with filter
		dataRows := sqlmock.NewRows([]string{
			"id", "number", "room_group_id", "note", "status",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(1, "101", 1, "Normal room", models.RoomStatusNormal, time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room` WHERE deleted_at = ? AND status = ? ORDER BY id DESC LIMIT ?")).
			WithArgs(defaultDeletedAt, status, limit).
			WillReturnRows(dataRows)

		// Preload RoomGroup query
		roomGroupRows := sqlmock.NewRows([]string{
			"id", "name", "peek_price", "off_peek_price", "description",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(1, "Standard", 100000, 80000, "Standard room", time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room_group` WHERE `room_group`.`id` = ? AND deleted_at = ?")).
			WithArgs(1, defaultDeletedAt).
			WillReturnRows(roomGroupRows)

		// When
		rooms, total, err := suite.repo.FindAll(suite.ctx, filter, offset, limit, sort)

		// Then
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), rooms, 1)
		assert.Equal(suite.T(), int64(1), total)
		assert.Equal(suite.T(), models.RoomStatusNormal, rooms[0].Status)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("필터링이 적용된다 - Search", func() {
		// Given
		filter := dto.RoomRepositoryFilter{
			Search: "101",
		}
		offset := 0
		limit := 10
		sort := ""
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		// Count query with search filter
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `room` WHERE deleted_at = ? AND number LIKE ?")).
			WithArgs(defaultDeletedAt, "%101%").
			WillReturnRows(countRows)

		// Data query with search filter
		dataRows := sqlmock.NewRows([]string{
			"id", "number", "room_group_id", "note", "status",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(1, "101", 1, "Room 101", models.RoomStatusNormal, time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room` WHERE deleted_at = ? AND number LIKE ? ORDER BY id DESC LIMIT ?")).
			WithArgs(defaultDeletedAt, "%101%", limit).
			WillReturnRows(dataRows)

		// Preload RoomGroup query
		roomGroupRows := sqlmock.NewRows([]string{
			"id", "name", "peek_price", "off_peek_price", "description",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(1, "Standard", 100000, 80000, "Standard room", time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room_group` WHERE `room_group`.`id` = ? AND deleted_at = ?")).
			WithArgs(1, defaultDeletedAt).
			WillReturnRows(roomGroupRows)

		// When
		rooms, total, err := suite.repo.FindAll(suite.ctx, filter, offset, limit, sort)

		// Then
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), rooms, 1)
		assert.Equal(suite.T(), int64(1), total)
		assert.Contains(suite.T(), rooms[0].Number, "101")
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("정렬이 적용된다", func() {
		// Given
		filter := dto.RoomRepositoryFilter{}
		offset := 0
		limit := 10
		sort := "number,asc"
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		// Count query
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `room` WHERE deleted_at = ?")).
			WithArgs(defaultDeletedAt).
			WillReturnRows(countRows)

		// Data query with custom sort
		dataRows := sqlmock.NewRows([]string{
			"id", "number", "room_group_id", "note", "status",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).
			AddRow(1, "101", 1, "Room 1", models.RoomStatusNormal, time.Now(), time.Now(), defaultDeletedAt, 1, 1).
			AddRow(2, "102", 1, "Room 2", models.RoomStatusNormal, time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room` WHERE deleted_at = ? ORDER BY number ASC LIMIT ?")).
			WithArgs(defaultDeletedAt, limit).
			WillReturnRows(dataRows)

		// Preload RoomGroup query
		roomGroupRows := sqlmock.NewRows([]string{
			"id", "name", "peek_price", "off_peek_price", "description",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(1, "Standard", 100000, 80000, "Standard room", time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room_group` WHERE `room_group`.`id` = ? AND deleted_at = ?")).
			WithArgs(1, defaultDeletedAt).
			WillReturnRows(roomGroupRows)

		// When
		rooms, total, err := suite.repo.FindAll(suite.ctx, filter, offset, limit, sort)

		// Then
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), rooms, 2)
		assert.Equal(suite.T(), int64(2), total)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})
}

// =====================================================
// Update Tests
// =====================================================

func (suite *RoomRepositoryTestSuite) TestUpdate() {
	suite.Run("객실 정보를 수정하면 업데이트된다", func() {
		// Given
		room := &models.Room{
			Number:      "101-UPDATED",
			RoomGroupID: 1,
			Note:        "Updated note",
			Status:      models.RoomStatusNormal,
		}
		room.ID = 1

		// Mock expectations
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `room`")).
			WithArgs(
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				sqlmock.AnyArg(), // deleted_at
				sqlmock.AnyArg(), // created_by
				sqlmock.AnyArg(), // updated_by
				room.Number,
				room.RoomGroupID,
				room.Note,
				room.Status,
				room.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		suite.mock.ExpectCommit()

		// When
		err := suite.repo.Update(suite.ctx, room)

		// Then
		assert.NoError(suite.T(), err)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("객실 상태를 변경할 수 있다", func() {
		// Given
		room := &models.Room{
			Number:      "101",
			RoomGroupID: 1,
			Note:        "Test room",
			Status:      models.RoomStatusInactive,
		}
		room.ID = 1

		// Mock expectations
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `room`")).
			WillReturnResult(sqlmock.NewResult(1, 1))
		suite.mock.ExpectCommit()

		// When
		err := suite.repo.Update(suite.ctx, room)

		// Then
		assert.NoError(suite.T(), err)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})
}

// =====================================================
// Delete Tests (Soft Delete)
// =====================================================

func (suite *RoomRepositoryTestSuite) TestDelete() {
	suite.Run("객실을 삭제하면 soft delete가 적용된다", func() {
		// Given
		roomID := uint(1)

		// Mock expectations - soft delete는 UPDATE 쿼리 사용
		// GORM Updates는 updated_at도 자동으로 업데이트함
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `room` SET `deleted_at`=?,`updated_at`=? WHERE id = ?")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), roomID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		suite.mock.ExpectCommit()

		// When
		err := suite.repo.Delete(suite.ctx, roomID)

		// Then
		assert.NoError(suite.T(), err)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("삭제 시 deleted_at이 현재 시간으로 설정된다", func() {
		// Given
		roomID := uint(2)

		// Mock expectations
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `room` SET `deleted_at`=?,`updated_at`=? WHERE id = ?")).
			WithArgs(
				sqlmock.AnyArg(), // deleted_at (현재 시간)
				sqlmock.AnyArg(), // updated_at (자동 업데이트)
				roomID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		suite.mock.ExpectCommit()

		// When
		err := suite.repo.Delete(suite.ctx, roomID)

		// Then
		assert.NoError(suite.T(), err)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("존재하지 않는 객실 삭제 시에도 에러가 발생하지 않는다", func() {
		// Given
		roomID := uint(999)

		// Mock expectations - 0 rows affected
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(regexp.QuoteMeta("UPDATE `room` SET `deleted_at`=?,`updated_at`=? WHERE id = ?")).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), roomID).
			WillReturnResult(sqlmock.NewResult(0, 0))
		suite.mock.ExpectCommit()

		// When
		err := suite.repo.Delete(suite.ctx, roomID)

		// Then
		assert.NoError(suite.T(), err)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})
}

// =====================================================
// Additional Helper Methods Tests
// =====================================================

func (suite *RoomRepositoryTestSuite) TestExistsByNumber() {
	suite.Run("객실 번호가 존재하면 true를 반환한다", func() {
		// Given
		number := "101"
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `room` WHERE number = ? AND deleted_at = ?")).
			WithArgs(number, defaultDeletedAt).
			WillReturnRows(countRows)

		// When
		exists, err := suite.repo.ExistsByNumber(suite.ctx, number, nil)

		// Then
		assert.NoError(suite.T(), err)
		assert.True(suite.T(), exists)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("객실 번호가 존재하지 않으면 false를 반환한다", func() {
		// Given
		number := "999"
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `room` WHERE number = ? AND deleted_at = ?")).
			WithArgs(number, defaultDeletedAt).
			WillReturnRows(countRows)

		// When
		exists, err := suite.repo.ExistsByNumber(suite.ctx, number, nil)

		// Then
		assert.NoError(suite.T(), err)
		assert.False(suite.T(), exists)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})

	suite.Run("특정 ID를 제외하고 중복 체크를 수행한다", func() {
		// Given
		number := "101"
		excludeID := uint(1)
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		// GORM wraps WHERE conditions in parentheses
		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `room` WHERE (number = ? AND deleted_at = ?) AND id != ?")).
			WithArgs(number, defaultDeletedAt, excludeID).
			WillReturnRows(countRows)

		// When
		exists, err := suite.repo.ExistsByNumber(suite.ctx, number, &excludeID)

		// Then
		assert.NoError(suite.T(), err)
		assert.False(suite.T(), exists)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})
}

func (suite *RoomRepositoryTestSuite) TestFindByNumber() {
	suite.Run("객실 번호로 조회하면 객실 정보가 반환된다", func() {
		// Given
		number := "101"
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		rows := sqlmock.NewRows([]string{
			"id", "number", "room_group_id", "note", "status",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(
			1, "101", 1, "Test room", models.RoomStatusNormal,
			time.Now(), time.Now(), defaultDeletedAt, 1, 1,
		)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room` WHERE number = ? AND deleted_at = ? ORDER BY `room`.`id` LIMIT ?")).
			WithArgs(number, defaultDeletedAt, 1).
			WillReturnRows(rows)

		// When
		result, err := suite.repo.FindByNumber(suite.ctx, number)

		// Then
		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), result)
		assert.Equal(suite.T(), number, result.Number)
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})
}

func (suite *RoomRepositoryTestSuite) TestFindByStatus() {
	suite.Run("상태별로 객실을 조회하면 해당 상태의 객실들이 반환된다", func() {
		// Given
		status := models.RoomStatusNormal
		defaultDeletedAt := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

		rows := sqlmock.NewRows([]string{
			"id", "number", "room_group_id", "note", "status",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).
			AddRow(1, "101", 1, "Room 1", models.RoomStatusNormal, time.Now(), time.Now(), defaultDeletedAt, 1, 1).
			AddRow(2, "102", 1, "Room 2", models.RoomStatusNormal, time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room` WHERE status = ? AND deleted_at = ? ORDER BY room_group_id, number")).
			WithArgs(status, defaultDeletedAt).
			WillReturnRows(rows)

		// Preload RoomGroup query
		roomGroupRows := sqlmock.NewRows([]string{
			"id", "name", "peek_price", "off_peek_price", "description",
			"created_at", "updated_at", "deleted_at", "created_by", "updated_by",
		}).AddRow(1, "Standard", 100000, 80000, "Standard room", time.Now(), time.Now(), defaultDeletedAt, 1, 1)

		suite.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `room_group` WHERE `room_group`.`id` = ? AND deleted_at = ?")).
			WithArgs(1, defaultDeletedAt).
			WillReturnRows(roomGroupRows)

		// When
		rooms, err := suite.repo.FindByStatus(suite.ctx, status)

		// Then
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), rooms, 2)
		for _, room := range rooms {
			assert.Equal(suite.T(), status, room.Status)
		}
		assert.NoError(suite.T(), suite.mock.ExpectationsWereMet())
	})
}

// AnyTime is a custom matcher for time.Time values in sqlmock
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestRoomRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RoomRepositoryTestSuite))
}
