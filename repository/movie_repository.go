package repository

import (
	"movie-review-api/models"

	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db}
}

func (r *MovieRepository) Create(movie *models.Movie) error {
	return r.db.Create(movie).Error
}

func (r *MovieRepository) FindByID(id uint) (*models.Movie, error) {
	var movie models.Movie
	err := r.db.Preload("Categories").Preload("Reviews").First(&movie, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &movie, err
}

func (r *MovieRepository) FindAll(page, limit int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64
	offset := (page - 1) * limit

	err := r.db.Model(&models.Movie{}).Count(&total).
		Preload("Categories").
		Preload("Reviews").
		Offset(offset).
		Limit(limit).
		Find(&movies).Error

	return movies, total, err
}

func (r *MovieRepository) Search(query string, page, limit int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64
	offset := (page - 1) * limit

	err := r.db.Model(&models.Movie{}).
		Where("title ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&total).
		Preload("Categories").
		Offset(offset).
		Limit(limit).
		Find(&movies).Error

	return movies, total, err
}

func (r *MovieRepository) FindByCategory(categoryName string, page, limit int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64
	offset := (page - 1) * limit

	err := r.db.Model(&models.Movie{}).
		Joins("JOIN movie_categories ON movie_categories.movie_id = movies.id").
		Joins("JOIN categories ON categories.id = movie_categories.category_id").
		Where("categories.name ILIKE ?", "%"+categoryName+"%").
		Count(&total).
		Preload("Categories").
		Offset(offset).
		Limit(limit).
		Find(&movies).Error

	return movies, total, err
}

func (r *MovieRepository) GetTopRated(limit int) ([]models.Movie, error) {
	var movies []models.Movie
	err := r.db.Model(&models.Movie{}).
		Joins("LEFT JOIN reviews ON reviews.movie_id = movies.id").
		Group("movies.id").
		Order("AVG(reviews.rating) DESC").
		Limit(limit).
		Preload("Categories").
		Find(&movies).Error

	return movies, err
}

func (r *MovieRepository) GetLatest(limit int) ([]models.Movie, error) {
	var movies []models.Movie
	err := r.db.Order("created_at DESC").
		Limit(limit).
		Preload("Categories").
		Find(&movies).Error

	return movies, err
}

func (r *MovieRepository) Update(movie *models.Movie) error {
	return r.db.Save(movie).Error
}

func (r *MovieRepository) Delete(id uint) error {
	return r.db.Delete(&models.Movie{}, id).Error
}
