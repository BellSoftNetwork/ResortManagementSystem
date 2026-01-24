package audit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetUserContext(t *testing.T) {
	t.Run("사용자 ID와 사용자명을 컨텍스트에 설정할 수 있다", func(t *testing.T) {
		ctx := context.Background()
		userID := uint(123)
		username := "testuser"

		newCtx := SetUserContext(ctx, &userID, username)

		userCtx := GetUserContext(newCtx)
		assert.Equal(t, userID, userCtx.UserID)
		assert.Equal(t, username, userCtx.Username)
	})

	t.Run("사용자 ID가 nil인 경우도 처리할 수 있다", func(t *testing.T) {
		ctx := context.Background()
		username := "testuser"

		newCtx := SetUserContext(ctx, nil, username)

		userCtx := GetUserContext(newCtx)
		assert.Equal(t, uint(0), userCtx.UserID)
		assert.Equal(t, username, userCtx.Username)
	})
}

func TestGetUserContext(t *testing.T) {
	t.Run("컨텍스트에서 사용자 정보를 가져올 수 있다", func(t *testing.T) {
		ctx := context.Background()
		userID := uint(456)
		username := "anotheruser"

		ctx = SetUserContext(ctx, &userID, username)
		userCtx := GetUserContext(ctx)

		assert.Equal(t, userID, userCtx.UserID)
		assert.Equal(t, username, userCtx.Username)
	})

	t.Run("사용자 컨텍스트가 없으면 빈 컨텍스트를 반환한다", func(t *testing.T) {
		ctx := context.Background()
		userCtx := GetUserContext(ctx)

		assert.Equal(t, uint(0), userCtx.UserID)
		assert.Empty(t, userCtx.Username)
	})
}

func TestHasUserContext(t *testing.T) {
	t.Run("사용자 컨텍스트가 있으면 true를 반환한다", func(t *testing.T) {
		ctx := context.Background()
		userID := uint(789)
		username := "user"

		ctx = SetUserContext(ctx, &userID, username)
		hasContext := HasUserContext(ctx)

		assert.True(t, hasContext)
	})

	t.Run("사용자 컨텍스트가 없으면 false를 반환한다", func(t *testing.T) {
		ctx := context.Background()
		hasContext := HasUserContext(ctx)

		assert.False(t, hasContext)
	})
}
