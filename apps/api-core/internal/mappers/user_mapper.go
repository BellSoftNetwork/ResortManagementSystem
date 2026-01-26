package mappers

import (
	"context"

	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/utils"
	pkgutils "gitlab.bellsoft.net/rms/api-core/pkg/utils"
)

// ToUserResponse converts a User model to UserResponse DTO
func ToUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:              user.ID,
		UserID:          user.UserID,
		Email:           utils.StringPtrToString(user.Email),
		Name:            user.Name,
		Status:          user.Status.String(),
		Role:            user.Role.String(),
		ProfileImageURL: pkgutils.GenerateGravatarURL(utils.StringPtrToString(user.Email)),
		CreatedAt:       dto.CustomTime{Time: user.CreatedAt},
		UpdatedAt:       dto.CustomTime{Time: user.UpdatedAt},
	}
}

// ToUserSummaryResponse converts a User model to UserSummaryResponse DTO
func ToUserSummaryResponse(user *models.User) dto.UserSummaryResponse {
	return dto.UserSummaryResponse{
		ID:              user.ID,
		UserID:          user.UserID,
		Email:           utils.StringPtrToString(user.Email),
		Name:            user.Name,
		ProfileImageURL: pkgutils.GenerateGravatarURL(utils.StringPtrToString(user.Email)),
	}
}

// GetUserSummaryHelper creates a helper function for retrieving user summaries
func GetUserSummaryHelper(getUserByID func(ctx context.Context, userID uint) (*models.User, error)) GetUserSummaryFunc {
	return func(ctx context.Context, userID uint) *dto.UserSummaryResponse {
		if userID == 0 {
			return nil
		}
		if user, err := getUserByID(ctx, userID); err == nil {
			summary := ToUserSummaryResponse(user)
			return &summary
		}
		return &dto.UserSummaryResponse{
			ID: userID,
		}
	}
}
