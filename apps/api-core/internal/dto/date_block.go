package dto

import (
	"fmt"
	"time"
)

type DateBlockResponse struct {
	ID        uint                 `json:"id"`
	StartDate JSONDate             `json:"startDate"`
	EndDate   JSONDate             `json:"endDate"`
	Reason    string               `json:"reason"`
	CreatedBy *UserSummaryResponse `json:"createdBy"`
	CreatedAt CustomTime           `json:"createdAt"`
}

type CreateDateBlockRequest struct {
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
	Reason    string `json:"reason" binding:"required,min=1,max=200"`
}

type UpdateDateBlockRequest struct {
	StartDate *string `json:"startDate"`
	EndDate   *string `json:"endDate"`
	Reason    *string `json:"reason" binding:"omitempty,min=1,max=200"`
}

func (r *CreateDateBlockRequest) Validate() error {
	startDate, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return fmt.Errorf("invalid startDate format, expected YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", r.EndDate)
	if err != nil {
		return fmt.Errorf("invalid endDate format, expected YYYY-MM-DD")
	}

	if startDate.After(endDate) {
		return fmt.Errorf("startDate must be before or equal to endDate")
	}

	return nil
}

type DateBlockFilter struct {
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}
