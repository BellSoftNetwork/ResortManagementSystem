package mappers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
)

func TestToPaymentMethodResponse(t *testing.T) {
	t.Run("결제 수단 모델을 PaymentMethodResponse DTO로 변환한다", func(t *testing.T) {
		now := time.Now()
		paymentMethod := &models.PaymentMethod{
			BaseTimeEntity: models.BaseTimeEntity{
				BaseEntity: models.BaseEntity{ID: 1},
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			Name:                     "신용카드",
			CommissionRate:           3.5,
			RequireUnpaidAmountCheck: true,
			IsDefaultSelect:          true,
			Status:                   models.PaymentMethodStatusActive,
		}

		result := ToPaymentMethodResponse(paymentMethod)

		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "신용카드", result.Name)
		assert.Equal(t, 3.5, result.CommissionRate)
		assert.True(t, result.RequireUnpaidAmountCheck)
		assert.True(t, result.IsDefaultSelect)
		assert.Equal(t, "ACTIVE", result.Status)
	})
}
