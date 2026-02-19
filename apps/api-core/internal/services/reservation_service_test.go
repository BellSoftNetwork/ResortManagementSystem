package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
)

// MockReservationRepository is a mock implementation of ReservationRepository
type MockReservationRepository struct {
	mock.Mock
}

func (m *MockReservationRepository) Create(ctx context.Context, reservation *models.Reservation) (*models.Reservation, error) {
	args := m.Called(ctx, reservation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reservation), args.Error(1)
}

func (m *MockReservationRepository) Update(ctx context.Context, reservation *models.Reservation) error {
	args := m.Called(ctx, reservation)
	return args.Error(0)
}

func (m *MockReservationRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockReservationRepository) FindByID(ctx context.Context, id uint) (*models.Reservation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reservation), args.Error(1)
}

func (m *MockReservationRepository) FindByIDWithDetails(ctx context.Context, id uint) (*models.Reservation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reservation), args.Error(1)
}

func (m *MockReservationRepository) FindAll(ctx context.Context, filter dto.ReservationRepositoryFilter, offset, limit int, sort string) ([]models.Reservation, int64, error) {
	args := m.Called(ctx, filter, offset, limit, sort)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Reservation), args.Get(1).(int64), args.Error(2)
}

func (m *MockReservationRepository) GetStatistics(ctx context.Context, startDate, endDate time.Time, periodType string) ([]repositories.ReservationStatistics, error) {
	args := m.Called(ctx, startDate, endDate, periodType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repositories.ReservationStatistics), args.Error(1)
}

func (m *MockReservationRepository) DeleteRooms(ctx context.Context, reservationID uint) error {
	args := m.Called(ctx, reservationID)
	return args.Error(0)
}

func (m *MockReservationRepository) FindLastReservationForRoom(ctx context.Context, roomID uint) (*models.Reservation, error) {
	args := m.Called(ctx, roomID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Reservation), args.Error(1)
}

type ReservationServiceTestSuite struct {
	suite.Suite
	ctx                   context.Context
	service               services.ReservationService
	mockReservationRepo   *MockReservationRepository
	mockRoomRepo          *MockRoomRepository
	mockPaymentMethodRepo *MockPaymentMethodRepository
}

func (suite *ReservationServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockReservationRepo = new(MockReservationRepository)
	suite.mockRoomRepo = new(MockRoomRepository)
	suite.mockPaymentMethodRepo = new(MockPaymentMethodRepository)

	suite.service = services.NewReservationService(
		suite.mockReservationRepo,
		suite.mockRoomRepo,
		suite.mockPaymentMethodRepo,
		nil,
	)
}

