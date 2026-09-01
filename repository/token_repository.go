package repository

import (
	"movie-review-api/models"

	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db}
}

func (r *TokenRepository) Create(token *models.Token) error {
	return r.db.Create(token).Error
}

func (r *TokenRepository) FindByToken(tokenString string) (*models.Token, error) {
	var token models.Token
	err := r.db.Where("token = ?", tokenString).Preload("User").First(&token).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &token, err
}

func (r *TokenRepository) RevokeToken(tokenString string) error {
	return r.db.Model(&models.Token{}).Where("token = ?", tokenString).Update("revoked", true).Error
}

func (r *TokenRepository) IsTokenValid(tokenString string) bool {
	var token models.Token
	err := r.db.Where("token = ? AND revoked = false AND expired = false", tokenString).First(&token).Error
	return err == nil
}

func (r *TokenRepository) DeleteExpiredTokens() error {
	return r.db.Where("expired = true").Delete(&models.Token{}).Error
}

func (r *TokenRepository) FindByUserID(userID uint) ([]models.Token, error) {
	var tokens []models.Token
	err := r.db.Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}
