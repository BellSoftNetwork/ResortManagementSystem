package services

import (
	"context"
	"errors"

	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
)

var (
	ErrPaymentMethodNotFound   = errors.New("존재하지 않는 결제 수단")
	ErrPaymentMethodNameExists = errors.New("이미 존재하는 결제 수단 이름")
	ErrPaymentMethodInUse      = errors.New("사용 중인 결제 수단")
)

type PaymentMethodService interface {
	GetByID(ctx context.Context, id uint) (*models.PaymentMethod, error)
	GetAll(ctx context.Context, page, size int, sort string) ([]models.PaymentMethod, int64, error)
	GetActive(ctx context.Context) ([]models.PaymentMethod, error)
	Create(ctx context.Context, paymentMethod *models.PaymentMethod) error
	Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.PaymentMethod, error)
	Delete(ctx context.Context, id uint) error
}

type paymentMethodService struct {
	paymentMethodRepo repositories.PaymentMethodRepository
}

func NewPaymentMethodService(paymentMethodRepo repositories.PaymentMethodRepository) PaymentMethodService {
	return &paymentMethodService{
		paymentMethodRepo: paymentMethodRepo,
	}
}

func (s *paymentMethodService) GetByID(ctx context.Context, id uint) (*models.PaymentMethod, error) {
	paymentMethod, err := s.paymentMethodRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrPaymentMethodNotFound
	}
	return paymentMethod, nil
}

func (s *paymentMethodService) GetAll(ctx context.Context, page, size int, sort string) ([]models.PaymentMethod, int64, error) {
	offset := page * size
	return s.paymentMethodRepo.FindAll(ctx, offset, size, sort)
}

func (s *paymentMethodService) GetActive(ctx context.Context) ([]models.PaymentMethod, error) {
	return s.paymentMethodRepo.FindActive(ctx)
}

func (s *paymentMethodService) Create(ctx context.Context, paymentMethod *models.PaymentMethod) error {
	exists, err := s.paymentMethodRepo.ExistsByName(ctx, paymentMethod.Name, nil)
	if err != nil {
		return err
	}
	if exists {
		return ErrPaymentMethodNameExists
	}

	_, err = s.paymentMethodRepo.Create(ctx, paymentMethod)
	return err
}

func (s *paymentMethodService) Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.PaymentMethod, error) {
	paymentMethod, err := s.paymentMethodRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrPaymentMethodNotFound
	}

	if name, ok := updates["name"].(string); ok {
		if name != paymentMethod.Name {
			exists, err := s.paymentMethodRepo.ExistsByName(ctx, name, &id)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, ErrPaymentMethodNameExists
			}
			paymentMethod.Name = name
		}
	}

	if commissionRate, ok := updates["commissionRate"].(float64); ok {
		paymentMethod.CommissionRate = commissionRate
	}

	if requireUnpaidAmountCheck, ok := updates["requireUnpaidAmountCheck"].(models.BitBool); ok {
		paymentMethod.RequireUnpaidAmountCheck = requireUnpaidAmountCheck
	}

	if isDefaultSelect, ok := updates["isDefaultSelect"].(models.BitBool); ok {
		// 기본 선택을 true로 변경하는 경우, 먼저 모든 다른 결제 수단의 기본 선택을 해제
		if bool(isDefaultSelect) {
			if err := s.paymentMethodRepo.ResetAllDefaultSelects(ctx); err != nil {
				return nil, err
			}
		}
		paymentMethod.IsDefaultSelect = isDefaultSelect
	}

	if status, ok := updates["status"].(models.PaymentMethodStatus); ok {
		paymentMethod.Status = status
	}

	if err := s.paymentMethodRepo.Update(ctx, paymentMethod); err != nil {
		return nil, err
	}

	return paymentMethod, nil
}

func (s *paymentMethodService) Delete(ctx context.Context, id uint) error {
	_, err := s.paymentMethodRepo.FindByID(ctx, id)
	if err != nil {
		return ErrPaymentMethodNotFound
	}

	return s.paymentMethodRepo.Delete(ctx, id)
}
