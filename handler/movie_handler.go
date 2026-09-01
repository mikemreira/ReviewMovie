package handler

import (
	"net/http"
	"strconv"
	"time"

	"movie-review-api/models"
	"movie-review-api/repository"
	"movie-review-api/service"
	"movie-review-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MovieHandler struct {
	movieService *service.MovieService
}

func NewMovieHandler(db *gorm.DB) *MovieHandler {
	movieRepo := repository.NewMovieRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	movieService := service.NewMovieService(movieRepo, categoryRepo)

	return &MovieHandler{movieService}
}

type CreateMovieRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date" binding:"required"` // Format: YYYY-MM-DD
	Director    string `json:"director"`
	Duration    int    `json:"duration"`
	Poster      string `json:"poster"`
	Categories  []uint `json:"categories"`
}

type UpdateMovieRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
	Director    string `json:"director"`
	Duration    int    `json:"duration"`
	Poster      string `json:"poster"`
	Categories  []uint `json:"categories"`
}

func (h *MovieHandler) ListMovies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	movies, total, err := h.movieService.ListMovies(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.PaginatedResponse{
		Success: true,
		Data:    movies,
		Page:    page,
		Limit:   limit,
		Total:   total,
	})
}

func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid movie ID"))
		return
	}

	movie, err := h.movieService.GetMovieByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	if movie == nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("movie not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(movie))
}

func (h *MovieHandler) SearchMovies(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("search query required"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	movies, total, err := h.movieService.SearchMovies(query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.PaginatedResponse{
		Success: true,
		Data:    movies,
		Page:    page,
		Limit:   limit,
		Total:   total,
	})
}

func (h *MovieHandler) GetMoviesByCategory(c *gin.Context) {
	category := c.Param("category")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	movies, total, err := h.movieService.GetMoviesByCategory(category, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.PaginatedResponse{
		Success: true,
		Data:    movies,
		Page:    page,
		Limit:   limit,
		Total:   total,
	})
}

func (h *MovieHandler) GetTopRatedMovies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	movies, err := h.movieService.GetTopRatedMovies(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(movies))
}

func (h *MovieHandler) GetLatestMovies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	movies, err := h.movieService.GetLatestMovies(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(movies))
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var req CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request"))
		return
	}

	releaseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid release date format, use YYYY-MM-DD"))
		return
	}

	movie := &models.Movie{
		Title:       req.Title,
		Description: req.Description,
		ReleaseDate: releaseDate,
		Director:    req.Director,
		Duration:    req.Duration,
		Poster:      req.Poster,
	}

	if err := h.movieService.CreateMovie(movie); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse(movie))
}

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid movie ID"))
		return
	}

	var req UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request"))
		return
	}

	movie, err := h.movieService.GetMovieByID(uint(id))
	if err != nil || movie == nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("movie not found"))
		return
	}

	if req.Title != "" {
		movie.Title = req.Title
	}
	if req.Description != "" {
		movie.Description = req.Description
	}
	if req.ReleaseDate != "" {
		releaseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
		if err == nil {
			movie.ReleaseDate = releaseDate
		}
	}
	if req.Director != "" {
		movie.Director = req.Director
	}
	if req.Duration > 0 {
		movie.Duration = req.Duration
	}
	if req.Poster != "" {
		movie.Poster = req.Poster
	}

	if err := h.movieService.UpdateMovie(movie); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(movie))
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid movie ID"))
		return
	}

	if err := h.movieService.DeleteMovie(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessMessage("movie deleted successfully"))
}

func (h *MovieHandler) GetCategories(c *gin.Context) {
	categories, err := h.movieService.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(categories))
}

func (h *MovieHandler) GetAverageRating(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid movie ID"))
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	reviewRepo := repository.NewReviewRepository(db)
	reviewService := service.NewReviewService(reviewRepo, repository.NewMovieRepository(db))

	rating, err := reviewService.GetAverageRating(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(gin.H{"average_rating": rating}))
}
