package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.bellsoft.net/rms/api-core/internal/audit"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
)

type MockAuditService struct {
	mock.Mock
}

func (m *MockAuditService) LogCreate(ctx context.Context, entity audit.Auditable) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockAuditService) LogUpdate(ctx context.Context, entity audit.Auditable, oldValues map[string]interface{}) error {
	args := m.Called(ctx, entity, oldValues)
	return args.Error(0)
}

func (m *MockAuditService) LogDelete(ctx context.Context, entity audit.Auditable) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockAuditService) GetHistory(ctx context.Context, entityType string, entityID uint, page, size int) ([]audit.AuditLog, int64, error) {
	args := m.Called(ctx, entityType, entityID, page, size)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]audit.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditService) GetAllHistory(ctx context.Context, filter audit.AuditLogFilter, page, size int) ([]audit.AuditLog, int64, error) {
	args := m.Called(ctx, filter, page, size)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]audit.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditService) GetByID(ctx context.Context, id uint) (*audit.AuditLog, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*audit.AuditLog), args.Error(1)
}

type ReservationServiceDateBlockTestSuite struct {
	suite.Suite
	ctx                   context.Context
	service               services.ReservationService
	mockReservationRepo   *MockReservationRepository
	mockRoomRepo          *MockRoomRepository
	mockPaymentMethodRepo *MockPaymentMethodRepository
	mockAuditService      *MockAuditService
	mockDateBlockRepo     *MockDateBlockRepository
}

func (s *ReservationServiceDateBlockTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockReservationRepo = new(MockReservationRepository)
	s.mockRoomRepo = new(MockRoomRepository)
	s.mockPaymentMethodRepo = new(MockPaymentMethodRepository)
	s.mockAuditService = new(MockAuditService)
	s.mockDateBlockRepo = new(MockDateBlockRepository)

	s.service = services.NewReservationService(
		s.mockReservationRepo,
		s.mockRoomRepo,
		s.mockPaymentMethodRepo,
		s.mockAuditService,
		s.mockDateBlockRepo,
	)
}

func (s *ReservationServiceDateBlockTestSuite) TestCreate_차단된_날짜에_예약_생성하면_ErrDateRangeBlocked를_반환한다() {
	// Given - 예약 가능한 결제수단이지만 날짜 범위가 차단된 예약 요청이 주어지면
	reservation := &models.Reservation{
		Name:            "홍길동",
		StayStartAt:     time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC),
		PaymentMethodID: 1,
	}
	paymentMethod := &models.PaymentMethod{Status: models.PaymentMethodStatusActive, CommissionRate: 0.025}
	paymentMethod.ID = 1

	s.mockPaymentMethodRepo.On("FindByID", s.ctx, uint(1)).Return(paymentMethod, nil)
	s.mockDateBlockRepo.On("IsDateRangeBlocked", s.ctx, reservation.StayStartAt, reservation.StayEndAt).Return(true, nil)

	// When - 예약 생성을 시도하면
	err := s.service.Create(s.ctx, reservation, []uint{1})

	// Then - ErrDateRangeBlocked를 반환하고 저장하지 않는다
	assert.Error(s.T(), err)
	assert.Equal(s.T(), services.ErrDateRangeBlocked, err)
	s.mockPaymentMethodRepo.AssertExpectations(s.T())
	s.mockDateBlockRepo.AssertExpectations(s.T())
	s.mockRoomRepo.AssertNotCalled(s.T(), "IsRoomAvailable", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	s.mockReservationRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

func (s *ReservationServiceDateBlockTestSuite) TestCreate_차단되지_않은_날짜에_예약_생성하면_정상_처리된다() {
	// Given - 차단되지 않은 날짜와 사용 가능한 객실로 예약 요청이 주어지면
	reservation := &models.Reservation{
		Name:            "김철수",
		StayStartAt:     time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2026, 6, 22, 0, 0, 0, 0, time.UTC),
		Price:           100000,
		PaymentMethodID: 2,
	}
	paymentMethod := &models.PaymentMethod{Status: models.PaymentMethodStatusActive, CommissionRate: 0.1}
	paymentMethod.ID = 2
	room := &models.Room{Number: "101"}
	room.ID = 1

	s.mockPaymentMethodRepo.On("FindByID", s.ctx, uint(2)).Return(paymentMethod, nil)
	s.mockDateBlockRepo.On("IsDateRangeBlocked", s.ctx, reservation.StayStartAt, reservation.StayEndAt).Return(false, nil)
	s.mockRoomRepo.On("IsRoomAvailable", s.ctx, uint(1), reservation.StayStartAt, reservation.StayEndAt, (*uint)(nil)).Return(true, nil)
	s.mockRoomRepo.On("FindByID", s.ctx, uint(1)).Return(room, nil)
	s.mockReservationRepo.On("Create", s.ctx, reservation).Return(reservation, nil)

	// When - 예약 생성을 시도하면
	err := s.service.Create(s.ctx, reservation, []uint{1})

	// Then - 정상 생성된다
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 10000, reservation.BrokerFee)
	assert.Len(s.T(), reservation.Rooms, 1)
	s.mockPaymentMethodRepo.AssertExpectations(s.T())
	s.mockDateBlockRepo.AssertExpectations(s.T())
	s.mockRoomRepo.AssertExpectations(s.T())
	s.mockReservationRepo.AssertExpectations(s.T())
}

