package internal

import (
	"net/http"
	"strings"

	"github.com/JonasLindermayr/FileBeam/types"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized!"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized!"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		JWT, err := DecodeToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		var foundToken types.Session
		if err := DB.First(&foundToken, "token = ?", tokenString).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not valid"})
			c.Abort()
			return
		}

		if JWT.UUID != foundToken.UserID {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not valid"})
			c.Abort()
			return
		}

		c.Set("uuid", JWT.UUID)
		c.Next()
	}
}
