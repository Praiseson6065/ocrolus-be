package main

import (
	"Praiseson6065/ocrolus-be/handlers"
	"Praiseson6065/ocrolus-be/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRouter(r *gin.Engine) {
	authRoutes := r.Group("/auth")
	authHandler := &handlers.AuthHandler{}
	{
		authRoutes.POST("/signup", authHandler.UserSignup)
		authRoutes.POST("/login", authHandler.UserLogin)
	}
}

func ApiRouter(r *gin.Engine) {
	apiRoutes := r.Group("/api")

	// User routes
	userHandler := &handlers.UserHandler{}
	userRoutes := apiRoutes.Group("/user", middleware.Authenicator())
	{
		userRoutes.GET("/", userHandler.GetUser)
		userRoutes.PUT("/", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.DeleteUser)
	}
	
	// Article routes
	articleHandler := &handlers.ArticleHandler{}

	// Public article routes (no authentication required)
	articleRoutes := apiRoutes.Group("/articles", middleware.OptionalAuthenticator())
	{
		// Public endpoints for articles (read-only)
		articleRoutes.GET("", articleHandler.ListArticles)
		articleRoutes.GET("/:id", articleHandler.GetArticle)
	}

	// Protected article routes (authentication required)
	authArticleRoutes := apiRoutes.Group("/articles", middleware.Authenicator())
	{
		// Create, update, delete (require authentication)
		authArticleRoutes.POST("", articleHandler.CreateArticle)
		authArticleRoutes.PUT("/:id", articleHandler.UpdateArticle)
		authArticleRoutes.DELETE("/:id", articleHandler.DeleteArticle)

		// User's recently viewed articles
		authArticleRoutes.GET("/recently-viewed", articleHandler.GetRecentlyViewedArticles)
	}
}
