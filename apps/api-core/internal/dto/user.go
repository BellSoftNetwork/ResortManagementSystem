package dto

type UserResponse struct {
	ID              uint       `json:"id"`
	UserID          string     `json:"userId"`
	Email           string     `json:"email"`
	Name            string     `json:"name"`
	Status          string     `json:"status"`
	Role            string     `json:"role"`
	ProfileImageURL string     `json:"profileImageUrl"`
	CreatedAt       CustomTime `json:"createdAt"`
	UpdatedAt       CustomTime `json:"updatedAt"`
}

type CreateUserRequest struct {
	UserID   string `json:"userId" binding:"required,min=4,max=30"`
	Email    string `json:"email" binding:"omitempty,email,max=100"`
	Name     string `json:"name" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required,min=8,max=20"`
	Status   string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	Role     string `json:"role" binding:"omitempty,oneof=NORMAL ADMIN SUPER_ADMIN"`
}

type UpdateUserRequest struct {
	Email  string  `json:"email" binding:"omitempty,email,max=100"`
	Name   *string `json:"name" binding:"omitempty,min=2,max=20"`
	Status *string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	Role   *string `json:"role" binding:"omitempty,oneof=NORMAL ADMIN SUPER_ADMIN"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

// UserSummaryResponse is a simplified user response for references
type UserSummaryResponse struct {
	ID              uint   `json:"id"`
	UserID          string `json:"userId"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	ProfileImageURL string `json:"profileImageUrl"`
}
