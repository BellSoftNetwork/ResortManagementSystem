package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.bellsoft.net/rms/api-core/internal/audit"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/mappers"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
)

var (
	ErrDateBlockNotFound       = errors.New("존재하지 않는 날짜 차단")
	ErrInvalidDateBlockRequest = errors.New("잘못된 날짜 차단 요청")
)

type DateBlockService interface {
	Create(ctx context.Context, req dto.CreateDateBlockRequest) (*dto.DateBlockResponse, error)
	Delete(ctx context.Context, id uint) error
	GetDateBlock(ctx context.Context, id uint) (*dto.DateBlockResponse, error)
	GetAll(ctx context.Context, filter dto.DateBlockFilter, page, size int) ([]dto.DateBlockResponse, int64, error)
	UpdateDateBlock(ctx context.Context, id uint, req dto.UpdateDateBlockRequest) (*dto.DateBlockResponse, error)
	IsDateRangeBlocked(ctx context.Context, startDate, endDate time.Time) (bool, error)
}

type dateBlockService struct {
	dateBlockRepo repositories.DateBlockRepository
	auditService  audit.AuditService
}

func NewDateBlockService(dateBlockRepo repositories.DateBlockRepository, auditService audit.AuditService) DateBlockService {
	return &dateBlockService{dateBlockRepo: dateBlockRepo, auditService: auditService}
}

func (s *dateBlockService) Create(ctx context.Context, req dto.CreateDateBlockRequest) (*dto.DateBlockResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}

	dateBlock := &models.DateBlock{
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    req.Reason,
	}

	created, err := s.dateBlockRepo.Create(ctx, dateBlock)
	if err != nil {
		return nil, err
	}

	result := mappers.ToDateBlockResponse(created)
	return &result, nil
}

func (s *dateBlockService) Delete(ctx context.Context, id uint) error {
	dateBlock, err := s.dateBlockRepo.FindByID(ctx, id)
	if err != nil {
		return ErrDateBlockNotFound
	}

	if err := s.dateBlockRepo.Delete(ctx, id); err != nil {
		return err
	}

	// Log deletion in audit — manual call required because soft delete bypasses GORM delete hooks
	_ = s.auditService.LogDelete(ctx, dateBlock)

	return nil
}

func (s *dateBlockService) GetAll(ctx context.Context, filter dto.DateBlockFilter, page, size int) ([]dto.DateBlockResponse, int64, error) {
	offset := page * size
	dateBlocks, total, err := s.dateBlockRepo.FindAll(ctx, filter, offset, size)
	if err != nil {
		return nil, 0, err
	}

	return mappers.ToDateBlockListResponse(dateBlocks), total, nil
}

func (s *dateBlockService) GetDateBlock(ctx context.Context, id uint) (*dto.DateBlockResponse, error) {
	dateBlock, err := s.dateBlockRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrDateBlockNotFound
	}

	result := mappers.ToDateBlockResponse(dateBlock)
	return &result, nil
}

func (s *dateBlockService) UpdateDateBlock(ctx context.Context, id uint, req dto.UpdateDateBlockRequest) (*dto.DateBlockResponse, error) {
	dateBlock, err := s.dateBlockRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrDateBlockNotFound
	}

	oldValues := dateBlock.GetAuditFields()

	if req.StartDate != nil {
		startDate, parseErr := time.Parse("2006-01-02", *req.StartDate)
		if parseErr != nil {
			return nil, fmt.Errorf("%w: invalid startDate format, expected YYYY-MM-DD", ErrInvalidDateBlockRequest)
		}
		dateBlock.StartDate = startDate
	}

	if req.EndDate != nil {
		endDate, parseErr := time.Parse("2006-01-02", *req.EndDate)
		if parseErr != nil {
			return nil, fmt.Errorf("%w: invalid endDate format, expected YYYY-MM-DD", ErrInvalidDateBlockRequest)
		}
		dateBlock.EndDate = endDate
	}

	if req.Reason != nil {
		dateBlock.Reason = *req.Reason
	}

	if dateBlock.StartDate.After(dateBlock.EndDate) {
		return nil, fmt.Errorf("%w: startDate must be before or equal to endDate", ErrInvalidDateBlockRequest)
	}

	if err := s.dateBlockRepo.Update(ctx, dateBlock); err != nil {
		return nil, err
	}

	if s.auditService != nil {
		_ = s.auditService.LogUpdate(ctx, dateBlock, oldValues)
	}

	result := mappers.ToDateBlockResponse(dateBlock)
	return &result, nil
}

func (s *dateBlockService) IsDateRangeBlocked(ctx context.Context, startDate, endDate time.Time) (bool, error) {
	return s.dateBlockRepo.IsDateRangeBlocked(ctx, startDate, endDate)
}
