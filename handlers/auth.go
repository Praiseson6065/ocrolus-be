package handlers

import (
	"Praiseson6065/ocrolus-be/database"
	"Praiseson6065/ocrolus-be/middleware"
	"Praiseson6065/ocrolus-be/models"
	"Praiseson6065/ocrolus-be/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct{}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type UserSignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) UserLogin(ctx *gin.Context) {

	var loginRequest LoginRequest
	if err := ctx.ShouldBindBodyWithJSON(&loginRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	hashedPwd, userId, err := database.GetPasswordByMail(ctx, loginRequest.Email)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !util.ComparePasswords(hashedPwd, loginRequest.Password) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid password",
		})
		return
	}
	token, err := middleware.GenerateToken(userId)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *AuthHandler) UserSignup(ctx *gin.Context) {
	var userSignupRequest UserSignupRequest

	if err := ctx.ShouldBindBodyWithJSON(&userSignupRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPwd := util.HashAndSalt(userSignupRequest.Password)

	id, err := database.CreateUser(ctx, &models.User{
		Name:     userSignupRequest.Name,
		Email:    userSignupRequest.Email,
		Password: hashedPwd,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Successfully signed up", "userId": id})

}
