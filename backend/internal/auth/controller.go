package auth

import (
	"net/http"

	"backend-travelku/internal/helper"
	userfeature "backend-travelku/internal/user"
	"backend-travelku/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Controller adalah HTTP handler untuk fitur auth.
type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

// Register godoc
// @Summary     Register user
// @Description Create a new user account
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body     RegisterRequest true "Register payload"
// @Success     201     {object} helper.SuccessResponse{data=user.UserResponse}
// @Failure     400     {object} helper.ErrorResponse
// @Router      /auth/register [post]
func (h *Controller) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid request body: " + err.Error())
		helper.ErrorJSON(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		logger.Error("Failed to register user: " + err.Error())
		helper.ErrorJSON(c, http.StatusBadRequest, "Failed to register user", err.Error())
		return
	}

	response := userfeature.ToUserResponse(user.ID, user.Name, user.Email, user.Phone, user.IsActive, user.LastLogin, user.CreatedAt, user.UpdatedAt)
	helper.SuccessJSON(c, http.StatusCreated, "User registered successfully", response)
}

// Login godoc
// @Summary     Login
// @Description Authenticate user and generate JWT token
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body     LoginRequest true "Login payload"
// @Success     200     {object} helper.SuccessResponse{data=auth.LoginResponse}
// @Failure     400     {object} helper.ErrorResponse
// @Failure     401     {object} helper.ErrorResponse
// @Router      /auth/login [post]
func (h *Controller) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid login request: " + err.Error())
		helper.ErrorJSON(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	token, user, err := h.service.Login(&req)
	if err != nil {
		logger.Warn("Login failed: " + err.Error())
		helper.ErrorJSON(c, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	userResp := userfeature.ToUserResponse(user.ID, user.Name, user.Email, user.Phone, user.IsActive, user.LastLogin, user.CreatedAt, user.UpdatedAt)
	response := LoginResponse{
		Token:     token,
		ExpiresIn: 24,
		TokenType: "Bearer",
		User:      userResp,
	}

	helper.SuccessJSON(c, http.StatusOK, "Login successful", response)
}
