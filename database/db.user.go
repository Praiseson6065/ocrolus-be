package database

import (
	"Praiseson6065/ocrolus-be/models"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUser(ctx *gin.Context, user *models.User) (string, error) {

	tx := db.WithContext(ctx).Create(user)
	if tx.Error != nil {
		return "", tx.Error
	}

	return user.ID, nil
}

func GetUserByID(ctx *gin.Context, id string) (*models.User, error) {
	var user models.User
	result := db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByEmail(ctx *gin.Context, email string) (*models.User, error) {
	var user models.User
	result := db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

func GetPasswordByMail(ctx *gin.Context, email string) (string, string, error) {
	var user models.User
	result := db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", "", errors.New("user not found")
		}
		return "", "", result.Error
	}
	return user.Password, user.ID, nil
}

func UpdateUser(ctx *gin.Context, user *models.User) (models.User, error) {
	var updatedUser models.User

	// Check if user exists
	result := db.WithContext(ctx).Where("id = ?", user.ID).First(&updatedUser)
	if result.Error != nil {
		return updatedUser, result.Error
	}

	// Update user fields
	result = db.WithContext(ctx).Model(&updatedUser).Updates(user)
	if result.Error != nil {
		return updatedUser, result.Error
	}

	// Fetch the updated user
	db.WithContext(ctx).Where("id = ?", user.ID).First(&updatedUser)
	return updatedUser, nil
}

func DeleteUser(ctx *gin.Context, id string) error {
	result := db.WithContext(ctx).Where("id = ?", id).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
