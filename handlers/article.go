package handlers

import (
	"net/http"
	"strconv"

	"Praiseson6065/ocrolus-be/database"
	"Praiseson6065/ocrolus-be/middleware"
	"Praiseson6065/ocrolus-be/models"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct{}

type CreateArticleRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Published bool   `json:"published"`
}

type UpdateArticleRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	Published bool   `json:"published"`
}

type ArticleResponse struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Published bool         `json:"published"`
	Author    UserResponse `json:"author,omitempty"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

// CreateArticle handles the creation of a new article
func (h *ArticleHandler) CreateArticle(ctx *gin.Context) {
	// Get user ID from context (set by authenticator middleware)
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	article := &models.Article{
		Title:     req.Title,
		Content:   req.Content,
		AuthorID:  userID,
		Published: req.Published,
	}

	createdArticleID, err := database.CreateArticle(ctx, article)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"articleId": createdArticleID,
	})
}

// GetArticle handles fetching a single article by ID
func (h *ArticleHandler) GetArticle(ctx *gin.Context) {
	id := ctx.Param("id")

	article, err := database.GetArticleByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// Record the view if user is authenticated
	userID := middleware.GetUserID(ctx)
	if userID != "" {
		// Ignoring errors for recently viewed as it's not critical
		_ = database.SaveRecentlyViewedArticle(ctx, userID, id)
	}

	response := ArticleResponse{
		ID:        article.ID,
		Title:     article.Title,
		Content:   article.Content,
		Published: article.Published,
		Author: UserResponse{
			ID:    article.Author.ID,
			Name:  article.Author.Name,
		},
		CreatedAt: article.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: article.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	ctx.JSON(http.StatusOK, response)
}

// ListArticles handles fetching a paginated list of articles
func (h *ArticleHandler) ListArticles(ctx *gin.Context) {
	// Get pagination parameters from query
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")
	onlyMine := ctx.Query("onlyMine")
	publishedOnly := ctx.Query("publishedOnly")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Get the authenticated user's ID
	userID := middleware.GetUserID(ctx)

	var articles []models.Article
	var total int64

	// Determine which articles to fetch
	if publishedOnly == "true" {
		// Public route - only show published articles
		articles, total, err = database.ListPublishedArticles(ctx, page, pageSize)
	} else if onlyMine == "true" && userID != "" {
		// User's own articles (published and unpublished)
		articles, total, err = database.ListArticles(ctx, page, pageSize, userID)
	} else if userID != "" {
		// Authenticated user can see all published articles + their own unpublished ones
		// For simplicity, we'll just show all published articles here
		articles, total, err = database.ListPublishedArticles(ctx, page, pageSize)
	} else {
		// Unauthenticated users can only see published articles
		articles, total, err = database.ListPublishedArticles(ctx, page, pageSize)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles: " + err.Error()})
		return
	}

	// Map to response objects
	responseArticles := make([]ArticleResponse, len(articles))
	for i, article := range articles {
		responseArticles[i] = ArticleResponse{
			ID:        article.ID,
			Title:     article.Title,
			Content:   article.Content,
			Published: article.Published,
			Author: UserResponse{
				ID:    article.Author.ID,
				Name:  article.Author.Name,
				Email: article.Author.Email,
			},
			CreatedAt: article.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: article.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"articles":    responseArticles,
		"totalCount":  total,
		"currentPage": page,
		"pageSize":    pageSize,
	})
}

// UpdateArticle handles updating an existing article
func (h *ArticleHandler) UpdateArticle(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := middleware.GetUserID(ctx)

	// Check if article exists
	existingArticle, err := database.GetArticleByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// Verify that the user is the author of the article
	if existingArticle.AuthorID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this article"})
		return
	}

	var req UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update fields if provided
	if req.Title != "" {
		existingArticle.Title = req.Title
	}
	if req.Content != "" {
		existingArticle.Content = req.Content
	}
	// Published is a boolean, so we always update it from the request
	existingArticle.Published = req.Published

	updatedArticle, err := database.UpdateArticle(ctx, existingArticle)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article: " + err.Error()})
		return
	}

	response := ArticleResponse{
		ID:        updatedArticle.ID,
		Title:     updatedArticle.Title,
		Content:   updatedArticle.Content,
		Published: updatedArticle.Published,
		Author: UserResponse{
			ID:    updatedArticle.Author.ID,
			Name:  updatedArticle.Author.Name,
			Email: updatedArticle.Author.Email,
		},
		CreatedAt: updatedArticle.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: updatedArticle.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteArticle handles the deletion of an article
func (h *ArticleHandler) DeleteArticle(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := middleware.GetUserID(ctx)

	// Check if article exists and user is the author
	existingArticle, err := database.GetArticleByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// Verify that the user is the author of the article
	if existingArticle.AuthorID != userID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this article"})
		return
	}

	if err := database.DeleteArticle(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article: " + err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetRecentlyViewedArticles returns a list of recently viewed articles for the authenticated user
func (h *ArticleHandler) GetRecentlyViewedArticles(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 20 {
		limit = 5
	}

	articles, err := database.GetRecentlyViewedArticles(ctx, userID, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve recently viewed articles: " + err.Error()})
		return
	}

	// Map to response objects
	responseArticles := make([]ArticleResponse, len(articles))
	for i, article := range articles {
		responseArticles[i] = ArticleResponse{
			ID:        article.ID,
			Title:     article.Title,
			Content:   article.Content,
			Published: article.Published,
			Author: UserResponse{
				ID:    article.Author.ID,
				Name:  article.Author.Name,
				Email: article.Author.Email,
			},
			CreatedAt: article.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: article.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"recentlyViewed": responseArticles,
	})
}
