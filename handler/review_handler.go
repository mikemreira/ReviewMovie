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

type ReviewHandler struct {
	reviewService *service.ReviewService
	db            *gorm.DB
}

func NewReviewHandler(db *gorm.DB) *ReviewHandler {
	reviewRepo := repository.NewReviewRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	reviewService := service.NewReviewService(reviewRepo, movieRepo)

	return &ReviewHandler{
		reviewService,
		db,
	}
}

type CreateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=10"`
	Comment string `json:"comment"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"min=1,max=10"`
	Comment string `json:"comment"`
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	movieID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid movie ID"))
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("user not authenticated"))
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request"))
		return
	}

	review := &models.Review{
		MovieID:   uint(movieID),
		UserID:    userID.(uint),
		Rating:    req.Rating,
		Comment:   req.Comment,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.reviewService.CreateReview(review); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse(review))
}

func (h *ReviewHandler) GetReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid review ID"))
		return
	}

	review, err := h.reviewService.GetReviewByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	if review == nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("review not found"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(review))
}

func (h *ReviewHandler) GetMovieReviews(c *gin.Context) {
	movieID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid movie ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.reviewService.GetMovieReviews(uint(movieID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.PaginatedResponse{
		Success: true,
		Data:    reviews,
		Page:    page,
		Limit:   limit,
		Total:   total,
	})
}

func (h *ReviewHandler) GetUserReviews(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid user ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.reviewService.GetUserReviews(uint(userID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.PaginatedResponse{
		Success: true,
		Data:    reviews,
		Page:    page,
		Limit:   limit,
		Total:   total,
	})
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid review ID"))
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("user not authenticated"))
		return
	}

	var req UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid request"))
		return
	}

	review, err := h.reviewService.GetReviewByID(uint(id))
	if err != nil || review == nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("review not found"))
		return
	}

	if req.Rating > 0 {
		review.Rating = req.Rating
	}
	if req.Comment != "" {
		review.Comment = req.Comment
	}
	review.UpdatedAt = time.Now()

	if err := h.reviewService.UpdateReview(review, userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(review))
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("invalid review ID"))
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("user not authenticated"))
		return
	}

	if err := h.reviewService.DeleteReview(uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessMessage("review deleted successfully"))
}
