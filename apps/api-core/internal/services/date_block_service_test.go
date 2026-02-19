package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.bellsoft.net/rms/api-core/internal/audit"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
)

type MockDateBlockRepository struct {
	mock.Mock
}

func (m *MockDateBlockRepository) Create(ctx context.Context, dateBlock *models.DateBlock) (*models.DateBlock, error) {
	args := m.Called(ctx, dateBlock)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DateBlock), args.Error(1)
}

func (m *MockDateBlockRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDateBlockRepository) Update(ctx context.Context, dateBlock *models.DateBlock) error {
	args := m.Called(ctx, dateBlock)
	return args.Error(0)
}

func (m *MockDateBlockRepository) FindByID(ctx context.Context, id uint) (*models.DateBlock, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.DateBlock), args.Error(1)
}

func (m *MockDateBlockRepository) FindAll(ctx context.Context, filter dto.DateBlockFilter, offset, limit int) ([]models.DateBlock, int64, error) {
	args := m.Called(ctx, filter, offset, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.DateBlock), args.Get(1).(int64), args.Error(2)
}

func (m *MockDateBlockRepository) IsDateRangeBlocked(ctx context.Context, startDate, endDate time.Time) (bool, error) {
	args := m.Called(ctx, startDate, endDate)
	return args.Bool(0), args.Error(1)
}

type DateBlockServiceTestSuite struct {
	suite.Suite
	ctx              context.Context
	mockRepo         *MockDateBlockRepository
	mockAuditService *MockAuditService
	service          services.DateBlockService
}

func (s *DateBlockServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockRepo = new(MockDateBlockRepository)
	s.mockAuditService = new(MockAuditService)
	s.service = services.NewDateBlockService(s.mockRepo, s.mockAuditService)
}

func (s *DateBlockServiceTestSuite) TestCreate_날짜_차단을_생성하면_저장된_차단을_반환한다() {
	// Given - 유효한 날짜 차단 생성 요청이 주어지면
	req := dto.CreateDateBlockRequest{
		StartDate: "2026-03-01",
		EndDate:   "2026-03-03",
		Reason:    "객실 점검",
	}

	created := &models.DateBlock{
		Reason:    "객실 점검",
		StartDate: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 3, 3, 0, 0, 0, 0, time.UTC),
	}
	created.ID = 10

	s.mockRepo.On("Create", s.ctx, mock.MatchedBy(func(block *models.DateBlock) bool {
		return block.Reason == req.Reason &&
			block.StartDate.Equal(created.StartDate) &&
			block.EndDate.Equal(created.EndDate)
	})).Return(created, nil)

	// When - 날짜 차단을 생성하면
	result, err := s.service.Create(s.ctx, req)

	// Then - 저장된 차단 정보가 반환된다
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), uint(10), result.ID)
	assert.Equal(s.T(), "객실 점검", result.Reason)
	assert.True(s.T(), result.StartDate.Time.Equal(created.StartDate))
	assert.True(s.T(), result.EndDate.Time.Equal(created.EndDate))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *DateBlockServiceTestSuite) TestCreate_시작일이_종료일보다_이후이면_에러를_반환한다() {
	// Given - 시작일이 종료일보다 이후인 요청이 주어지면
	req := dto.CreateDateBlockRequest{
		StartDate: "2026-03-05",
		EndDate:   "2026-03-03",
		Reason:    "운영 중단",
	}

	// When - 날짜 차단 생성을 시도하면
	result, err := s.service.Create(s.ctx, req)

	// Then - 검증 에러가 반환되고 저장은 호출되지 않는다
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.EqualError(s.T(), err, "startDate must be before or equal to endDate")
	s.mockRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

func (s *DateBlockServiceTestSuite) TestCreate_단일_날짜_차단시작일_종료일이_가능하다() {
	// Given - 시작일과 종료일이 같은 단일 날짜 요청이 주어지면
	req := dto.CreateDateBlockRequest{
		StartDate: "2026-04-01",
		EndDate:   "2026-04-01",
		Reason:    "당일 행사",
	}

	created := &models.DateBlock{
		Reason:    "당일 행사",
		StartDate: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
	}
	created.ID = 11

	s.mockRepo.On("Create", s.ctx, mock.MatchedBy(func(block *models.DateBlock) bool {
		return block.StartDate.Equal(block.EndDate)
	})).Return(created, nil)

	// When - 날짜 차단을 생성하면
	result, err := s.service.Create(s.ctx, req)

	// Then - 정상적으로 생성된다
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), uint(11), result.ID)
	assert.True(s.T(), result.StartDate.Time.Equal(result.EndDate.Time))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *DateBlockServiceTestSuite) TestDelete_차단_삭제_시_존재하지_않는_ID면_에러를_반환한다() {
	// Given - 존재하지 않는 날짜 차단 ID가 주어지면
	s.mockRepo.On("FindByID", s.ctx, uint(999)).Return(nil, errors.New("record not found"))

	// When - 삭제를 시도하면
	err := s.service.Delete(s.ctx, 999)

	// Then - ErrDateBlockNotFound를 반환한다
	assert.Error(s.T(), err)
	assert.Equal(s.T(), services.ErrDateBlockNotFound, err)
	s.mockRepo.AssertExpectations(s.T())
	s.mockRepo.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
	s.mockAuditService.AssertNotCalled(s.T(), "LogDelete", mock.Anything, mock.Anything)
}

