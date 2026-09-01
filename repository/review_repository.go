package repository

import (
	"movie-review-api/models"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) FindByID(id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("User").Preload("Movie").First(&review, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &review, err
}

func (r *ReviewRepository) FindByMovieID(movieID uint, page, limit int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64
	offset := (page - 1) * limit

	err := r.db.Where("movie_id = ?", movieID).
		Count(&total).
		Preload("User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&reviews).Error

	return reviews, total, err
}

func (r *ReviewRepository) FindByUserID(userID uint, page, limit int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64
	offset := (page - 1) * limit

	err := r.db.Where("user_id = ?", userID).
		Count(&total).
		Preload("Movie").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&reviews).Error

	return reviews, total, err
}

func (r *ReviewRepository) GetAverageRating(movieID uint) (float64, error) {
	var avg float64
	err := r.db.Model(&models.Review{}).
		Where("movie_id = ?", movieID).
		Select("COALESCE(AVG(rating), 0)").
		Row().
		Scan(&avg)

	return avg, err
}

func (r *ReviewRepository) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

func (r *ReviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}

func (r *ReviewRepository) FindByMovieAndUser(movieID, userID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("movie_id = ? AND user_id = ?", movieID, userID).First(&review).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &review, err
}