func (s *ReservationServiceDateBlockTestSuite) TestUpdate_예약_수정_시_날짜를_차단_범위로_변경하면_에러를_반환한다() {
	// Given - 기존 예약의 날짜를 차단된 범위로 변경하는 수정 요청이 주어지면
	existingReservation := &models.Reservation{
		Name:            "이영희",
		StayStartAt:     time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC),
		PaymentMethodID: 1,
	}
	existingReservation.ID = 10

	updates := map[string]interface{}{
		"stayStartAt": time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC),
		"stayEndAt":   time.Date(2026, 7, 7, 0, 0, 0, 0, time.UTC),
	}

	s.mockReservationRepo.On("FindByIDWithDetails", s.ctx, uint(10)).Return(existingReservation, nil)
	s.mockDateBlockRepo.On("IsDateRangeBlocked", s.ctx, updates["stayStartAt"], updates["stayEndAt"]).Return(true, nil)

	// When - 예약 수정을 시도하면
	result, err := s.service.Update(s.ctx, 10, updates, nil, false)

	// Then - ErrDateRangeBlocked를 반환한다
	assert.Error(s.T(), err)
	assert.Equal(s.T(), services.ErrDateRangeBlocked, err)
	assert.Nil(s.T(), result)
	s.mockReservationRepo.AssertExpectations(s.T())
	s.mockDateBlockRepo.AssertExpectations(s.T())
	s.mockReservationRepo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
}

func (s *ReservationServiceDateBlockTestSuite) TestUpdate_예약_수정_시_날짜를_변경하지_않으면_차단_검증을_하지_않는다() {
	// Given - 날짜 변경 없이 메모만 변경하는 수정 요청이 주어지면
	existingReservation := &models.Reservation{
		Name:            "박민수",
		StayStartAt:     time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		StayEndAt:       time.Date(2026, 8, 3, 0, 0, 0, 0, time.UTC),
		PaymentMethodID: 1,
	}
	existingReservation.ID = 11

	updates := map[string]interface{}{
		"note": "요청사항 수정",
	}

	s.mockReservationRepo.On("FindByIDWithDetails", s.ctx, uint(11)).Return(existingReservation, nil).Once()
	s.mockReservationRepo.On("Update", s.ctx, existingReservation).Return(nil).Once()
	s.mockReservationRepo.On("FindByIDWithDetails", s.ctx, uint(11)).Return(existingReservation, nil).Once()

	// When - 예약 수정을 시도하면
	result, err := s.service.Update(s.ctx, 11, updates, nil, false)

	// Then - 정상 수정되고 날짜 차단 검증은 수행하지 않는다
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "요청사항 수정", result.Note)
	s.mockReservationRepo.AssertExpectations(s.T())
	s.mockDateBlockRepo.AssertNotCalled(s.T(), "IsDateRangeBlocked", mock.Anything, mock.Anything, mock.Anything)
}

func TestReservationServiceDateBlockTestSuite(t *testing.T) {
	suite.Run(t, new(ReservationServiceDateBlockTestSuite))
}
