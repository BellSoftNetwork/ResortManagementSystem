package mappers

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
)

func TestToUserResponse(t *testing.T) {
	t.Run("사용자 모델을 UserResponse DTO로 변환한다", func(t *testing.T) {
		email := "test@example.com"
		now := time.Now()
		user := &models.User{
			BaseTimeEntity: models.BaseTimeEntity{
				BaseEntity: models.BaseEntity{ID: 1},
				CreatedAt:  now,
				UpdatedAt:  now,
			},
			UserID: "testuser",
			Email:  &email,
			Name:   "Test User",
			Status: models.UserStatusActive,
			Role:   models.UserRoleNormal,
		}

		result := ToUserResponse(user)

		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "testuser", result.UserID)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, "Test User", result.Name)
		assert.Equal(t, "ACTIVE", result.Status)
		assert.Equal(t, "NORMAL", result.Role)
		assert.NotEmpty(t, result.ProfileImageURL)
	})
}

func TestToUserSummaryResponse(t *testing.T) {
	t.Run("사용자 모델을 UserSummaryResponse DTO로 변환한다", func(t *testing.T) {
		email := "test@example.com"
		user := &models.User{
			BaseTimeEntity: models.BaseTimeEntity{
				BaseEntity: models.BaseEntity{ID: 1},
			},
			UserID: "testuser",
			Email:  &email,
			Name:   "Test User",
		}

		result := ToUserSummaryResponse(user)

		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "testuser", result.UserID)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, "Test User", result.Name)
		assert.NotEmpty(t, result.ProfileImageURL)
	})
}

func TestGetUserSummaryHelper(t *testing.T) {
	t.Run("userID가 0이면 nil을 반환한다", func(t *testing.T) {
		getUserByID := func(ctx context.Context, userID uint) (*models.User, error) {
			return nil, nil
		}
		helper := GetUserSummaryHelper(getUserByID)

		result := helper(context.Background(), 0)

		assert.Nil(t, result)
	})

	t.Run("사용자를 찾으면 UserSummaryResponse를 반환한다", func(t *testing.T) {
		email := "test@example.com"
		getUserByID := func(ctx context.Context, userID uint) (*models.User, error) {
			return &models.User{
				BaseTimeEntity: models.BaseTimeEntity{
					BaseEntity: models.BaseEntity{ID: userID},
				},
				UserID: "testuser",
				Email:  &email,
				Name:   "Test User",
			}, nil
		}
		helper := GetUserSummaryHelper(getUserByID)

		result := helper(context.Background(), 1)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "testuser", result.UserID)
		assert.Equal(t, "test@example.com", result.Email)
	})

	t.Run("사용자를 찾지 못하면 ID만 포함된 응답을 반환한다", func(t *testing.T) {
		getUserByID := func(ctx context.Context, userID uint) (*models.User, error) {
			return nil, assert.AnError
		}
		helper := GetUserSummaryHelper(getUserByID)

		result := helper(context.Background(), 1)

		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Empty(t, result.UserID)
	})
}
