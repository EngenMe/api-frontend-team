package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/EngenMe/api-frontend-team/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	return (func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(tokenString, " ")
		fmt.Println(tokenParts)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication token 1"})
			ctx.Abort()
			return
		}

		tokenString = tokenParts[1]

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication token 2"})
			ctx.Abort()
			return
		}
		ctx.Set("userId", claims["user_id"])
		ctx.Next()
	})

}