func (s *DateBlockServiceTestSuite) TestDelete_날짜_차단_삭제_시_감사_로그가_수동_기록된다() {
	// Given - 존재하는 날짜 차단 ID와 삭제 가능한 상태가 주어지면
	dateBlock := &models.DateBlock{
		StartDate: time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC),
		Reason:    "점검",
	}
	dateBlock.ID = 101

	s.mockRepo.On("FindByID", s.ctx, uint(101)).Return(dateBlock, nil)
	s.mockRepo.On("Delete", s.ctx, uint(101)).Return(nil)
	s.mockAuditService.On("LogDelete", s.ctx, dateBlock).Return(nil)

	// When - 날짜 차단 삭제를 실행하면
	err := s.service.Delete(s.ctx, 101)

	// Then - 삭제 후 감사 로그가 수동으로 기록된다
	assert.NoError(s.T(), err)
	s.mockRepo.AssertExpectations(s.T())
	s.mockAuditService.AssertExpectations(s.T())
}

func (s *DateBlockServiceTestSuite) TestDelete_날짜_차단_삭제_실패_시_감사_로그가_기록되지_않는다() {
	// Given - 날짜 차단 조회는 성공하지만 삭제가 실패하면
	dateBlock := &models.DateBlock{
		StartDate: time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 12, 2, 0, 0, 0, 0, time.UTC),
		Reason:    "긴급 점검",
	}
	dateBlock.ID = 202

	s.mockRepo.On("FindByID", s.ctx, uint(202)).Return(dateBlock, nil)
	s.mockRepo.On("Delete", s.ctx, uint(202)).Return(errors.New("delete failed"))

	// When - 삭제를 실행하면
	err := s.service.Delete(s.ctx, 202)

	// Then - 에러를 반환하고 감사 로그는 기록되지 않는다
	assert.Error(s.T(), err)
	assert.EqualError(s.T(), err, "delete failed")
	s.mockRepo.AssertExpectations(s.T())
	s.mockAuditService.AssertNotCalled(s.T(), "LogDelete", mock.Anything, mock.Anything)
}

func (s *DateBlockServiceTestSuite) TestGetAll_날짜_범위로_차단_목록을_조회한다() {
	// Given - 날짜 범위 필터가 주어지면
	filter := dto.DateBlockFilter{
		StartDate: "2026-05-01",
		EndDate:   "2026-05-31",
	}

	list := []models.DateBlock{
		{StartDate: time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC), Reason: "정기 점검"},
		{StartDate: time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC), EndDate: time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC), Reason: "행사"},
	}
	list[0].ID = 1
	list[1].ID = 2

	s.mockRepo.On("FindAll", s.ctx, filter, 10, 10).Return(list, int64(2), nil)

	// When - 페이지 1, 크기 10으로 목록을 조회하면
	result, total, err := s.service.GetAll(s.ctx, filter, 1, 10)

	// Then - 필터링된 차단 목록이 반환된다
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(2), total)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), uint(1), result[0].ID)
	assert.Equal(s.T(), "정기 점검", result[0].Reason)
	assert.Equal(s.T(), uint(2), result[1].ID)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *DateBlockServiceTestSuite) TestGetDateBlock_날짜_차단_단건_조회_정상_조회() {
	// Given - 존재하는 날짜 차단 ID가 주어지면
	dateBlock := &models.DateBlock{
		StartDate: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC),
		Reason:    "여름 성수기 운영",
	}
	dateBlock.ID = 501

	s.mockRepo.On("FindByID", s.ctx, uint(501)).Return(dateBlock, nil)

	// When - 단건 조회를 실행하면
	result, err := s.service.GetDateBlock(s.ctx, 501)

	// Then - 날짜 차단 정보가 반환된다
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), uint(501), result.ID)
	assert.Equal(s.T(), "여름 성수기 운영", result.Reason)
	assert.True(s.T(), result.StartDate.Time.Equal(dateBlock.StartDate))
	assert.True(s.T(), result.EndDate.Time.Equal(dateBlock.EndDate))
	s.mockRepo.AssertExpectations(s.T())
}