func (suite *ReservationServiceTestSuite) TestGetByID() {
	// Given - 예약이 등록된 상황에서
	reservation := &models.Reservation{
		Name:            "홍길동",
		Phone:           "010-1234-5678",
		PeopleCount:     2,
		StayStartAt:     time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC),
		Price:           200000,
		Status:          models.ReservationStatusNormal,
		Type:            models.ReservationTypeStay,
		PaymentMethodID: 1,
	}
	reservation.ID = 1

	suite.mockReservationRepo.On("FindByID", suite.ctx, uint(1)).Return(reservation, nil)

	// When - 특정 예약 정보를 조회하면
	result, err := suite.service.GetByID(suite.ctx, 1)

	// Then - 정상적으로 조회된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), reservation.Name, result.Name)
	assert.Equal(suite.T(), reservation.Phone, result.Phone)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetByID_NotFound() {
	// Given - 존재하지 않는 예약 ID로
	suite.mockReservationRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, errors.New("record not found"))

	// When - 조회하면
	result, err := suite.service.GetByID(suite.ctx, 999)

	// Then - ErrReservationNotFound 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrReservationNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetAll_Empty() {
	// Given - 예약 정보가 없는 상황에서
	filter := dto.ReservationRepositoryFilter{}
	suite.mockReservationRepo.On("FindAll", suite.ctx, filter, 0, 10, "").Return([]models.Reservation{}, int64(0), nil)

	// When - 전체 예약 정보를 조회하면 (page 0부터 시작 - Spring Boot 호환)
	reservations, total, err := suite.service.GetAll(suite.ctx, filter, 0, 10, "")

	// Then - 빈 예약 목록이 반환된다
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), reservations)
	assert.Equal(suite.T(), int64(0), total)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetAll_Pagination() {
	// Given - 여러 예약이 있는 상황에서
	filter := dto.ReservationRepositoryFilter{}
	mockReservations := []models.Reservation{
		{Name: "예약1"},
		{Name: "예약2"},
		{Name: "예약3"},
		{Name: "예약4"},
		{Name: "예약5"},
	}

	// page 0, size 5로 조회할 때 offset 0으로 계산되어야 함
	suite.mockReservationRepo.On("FindAll", suite.ctx, filter, 0, 5, "").Return(mockReservations[:5], int64(15), nil).Once()

	// page 1, size 5로 조회할 때 offset 5로 계산되어야 함
	suite.mockReservationRepo.On("FindAll", suite.ctx, filter, 5, 5, "").Return(mockReservations[:5], int64(15), nil).Once()

	// page 2, size 5로 조회할 때 offset 10으로 계산되어야 함
	suite.mockReservationRepo.On("FindAll", suite.ctx, filter, 10, 5, "").Return(mockReservations[:5], int64(15), nil).Once()

	// When - 각 페이지를 조회하면
	// 첫 번째 페이지 (page 0)
	reservations1, total1, err1 := suite.service.GetAll(suite.ctx, filter, 0, 5, "")
	assert.NoError(suite.T(), err1)
	assert.Len(suite.T(), reservations1, 5)
	assert.Equal(suite.T(), int64(15), total1)

	// 두 번째 페이지 (page 1)
	reservations2, total2, err2 := suite.service.GetAll(suite.ctx, filter, 1, 5, "")
	assert.NoError(suite.T(), err2)
	assert.Len(suite.T(), reservations2, 5)
	assert.Equal(suite.T(), int64(15), total2)

	// 세 번째 페이지 (page 2)
	reservations3, total3, err3 := suite.service.GetAll(suite.ctx, filter, 2, 5, "")
	assert.NoError(suite.T(), err3)
	assert.Len(suite.T(), reservations3, 5)
	assert.Equal(suite.T(), int64(15), total3)

	// Then - 각 페이지에 대해 올바른 offset으로 repository가 호출되었는지 확인
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestCreate() {
	// Given - 새로운 예약 정보가 주어지면
	paymentMethod := &models.PaymentMethod{
		Name:           "신용카드",
		CommissionRate: 0.025,
		Status:         models.PaymentMethodStatusActive,
	}
	paymentMethod.ID = 1

	newReservation := &models.Reservation{
		Name:            "홍길동",
		Phone:           "010-1234-5678",
		PeopleCount:     2,
		StayStartAt:     time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC),
		Price:           200000,
		Status:          models.ReservationStatusNormal,
		Type:            models.ReservationTypeStay,
		PaymentMethodID: 1,
	}

	createdReservation := &models.Reservation{
		Name:            "홍길동",
		Phone:           "010-1234-5678",
		PeopleCount:     2,
		StayStartAt:     time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC),
		Price:           200000,
		BrokerFee:       5000, // 200000 * 2.5%
		Status:          models.ReservationStatusNormal,
		Type:            models.ReservationTypeStay,
		PaymentMethodID: 1,
		Rooms: []models.ReservationRoom{
			{RoomID: 1},
			{RoomID: 2},
		},
	}
	createdReservation.ID = 1

	roomIDs := []uint{1, 2}

	// 결제 수단 확인
	suite.mockPaymentMethodRepo.On("FindByID", suite.ctx, uint(1)).Return(paymentMethod, nil)
	// 객실 가용성 확인
	suite.mockRoomRepo.On("IsRoomAvailable", suite.ctx, uint(1), newReservation.StayStartAt, newReservation.StayEndAt, (*uint)(nil)).Return(true, nil)
	suite.mockRoomRepo.On("IsRoomAvailable", suite.ctx, uint(2), newReservation.StayStartAt, newReservation.StayEndAt, (*uint)(nil)).Return(true, nil)
	// 객실 정보 로드
	suite.mockRoomRepo.On("FindByID", suite.ctx, uint(1)).Return(&models.Room{Number: "101"}, nil)
	suite.mockRoomRepo.On("FindByID", suite.ctx, uint(2)).Return(&models.Room{Number: "102"}, nil)
	// 생성
	suite.mockReservationRepo.On("Create", suite.ctx, newReservation).Return(createdReservation, nil)

	// When - 예약을 생성하면
	err := suite.service.Create(suite.ctx, newReservation, roomIDs)

	// Then - 정상적으로 생성된다
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 5000, newReservation.BrokerFee) // 수수료가 계산됨
	assert.Len(suite.T(), newReservation.Rooms, 2)
	assert.NotNil(suite.T(), newReservation.PaymentMethod)
	assert.Equal(suite.T(), "신용카드", newReservation.PaymentMethod.Name)
	assert.NotNil(suite.T(), newReservation.Rooms[0].Room)
	assert.Equal(suite.T(), "101", newReservation.Rooms[0].Room.Number)
	suite.mockPaymentMethodRepo.AssertExpectations(suite.T())
	suite.mockRoomRepo.AssertExpectations(suite.T())
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestCreate_InvalidDateRange() {
	// Given - 잘못된 날짜 범위로 예약 시도
	newReservation := &models.Reservation{
		Name:            "홍길동",
		StayStartAt:     time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC), // 시작일이 종료일보다 늦음
		PaymentMethodID: 1,
	}

	// When - 예약을 생성하면
	err := suite.service.Create(suite.ctx, newReservation, []uint{})

	// Then - 날짜 범위 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrInvalidDateRange, err)
}

