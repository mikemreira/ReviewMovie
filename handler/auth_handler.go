package handler

import (
	"net/http"

	"movie-review-api/config"
	"movie-review-api/models"
	"movie-review-api/repository"
	"movie-review-api/service"
	"movie-review-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret, cfg.JWTExpirationMs)
	authService := service.NewAuthService(userRepo, tokenRepo, jwtUtil)

	return &AuthHandler{authService}
}

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request"))
		return
	}

	user, err := h.authService.Signup(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse(user))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request"))
		return
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
		return
	}

	// Get user info
	user, _ := repository.NewUserRepository(c.MustGet("db").(*gorm.DB)).FindByUsername(req.Username)

	response := LoginResponse{
		Token: token,
		User:  user,
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(response))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("authorization header missing"))
		return
	}

	token := authHeader[7:] // Remove "Bearer " prefix
	if err := h.authService.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessMessage("logged out successfully"))
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("user not found in token"))
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("authorization header missing"))
		return
	}

	token := authHeader[7:] // Remove "Bearer " prefix
	user, err := h.authService.GetUserFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error()))
		return
	}

	// Verify the user ID matches
	if user.ID != userID.(uint) {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("unauthorized"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(user))
}