func (s *DateBlockServiceTestSuite) TestGetDateBlock_날짜_차단_단건_조회_존재하지_않는_ID() {
	// Given - 존재하지 않는 날짜 차단 ID가 주어지면
	s.mockRepo.On("FindByID", s.ctx, uint(99999)).Return(nil, errors.New("record not found"))

	// When - 단건 조회를 실행하면
	result, err := s.service.GetDateBlock(s.ctx, 99999)

	// Then - ErrDateBlockNotFound를 반환한다
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), services.ErrDateBlockNotFound, err)
	s.mockRepo.AssertExpectations(s.T())
}

func (s *DateBlockServiceTestSuite) TestUpdateDateBlock_날짜_차단_수정_사유_변경() {
	// Given - 기존 날짜 차단과 사유 변경 요청이 주어지면
	dateBlock := &models.DateBlock{
		StartDate: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC),
		Reason:    "기존 사유",
	}
	dateBlock.ID = 601

	newReason := "수정된 사유"
	req := dto.UpdateDateBlockRequest{Reason: &newReason}

	s.mockRepo.On("FindByID", s.ctx, uint(601)).Return(dateBlock, nil)
	s.mockRepo.On("Update", s.ctx, mock.MatchedBy(func(block *models.DateBlock) bool {
		return block.ID == 601 &&
			block.Reason == "수정된 사유" &&
			block.StartDate.Equal(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)) &&
			block.EndDate.Equal(time.Date(2026, 8, 2, 0, 0, 0, 0, time.UTC))
	})).Return(nil)
	s.mockAuditService.On("LogUpdate", s.ctx, dateBlock, mock.Anything).Return(nil)

	// When - 날짜 차단 수정을 실행하면
	result, err := s.service.UpdateDateBlock(s.ctx, 601, req)

	// Then - 사유가 변경된 결과가 반환된다
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), uint(601), result.ID)
	assert.Equal(s.T(), "수정된 사유", result.Reason)
	s.mockRepo.AssertExpectations(s.T())
	s.mockAuditService.AssertExpectations(s.T())
}

func (s *DateBlockServiceTestSuite) TestUpdateDateBlock_날짜_차단_수정_날짜_범위_변경() {
	// Given - 기존 날짜 차단과 날짜 범위 변경 요청이 주어지면
	dateBlock := &models.DateBlock{
		StartDate: time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 9, 12, 0, 0, 0, 0, time.UTC),
		Reason:    "시설 점검",
	}
	dateBlock.ID = 602

	newStart := "2026-09-11"
	newEnd := "2026-09-13"
	req := dto.UpdateDateBlockRequest{StartDate: &newStart, EndDate: &newEnd}

	s.mockRepo.On("FindByID", s.ctx, uint(602)).Return(dateBlock, nil)
	s.mockRepo.On("Update", s.ctx, mock.MatchedBy(func(block *models.DateBlock) bool {
		return block.ID == 602 &&
			block.Reason == "시설 점검" &&
			block.StartDate.Equal(time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC)) &&
			block.EndDate.Equal(time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC))
	})).Return(nil)
	s.mockAuditService.On("LogUpdate", s.ctx, dateBlock, mock.Anything).Return(nil)

	// When - 날짜 차단 수정을 실행하면
	result, err := s.service.UpdateDateBlock(s.ctx, 602, req)

	// Then - 날짜 범위가 변경된 결과가 반환된다
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), uint(602), result.ID)
	assert.True(s.T(), result.StartDate.Time.Equal(time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC)))
	assert.True(s.T(), result.EndDate.Time.Equal(time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC)))
	s.mockRepo.AssertExpectations(s.T())
	s.mockAuditService.AssertExpectations(s.T())
}

func (s *DateBlockServiceTestSuite) TestUpdateDateBlock_날짜_차단_수정_존재하지_않는_ID() {
	// Given - 존재하지 않는 날짜 차단 ID가 주어지면
	newReason := "수정 시도"
	req := dto.UpdateDateBlockRequest{Reason: &newReason}

	s.mockRepo.On("FindByID", s.ctx, uint(99998)).Return(nil, errors.New("record not found"))

	// When - 날짜 차단 수정을 실행하면
	result, err := s.service.UpdateDateBlock(s.ctx, 99998, req)

	// Then - ErrDateBlockNotFound를 반환하고 수정/감사로그는 호출되지 않는다
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.Equal(s.T(), services.ErrDateBlockNotFound, err)
	s.mockRepo.AssertExpectations(s.T())
	s.mockRepo.AssertNotCalled(s.T(), "Update", mock.Anything, mock.Anything)
	s.mockAuditService.AssertNotCalled(s.T(), "LogUpdate", mock.Anything, mock.Anything, mock.Anything)
}

