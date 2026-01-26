package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gorm.io/gorm"
)

// MockRoomRepository is a mock implementation of RoomRepository
type MockRoomRepository struct {
	mock.Mock
}

func (m *MockRoomRepository) Create(ctx context.Context, room *models.Room) (*models.Room, error) {
	args := m.Called(ctx, room)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Room), args.Error(1)
}

func (m *MockRoomRepository) Update(ctx context.Context, room *models.Room) error {
	args := m.Called(ctx, room)
	return args.Error(0)
}

func (m *MockRoomRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoomRepository) FindByID(ctx context.Context, id uint) (*models.Room, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Room), args.Error(1)
}

func (m *MockRoomRepository) FindByIDWithGroup(ctx context.Context, id uint) (*models.Room, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Room), args.Error(1)
}

func (m *MockRoomRepository) FindByNumber(ctx context.Context, number string) (*models.Room, error) {
	args := m.Called(ctx, number)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Room), args.Error(1)
}

func (m *MockRoomRepository) FindAll(ctx context.Context, filter dto.RoomRepositoryFilter, offset, limit int, sort string) ([]models.Room, int64, error) {
	args := m.Called(ctx, filter, offset, limit, sort)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Room), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomRepository) ExistsByNumber(ctx context.Context, number string, excludeID *uint) (bool, error) {
	args := m.Called(ctx, number, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoomRepository) FindAvailableRooms(ctx context.Context, startDate, endDate time.Time, excludeReservationID *uint) ([]models.Room, error) {
	args := m.Called(ctx, startDate, endDate, excludeReservationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Room), args.Error(1)
}

func (m *MockRoomRepository) IsRoomAvailable(ctx context.Context, roomID uint, startDate, endDate time.Time, excludeReservationID *uint) (bool, error) {
	args := m.Called(ctx, roomID, startDate, endDate, excludeReservationID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoomRepository) FindByStatus(ctx context.Context, status models.RoomStatus) ([]models.Room, error) {
	args := m.Called(ctx, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Room), args.Error(1)
}

// MockRoomGroupRepository is a mock implementation of RoomGroupRepository
type MockRoomGroupRepository struct {
	mock.Mock
}

func (m *MockRoomGroupRepository) Create(ctx context.Context, roomGroup *models.RoomGroup) (*models.RoomGroup, error) {
	args := m.Called(ctx, roomGroup)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomGroup), args.Error(1)
}

func (m *MockRoomGroupRepository) Update(ctx context.Context, roomGroup *models.RoomGroup) error {
	args := m.Called(ctx, roomGroup)
	return args.Error(0)
}

func (m *MockRoomGroupRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoomGroupRepository) FindByID(ctx context.Context, id uint) (*models.RoomGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomGroup), args.Error(1)
}

func (m *MockRoomGroupRepository) FindByName(ctx context.Context, name string) (*models.RoomGroup, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomGroup), args.Error(1)
}

func (m *MockRoomGroupRepository) FindAll(ctx context.Context, offset, limit int) ([]models.RoomGroup, int64, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.RoomGroup), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoomGroupRepository) ExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error) {
	args := m.Called(ctx, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoomGroupRepository) FindByIDWithRooms(ctx context.Context, id uint, roomStatus *models.RoomStatus) (*models.RoomGroup, error) {
	args := m.Called(ctx, id, roomStatus)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomGroup), args.Error(1)
}

func (m *MockRoomGroupRepository) FindByIDWithFilteredRooms(ctx context.Context, id uint, filter repositories.RoomGroupRoomFilter) (*models.RoomGroup, error) {
	args := m.Called(ctx, id, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomGroup), args.Error(1)
}

func (m *MockRoomGroupRepository) FindByIDWithUsers(ctx context.Context, id uint) (*models.RoomGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RoomGroup), args.Error(1)
}

func (m *MockRoomGroupRepository) FindAllWithUsers(ctx context.Context, offset, limit int, sort string) ([]models.RoomGroup, int64, error) {
	args := m.Called(ctx, offset, limit, sort)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.RoomGroup), args.Get(1).(int64), args.Error(2)
}

type RoomServiceTestSuite struct {
	suite.Suite
	ctx               context.Context
	service           services.RoomService
	mockRoomRepo      *MockRoomRepository
	mockRoomGroupRepo *MockRoomGroupRepository
}

