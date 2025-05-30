package handlers

import (
	"net/http"

	"Praiseson6065/ocrolus-be/database"
	"Praiseson6065/ocrolus-be/middleware"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	id := middleware.GetUserID(ctx)

	user, err := database.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	id := middleware.GetUserID(ctx)
	// Check if user exists
	existingUser, err := database.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update fields if provided
	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.Email != "" {
		// Check if new email is already taken by another user
		if req.Email != existingUser.Email {
			userWithEmail, _ := database.GetUserByEmail(ctx, req.Email)
			if userWithEmail != nil {
				ctx.JSON(http.StatusConflict, gin.H{"error": "Email is already in use"})
				return
			}
		}
		existingUser.Email = req.Email
	}
	if req.Password != "" {
		existingUser.Password = req.Password // Again, password should be hashed in real app
	}

	updatedUser, err := database.UpdateUser(ctx, existingUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	response := UserResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	id := middleware.GetUserID(ctx)

	if err := database.DeleteUser(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
