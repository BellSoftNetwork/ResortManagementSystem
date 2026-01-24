package dto

type PaymentMethodResponse struct {
	ID                       uint       `json:"id"`
	Name                     string     `json:"name"`
	CommissionRate           float64    `json:"commissionRate"`
	RequireUnpaidAmountCheck bool       `json:"requireUnpaidAmountCheck"`
	IsDefaultSelect          bool       `json:"isDefaultSelect"`
	Status                   string     `json:"status"`
	CreatedAt                CustomTime `json:"createdAt"`
	UpdatedAt                CustomTime `json:"updatedAt"`
}

type CreatePaymentMethodRequest struct {
	Name                     string  `json:"name" binding:"required,min=2,max=20"`
	CommissionRate           float64 `json:"commissionRate" binding:"min=0,max=1"`
	RequireUnpaidAmountCheck bool    `json:"requireUnpaidAmountCheck"`
}

type UpdatePaymentMethodRequest struct {
	Name                     *string  `json:"name" binding:"omitempty,min=1,max=20"`
	CommissionRate           *float64 `json:"commissionRate" binding:"omitempty,min=0,max=1"`
	RequireUnpaidAmountCheck *bool    `json:"requireUnpaidAmountCheck"`
	IsDefaultSelect          *bool    `json:"isDefaultSelect"`
	Status                   *string  `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}
