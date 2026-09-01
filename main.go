package main

import (
	"log"
	"movie-review-api/config"
	"movie-review-api/handler"
	"movie-review-api/middleware"
	"movie-review-api/models"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db := config.ConnectDB(cfg)

	// Auto migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Token{},
		&models.Movie{},
		&models.Review{},
		&models.Category{},
	); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Initialize handlers
	authHandler := handler.NewAuthHandler(db, cfg)
	movieHandler := handler.NewMovieHandler(db)
	reviewHandler := handler.NewReviewHandler(db)

	// Public routes
	public := r.Group("/api/auth")
	{
		public.POST("/signup", authHandler.Signup)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes - Auth
	auth := r.Group("/api/auth")
	auth.Use(middleware.AuthMiddleware(cfg))
	{
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", authHandler.GetMe)
	}

	// Public movie routes
	movies := r.Group("/api/movies")
	{
		movies.GET("", movieHandler.ListMovies)
		movies.GET("/:id", movieHandler.GetMovieByID)
		movies.GET("/category/:category", movieHandler.GetMoviesByCategory)
		movies.GET("/search", movieHandler.SearchMovies)
		movies.GET("/top-rated", movieHandler.GetTopRatedMovies)
		movies.GET("/latest", movieHandler.GetLatestMovies)
		movies.GET("/:id/rating", movieHandler.GetAverageRating)
		movies.GET("/:id/reviews", reviewHandler.GetMovieReviews)
	}

	// Protected movie routes - Admin
	adminMovies := r.Group("/api/movies")
	adminMovies.Use(middleware.AuthMiddleware(cfg))
	{
		adminMovies.POST("", movieHandler.CreateMovie)
		adminMovies.PUT("/:id", movieHandler.UpdateMovie)
		adminMovies.DELETE("/:id", movieHandler.DeleteMovie)
	}

	// Protected review routes
	reviews := r.Group("/api/reviews")
	reviews.Use(middleware.AuthMiddleware(cfg))
	{
		reviews.GET("/:id", reviewHandler.GetReview)
		reviews.PUT("/:id", reviewHandler.UpdateReview)
		reviews.DELETE("/:id", reviewHandler.DeleteReview)
	}

	// Protected movie reviews
	movieReviews := r.Group("/api/movies/:id/reviews")
	movieReviews.Use(middleware.AuthMiddleware(cfg))
	{
		movieReviews.POST("", reviewHandler.CreateReview)
	}

	// User reviews
	userReviews := r.Group("/api/users/:id/reviews")
	{
		userReviews.GET("", reviewHandler.GetUserReviews)
	}

	// Categories
	categories := r.Group("/api/categories")
	{
		categories.GET("", movieHandler.GetCategories)
	}

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