func (suite *ReservationServiceTestSuite) TestCreate_RoomNotAvailable() {
	// Given - 이미 예약된 객실로 예약 시도
	paymentMethod := &models.PaymentMethod{
		Name:   "신용카드",
		Status: models.PaymentMethodStatusActive,
	}
	paymentMethod.ID = 1

	newReservation := &models.Reservation{
		Name:            "홍길동",
		StayStartAt:     time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC),
		PaymentMethodID: 1,
	}

	roomIDs := []uint{1}

	// 결제 수단 확인
	suite.mockPaymentMethodRepo.On("FindByID", suite.ctx, uint(1)).Return(paymentMethod, nil)
	// 객실 가용성 확인 - 사용 불가
	suite.mockRoomRepo.On("IsRoomAvailable", suite.ctx, uint(1), newReservation.StayStartAt, newReservation.StayEndAt, (*uint)(nil)).Return(false, nil)

	// When - 예약을 생성하면
	err := suite.service.Create(suite.ctx, newReservation, roomIDs)

	// Then - 객실 사용 불가 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomNotAvailable, err)
	suite.mockPaymentMethodRepo.AssertExpectations(suite.T())
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestUpdate() {
	// Given - 예약이 등록된 상황에서
	existingReservation := &models.Reservation{
		Name:            "홍길동",
		Phone:           "010-1234-5678",
		PeopleCount:     2,
		StayStartAt:     time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC),
		Price:           200000,
		Status:          models.ReservationStatusNormal,
		PaymentMethodID: 1,
		Rooms: []models.ReservationRoom{
			{RoomID: 1},
		},
	}
	existingReservation.ID = 1

	updates := map[string]interface{}{
		"name":        "김철수",
		"phone":       "010-9876-5432",
		"peopleCount": 3,
		"note":        "특별 요청 사항",
	}

	suite.mockReservationRepo.On("FindByIDWithDetails", suite.ctx, uint(1)).Return(existingReservation, nil).Once()
	suite.mockReservationRepo.On("Update", suite.ctx, existingReservation).Return(nil)
	suite.mockReservationRepo.On("FindByIDWithDetails", suite.ctx, uint(1)).Return(existingReservation, nil).Once()

	// When - 예약 정보 수정을 시도하면
	result, err := suite.service.Update(suite.ctx, 1, updates, nil, false)

	// Then - 정상적으로 수정된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "김철수", result.Name)
	assert.Equal(suite.T(), "010-9876-5432", result.Phone)
	assert.Equal(suite.T(), 3, result.PeopleCount)
	assert.Equal(suite.T(), "특별 요청 사항", result.Note)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestUpdate_NotFound() {
	// Given - 존재하지 않는 예약에 대해
	updates := map[string]interface{}{}

	suite.mockReservationRepo.On("FindByIDWithDetails", suite.ctx, uint(999)).Return(nil, errors.New("record not found"))

	// When - 수정 시도 시
	result, err := suite.service.Update(suite.ctx, 999, updates, nil, false)

	// Then - 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrReservationNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestUpdate_StatusToCanceled() {
	// Given - 예약을 취소 상태로 변경하려는 상황에서
	existingReservation := &models.Reservation{
		Name:            "홍길동",
		Status:          models.ReservationStatusNormal,
		StayStartAt:     time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC),
		PaymentMethodID: 1,
	}
	existingReservation.ID = 1

	updates := map[string]interface{}{
		"status": models.ReservationStatusCancel,
	}

	suite.mockReservationRepo.On("FindByIDWithDetails", suite.ctx, uint(1)).Return(existingReservation, nil).Once()
	suite.mockReservationRepo.On("Update", suite.ctx, existingReservation).Return(nil)
	suite.mockReservationRepo.On("FindByIDWithDetails", suite.ctx, uint(1)).Return(existingReservation, nil).Once()

	// When - 상태를 취소로 변경하면
	result, err := suite.service.Update(suite.ctx, 1, updates, nil, false)

	// Then - 정상적으로 수정되고 취소 시간이 설정된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), models.ReservationStatusCancel, result.Status)
	assert.NotNil(suite.T(), result.CanceledAt)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestUpdate_결제수단변경시_PaymentMethodID가_올바르게_저장된다() {
	// Given - 결제수단이 등록된 예약에서
	now := time.Now()
	existingReservation := &models.Reservation{
		PaymentMethodID: 1,
		PaymentMethod: &models.PaymentMethod{
			CommissionRate: 0.1,
			Status:         models.PaymentMethodStatusActive,
		},
		Price:       100000,
		StayStartAt: now.Add(24 * time.Hour),
		StayEndAt:   now.Add(48 * time.Hour),
	}
	existingReservation.ID = 1
	existingReservation.PaymentMethod.ID = 1

	newPaymentMethod := &models.PaymentMethod{
		CommissionRate: 0.2,
		Status:         models.PaymentMethodStatusActive,
	}
	newPaymentMethod.ID = 3

	updates := map[string]interface{}{
		"paymentMethodId": uint(3),
	}

	// 1. 기존 예약 조회 (FindByIDWithDetails)
	suite.mockReservationRepo.On("FindByIDWithDetails", suite.ctx, uint(1)).Return(existingReservation, nil).Once()
	// 2. 새 결제수단 조회
	suite.mockPaymentMethodRepo.On("FindByID", suite.ctx, uint(3)).Return(newPaymentMethod, nil).Once()
	// 3. 업데이트 수행
	suite.mockReservationRepo.On("Update", suite.ctx, existingReservation).Return(nil).Once()
	// 4. 업데이트 후 다시 조회 (FindByIDWithDetails)
	suite.mockReservationRepo.On("FindByIDWithDetails", suite.ctx, uint(1)).Return(existingReservation, nil).Once()

	// When - 결제수단을 변경하면
	result, err := suite.service.Update(suite.ctx, 1, updates, nil, false)

	// Then - PaymentMethodID가 변경되어야 하고, preloaded된 PaymentMethod는 초기화되어야 한다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), uint(3), existingReservation.PaymentMethodID)
	assert.Equal(suite.T(), 20000, existingReservation.BrokerFee) // 100000 * 0.2

	// 이 부분에서 실패할 것으로 예상됨 (버그: PaymentMethod가 nil이 아니면 GORM이 예전 값을 다시 쓸 수 있음)
	// 서비스 로직에서 reservation.PaymentMethod = nil 처리가 누락되어 있음
	assert.Nil(suite.T(), existingReservation.PaymentMethod, "결제수단 ID가 변경될 때 기존 preloaded 객체는 제거되어야 함")

	suite.mockReservationRepo.AssertExpectations(suite.T())
	suite.mockPaymentMethodRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestDelete() {
	// Given - 예약이 등록된 상황에서
	reservation := &models.Reservation{
		Name: "홍길동",
	}
	reservation.ID = 1

	suite.mockReservationRepo.On("FindByID", suite.ctx, uint(1)).Return(reservation, nil)
	suite.mockReservationRepo.On("Delete", suite.ctx, uint(1)).Return(nil)

	// When - 예약 삭제를 시도하면
	err := suite.service.Delete(suite.ctx, 1)

	// Then - 정상적으로 삭제된다
	assert.NoError(suite.T(), err)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestDelete_NotFound() {
	// Given - 존재하지 않는 예약에 대해
	suite.mockReservationRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, errors.New("record not found"))

	// When - 삭제 시도 시
	err := suite.service.Delete(suite.ctx, 999)

	// Then - 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrReservationNotFound, err)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetAvailableRooms() {
	// Given - 특정 기간의 예약 가능 객실 조회
	startDate := time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC)

	availableRooms := []models.Room{
		{
			Number: "101",
			Status: models.RoomStatusNormal,
		},
		{
			Number: "102",
			Status: models.RoomStatusNormal,
		},
	}
	availableRooms[0].ID = 1
	availableRooms[1].ID = 2

	suite.mockRoomRepo.On("FindAvailableRooms", suite.ctx, startDate, endDate, (*uint)(nil)).Return(availableRooms, nil)

	// When - 예약 가능 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 예약 가능한 객실 목록이 반환된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), "101", result[0].Number)
	assert.Equal(suite.T(), "102", result[1].Number)
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetAvailableRooms_InvalidDateRange() {
	// Given - 잘못된 날짜 범위로 조회
	startDate := time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC) // 시작일이 종료일보다 늦음

	// When - 예약 가능 객실을 조회하면
	result, err := suite.service.GetAvailableRooms(suite.ctx, startDate, endDate, nil)

	// Then - 날짜 범위 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrInvalidDateRange, err)
	assert.Nil(suite.T(), result)
}

