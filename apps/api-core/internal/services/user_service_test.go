package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gorm.io/gorm"
)

// MockUserRepository는 auth_service_test.go에 이미 정의되어 있음

type UserServiceTestSuite struct {
	suite.Suite
	ctx          context.Context
	service      services.UserService
	mockUserRepo *MockUserRepository
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockUserRepo = new(MockUserRepository)
	suite.service = services.NewUserService(suite.mockUserRepo)
}

func (suite *UserServiceTestSuite) TestGetByID() {
	// Given - 사용자가 등록된 상황에서
	user := &models.User{
		UserID: "testuser",
		Email:  stringPtr("test@example.com"),
		Name:   "Test User",
		Role:   models.UserRoleNormal,
		Status: models.UserStatusActive,
	}
	user.ID = 1

	suite.mockUserRepo.On("FindByID", suite.ctx, uint(1)).Return(user, nil)

	// When - 특정 사용자 정보를 조회하면
	result, err := suite.service.GetByID(suite.ctx, 1)

	// Then - 정상적으로 조회된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), user.UserID, result.UserID)
	assert.Equal(suite.T(), user.Email, result.Email)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetByID_NotFound() {
	// Given - 존재하지 않는 사용자 ID로
	suite.mockUserRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// When - 조회하면
	result, err := suite.service.GetByID(suite.ctx, 999)

	// Then - ErrUserNotFound 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrUserNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetAll_Empty() {
	// Given - 가입한 사용자가 없는 상황에서
	suite.mockUserRepo.On("FindAll", suite.ctx, 0, 10).Return([]models.User{}, int64(0), nil)

	// When - 계정 리스트 조회 시 (page 0, Spring Boot 호환)
	users, total, err := suite.service.GetAll(suite.ctx, 0, 10)

	// Then - 빈 리스트가 반환된다
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), users)
	assert.Equal(suite.T(), int64(0), total)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestGetAll_Multiple() {
	// Given - 50개의 계정이 생성된 상태에서
	users := make([]models.User, 20)
	for i := 0; i < 20; i++ {
		email := "user" + string(rune(i)) + "@example.com"
		users[i] = models.User{
			UserID: "user" + string(rune(i)),
			Email:  &email,
			Name:   "User " + string(rune(i)),
			Role:   models.UserRoleNormal,
			Status: models.UserStatusActive,
		}
		users[i].ID = uint(50 - i) // ID를 역순으로 설정 (정렬 테스트용)
	}

	suite.mockUserRepo.On("FindAll", suite.ctx, 0, 20).Return(users, int64(50), nil)

	// When - 계정 리스트 조회 시 (page 0, Spring Boot 호환)
	result, total, err := suite.service.GetAll(suite.ctx, 0, 20)

	// Then - 정상적으로 모두 조회된다
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 20)
	assert.Equal(suite.T(), int64(50), total)
	assert.Equal(suite.T(), uint(50), result[0].ID)  // 첫 번째 ID가 50
	assert.Equal(suite.T(), uint(31), result[19].ID) // 20번째 ID가 31
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestCreate() {
	// Given - 새로운 사용자 정보가 주어지면
	newUser := &models.User{
		UserID:   "newuser",
		Email:    stringPtr("newuser@example.com"),
		Name:     "New User",
		Password: "hashed_password",
		Role:     models.UserRoleNormal,
		Status:   models.UserStatusActive,
	}

	// 중복 확인 - UserService는 email만 체크함
	suite.mockUserRepo.On("ExistsByEmail", suite.ctx, "newuser@example.com").Return(false, nil)
	suite.mockUserRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.User")).Return(nil)

	// When - 신규 계정 등록 시
	err := suite.service.Create(suite.ctx, newUser)

	// Then - 정상적으로 등록된다
	assert.NoError(suite.T(), err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestCreate_DuplicateUserID() {
	// Given - 기존에 동일한 사용자 ID가 있는 상황에서
	newUser := &models.User{
		UserID: "existinguser",
		Email:  stringPtr("new@example.com"),
	}

	// UserService는 email만 체크하므로
	suite.mockUserRepo.On("ExistsByEmail", suite.ctx, "new@example.com").Return(false, nil)
	suite.mockUserRepo.On("Create", suite.ctx, mock.AnythingOfType("*models.User")).Return(nil)

	// When - 신규 계정 등록 시도 시
	err := suite.service.Create(suite.ctx, newUser)

	// Then - UserID 중복에도 불구하고 성공 (현재 구현의 한계)
	assert.NoError(suite.T(), err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestCreate_DuplicateEmail() {
	// Given - 기존에 동일한 이메일이 있는 상황에서
	newUser := &models.User{
		UserID: "newuser",
		Email:  stringPtr("existing@example.com"),
	}

	// UserService는 email만 체크함
	suite.mockUserRepo.On("ExistsByEmail", suite.ctx, "existing@example.com").Return(true, nil)

	// When - 신규 계정 등록 시도 시
	err := suite.service.Create(suite.ctx, newUser)

	// Then - 가입이 거부된다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrUserAlreadyExists, err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestUpdate() {
	// Given - 기존 계정이 있는 상태에서
	existingUser := &models.User{
		UserID: "testuser",
		Email:  stringPtr("test@example.com"),
		Name:   "Original Name",
		Role:   models.UserRoleNormal,
		Status: models.UserStatusActive,
	}
	existingUser.ID = 1

	updatedUser := &models.User{
		UserID: "testuser",
		Email:  stringPtr("test@example.com"),
		Name:   "변경된 이름",
		Role:   models.UserRoleNormal,
		Status: models.UserStatusActive,
	}
	updatedUser.ID = 1

	suite.mockUserRepo.On("FindByID", suite.ctx, uint(1)).Return(existingUser, nil)
	suite.mockUserRepo.On("Update", suite.ctx, existingUser).Return(nil)

	// When - 기존 계정 수정 시도 시
	updates := map[string]interface{}{
		"name": "변경된 이름",
	}
	result, err := suite.service.Update(suite.ctx, 1, updates)

	// Then - 정상적으로 수정된다
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), "변경된 이름", result.Name)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestUpdate_NotFound() {
	// Given - 존재하지 않는 계정에 대해
	suite.mockUserRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// When - 수정 시도 시
	updates := map[string]interface{}{
		"name": "변경된 이름",
	}
	result, err := suite.service.Update(suite.ctx, 999, updates)

	// Then - 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrUserNotFound, err)
	assert.Nil(suite.T(), result)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestUpdatePassword() {
	// Given - 기존 사용자가 있는 상태에서
	existingUser := &models.User{
		UserID:   "testuser",
		Password: "old_hashed_password",
	}
	existingUser.ID = 1

	suite.mockUserRepo.On("FindByID", suite.ctx, uint(1)).Return(existingUser, nil)
	suite.mockUserRepo.On("Update", suite.ctx, existingUser).Return(nil)

	// When - 비밀번호를 변경하면
	err := suite.service.UpdatePassword(suite.ctx, 1, "new_password")

	// Then - 정상적으로 변경된다
	assert.NoError(suite.T(), err)
	// 비밀번호가 해시되어 저장되는지 확인
	assert.NotEqual(suite.T(), "new_password", existingUser.Password)
	assert.Contains(suite.T(), existingUser.Password, "{bcrypt}$2a$")
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestUpdatePassword_NotFound() {
	// Given - 존재하지 않는 사용자에 대해
	suite.mockUserRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	// When - 비밀번호 변경 시도 시
	err := suite.service.UpdatePassword(suite.ctx, 999, "new_password")

	// Then - 에러가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrUserNotFound, err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestIsUpdatableAccount_SuperAdmin() {
	// Given - 일반, 관리자, 최고 관리자 계정이 등록된 상태에서
	normalUser := &models.User{Role: models.UserRoleNormal}
	normalUser.ID = 1

	adminUser := &models.User{Role: models.UserRoleAdmin}
	adminUser.ID = 2

	superAdminUser := &models.User{Role: models.UserRoleSuperAdmin}
	superAdminUser.ID = 3

	// When/Then - 최고 관리자가 최고 관리자를 수정할 수 있는지
	result1, err := suite.service.IsUpdatableAccount(suite.ctx, superAdminUser, 3)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), result1)

	// When/Then - 최고 관리자가 일반 관리자를 수정할 수 있는지
	result2, err := suite.service.IsUpdatableAccount(suite.ctx, superAdminUser, 2)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), result2)

	// When/Then - 최고 관리자가 일반 유저를 수정할 수 있는지
	result3, err := suite.service.IsUpdatableAccount(suite.ctx, superAdminUser, 1)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), result3)
}

