package service

import (
	"errors"
	"movie-review-api/models"
	"movie-review-api/repository"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	movieRepo  *repository.MovieRepository
}

func NewReviewService(reviewRepo *repository.ReviewRepository, movieRepo *repository.MovieRepository) *ReviewService {
	return &ReviewService{
		reviewRepo,
		movieRepo,
	}
}

func (s *ReviewService) CreateReview(review *models.Review) error {
	// Validate rating
	if review.Rating < 1 || review.Rating > 10 {
		return errors.New("rating must be between 1 and 10")
	}

	// Check if movie exists
	movie, err := s.movieRepo.FindByID(review.MovieID)
	if err != nil {
		return err
	}
	if movie == nil {
		return errors.New("movie not found")
	}

	return s.reviewRepo.Create(review)
}

func (s *ReviewService) GetReviewByID(id uint) (*models.Review, error) {
	review, err := s.reviewRepo.FindByID(id)
	if review == nil {
		return nil, errors.New("review not found")
	}
	return review, err
}

func (s *ReviewService) GetMovieReviews(movieID uint, page, limit int) ([]models.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Check if movie exists
	movie, err := s.movieRepo.FindByID(movieID)
	if err != nil {
		return nil, 0, err
	}
	if movie == nil {
		return nil, 0, errors.New("movie not found")
	}

	return s.reviewRepo.FindByMovieID(movieID, page, limit)
}

func (s *ReviewService) GetUserReviews(userID uint, page, limit int) ([]models.Review, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	return s.reviewRepo.FindByUserID(userID, page, limit)
}

func (s *ReviewService) GetAverageRating(movieID uint) (float64, error) {
	// Check if movie exists
	movie, err := s.movieRepo.FindByID(movieID)
	if err != nil {
		return 0, err
	}
	if movie == nil {
		return 0, errors.New("movie not found")
	}

	return s.reviewRepo.GetAverageRating(movieID)
}

func (s *ReviewService) UpdateReview(review *models.Review, userID uint) error {
	// Verify user owns the review
	existing, err := s.reviewRepo.FindByID(review.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("review not found")
	}
	if existing.UserID != userID {
		return errors.New("unauthorized: can only update own review")
	}

	// Validate rating
	if review.Rating < 1 || review.Rating > 10 {
		return errors.New("rating must be between 1 and 10")
	}

	return s.reviewRepo.Update(review)
}

func (s *ReviewService) DeleteReview(reviewID, userID uint) error {
	// Verify user owns the review
	review, err := s.reviewRepo.FindByID(reviewID)
	if err != nil {
		return err
	}
	if review == nil {
		return errors.New("review not found")
	}
	if review.UserID != userID {
		return errors.New("unauthorized: can only delete own review")
	}

	return s.reviewRepo.Delete(reviewID)
}