func (suite *ReservationServiceTestSuite) TestGetAll_DateFiltering() {
	// Given - 다양한 날짜 범위의 예약들이 있는 상황에서
	filter := dto.ReservationRepositoryFilter{
		StartDate: &time.Time{}, // 2025-08-17로 설정할 예정
		EndDate:   &time.Time{}, // 2025-11-17로 설정할 예정
	}

	// 필터 날짜 설정
	startDate := time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 11, 17, 0, 0, 0, 0, time.UTC)
	filter.StartDate = &startDate
	filter.EndDate = &endDate

	// api-legacy 로직과 동일하게 필터링되어야 하는 예약들
	expectedReservations := []models.Reservation{
		{Name: "예약1"}, // 2025-08-01 ~ 2025-08-20 (체크아웃이 필터 시작일 이후)
		{Name: "예약2"}, // 2025-09-01 ~ 2025-09-05 (완전히 필터 기간 내)
		{Name: "예약3"}, // 2025-11-10 ~ 2025-11-20 (체크인이 필터 종료일 이전)
		{Name: "예약4"}, // 2025-07-01 ~ 2025-12-31 (필터 기간을 포함)
	}

	suite.mockReservationRepo.On("FindAll", suite.ctx, filter, 0, 15, "").Return(expectedReservations, int64(4), nil)

	// When - 날짜 필터로 예약을 조회하면
	reservations, total, err := suite.service.GetAll(suite.ctx, filter, 0, 15, "")

	// Then - api-legacy와 동일한 결과가 반환된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), reservations, 4)
	assert.Equal(suite.T(), int64(4), total)
	assert.Equal(suite.T(), "예약1", reservations[0].Name)
	assert.Equal(suite.T(), "예약2", reservations[1].Name)
	assert.Equal(suite.T(), "예약3", reservations[2].Name)
	assert.Equal(suite.T(), "예약4", reservations[3].Name)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetAll_DateFiltering_OnlyStartDate() {
	// Given - 시작일만 필터가 있는 상황에서
	filter := dto.ReservationRepositoryFilter{}
	startDate := time.Date(2025, 8, 17, 0, 0, 0, 0, time.UTC)
	filter.StartDate = &startDate

	expectedReservations := []models.Reservation{
		{Name: "예약1"}, // 시작일 이후의 예약들
		{Name: "예약2"},
	}

	suite.mockReservationRepo.On("FindAll", suite.ctx, filter, 0, 15, "").Return(expectedReservations, int64(2), nil)

	// When - 시작일 필터로 예약을 조회하면
	reservations, total, err := suite.service.GetAll(suite.ctx, filter, 0, 15, "")

	// Then - 시작일 이후의 예약들이 반환된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), reservations, 2)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetAll_DateFiltering_OnlyEndDate() {
	// Given - 종료일만 필터가 있는 상황에서
	filter := dto.ReservationRepositoryFilter{}
	endDate := time.Date(2025, 11, 17, 0, 0, 0, 0, time.UTC)
	filter.EndDate = &endDate

	expectedReservations := []models.Reservation{
		{Name: "예약1"}, // 종료일 이전의 예약들
		{Name: "예약2"},
		{Name: "예약3"},
	}

	suite.mockReservationRepo.On("FindAll", suite.ctx, filter, 0, 15, "").Return(expectedReservations, int64(3), nil)

	// When - 종료일 필터로 예약을 조회하면
	reservations, total, err := suite.service.GetAll(suite.ctx, filter, 0, 15, "")

	// Then - 종료일 이전의 예약들이 반환된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), reservations, 3)
	assert.Equal(suite.T(), int64(3), total)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestUpdate_기간변경과_객실재배정시_충돌_실패() {
	// Given - 객실이 배정된 예약을 충돌하는 기간으로 변경하고 같은 객실 재배정 시도
	existingReservation := &models.Reservation{
		Name:            "홍길동",
		StayStartAt:     time.Date(2023, 10, 31, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC),
		Status:          models.ReservationStatusNormal,
		PaymentMethodID: 1,
		Rooms: []models.ReservationRoom{
			{RoomID: 1},
		},
	}
	existingReservation.ID = 1

	newStayStartAt := time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC)
	newStayEndAt := time.Date(2023, 11, 2, 0, 0, 0, 0, time.UTC)

	updates := map[string]interface{}{
		"stayStartAt": newStayStartAt,
		"stayEndAt":   newStayEndAt,
	}

	suite.mockReservationRepo.On("FindByIDWithDetails", suite.ctx, uint(1)).Return(existingReservation, nil)
	excludeID := uint(1)
	suite.mockRoomRepo.On("IsRoomAvailable", suite.ctx, uint(1), newStayStartAt, newStayEndAt, &excludeID).Return(false, nil)

	// When - 충돌하는 기간으로 수정하고 객실도 재배정(hasRoomsUpdate=true)하면
	result, err := suite.service.Update(suite.ctx, 1, updates, []uint{1}, true)

	// Then - 중복 객실 배정 예외가 발생하면서 수정되지 않는다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomNotAvailable, err)
	assert.Nil(suite.T(), result)
	suite.mockReservationRepo.AssertExpectations(suite.T())
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestUpdate_객실변경시_새객실_충돌_실패() {
	// Given - 객실이 배정된 예약에서 이미 예약된 다른 객실로 변경 시도
	existingReservation := &models.Reservation{
		Name:            "홍길동",
		StayStartAt:     time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2023, 11, 5, 0, 0, 0, 0, time.UTC),
		Status:          models.ReservationStatusNormal,
		PaymentMethodID: 1,
		Rooms: []models.ReservationRoom{
			{RoomID: 1},
		},
	}
	existingReservation.ID = 1

	newRoomIDs := []uint{2}

	suite.mockReservationRepo.On("FindByIDWithDetails", suite.ctx, uint(1)).Return(existingReservation, nil)
	excludeID := uint(1)
	suite.mockRoomRepo.On("IsRoomAvailable", suite.ctx, uint(2), existingReservation.StayStartAt, existingReservation.StayEndAt, &excludeID).Return(false, nil)

	// When - 충돌하는 객실로 변경을 시도하면
	result, err := suite.service.Update(suite.ctx, 1, nil, newRoomIDs, true)

	// Then - 중복 객실 배정 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrRoomNotAvailable, err)
	assert.Nil(suite.T(), result)
	suite.mockReservationRepo.AssertExpectations(suite.T())
	suite.mockRoomRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetStatistics_월별통계() {
	// Given - 10월~12월 예약이 있는 상황에서
	startDate := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	expectedStats := []repositories.ReservationStatistics{
		{Period: "2023-10", TotalRevenue: 100000, ReservationCount: 1, TotalGuests: 2},
		{Period: "2023-11", TotalRevenue: 500000, ReservationCount: 2, TotalGuests: 7},
		{Period: "2023-12", TotalRevenue: 400000, ReservationCount: 1, TotalGuests: 5},
	}

	suite.mockReservationRepo.On("GetStatistics", suite.ctx, startDate, endDate, "MONTHLY").Return(expectedStats, nil)

	// When - MONTHLY 타입으로 예약 통계를 조회하면
	stats, err := suite.service.GetStatistics(suite.ctx, startDate, endDate, "MONTHLY")

	// Then - 월별 통계 데이터가 정상적으로 조회된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), stats, 3)
	assert.Equal(suite.T(), "2023-10", stats[0].Period)
	assert.Equal(suite.T(), float64(100000), stats[0].TotalRevenue)
	assert.Equal(suite.T(), "2023-11", stats[1].Period)
	assert.Equal(suite.T(), float64(500000), stats[1].TotalRevenue)
	assert.Equal(suite.T(), "2023-12", stats[2].Period)
	assert.Equal(suite.T(), float64(400000), stats[2].TotalRevenue)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetStatistics_일별통계() {
	// Given - 11월 초 예약이 있는 상황에서
	startDate := time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 11, 5, 0, 0, 0, 0, time.UTC)

	expectedStats := []repositories.ReservationStatistics{
		{Period: "2023-11-01", TotalRevenue: 200000, ReservationCount: 1, TotalGuests: 3},
		{Period: "2023-11-30", TotalRevenue: 300000, ReservationCount: 1, TotalGuests: 4},
	}

	suite.mockReservationRepo.On("GetStatistics", suite.ctx, startDate, endDate, "DAILY").Return(expectedStats, nil)

	// When - DAILY 타입으로 예약 통계를 조회하면
	stats, err := suite.service.GetStatistics(suite.ctx, startDate, endDate, "DAILY")

	// Then - 일별 통계 데이터가 정상적으로 조회된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), stats, 2)
	assert.Equal(suite.T(), "2023-11-01", stats[0].Period)
	assert.Equal(suite.T(), float64(200000), stats[0].TotalRevenue)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetByIDWithDetails() {
	// Given - 예약이 등록된 상황에서
	reservation := &models.Reservation{
		Name:            "홍길동",
		Phone:           "010-1234-5678",
		PeopleCount:     2,
		StayStartAt:     time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC),
		Price:           200000,
		Status:          models.ReservationStatusNormal,
		Type:            models.ReservationTypeStay,
		PaymentMethodID: 1,
		Rooms: []models.ReservationRoom{
			{RoomID: 1},
			{RoomID: 2},
		},
	}
	reservation.ID = 1

	suite.mockReservationRepo.On("FindByIDWithDetails", suite.ctx, uint(1)).Return(reservation, nil)

	// When - 특정 예약의 상세 정보를 조회하면
	result, err := suite.service.GetByIDWithDetails(suite.ctx, 1)

	// Then - 정상적으로 조회되고 연관 데이터도 포함된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), reservation.Name, result.Name)
	assert.Equal(suite.T(), reservation.Phone, result.Phone)
	assert.Len(suite.T(), result.Rooms, 2)
	assert.Equal(suite.T(), uint(1), result.Rooms[0].RoomID)
	assert.Equal(suite.T(), uint(2), result.Rooms[1].RoomID)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetByIDWithDetails_NotFound() {
	// Given - 존재하지 않는 예약 ID로
	suite.mockReservationRepo.On("FindByIDWithDetails", suite.ctx, uint(999)).Return(nil, errors.New("record not found"))

	// When - 상세 정보를 조회하면
	result, err := suite.service.GetByIDWithDetails(suite.ctx, 999)

	// Then - ErrReservationNotFound 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrReservationNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetLastReservationForRoom() {
	// Given - 특정 객실에 예약이 있는 상황에서
	lastReservation := &models.Reservation{
		Name:        "최근 예약",
		Phone:       "010-1111-2222",
		StayStartAt: time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC),
		StayEndAt:   time.Date(2024, 3, 27, 0, 0, 0, 0, time.UTC),
		Status:      models.ReservationStatusNormal,
		Rooms: []models.ReservationRoom{
			{RoomID: 1},
		},
	}
	lastReservation.ID = 5

	suite.mockReservationRepo.On("FindLastReservationForRoom", suite.ctx, uint(1)).Return(lastReservation, nil)

	// When - 해당 객실의 마지막 예약을 조회하면
	result, err := suite.service.GetLastReservationForRoom(suite.ctx, 1)

	// Then - 정상적으로 조회된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), lastReservation.Name, result.Name)
	assert.Equal(suite.T(), lastReservation.ID, result.ID)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func (suite *ReservationServiceTestSuite) TestGetLastReservationForRoom_NoReservation() {
	// Given - 예약이 없는 객실에 대해
	suite.mockReservationRepo.On("FindLastReservationForRoom", suite.ctx, uint(999)).Return(nil, errors.New("record not found"))

	// When - 마지막 예약을 조회하면
	result, err := suite.service.GetLastReservationForRoom(suite.ctx, 999)

	// Then - 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	suite.mockReservationRepo.AssertExpectations(suite.T())
}

func TestReservationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ReservationServiceTestSuite))
}
