package database

import (
	"Praiseson6065/ocrolus-be/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateArticle(ctx *gin.Context, article *models.Article) (string, error) {
	tx := db.WithContext(ctx).Create(article)
	if tx.Error != nil {
		return "", tx.Error
	}

	return article.ID, nil
}

func GetArticleByID(ctx *gin.Context, id string) (*models.Article, error) {
	var article models.Article
	result := db.WithContext(ctx).Preload("Author").Where("id = ?", id).First(&article)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("article not found")
		}
		return nil, result.Error
	}
	return &article, nil
}

func ListArticles(ctx *gin.Context, page, pageSize int, authorID string) ([]models.Article, int64, error) {
	var articles []models.Article
	var count int64
	query := db.WithContext(ctx).Model(&models.Article{})

	// Filter by author if specified
	if authorID != "" {
		query = query.Where("author_id = ?", authorID)
	}

	// Count total articles matching the filter
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and fetch articles with author information
	offset := (page - 1) * pageSize
	result := query.Preload("Author").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, count, nil
}

func ListPublishedArticles(ctx *gin.Context, page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var count int64
	query := db.WithContext(ctx).Model(&models.Article{}).Where("published = ?", true)

	// Count total published articles
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and fetch articles with author information
	offset := (page - 1) * pageSize
	result := query.Preload("Author").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&articles)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return articles, count, nil
}

func UpdateArticle(ctx *gin.Context, article *models.Article) (*models.Article, error) {
	var updatedArticle models.Article

	// Check if article exists
	result := db.WithContext(ctx).Where("id = ?", article.ID).First(&updatedArticle)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update article fields
	result = db.WithContext(ctx).Model(&updatedArticle).Updates(map[string]interface{}{
		"title":     article.Title,
		"content":   article.Content,
		"published": article.Published,
	})
	if result.Error != nil {
		return nil, result.Error
	}

	// Fetch the updated article with author information
	db.WithContext(ctx).Preload("Author").Where("id = ?", article.ID).First(&updatedArticle)
	return &updatedArticle, nil
}

func DeleteArticle(ctx *gin.Context, id string) error {
	result := db.WithContext(ctx).Where("id = ?", id).Delete(&models.Article{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("article not found")
	}
	return nil
}

func SaveRecentlyViewedArticle(ctx *gin.Context, userID, articleID string) error {
	// Check if the article exists first
	var article models.Article
	if err := db.WithContext(ctx).Where("id = ?", articleID).First(&article).Error; err != nil {
		return err
	}

	// Check if a record already exists for this user and article
	var existingRecord models.RecentlyViewedArticle
	result := db.WithContext(ctx).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		First(&existingRecord)

	if result.Error == nil {
		// Record exists, update ViewedAt
		return db.WithContext(ctx).
			Model(&existingRecord).
			Update("viewed_at", db.NowFunc()).
			Error
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Record doesn't exist, create a new one
		recentlyViewed := &models.RecentlyViewedArticle{
			UserID:    userID,
			ArticleID: articleID,
			ViewedAt:  db.NowFunc(),
		}
		return db.WithContext(ctx).Create(recentlyViewed).Error
	} else {
		// Some other error occurred
		return result.Error
	}
}

func GetRecentlyViewedArticles(ctx *gin.Context, userID string, limit int) ([]models.Article, error) {
	var recentArticles []models.Article

	// Use a subquery to get the most recent viewed articles by the user
	subQuery := db.WithContext(ctx).
		Model(&models.RecentlyViewedArticle{}).
		Select("article_id, MAX(viewed_at) as last_viewed").
		Where("user_id = ?", userID).
		Group("article_id").
		Order("last_viewed DESC").
		Limit(limit)

	// Join with articles to get the full article details
	err := db.WithContext(ctx).Table("(?) as rv", subQuery).
		Joins("JOIN articles ON articles.id = rv.article_id").
		Preload("Author").
		Find(&recentArticles).Error

	if err != nil {
		return nil, err
	}

	return recentArticles, nil
}
