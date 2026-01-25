package services

import (
	"context"
	"errors"

	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("존재하지 않는 사용자")
)

type UserService interface {
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetAll(ctx context.Context, page, size int) ([]models.User, int64, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.User, error)
	UpdatePassword(ctx context.Context, id uint, newPassword string) error
	IsUpdatableAccount(ctx context.Context, requestUser *models.User, targetUserID uint) (bool, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) GetAll(ctx context.Context, page, size int) ([]models.User, int64, error) {
	offset := page * size
	return s.userRepo.FindAll(ctx, offset, size)
}

func (s *userService) Create(ctx context.Context, user *models.User) error {
	if user.Email != nil && *user.Email != "" {
		exists, err := s.userRepo.ExistsByEmail(ctx, *user.Email)
		if err != nil {
			return err
		}
		if exists {
			return ErrUserAlreadyExists
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = "{bcrypt}" + string(hashedPassword)

	return s.userRepo.Create(ctx, user)
}

func (s *userService) Update(ctx context.Context, id uint, updates map[string]interface{}) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if name, ok := updates["name"].(string); ok {
		user.Name = name
	}

	if email, ok := updates["email"].(string); ok {
		if email != "" {
			exists, err := s.userRepo.ExistsByEmail(ctx, email)
			if err != nil {
				return nil, err
			}
			if exists && (user.Email == nil || *user.Email != email) {
				return nil, ErrUserAlreadyExists
			}
			user.Email = &email
		} else {
			user.Email = nil
		}
	}

	if status, ok := updates["status"].(models.UserStatus); ok {
		user.Status = status
	}

	if role, ok := updates["role"].(models.UserRole); ok {
		user.Role = role
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdatePassword(ctx context.Context, id uint, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = "{bcrypt}" + string(hashedPassword)

	return s.userRepo.Update(ctx, user)
}

func (s *userService) IsUpdatableAccount(ctx context.Context, requestUser *models.User, targetUserID uint) (bool, error) {
	if requestUser.Role == models.UserRoleSuperAdmin {
		return true, nil
	}

	targetUser, err := s.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return false, ErrUserNotFound
	}

	return targetUser.Role < models.UserRoleAdmin, nil
}
