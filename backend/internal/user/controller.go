package user

import (
	"net/http"

	"backend-travelku/internal/helper"
	"backend-travelku/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Controller adalah HTTP handler untuk fitur user.
type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// GetUser godoc
// @Summary     Get user by ID
// @Description Retrieve a user by id
// @Tags        Users
// @Produce     json
// @Security    BearerAuth
// @Param       id path     string true "User ID"
// @Success     200 {object} helper.SuccessResponse
// @Failure     404 {object} helper.ErrorResponse
// @Router      /users/{id} [get]
func (h *Controller) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.service.GetUserByID(id)
	if err != nil {
		logger.Warn("User not found: " + id)
		helper.ErrorJSON(c, http.StatusNotFound, "User not found", "")
		return
	}

	response := ToUserResponse(user.ID, user.Name, user.Email, user.Phone, user.IsActive, user.LastLogin, user.CreatedAt, user.UpdatedAt)
	helper.SuccessJSON(c, http.StatusOK, "User retrieved successfully", response)
}

// GetAllUsers godoc
// @Summary     List users
// @Description Retrieve all users
// @Tags        Users
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} helper.SuccessResponse
// @Failure     500 {object} helper.ErrorResponse
// @Router      /users [get]
func (h *Controller) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		logger.Error("Failed to get users: " + err.Error())
		helper.ErrorJSON(c, http.StatusInternalServerError, "Failed to get users", err.Error())
		return
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(user.ID, user.Name, user.Email, user.Phone, user.IsActive, user.LastLogin, user.CreatedAt, user.UpdatedAt)
	}

	response := UserListResponse{Count: len(userResponses), Users: userResponses}
	helper.SuccessJSON(c, http.StatusOK, "Users retrieved successfully", response)
}

// UpdateProfile godoc
// @Summary     Update user profile
// @Description Update name, email, and phone
// @Tags        Users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string               true "User ID"
// @Param       request body     UpdateProfileRequest true "Update profile payload"
// @Success     200 {object} helper.SuccessResponse
// @Failure     400 {object} helper.ErrorResponse
// @Router      /users/{id} [put]
func (h *Controller) UpdateProfile(c *gin.Context) {
	id := c.Param("id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request body: " + err.Error())
		helper.ErrorJSON(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := h.service.UpdateProfile(id, &req)
	if err != nil {
		logger.Error("Failed to update profile: " + err.Error())
		helper.ErrorJSON(c, http.StatusBadRequest, "Failed to update profile", err.Error())
		return
	}

	response := ToUserResponse(user.ID, user.Name, user.Email, user.Phone, user.IsActive, user.LastLogin, user.CreatedAt, user.UpdatedAt)
	helper.SuccessJSON(c, http.StatusOK, "Profile updated successfully", response)
}

// ChangePassword godoc
// @Summary     Change password
// @Description Change user's password
// @Tags        Users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id      path     string                true "User ID"
// @Param       request body     ChangePasswordRequest true "Change password payload"
// @Success     200 {object} helper.SuccessResponse
// @Failure     400 {object} helper.ErrorResponse
// @Router      /users/{id}/change-password [post]
func (h *Controller) ChangePassword(c *gin.Context) {
	id := c.Param("id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request body: " + err.Error())
		helper.ErrorJSON(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	if err := h.service.ChangePassword(id, &req); err != nil {
		logger.Warn("Failed to change password: " + err.Error())
		helper.ErrorJSON(c, http.StatusBadRequest, "Failed to change password", err.Error())
		return
	}

	helper.SuccessEmptyJSON(c, http.StatusOK, "Password changed successfully")
}

// DeleteUser godoc
// @Summary     Delete user
// @Description Delete user by id
// @Tags        Users
// @Produce     json
// @Security    BearerAuth
// @Param       id path     string true "User ID"
// @Success     200 {object} helper.SuccessResponse
// @Failure     404 {object} helper.ErrorResponse
// @Router      /users/{id} [delete]
func (h *Controller) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteUser(id); err != nil {
		logger.Error("Failed to delete user: " + err.Error())
		helper.ErrorJSON(c, http.StatusNotFound, "Failed to delete user", err.Error())
		return
	}

	helper.SuccessEmptyJSON(c, http.StatusOK, "User deleted successfully")
}
