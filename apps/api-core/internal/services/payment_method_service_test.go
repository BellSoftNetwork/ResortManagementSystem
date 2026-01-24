package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
)

// MockPaymentMethodRepository is a mock implementation of PaymentMethodRepository
type MockPaymentMethodRepository struct {
	mock.Mock
}

func (m *MockPaymentMethodRepository) FindByID(ctx context.Context, id uint) (*models.PaymentMethod, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodRepository) FindAll(ctx context.Context, offset, limit int, sort string) ([]models.PaymentMethod, int64, error) {
	args := m.Called(ctx, offset, limit, sort)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.PaymentMethod), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentMethodRepository) FindByName(ctx context.Context, name string) (*models.PaymentMethod, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodRepository) FindActive(ctx context.Context) ([]models.PaymentMethod, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodRepository) Create(ctx context.Context, paymentMethod *models.PaymentMethod) (*models.PaymentMethod, error) {
	args := m.Called(ctx, paymentMethod)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentMethodRepository) Update(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	args := m.Called(ctx, paymentMethod)
	return args.Error(0)
}

func (m *MockPaymentMethodRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPaymentMethodRepository) ExistsByName(ctx context.Context, name string, excludeID *uint) (bool, error) {
	args := m.Called(ctx, name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockPaymentMethodRepository) ResetAllDefaultSelects(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type PaymentMethodServiceTestSuite struct {
	suite.Suite
	ctx      context.Context
	service  services.PaymentMethodService
	mockRepo *MockPaymentMethodRepository
}

func (suite *PaymentMethodServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockRepo = new(MockPaymentMethodRepository)
	suite.service = services.NewPaymentMethodService(suite.mockRepo)
}

func (suite *PaymentMethodServiceTestSuite) TestGetByID() {
	// Given - 결제 수단이 등록된 상황에서
	paymentMethod := &models.PaymentMethod{
		BaseTimeEntity: models.BaseTimeEntity{
			BaseEntity: models.BaseEntity{ID: 1},
		},
		Name:                     "신용카드",
		CommissionRate:           0.025,
		RequireUnpaidAmountCheck: models.BitBool(true),
		IsDefaultSelect:          models.BitBool(false),
		Status:                   models.PaymentMethodStatusActive,
	}

	suite.mockRepo.On("FindByID", suite.ctx, uint(1)).Return(paymentMethod, nil)

	// When - 특정 결제 수단 정보를 조회하면
	result, err := suite.service.GetByID(suite.ctx, 1)

	// Then - 정상적으로 조회된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), paymentMethod.Name, result.Name)
	assert.Equal(suite.T(), paymentMethod.CommissionRate, result.CommissionRate)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *PaymentMethodServiceTestSuite) TestGetByID_NotFound() {
	// Given - 결제 수단이 없는 상황에서
	suite.mockRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, assert.AnError)

	// When - 존재하지 않는 결제 수단 정보를 조회하면
	result, err := suite.service.GetByID(suite.ctx, 999)

	// Then - 조회 불가 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrPaymentMethodNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *PaymentMethodServiceTestSuite) TestGetAll() {
	// Given - 결제 수단이 등록된 상황에서
	paymentMethods := []models.PaymentMethod{
		{
			BaseTimeEntity: models.BaseTimeEntity{
				BaseEntity: models.BaseEntity{ID: 1},
			},
			Name:           "신용카드",
			CommissionRate: 2.5,
		},
		{
			BaseTimeEntity: models.BaseTimeEntity{
				BaseEntity: models.BaseEntity{ID: 2},
			},
			Name:           "체크카드",
			CommissionRate: 0.015,
		},
	}

	page := 0 // Spring Boot 호환성을 위해 0부터 시작
	size := 20
	offset := 0

	suite.mockRepo.On("FindAll", suite.ctx, offset, size, "").Return(paymentMethods, int64(2), nil)

	// When - 결제 수단 리스트를 조회하면
	result, total, err := suite.service.GetAll(suite.ctx, page, size, "")

	// Then - 모든 결제 수단이 정상적으로 조회된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	assert.Equal(suite.T(), int64(2), total)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *PaymentMethodServiceTestSuite) TestGetActive() {
	// Given - 활성화된 결제 수단이 등록된 상황에서
	paymentMethods := []models.PaymentMethod{
		{
			BaseTimeEntity: models.BaseTimeEntity{
				BaseEntity: models.BaseEntity{ID: 1},
			},
			Name:   "신용카드",
			Status: models.PaymentMethodStatusActive,
		},
		{
			BaseTimeEntity: models.BaseTimeEntity{
				BaseEntity: models.BaseEntity{ID: 2},
			},
			Name:   "체크카드",
			Status: models.PaymentMethodStatusActive,
		},
	}

	suite.mockRepo.On("FindActive", suite.ctx).Return(paymentMethods, nil)

	// When - 활성화된 결제 수단 리스트를 조회하면
	result, err := suite.service.GetActive(suite.ctx)

	// Then - 활성화된 결제 수단만 조회된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *PaymentMethodServiceTestSuite) TestCreate() {
	// Given - 새로운 결제 수단 정보가 주어지면
	paymentMethod := &models.PaymentMethod{
		Name:                     "네이버페이",
		CommissionRate:           0.1,
		RequireUnpaidAmountCheck: models.BitBool(false),
		IsDefaultSelect:          models.BitBool(false),
		Status:                   models.PaymentMethodStatusActive,
	}

	createdPaymentMethod := &models.PaymentMethod{
		BaseTimeEntity: models.BaseTimeEntity{
			BaseEntity: models.BaseEntity{ID: 1},
		},
		Name:                     "네이버페이",
		CommissionRate:           0.1,
		RequireUnpaidAmountCheck: models.BitBool(false),
		IsDefaultSelect:          models.BitBool(false),
		Status:                   models.PaymentMethodStatusActive,
	}

	// 중복 체크
	suite.mockRepo.On("ExistsByName", suite.ctx, "네이버페이", (*uint)(nil)).Return(false, nil)
	suite.mockRepo.On("Create", suite.ctx, paymentMethod).Return(createdPaymentMethod, nil)

	// When - 결제 수단을 생성하면
	err := suite.service.Create(suite.ctx, paymentMethod)

	// Then - 정상적으로 생성된다
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *PaymentMethodServiceTestSuite) TestCreate_DuplicateName() {
	// Given - 이미 존재하는 이름으로 결제 수단 생성 시도
	paymentMethod := &models.PaymentMethod{
		Name: "신용카드",
	}

	// 중복 체크 - 이미 존재함
	suite.mockRepo.On("ExistsByName", suite.ctx, "신용카드", (*uint)(nil)).Return(true, nil)

	// When - 동일한 결제 수단 생성을 시도하면
	err := suite.service.Create(suite.ctx, paymentMethod)

	// Then - 중복 생성으로 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrPaymentMethodNameExists, err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *PaymentMethodServiceTestSuite) TestUpdate() {
	// Given - 결제 수단이 등록된 상황에서
	paymentMethod := &models.PaymentMethod{
		BaseTimeEntity: models.BaseTimeEntity{
			BaseEntity: models.BaseEntity{ID: 1},
		},
		Name:                     "신용카드",
		CommissionRate:           0.025,
		RequireUnpaidAmountCheck: models.BitBool(true),
		IsDefaultSelect:          models.BitBool(false),
		Status:                   models.PaymentMethodStatusActive,
	}

	updates := map[string]interface{}{
		"name":            "BSN",
		"commissionRate":  0.2,
		"isDefaultSelect": models.BitBool(true),
	}

	suite.mockRepo.On("FindByID", suite.ctx, uint(1)).Return(paymentMethod, nil)

	// 이름 중복 체크 - ID를 제외하고 체크
	id := uint(1)
	suite.mockRepo.On("ExistsByName", suite.ctx, "BSN", &id).Return(false, nil)

	// 기본 선택을 true로 변경할 때, 먼저 모든 다른 결제 수단의 기본 선택을 해제
	suite.mockRepo.On("ResetAllDefaultSelects", suite.ctx).Return(nil)

	// 현재 결제 수단 업데이트
	suite.mockRepo.On("Update", suite.ctx, paymentMethod).Return(nil)

	// When - 등록한 결제 수단 정보 수정을 시도하면
	result, err := suite.service.Update(suite.ctx, 1, updates)

	// Then - 정상적으로 수정된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "BSN", result.Name)
	assert.Equal(suite.T(), 0.2, result.CommissionRate)
	assert.True(suite.T(), bool(result.IsDefaultSelect))
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *PaymentMethodServiceTestSuite) TestUpdate_NotFound() {
	// Given - 존재하지 않는 결제 수단에 대해
	updates := map[string]interface{}{}

	suite.mockRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, assert.AnError)

	// When - 존재하지 않는 결제 수단 정보 수정을 시도하면
	result, err := suite.service.Update(suite.ctx, 999, updates)

	// Then - 조회 불가 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrPaymentMethodNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *PaymentMethodServiceTestSuite) TestDelete() {
	// Given - 결제 수단이 등록된 상황에서
	paymentMethod := &models.PaymentMethod{
		BaseTimeEntity: models.BaseTimeEntity{
			BaseEntity: models.BaseEntity{ID: 1},
		},
		Name: "신용카드",
	}

	suite.mockRepo.On("FindByID", suite.ctx, uint(1)).Return(paymentMethod, nil)
	suite.mockRepo.On("Delete", suite.ctx, uint(1)).Return(nil)

	// When - 등록한 결제 수단 정보 삭제를 시도하면
	err := suite.service.Delete(suite.ctx, 1)

	// Then - 정상적으로 삭제된다
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *PaymentMethodServiceTestSuite) TestDelete_NotFound() {
	// Given - 존재하지 않는 결제 수단에 대해
	suite.mockRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, assert.AnError)

	// When - 존재하지 않는 결제 수단 정보 삭제를 시도하면
	err := suite.service.Delete(suite.ctx, 999)

	// Then - 조회 불가 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrPaymentMethodNotFound, err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *PaymentMethodServiceTestSuite) TestUpdate_DefaultSelect_OnlyOneCanBeDefault() {
	// Given - 결제 수단이 등록된 상황에서
	paymentMethod := &models.PaymentMethod{
		BaseTimeEntity: models.BaseTimeEntity{
			BaseEntity: models.BaseEntity{ID: 1},
		},
		Name:                     "신용카드",
		CommissionRate:           0.025,
		RequireUnpaidAmountCheck: models.BitBool(true),
		IsDefaultSelect:          models.BitBool(false),
		Status:                   models.PaymentMethodStatusActive,
	}

	updates := map[string]interface{}{
		"isDefaultSelect": models.BitBool(true),
	}

	suite.mockRepo.On("FindByID", suite.ctx, uint(1)).Return(paymentMethod, nil)

	// 기본 선택을 true로 변경할 때, 먼저 모든 다른 결제 수단의 기본 선택을 해제
	suite.mockRepo.On("ResetAllDefaultSelects", suite.ctx).Return(nil)

	// 현재 결제 수단 업데이트
	suite.mockRepo.On("Update", suite.ctx, paymentMethod).Return(nil)

	// When - 결제 수단을 기본 선택으로 설정하면
	result, err := suite.service.Update(suite.ctx, 1, updates)

	// Then - 정상적으로 수정되고, 다른 결제 수단들의 기본 선택이 먼저 해제된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.True(suite.T(), bool(result.IsDefaultSelect))
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *PaymentMethodServiceTestSuite) TestUpdate_DefaultSelect_FalseDoesNotResetOthers() {
	// Given - 결제 수단이 등록된 상황에서
	paymentMethod := &models.PaymentMethod{
		BaseTimeEntity: models.BaseTimeEntity{
			BaseEntity: models.BaseEntity{ID: 1},
		},
		Name:                     "신용카드",
		CommissionRate:           0.025,
		RequireUnpaidAmountCheck: models.BitBool(true),
		IsDefaultSelect:          models.BitBool(true),
		Status:                   models.PaymentMethodStatusActive,
	}

	updates := map[string]interface{}{
		"isDefaultSelect": models.BitBool(false),
	}

	suite.mockRepo.On("FindByID", suite.ctx, uint(1)).Return(paymentMethod, nil)

	// 기본 선택을 false로 변경할 때는 ResetAllDefaultSelects를 호출하지 않음
	// 현재 결제 수단만 업데이트
	suite.mockRepo.On("Update", suite.ctx, paymentMethod).Return(nil)

	// When - 결제 수단의 기본 선택을 해제하면
	result, err := suite.service.Update(suite.ctx, 1, updates)

	// Then - 정상적으로 수정되고, 다른 결제 수단들에는 영향을 주지 않는다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.False(suite.T(), bool(result.IsDefaultSelect))
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestPaymentMethodServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentMethodServiceTestSuite))
}
