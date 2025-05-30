package models

import (
	"time"

	"gorm.io/gorm"
)

type RecentlyViewedArticle struct {
	gorm.Model
	ID        string    `gorm:"primaryKey;<-:create" json:"id"`
	ArticleID string    `json:"article_id" gorm:"not null;index"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	Article   Article   `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ViewedAt  time.Time `json:"viewed_at" gorm:"not null;default:current_timestamp"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
