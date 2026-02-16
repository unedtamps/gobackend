package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/unedtamps/gobackend/pkg/utils"
)

func GetUserID(c *gin.Context) (utils.ULID, bool) {
	val, exists := c.Get("userID")
	if !exists {
		return utils.ULID{}, false
	}
	userID, ok := val.(utils.ULID)
	return userID, ok
}

func JWT(jwtGen *utils.JWTGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		claims, err := jwtGen.ParseClaims(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		userID, err := utils.ParseULID(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		c.Set("userID", userID)
		c.Set("userEmail", claims.Email)
		c.Next()
	}
}
