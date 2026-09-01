package service

import (
	"errors"
	"movie-review-api/models"
	"movie-review-api/repository"
)

type MovieService struct {
	movieRepo    *repository.MovieRepository
	categoryRepo *repository.CategoryRepository
}

func NewMovieService(movieRepo *repository.MovieRepository, categoryRepo *repository.CategoryRepository) *MovieService {
	return &MovieService{
		movieRepo,
		categoryRepo,
	}
}

func (s *MovieService) CreateMovie(movie *models.Movie) error {
	return s.movieRepo.Create(movie)
}

func (s *MovieService) GetMovieByID(id uint) (*models.Movie, error) {
	movie, err := s.movieRepo.FindByID(id)
	if movie == nil {
		return nil, errors.New("movie not found")
	}
	return movie, err
}

func (s *MovieService) ListMovies(page, limit int) ([]models.Movie, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.movieRepo.FindAll(page, limit)
}

func (s *MovieService) SearchMovies(query string, page, limit int) ([]models.Movie, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.movieRepo.Search(query, page, limit)
}

func (s *MovieService) GetMoviesByCategory(category string, page, limit int) ([]models.Movie, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.movieRepo.FindByCategory(category, page, limit)
}

func (s *MovieService) GetTopRatedMovies(limit int) ([]models.Movie, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.movieRepo.GetTopRated(limit)
}

func (s *MovieService) GetLatestMovies(limit int) ([]models.Movie, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return s.movieRepo.GetLatest(limit)
}

func (s *MovieService) UpdateMovie(movie *models.Movie) error {
	return s.movieRepo.Update(movie)
}

func (s *MovieService) DeleteMovie(id uint) error {
	return s.movieRepo.Delete(id)
}

func (s *MovieService) GetCategories() ([]models.Category, error) {
	return s.categoryRepo.FindAll()
}
