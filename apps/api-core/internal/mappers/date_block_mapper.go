package mappers

import (
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
)

func ToDateBlockResponse(model *models.DateBlock) dto.DateBlockResponse {
	response := dto.DateBlockResponse{
		ID:        model.ID,
		StartDate: dto.JSONDate{Time: model.StartDate},
		EndDate:   dto.JSONDate{Time: model.EndDate},
		Reason:    model.Reason,
		CreatedAt: dto.CustomTime{Time: model.CreatedAt},
	}

	if model.CreatedByUser != nil {
		email := ""
		if model.CreatedByUser.Email != nil {
			email = *model.CreatedByUser.Email
		}
		response.CreatedBy = &dto.UserSummaryResponse{
			ID:              model.CreatedByUser.ID,
			UserID:          model.CreatedByUser.UserID,
			Email:           email,
			Name:            model.CreatedByUser.Name,
			ProfileImageURL: "",
		}
	}

	return response
}

func ToDateBlockListResponse(models []models.DateBlock) []dto.DateBlockResponse {
	responses := make([]dto.DateBlockResponse, len(models))
	for i, model := range models {
		responses[i] = ToDateBlockResponse(&model)
	}
	return responses
}
