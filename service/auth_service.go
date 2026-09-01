package service

import (
	"errors"
	"movie-review-api/models"
	"movie-review-api/repository"
	"movie-review-api/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.TokenRepository
	jwtUtil   *utils.JWTUtil
}

func NewAuthService(userRepo *repository.UserRepository, tokenRepo *repository.TokenRepository, jwtUtil *utils.JWTUtil) *AuthService {
	return &AuthService{
		userRepo,
		tokenRepo,
		jwtUtil,
	}
}

func (s *AuthService) Signup(username, email, password string) (*models.User, error) {
	// Check if user exists
	if s.userRepo.ExistsByUsername(username) {
		return nil, errors.New("username already exists")
	}

	if s.userRepo.ExistsByEmail(email) {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	// Find user
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	// Generate token
	token, err := s.jwtUtil.GenerateToken(user)
	if err != nil {
		return "", err
	}

	// Store token in database
	tokenEntity := &models.Token{
		UserID:    user.ID,
		Token:     token,
		TokenType: "BEARER",
		Expired:   false,
		Revoked:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.tokenRepo.Create(tokenEntity); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Logout(tokenString string) error {
	return s.tokenRepo.RevokeToken(tokenString)
}

func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	// Validate token
	claims, err := s.jwtUtil.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Check if token is revoked
	if !s.tokenRepo.IsTokenValid(tokenString) {
		return nil, errors.New("token is revoked")
	}

	// Extract user ID from claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Get user
	user, err := s.userRepo.FindByID(uint(userID))
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