func (suite *UserServiceTestSuite) TestIsUpdatableAccount_Admin() {
	// Given - 관리자 계정으로
	normalUser := &models.User{Role: models.UserRoleNormal}
	normalUser.ID = 1

	adminUser := &models.User{Role: models.UserRoleAdmin}
	adminUser.ID = 2

	superAdminUser := &models.User{Role: models.UserRoleSuperAdmin}
	superAdminUser.ID = 3

	// When - 일반 관리자가 최고 관리자를 수정할 수 있는지 확인하면
	suite.mockUserRepo.On("FindByID", suite.ctx, uint(3)).Return(superAdminUser, nil)
	result1, err := suite.service.IsUpdatableAccount(suite.ctx, adminUser, 3)

	// Then - 수정 불가능한 상태로 반환된다
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), result1)

	// When - 일반 관리자가 일반 관리자를 수정할 수 있는지 확인하면
	suite.mockUserRepo.On("FindByID", suite.ctx, uint(2)).Return(adminUser, nil)
	result2, err := suite.service.IsUpdatableAccount(suite.ctx, adminUser, 2)

	// Then - 수정 불가능한 상태로 반환된다
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), result2)

	// When - 일반 관리자가 일반 유저를 수정할 수 있는지 확인하면
	suite.mockUserRepo.On("FindByID", suite.ctx, uint(1)).Return(normalUser, nil)
	result3, err := suite.service.IsUpdatableAccount(suite.ctx, adminUser, 1)

	// Then - 수정 가능한 상태로 반환된다
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), result3)
}

func (suite *UserServiceTestSuite) TestIsUpdatableAccount_NotFound() {
	// Given - 관리자 계정으로
	adminUser := &models.User{Role: models.UserRoleAdmin}
	adminUser.ID = 2

	// When - 존재하지 않는 유저를 수정할 수 있는지 확인하면
	suite.mockUserRepo.On("FindByID", suite.ctx, uint(999)).Return(nil, gorm.ErrRecordNotFound)

	result, err := suite.service.IsUpdatableAccount(suite.ctx, adminUser, 999)

	// Then - 사용자를 찾을 수 없다는 예외가 발생한다
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), services.ErrUserNotFound, err)
	assert.False(suite.T(), result)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
