package utils

import (
	"context"

	"gorm.io/gorm"
)

const userIDKey = "userID"

// SetUserIDInContext sets the user ID in the context for GORM operations
func SetUserIDInContext(db *gorm.DB, userID uint) *gorm.DB {
	ctx := context.WithValue(db.Statement.Context, userIDKey, userID)
	return db.WithContext(ctx)
}

// GetUserIDFromContext retrieves the user ID from the context
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(userIDKey).(uint)
	return userID, ok
}
