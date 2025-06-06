package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenicator() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		authorization := ctx.GetHeader("Authorization")
		if len(authorization) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		fields := strings.Fields(authorization)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is invalid"})
			return
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization type is invalid"})
		}

		encodedToken := fields[1]
		userId, err := ValidateToken(encodedToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid Token"})
		}

		ctx.Set("userId", userId)

		ctx.Next()

	}

}

func OptionalAuthenticator() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		authorization := ctx.GetHeader("Authorization")
		if len(authorization) > 0 {
			fields := strings.Fields(authorization)
			if len(fields) == 2 && strings.ToLower(fields[0]) == "bearer" {
				encodedToken := fields[1]
				userId, err := ValidateToken(encodedToken)
				if err == nil {
					ctx.Set("userId", userId)
				}
			}
		}
		ctx.Next()
	}

}

func GetUserID(ctx *gin.Context) string {
	return ctx.GetString("userId")
}
