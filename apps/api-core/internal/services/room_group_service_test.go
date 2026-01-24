package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gorm.io/gorm"
)

type RoomGroupServiceTestSuite struct {
	suite.Suite
	ctx      context.Context
	service  services.RoomGroupService
	mockRepo *MockRoomGroupRepository
}

func (suite *RoomGroupServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockRepo = new(MockRoomGroupRepository)
	suite.service = services.NewRoomGroupService(suite.mockRepo)
}

func (suite *RoomGroupServiceTestSuite) TestGetByID() {
	// Given - 객실 그룹이 등록된 상황에서
	roomGroup := &models.RoomGroup{
		Name:         "Standard",
		PeekPrice:    100000,
		OffPeekPrice: 80000,
		Description:  "Standard room group",
	}
	roomGroup.ID = 1

	suite.mockRepo.On("FindByID", suite.ctx, uint(1)).Return(roomGroup, nil)

	// When - 특정 객실 그룹 정보를 조회하면
	result, err := suite.service.GetByID(suite.ctx, 1)

	// Then - 정상적으로 조회된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), roomGroup.Name, result.Name)
	assert.Equal(suite.T(), roomGroup.PeekPrice, result.PeekPrice)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestGetByID_NotFound() {
	// Given - 존재하지 않는 객실 그룹 ID로
	suite.mockRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// When - 조회하면
	result, err := suite.service.GetByID(suite.ctx, 999)

	// Then - ErrRoomGroupNotFound 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomGroupNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestGetByIDWithRooms() {
	// Given - 객실 그룹과 객실들이 있는 상황에서
	rooms := []models.Room{
		{
			Number: "101",
			Status: models.RoomStatusNormal,
		},
		{
			Number: "102",
			Status: models.RoomStatusNormal,
		},
	}
	rooms[0].ID = 1
	rooms[1].ID = 2

	roomGroup := &models.RoomGroup{
		Name:         "Standard",
		PeekPrice:    100000,
		OffPeekPrice: 80000,
		Description:  "Standard room group",
		Rooms:        rooms,
	}
	roomGroup.ID = 1

	statusNormal := models.RoomStatusNormal
	suite.mockRepo.On("FindByIDWithRooms", suite.ctx, uint(1), &statusNormal).Return(roomGroup, nil)

	// When - 활성 상태의 객실과 함께 조회하면
	result, err := suite.service.GetByIDWithRooms(suite.ctx, 1, &statusNormal)

	// Then - 객실 그룹과 활성 객실이 함께 조회된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), roomGroup.Name, result.Name)
	assert.Len(suite.T(), result.Rooms, 2)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestGetByIDWithFilteredRooms_날짜필터링으로_예약가능한_객실만_조회() {
	// Given - 객실 그룹에 2개의 객실이 있고, 1개만 해당 기간에 예약 가능한 상황
	stayStartAt := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
	stayEndAt := time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC)
	statusNormal := models.RoomStatusNormal

	availableRooms := []models.Room{
		{
			Number: "502",
			Status: models.RoomStatusNormal,
		},
	}
	availableRooms[0].ID = 2

	roomGroup := &models.RoomGroup{
		Name:         "Standard",
		PeekPrice:    100000,
		OffPeekPrice: 80000,
		Description:  "Standard room group",
		Rooms:        availableRooms,
	}
	roomGroup.ID = 9

	filter := repositories.RoomGroupRoomFilter{
		RoomStatus:  &statusNormal,
		StayStartAt: &stayStartAt,
		StayEndAt:   &stayEndAt,
	}

	suite.mockRepo.On("FindByIDWithFilteredRooms", suite.ctx, uint(9), filter).Return(roomGroup, nil)

	// When - 날짜 필터와 함께 조회하면
	result, err := suite.service.GetByIDWithFilteredRooms(suite.ctx, 9, filter)

	// Then - 예약 가능한 객실만 포함되어 조회된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result.Rooms, 1)
	assert.Equal(suite.T(), "502", result.Rooms[0].Number)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestGetByIDWithFilteredRooms_예약수정시_현재예약_제외() {
	// Given - 현재 예약의 객실도 포함하여 조회해야 하는 상황 (예약 수정 시)
	stayStartAt := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
	stayEndAt := time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC)
	statusNormal := models.RoomStatusNormal
	excludeReservationID := uint(123)

	rooms := []models.Room{
		{
			Number: "501",
			Status: models.RoomStatusNormal,
		},
		{
			Number: "502",
			Status: models.RoomStatusNormal,
		},
	}
	rooms[0].ID = 1
	rooms[1].ID = 2

	roomGroup := &models.RoomGroup{
		Name:  "Standard",
		Rooms: rooms,
	}
	roomGroup.ID = 9

	filter := repositories.RoomGroupRoomFilter{
		RoomStatus:           &statusNormal,
		StayStartAt:          &stayStartAt,
		StayEndAt:            &stayEndAt,
		ExcludeReservationID: &excludeReservationID,
	}

	suite.mockRepo.On("FindByIDWithFilteredRooms", suite.ctx, uint(9), filter).Return(roomGroup, nil)

	// When - excludeReservationId와 함께 조회하면
	result, err := suite.service.GetByIDWithFilteredRooms(suite.ctx, 9, filter)

	// Then - 현재 예약을 제외하고 사용 가능한 객실이 조회된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Len(suite.T(), result.Rooms, 2)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestGetByIDWithFilteredRooms_모든객실_예약불가() {
	// Given - 해당 기간에 모든 객실이 예약 불가능한 상황
	stayStartAt := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
	stayEndAt := time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC)
	statusNormal := models.RoomStatusNormal

	roomGroup := &models.RoomGroup{
		Name:  "Standard",
		Rooms: []models.Room{},
	}
	roomGroup.ID = 9

	filter := repositories.RoomGroupRoomFilter{
		RoomStatus:  &statusNormal,
		StayStartAt: &stayStartAt,
		StayEndAt:   &stayEndAt,
	}

	suite.mockRepo.On("FindByIDWithFilteredRooms", suite.ctx, uint(9), filter).Return(roomGroup, nil)

	// When - 날짜 필터와 함께 조회하면
	result, err := suite.service.GetByIDWithFilteredRooms(suite.ctx, 9, filter)

	// Then - 빈 객실 목록이 반환된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Empty(suite.T(), result.Rooms)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestGetByIDWithFilteredRooms_NotFound() {
	// Given - 존재하지 않는 객실 그룹
	stayStartAt := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
	stayEndAt := time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC)
	statusNormal := models.RoomStatusNormal

	filter := repositories.RoomGroupRoomFilter{
		RoomStatus:  &statusNormal,
		StayStartAt: &stayStartAt,
		StayEndAt:   &stayEndAt,
	}

	suite.mockRepo.On("FindByIDWithFilteredRooms", suite.ctx, uint(999), filter).Return(nil, gorm.ErrRecordNotFound)

	// When - 존재하지 않는 그룹 ID로 조회하면
	result, err := suite.service.GetByIDWithFilteredRooms(suite.ctx, 999, filter)

	// Then - ErrRoomGroupNotFound 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomGroupNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestGetAll_Empty() {
	// Given - 객실 그룹 정보가 없는 상황에서
	suite.mockRepo.On("FindAll", suite.ctx, 0, 10).Return([]models.RoomGroup{}, int64(0), nil)

	// When - 전체 객실 그룹 정보를 조회하면 (page 0부터 시작 - Spring Boot 호환)
	roomGroups, total, err := suite.service.GetAll(suite.ctx, 0, 10)

	// Then - 빈 객실 그룹 목록이 반환된다
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), roomGroups)
	assert.Equal(suite.T(), int64(0), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestGetAll_Multiple() {
	// Given - 객실 그룹이 10개 등록된 상황에서
	roomGroups := make([]models.RoomGroup, 10)
	for i := 0; i < 10; i++ {
		roomGroups[i] = models.RoomGroup{
			Name:         "Group" + string(rune(i)),
			PeekPrice:    100000 + i*10000,
			OffPeekPrice: 80000 + i*10000,
			Description:  "Test room group " + string(rune(i)),
		}
		roomGroups[i].ID = uint(i + 1)
	}

	suite.mockRepo.On("FindAll", suite.ctx, 0, 10).Return(roomGroups, int64(10), nil)

	// When - 전체 객실 그룹 정보를 조회하면 (page 0부터 시작 - Spring Boot 호환)
	result, total, err := suite.service.GetAll(suite.ctx, 0, 10)

	// Then - 10개의 객실 그룹 정보가 반환된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 10)
	assert.Equal(suite.T(), int64(10), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestGetByIDWithUsers() {
	// Given - 생성자/수정자 정보가 포함된 객실 그룹이 있는 상황에서
	createdBy := &models.User{
		UserID: "creator",
		Name:   "Creator User",
	}
	createdBy.ID = 1

	updatedBy := &models.User{
		UserID: "updater",
		Name:   "Updater User",
	}
	updatedBy.ID = 2

	roomGroup := &models.RoomGroup{
		Name:          "Standard",
		PeekPrice:     100000,
		OffPeekPrice:  80000,
		Description:   "Standard room group",
		CreatedByUser: createdBy,
		UpdatedByUser: updatedBy,
	}
	roomGroup.ID = 1
	roomGroup.CreatedBy = 1
	roomGroup.UpdatedBy = 2

	suite.mockRepo.On("FindByIDWithUsers", suite.ctx, uint(1)).Return(roomGroup, nil)

	// When - 사용자 정보와 함께 조회하면
	result, err := suite.service.GetByIDWithUsers(suite.ctx, 1)

	// Then - 객실 그룹과 사용자 정보가 함께 조회된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.NotNil(suite.T(), result.CreatedByUser)
	assert.NotNil(suite.T(), result.UpdatedByUser)
	assert.Equal(suite.T(), "creator", result.CreatedByUser.UserID)
	assert.Equal(suite.T(), "updater", result.UpdatedByUser.UserID)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestCreate() {
	// Given - 새로운 객실 그룹 정보가 주어지면
	newRoomGroup := &models.RoomGroup{
		Name:         "Deluxe",
		PeekPrice:    150000,
		OffPeekPrice: 120000,
		Description:  "Deluxe room group",
	}

	createdRoomGroup := &models.RoomGroup{
		Name:         "Deluxe",
		PeekPrice:    150000,
		OffPeekPrice: 120000,
		Description:  "Deluxe room group",
	}
	createdRoomGroup.ID = 1

	// 중복 체크
	suite.mockRepo.On("ExistsByName", suite.ctx, "Deluxe", (*uint)(nil)).Return(false, nil)
	// 생성
	suite.mockRepo.On("Create", suite.ctx, newRoomGroup).Return(createdRoomGroup, nil)

	// When - 객실 그룹을 생성하면
	err := suite.service.Create(suite.ctx, newRoomGroup)

	// Then - 정상적으로 생성된다
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestCreate_DuplicateName() {
	// Given - 이미 존재하는 이름으로 객실 그룹 생성 시도
	newRoomGroup := &models.RoomGroup{
		Name: "Standard",
	}

	// 중복 체크 - 이미 존재함
	suite.mockRepo.On("ExistsByName", suite.ctx, "Standard", (*uint)(nil)).Return(true, nil)

	// When - 동일한 이름으로 생성을 시도하면
	err := suite.service.Create(suite.ctx, newRoomGroup)

	// Then - 중복 생성으로 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomGroupNameExists, err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestUpdate() {
	// Given - 객실 그룹이 등록된 상황에서
	existingRoomGroup := &models.RoomGroup{
		Name:         "Standard",
		PeekPrice:    100000,
		OffPeekPrice: 80000,
		Description:  "Original description",
	}
	existingRoomGroup.ID = 1

	updates := map[string]interface{}{
		"name":        "UPDATED",
		"peekPrice":   120000,
		"description": "Updated description",
	}

	suite.mockRepo.On("FindByID", suite.ctx, uint(1)).Return(existingRoomGroup, nil)
	// 이름 중복 체크 - ID를 제외하고 체크
	id := uint(1)
	suite.mockRepo.On("ExistsByName", suite.ctx, "UPDATED", &id).Return(false, nil)
	// 업데이트
	suite.mockRepo.On("Update", suite.ctx, existingRoomGroup).Return(nil)

	// When - 객실 그룹 정보 수정을 시도하면
	result, err := suite.service.Update(suite.ctx, 1, updates)

	// Then - 정상적으로 수정된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "UPDATED", result.Name)
	assert.Equal(suite.T(), 120000, result.PeekPrice)
	assert.Equal(suite.T(), "Updated description", result.Description)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestUpdate_NotFound() {
	// Given - 존재하지 않는 객실 그룹에 대해
	updates := map[string]interface{}{}

	suite.mockRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// When - 수정 시도 시
	result, err := suite.service.Update(suite.ctx, 999, updates)

	// Then - 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomGroupNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestUpdate_DuplicateName() {
	// Given - 다른 객실 그룹이 사용 중인 이름으로 변경 시도
	existingRoomGroup := &models.RoomGroup{
		Name: "Standard",
	}
	existingRoomGroup.ID = 1

	updates := map[string]interface{}{
		"name": "Deluxe", // 이미 존재하는 이름
	}

	suite.mockRepo.On("FindByID", suite.ctx, uint(1)).Return(existingRoomGroup, nil)
	// 이름 중복 체크 - 이미 존재함
	id := uint(1)
	suite.mockRepo.On("ExistsByName", suite.ctx, "Deluxe", &id).Return(true, nil)

	// When - 중복된 이름으로 수정을 시도하면
	result, err := suite.service.Update(suite.ctx, 1, updates)

	// Then - 중복 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomGroupNameExists, err)
	assert.Nil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestDelete() {
	// Given - 객실이 없는 객실 그룹이 등록된 상황에서
	roomGroup := &models.RoomGroup{
		Name:  "Standard",
		Rooms: []models.Room{}, // 빈 객실 목록
	}
	roomGroup.ID = 1

	suite.mockRepo.On("FindByIDWithRooms", suite.ctx, uint(1), (*models.RoomStatus)(nil)).Return(roomGroup, nil)
	suite.mockRepo.On("Delete", suite.ctx, uint(1)).Return(nil)

	// When - 객실 그룹 삭제를 시도하면
	err := suite.service.Delete(suite.ctx, 1)

	// Then - 정상적으로 삭제된다
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestDelete_NotFound() {
	// Given - 존재하지 않는 객실 그룹에 대해
	suite.mockRepo.On("FindByIDWithRooms", suite.ctx, uint(999), (*models.RoomStatus)(nil)).Return(nil, gorm.ErrRecordNotFound)

	// When - 삭제 시도 시
	err := suite.service.Delete(suite.ctx, 999)

	// Then - 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomGroupNotFound, err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *RoomGroupServiceTestSuite) TestDelete_HasRooms() {
	// Given - 객실이 있는 객실 그룹이 등록된 상황에서
	rooms := []models.Room{
		{
			Number: "101",
			Status: models.RoomStatusInactive,
		},
	}
	rooms[0].ID = 1

	roomGroup := &models.RoomGroup{
		Name:  "Standard",
		Rooms: rooms,
	}
	roomGroup.ID = 1

	suite.mockRepo.On("FindByIDWithRooms", suite.ctx, uint(1), (*models.RoomStatus)(nil)).Return(roomGroup, nil)

	// When - 객실 그룹 삭제를 시도하면
	err := suite.service.Delete(suite.ctx, 1)

	// Then - 객실이 존재하여 삭제할 수 없다는 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomGroupHasRooms, err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestRoomGroupServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RoomGroupServiceTestSuite))
}
