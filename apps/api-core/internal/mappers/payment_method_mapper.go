package mappers

import (
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
)

// ToPaymentMethodResponse converts a PaymentMethod model to PaymentMethodResponse DTO
func ToPaymentMethodResponse(paymentMethod *models.PaymentMethod) dto.PaymentMethodResponse {
	return dto.PaymentMethodResponse{
		ID:                       paymentMethod.ID,
		Name:                     paymentMethod.Name,
		CommissionRate:           paymentMethod.CommissionRate,
		RequireUnpaidAmountCheck: bool(paymentMethod.RequireUnpaidAmountCheck),
		IsDefaultSelect:          bool(paymentMethod.IsDefaultSelect),
		Status:                   paymentMethod.Status.String(),
		CreatedAt:                dto.CustomTime{Time: paymentMethod.CreatedAt},
		UpdatedAt:                dto.CustomTime{Time: paymentMethod.UpdatedAt},
	}
}