func (s *DateBlockServiceTestSuite) TestUpdateDateBlock_날짜_차단_수정_감사로그_기록() {
	// Given - 수정 가능한 날짜 차단과 변경 요청이 주어지면
	dateBlock := &models.DateBlock{
		StartDate: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 10, 2, 0, 0, 0, 0, time.UTC),
		Reason:    "기존 감사로그 사유",
	}
	dateBlock.ID = 603

	newReason := "감사로그 대상 사유"
	req := dto.UpdateDateBlockRequest{Reason: &newReason}

	s.mockRepo.On("FindByID", s.ctx, uint(603)).Return(dateBlock, nil)
	s.mockRepo.On("Update", s.ctx, mock.AnythingOfType("*models.DateBlock")).Return(nil)
	s.mockAuditService.On("LogUpdate", s.ctx, dateBlock, mock.MatchedBy(func(oldValues map[string]interface{}) bool {
		return oldValues != nil && oldValues["reason"] == "기존 감사로그 사유"
	})).Return(nil)

	// When - 날짜 차단 수정을 실행하면
	result, err := s.service.UpdateDateBlock(s.ctx, 603, req)

	// Then - 수정 후 감사 로그가 수동으로 기록된다
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "감사로그 대상 사유", result.Reason)
	s.mockRepo.AssertExpectations(s.T())
	s.mockAuditService.AssertExpectations(s.T())
}

func (s *DateBlockServiceTestSuite) TestDateBlock_Auditable_인터페이스를_구현한다() {
	// Given - 감사 대상 DateBlock 엔티티가 주어지면
	dateBlock := &models.DateBlock{
		StartDate: time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2026, 9, 3, 0, 0, 0, 0, time.UTC),
		Reason:    "시설 공사",
	}
	dateBlock.ID = 303

	// When - Auditable 메서드를 호출하면
	entityType := dateBlock.GetAuditEntityType()
	entityID := dateBlock.GetAuditEntityID()
	fields := dateBlock.GetAuditFields()

	// Then - entity type/id/fields를 올바르게 반환한다
	assert.Equal(s.T(), "date_block", entityType)
	assert.Equal(s.T(), uint(303), entityID)
	assert.Equal(s.T(), "2026-09-01", fields["startDate"])
	assert.Equal(s.T(), "2026-09-03", fields["endDate"])
	assert.Equal(s.T(), "시설 공사", fields["reason"])
}

func (s *DateBlockServiceTestSuite) TestGetDateBlockHistory_날짜_차단_이력_조회가_페이지네이션과_함께_동작한다() {
	// Given - date_block 감사 이력이 페이지네이션 파라미터와 함께 주어지면
	historyService := services.NewHistoryService(s.mockAuditService, nil)
	newValues, marshalErr := json.Marshal(map[string]interface{}{
		"id":        404,
		"startDate": "2026-11-01",
		"endDate":   "2026-11-03",
		"reason":    "통합테스트",
		"createdBy": float64(0),
	})
	assert.NoError(s.T(), marshalErr)

	logs := []audit.AuditLog{
		{
			EntityType: "date_block",
			EntityID:   404,
			Action:     audit.ActionCreate,
			NewValues:  newValues,
			Username:   "testadmin",
			CreatedAt:  time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC),
		},
	}

	s.mockAuditService.On("GetHistory", s.ctx, "date_block", uint(404), 0, 10).Return(logs, int64(1), nil)

	// When - 날짜 차단 이력 조회를 요청하면
	result, total, err := historyService.GetDateBlockHistory(s.ctx, 404, 0, 10)

	// Then - 페이지네이션과 함께 이력 결과를 반환한다
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), total)
	assert.Len(s.T(), result, 1)
	assert.Equal(s.T(), uint(404), result[0].Entity.ID)
	assert.Equal(s.T(), "통합테스트", result[0].Entity.Reason)
	assert.Equal(s.T(), "testadmin", result[0].HistoryUsername)
	s.mockAuditService.AssertExpectations(s.T())
}

func TestDateBlockService(t *testing.T) {
	suite.Run(t, new(DateBlockServiceTestSuite))
}
