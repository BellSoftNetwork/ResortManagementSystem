package handlers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.bellsoft.net/rms/api-core/internal/dto"
	"gitlab.bellsoft.net/rms/api-core/internal/middleware"
	"gitlab.bellsoft.net/rms/api-core/internal/models"
	"gitlab.bellsoft.net/rms/api-core/internal/services"
	"gitlab.bellsoft.net/rms/api-core/pkg/response"
	"gitlab.bellsoft.net/rms/api-core/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "존재하지 않는 사용자")
		return
	}

	userResponse := h.toUserResponse(user)
	response.Success(c, userResponse)
}

func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	user, err := h.userService.Update(c.Request.Context(), userID, updates)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			response.Conflict(c, "이미 존재하는 이메일")
			return
		}
		response.InternalServerError(c, "사용자 정보 수정 실패")
		return
	}

	userResponse := h.toUserResponse(user)
	response.Success(c, userResponse)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "로그인 필요")
		return
	}

	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "존재하지 않는 사용자")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		response.BadRequest(c, "현재 비밀번호가 일치하지 않습니다")
		return
	}

	if err := h.userService.UpdatePassword(c.Request.Context(), userID, req.NewPassword); err != nil {
		response.InternalServerError(c, "비밀번호 변경 실패")
		return
	}

	response.Success(c, map[string]string{"message": "Password updated successfully"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var query dto.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "잘못된 쿼리 파라미터", err.Error())
		return
	}

	users, total, err := h.userService.GetAll(c.Request.Context(), query.Page, query.Size)
	if err != nil {
		response.InternalServerError(c, "사용자 목록 조회 실패")
		return
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = h.toUserResponse(&user)
	}

	totalPages := int(total) / query.Size
	if int(total)%query.Size > 0 {
		totalPages++
	}

	pagination := &response.Pagination{
		Page:          query.Page,
		Size:          query.Size,
		TotalPages:    totalPages,
		TotalElements: total,
	}

	// User 리스트는 현재 필터를 지원하지 않으므로 빈 필터 객체 반환
	filterResponse := map[string]interface{}{}

	response.SuccessListWithFilter(c, userResponses, pagination, filterResponse)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	user := &models.User{
		UserID:   req.UserID,
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	if req.Status != "" {
		switch req.Status {
		case "ACTIVE":
			user.Status = models.UserStatusActive
		case "INACTIVE":
			user.Status = models.UserStatusInactive
		}
	} else {
		user.Status = models.UserStatusActive
	}

	if req.Role != "" {
		switch req.Role {
		case "NORMAL":
			user.Role = models.UserRoleNormal
		case "ADMIN":
			user.Role = models.UserRoleAdmin
		case "SUPER_ADMIN":
			user.Role = models.UserRoleSuperAdmin
		}
	} else {
		user.Role = models.UserRoleNormal
	}

	if err := h.userService.Create(c.Request.Context(), user); err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			response.Conflict(c, "이미 존재하는 사용자")
			return
		}
		response.InternalServerError(c, "사용자 등록 실패")
		return
	}

	userResponse := h.toUserResponse(user)
	response.Created(c, userResponse)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "잘못된 사용자 ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "잘못된 요청", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Status != nil {
		switch *req.Status {
		case "ACTIVE":
			updates["status"] = models.UserStatusActive
		case "INACTIVE":
			updates["status"] = models.UserStatusInactive
		}
	}
	if req.Role != nil {
		switch *req.Role {
		case "NORMAL":
			updates["role"] = models.UserRoleNormal
		case "ADMIN":
			updates["role"] = models.UserRoleAdmin
		case "SUPER_ADMIN":
			updates["role"] = models.UserRoleSuperAdmin
		}
	}

	user, err := h.userService.Update(c.Request.Context(), uint(id), updates)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			response.NotFound(c, "존재하지 않는 사용자")
			return
		}
		if errors.Is(err, services.ErrUserAlreadyExists) {
			response.Conflict(c, "이미 존재하는 이메일")
			return
		}
		response.InternalServerError(c, "사용자 정보 수정 실패")
		return
	}

	userResponse := h.toUserResponse(user)
	response.Success(c, userResponse)
}

func (h *UserHandler) toUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:              user.ID,
		UserID:          user.UserID,
		Email:           user.Email,
		Name:            user.Name,
		Status:          user.Status.String(),
		Role:            user.Role.String(),
		ProfileImageURL: utils.GenerateGravatarURL(user.Email),
		CreatedAt:       dto.CustomTime{Time: user.CreatedAt},
		UpdatedAt:       dto.CustomTime{Time: user.UpdatedAt},
	}
}