func (suite *RoomServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockRoomRepo = new(MockRoomRepository)
	suite.mockRoomGroupRepo = new(MockRoomGroupRepository)
	suite.service = services.NewRoomService(suite.mockRoomRepo, suite.mockRoomGroupRepo, nil)
}

func (suite *RoomServiceTestSuite) TestGetByID() {
	// Given - 객실이 등록된 상황에서
	room := &models.Room{
		Number:      "101",
		Note:        "Test room",
		Status:      models.RoomStatusNormal,
		RoomGroupID: 1,
	}
	room.ID = 1

	suite.mockRoomRepo.On("FindByID", suite.ctx, uint(1)).Return(room, nil)

	// When - 특정 객실 정보를 조회하면
	result, err := suite.service.GetByID(suite.ctx, 1)

	// Then - 정상적으로 조회된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), room.Number, result.Number)
	assert.Equal(suite.T(), room.Note, result.Note)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetByID_NotFound() {
	// Given - 존재하지 않는 객실 ID로
	suite.mockRoomRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// When - 조회하면
	result, err := suite.service.GetByID(suite.ctx, 999)

	// Then - ErrRoomNotFound 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetByIDWithGroup() {
	// Given - 객실과 그룹 정보가 있는 상황에서
	roomGroup := &models.RoomGroup{
		Name:         "Standard",
		PeekPrice:    100000,
		OffPeekPrice: 80000,
		Description:  "Standard room group",
	}
	roomGroup.ID = 1

	room := &models.Room{
		Number:      "101",
		Note:        "Test room",
		Status:      models.RoomStatusNormal,
		RoomGroupID: 1,
		RoomGroup:   roomGroup,
	}
	room.ID = 1

	suite.mockRoomRepo.On("FindByIDWithGroup", suite.ctx, uint(1)).Return(room, nil)

	// When - 그룹 정보와 함께 객실을 조회하면
	result, err := suite.service.GetByIDWithGroup(suite.ctx, 1)

	// Then - 객실과 그룹 정보가 함께 조회된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), room.Number, result.Number)
	assert.NotNil(suite.T(), result.RoomGroup)
	assert.Equal(suite.T(), "Standard", result.RoomGroup.Name)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAll_Empty() {
	// Given - 객실 정보가 없는 상황에서
	filter := dto.RoomRepositoryFilter{}
	suite.mockRoomRepo.On("FindAll", suite.ctx, filter, 0, 10, "").Return([]models.Room{}, int64(0), nil)

	// When - 전체 객실 정보를 조회하면 (page 0부터 시작 - Spring Boot 호환)
	rooms, total, err := suite.service.GetAll(suite.ctx, filter, 0, 10, "")

	// Then - 빈 객실 목록이 반환된다
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), rooms)
	assert.Equal(suite.T(), int64(0), total)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAll_Multiple() {
	// Given - 객실이 10개 등록된 상황에서
	rooms := make([]models.Room, 10)
	for i := 0; i < 10; i++ {
		rooms[i] = models.Room{
			Number:      string(rune('1'+i)) + "01",
			Note:        "Test room " + string(rune(i)),
			Status:      models.RoomStatusNormal,
			RoomGroupID: 1,
		}
		rooms[i].ID = uint(i + 1)
	}

	filter := dto.RoomRepositoryFilter{}
	suite.mockRoomRepo.On("FindAll", suite.ctx, filter, 0, 10, "").Return(rooms, int64(10), nil)

	// When - 전체 객실 정보를 조회하면 (page 0부터 시작 - Spring Boot 호환)
	result, total, err := suite.service.GetAll(suite.ctx, filter, 0, 10, "")

	// Then - 10개의 객실 정보가 반환된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 10)
	assert.Equal(suite.T(), int64(10), total)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAll_WithFilter() {
	// Given - 비활성 상태의 객실이 있는 상황에서
	activeRoom := models.Room{
		Number: "101",
		Status: models.RoomStatusNormal,
	}
	activeRoom.ID = 1

	statusNormal := models.RoomStatusNormal
	filter := dto.RoomRepositoryFilter{
		Status: &statusNormal,
	}
	suite.mockRoomRepo.On("FindAll", suite.ctx, filter, 0, 10, "").Return([]models.Room{activeRoom}, int64(1), nil)

	// When - 활성 상태의 객실만 조회하면 (page 0부터 시작 - Spring Boot 호환)
	result, total, err := suite.service.GetAll(suite.ctx, filter, 0, 10, "")

	// Then - 활성 객실만 반환된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), int64(1), total)
	assert.Equal(suite.T(), models.RoomStatusNormal, result[0].Status)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestCreate() {
	// Given - 새로운 객실 정보가 주어지면
	roomGroup := &models.RoomGroup{
		Name: "Standard",
	}
	roomGroup.ID = 1

	newRoom := &models.Room{
		Number:      "102",
		Note:        "New room",
		Status:      models.RoomStatusNormal,
		RoomGroupID: 1,
	}

	createdRoom := &models.Room{
		Number:      "102",
		Note:        "New room",
		Status:      models.RoomStatusNormal,
		RoomGroupID: 1,
	}
	createdRoom.ID = 1

	// 중복 체크
	suite.mockRoomRepo.On("ExistsByNumber", suite.ctx, "102", (*uint)(nil)).Return(false, nil)
	// 룸 그룹 존재 확인
	suite.mockRoomGroupRepo.On("FindByID", suite.ctx, uint(1)).Return(roomGroup, nil)
	// 생성
	suite.mockRoomRepo.On("Create", suite.ctx, newRoom).Return(createdRoom, nil)

	// When - 객실을 생성하면
	err := suite.service.Create(suite.ctx, newRoom)

	// Then - 정상적으로 생성된다
	assert.NoError(suite.T(), err)
	suite.mockRoomRepo.AssertExpectations(suite.T())
	suite.mockRoomGroupRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestCreate_DuplicateNumber() {
	// Given - 이미 존재하는 번호로 객실 생성 시도
	newRoom := &models.Room{
		Number:      "101",
		RoomGroupID: 1,
	}

	// 중복 체크 - 이미 존재함
	suite.mockRoomRepo.On("ExistsByNumber", suite.ctx, "101", (*uint)(nil)).Return(true, nil)

	// When - 동일한 객실 번호로 생성을 시도하면
	err := suite.service.Create(suite.ctx, newRoom)

	// Then - 중복 생성으로 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomNumberExists, err)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestCreate_RoomGroupNotFound() {
	// Given - 존재하지 않는 룸 그룹으로 객실 생성 시도
	newRoom := &models.Room{
		Number:      "102",
		RoomGroupID: 999,
	}

	// 중복 체크 통과
	suite.mockRoomRepo.On("ExistsByNumber", suite.ctx, "102", (*uint)(nil)).Return(false, nil)
	// 룸 그룹 조회 실패
	suite.mockRoomGroupRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// When - 존재하지 않는 룸 그룹으로 생성을 시도하면
	err := suite.service.Create(suite.ctx, newRoom)

	// Then - 룸 그룹을 찾을 수 없다는 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomGroupNotFound, err)
	suite.mockRoomRepo.AssertExpectations(suite.T())
	suite.mockRoomGroupRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestUpdate() {
	// Given - 객실이 등록된 상황에서
	existingRoom := &models.Room{
		Number:      "101",
		Note:        "Original note",
		Status:      models.RoomStatusNormal,
		RoomGroupID: 1,
	}
	existingRoom.ID = 1

	updates := map[string]interface{}{
		"number": "UPDATED",
		"note":   "Updated note",
	}

	suite.mockRoomRepo.On("FindByID", suite.ctx, uint(1)).Return(existingRoom, nil)
	// 번호 중복 체크 - ID를 제외하고 체크
	id := uint(1)
	suite.mockRoomRepo.On("ExistsByNumber", suite.ctx, "UPDATED", &id).Return(false, nil)
	// 업데이트
	suite.mockRoomRepo.On("Update", suite.ctx, existingRoom).Return(nil)

	// When - 객실 정보 수정을 시도하면
	result, err := suite.service.Update(suite.ctx, 1, updates)

	// Then - 정상적으로 수정된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "UPDATED", result.Number)
	assert.Equal(suite.T(), "Updated note", result.Note)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestUpdate_NotFound() {
	// Given - 존재하지 않는 객실에 대해
	updates := map[string]interface{}{}

	suite.mockRoomRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// When - 수정 시도 시
	result, err := suite.service.Update(suite.ctx, 999, updates)

	// Then - 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestUpdate_DuplicateNumber() {
	// Given - 다른 객실이 사용 중인 번호로 변경 시도
	existingRoom := &models.Room{
		Number:      "101",
		RoomGroupID: 1,
	}
	existingRoom.ID = 1

	updates := map[string]interface{}{
		"number": "102", // 이미 존재하는 번호
	}

	suite.mockRoomRepo.On("FindByID", suite.ctx, uint(1)).Return(existingRoom, nil)
	// 번호 중복 체크 - 이미 존재함
	id := uint(1)
	suite.mockRoomRepo.On("ExistsByNumber", suite.ctx, "102", &id).Return(true, nil)

	// When - 중복된 번호로 수정을 시도하면
	result, err := suite.service.Update(suite.ctx, 1, updates)

	// Then - 중복 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomNumberExists, err)
	assert.Nil(suite.T(), result)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestUpdate_RoomGroupChange() {
	// Given - 룸 그룹을 변경하려는 상황에서
	existingRoom := &models.Room{
		Number:      "101",
		RoomGroupID: 1,
	}
	existingRoom.ID = 1

	newRoomGroup := &models.RoomGroup{
		Name: "Deluxe",
	}
	newRoomGroup.ID = 2

	updates := map[string]interface{}{
		"room_group_id": uint(2),
	}

	suite.mockRoomRepo.On("FindByID", suite.ctx, uint(1)).Return(existingRoom, nil)
	// 새 룸 그룹 존재 확인
	suite.mockRoomGroupRepo.On("FindByID", suite.ctx, uint(2)).Return(newRoomGroup, nil)
	// 업데이트
	suite.mockRoomRepo.On("Update", suite.ctx, existingRoom).Return(nil)

	// When - 룸 그룹을 변경하면
	result, err := suite.service.Update(suite.ctx, 1, updates)

	// Then - 정상적으로 변경된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), uint(2), result.RoomGroupID)
	suite.mockRoomRepo.AssertExpectations(suite.T())
	suite.mockRoomGroupRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestDelete() {
	// Given - 객실이 등록된 상황에서
	room := &models.Room{
		Number: "101",
	}
	room.ID = 1

	suite.mockRoomRepo.On("FindByID", suite.ctx, uint(1)).Return(room, nil)
	suite.mockRoomRepo.On("Delete", suite.ctx, uint(1)).Return(nil)

	// When - 객실 삭제를 시도하면
	err := suite.service.Delete(suite.ctx, 1)

	// Then - 정상적으로 삭제된다
	assert.NoError(suite.T(), err)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestDelete_NotFound() {
	// Given - 존재하지 않는 객실에 대해
	suite.mockRoomRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// When - 삭제 시도 시
	err := suite.service.Delete(suite.ctx, 999)

	// Then - 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomNotFound, err)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_기간내_예약가능한_객실_조회() {
	// Given - 예약 가능한 객실이 4개 있는 상황에서
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 11, 0, 0, 0, 0, time.UTC)

	availableRooms := []models.Room{
		{Number: "101", Status: models.RoomStatusNormal, RoomGroupID: 1},
		{Number: "102", Status: models.RoomStatusNormal, RoomGroupID: 1},
		{Number: "201", Status: models.RoomStatusNormal, RoomGroupID: 2},
		{Number: "202", Status: models.RoomStatusNormal, RoomGroupID: 2},
	}
	for i := range availableRooms {
		availableRooms[i].ID = uint(i + 1)
	}

	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return(availableRooms, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 4개의 객실이 반환된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 4)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_예약수정시_현재예약_제외() {
	// Given - 특정 예약을 수정하려는 상황에서
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)
	excludeReservationID := uint(1)

	availableRoom := models.Room{Number: "101", Status: models.RoomStatusNormal, RoomGroupID: 1}
	availableRoom.ID = 1

	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, &excludeReservationID).Return([]models.Room{availableRoom}, nil)

	// When - 현재 예약을 제외하고 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, &excludeReservationID)

	// Then - 현재 예약의 객실도 사용 가능으로 나온다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_모든객실_예약불가() {
	// Given - 모든 객실이 예약된 상황에서
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)

	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{}, nil)

	// When - 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 빈 목록이 반환된다
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

// =====================================================
// 희망 기간 외 예약이 잡혀있어 예약이 가능한 객실 시나리오
// (api-legacy RoomServiceTest 동기화)
// =====================================================

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_희망기간외_예약_시나리오0_기존예약이_희망시작일에_끝남() {
	// Given - 기존 예약이 희망 시작일에 끝나는 상황에서
	// 기존 예약 기간: ##=  (2023-11-09 ~ 2023-11-10)
	// 희망 예약 기간: =##  (2023-11-10 ~ 2023-11-11)
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 11, 0, 0, 0, 0, time.UTC)

	availableRoom := models.Room{
		Number:      "101",
		Note:        "[0] 기존: ##= / 희망: =## (기존 예약이 희망 시작일에 끝남)",
		Status:      models.RoomStatusNormal,
		RoomGroupID: 1,
	}
	availableRoom.ID = 1

	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{availableRoom}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 가능하다 (체크아웃 날짜와 체크인 날짜가 같으면 겹치지 않음)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_희망기간외_예약_시나리오1_기존예약이_희망종료일에_시작() {
	// Given - 기존 예약이 희망 종료일에 시작하는 상황에서
	// 기존 예약 기간: =##  (2023-11-11 ~ 2023-11-12)
	// 희망 예약 기간: ##=  (2023-11-10 ~ 2023-11-11)
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 11, 0, 0, 0, 0, time.UTC)

	availableRoom := models.Room{
		Number:      "102",
		Note:        "[1] 기존: =## / 희망: ##= (기존 예약이 희망 종료일에 시작)",
		Status:      models.RoomStatusNormal,
		RoomGroupID: 1,
	}
	availableRoom.ID = 2

	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{availableRoom}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 가능하다 (체크아웃 날짜와 체크인 날짜가 같으면 겹치지 않음)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_희망기간외_예약_시나리오2_기존예약2개가_희망기간_앞뒤() {
	// Given - 기존 예약 2개가 희망 기간 앞뒤에 있는 상황에서
	// 기존 예약 기간: ##@@  (2023-11-09~10, 2023-11-11~12)
	// 희망 예약 기간: =##=  (2023-11-10 ~ 2023-11-11)
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 11, 0, 0, 0, 0, time.UTC)

	availableRoom := models.Room{
		Number:      "103",
		Note:        "[2] 기존: ##@@ / 희망: =##= (기존 예약 2개가 희망 기간 앞뒤에)",
		Status:      models.RoomStatusNormal,
		RoomGroupID: 1,
	}
	availableRoom.ID = 3

	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{availableRoom}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 가능하다 (양쪽 예약 모두 희망 기간과 겹치지 않음)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_희망기간외_예약_시나리오3_기존예약없음() {
	// Given - 기존 예약이 없는 상황에서
	// 기존 예약 기간: ====  (없음)
	// 희망 예약 기간: =##=  (2023-11-10 ~ 2023-11-11)
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 11, 0, 0, 0, 0, time.UTC)

	availableRoom := models.Room{
		Number:      "104",
		Note:        "[3] 기존: ==== / 희망: =##= (기존 예약 없음)",
		Status:      models.RoomStatusNormal,
		RoomGroupID: 1,
	}
	availableRoom.ID = 4

	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{availableRoom}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 가능하다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

// =====================================================
// 희망 기간 내 예약이 잡혀있어 예약이 불가능한 객실 시나리오
// (api-legacy RoomServiceTest 동기화)
// =====================================================

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_희망기간내_예약_시나리오0_before_overlap() {
	// Given - 기존 예약이 희망 시작 전에 시작하여 희망 기간 내로 끝나는 상황
	// 기존 예약 기간: ###=  (2023-11-09 ~ 2023-11-11)
	// 희망 예약 기간: =###  (2023-11-10 ~ 2023-11-20)
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)

	// 예약 불가능 - 빈 목록 반환
	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 불가능하다 (before overlap)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_희망기간내_예약_시나리오1_exact_match() {
	// Given - 기존 예약과 희망 예약이 정확히 일치하는 상황
	// 기존 예약 기간: ###  (2023-11-10 ~ 2023-11-20)
	// 희망 예약 기간: ###  (2023-11-10 ~ 2023-11-20)
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)

	// 예약 불가능 - 빈 목록 반환
	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 불가능하다 (exact match)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_희망기간내_예약_시나리오2_partial_start() {
	// Given - 기존 예약이 희망 시작일에 시작하여 희망 기간 내에 끝나는 상황
	// 기존 예약 기간: ##=  (2023-11-10 ~ 2023-11-11)
	// 희망 예약 기간: ###  (2023-11-10 ~ 2023-11-20)
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)

	// 예약 불가능 - 빈 목록 반환
	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 불가능하다 (partial start overlap)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_희망기간내_예약_시나리오3_contained() {
	// Given - 기존 예약이 희망 기간 안에 완전히 포함된 상황
	// 기존 예약 기간: =###=  (2023-11-11 ~ 2023-11-19)
	// 희망 예약 기간: #####  (2023-11-10 ~ 2023-11-20)
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)

	// 예약 불가능 - 빈 목록 반환
	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 불가능하다 (contained)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_희망기간내_예약_시나리오4_partial_end() {
	// Given - 기존 예약이 희망 기간 내에 시작하여 희망 종료일에 끝나는 상황
	// 기존 예약 기간: =##  (2023-11-19 ~ 2023-11-20)
	// 희망 예약 기간: ###  (2023-11-10 ~ 2023-11-20)
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)

	// 예약 불가능 - 빈 목록 반환
	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 불가능하다 (partial end overlap)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_희망기간내_예약_시나리오5_after_overlap() {
	// Given - 기존 예약이 희망 기간 내에 시작하여 희망 종료일 이후에 끝나는 상황
	// 기존 예약 기간: =###  (2023-11-19 ~ 2023-11-21)
	// 희망 예약 기간: ###=  (2023-11-10 ~ 2023-11-20)
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)

	// 예약 불가능 - 빈 목록 반환
	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 불가능하다 (after overlap)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_희망기간내_예약_시나리오6_wrap() {
	// Given - 기존 예약이 희망 기간을 완전히 감싸는 상황
	// 기존 예약 기간: #####  (2023-11-01 ~ 2023-11-30)
	// 희망 예약 기간: =###=  (2023-11-10 ~ 2023-11-20)
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)

	// 예약 불가능 - 빈 목록 반환
	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 불가능하다 (wrap)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

// =====================================================
// 특수 케이스 시나리오
// (api-legacy RoomServiceTest 동기화)
// =====================================================

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_객실미배정_예약은_가용성에_영향없음() {
	// Given - 객실이 배정되지 않은 예약만 있는 상황에서
	// 희망 예약 기간: 2023-11-10 ~ 2023-11-20
	// 기존에 2023-11-01 ~ 2023-11-30 기간의 예약이 있지만 객실 미배정
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)

	availableRoom := models.Room{
		Number:      "101",
		Note:        "객실 미배정 예약은 가용성에 영향 없음",
		Status:      models.RoomStatusNormal,
		RoomGroupID: 1,
	}
	availableRoom.ID = 1

	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return([]models.Room{availableRoom}, nil)

	// When - 해당 기간에 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 객실이 사용 가능하다 (객실 미배정 예약은 가용성에 영향 없음)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "101", result[0].Number)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *RoomServiceTestSuite) TestGetAvailableRooms_예약수정시_해당예약_객실_사용가능() {
	// Given - 특정 예약을 수정하려는 상황에서
	// 첫 번째 예약이 희망 기간과 겹치지만, 해당 예약을 수정하려는 경우
	startDate := time.Date(2023, 11, 10, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 20, 0, 0, 0, 0, time.UTC)
	excludeReservationID := uint(1)

	// 원래는 예약 불가능하지만, 수정 중인 예약을 제외하면 사용 가능
	availableRoom := models.Room{
		Number:      "101",
		Note:        "수정 중인 예약의 객실",
		Status:      models.RoomStatusNormal,
		RoomGroupID: 1,
	}
	availableRoom.ID = 1

	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, &excludeReservationID).Return([]models.Room{availableRoom}, nil)

	// When - 해당 예약을 제외하고 예약 가능한 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, &excludeReservationID)

	// Then - 해당 예약에 배정된 객실이 사용 가능으로 나온다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 1)
	assert.Equal(suite.T(), "101", result[0].Number)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func TestRoomServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RoomServiceTestSuite))
}
