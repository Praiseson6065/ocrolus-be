package models

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = "U" + strings.Replace(uuid.New().String(), "-", "", -1)
	return
}

func (article *Article) BeforeCreate(tx *gorm.DB) (err error) {
	article.ID = "AR" + strings.Replace(uuid.New().String(), "-", "", -1)
	return
}

func (rva *RecentlyViewedArticle) BeforeCreate(tx *gorm.DB) (err error) {
	rva.ID = "RV" + strings.Replace(uuid.New().String(), "-", "", -1)
	return
}
