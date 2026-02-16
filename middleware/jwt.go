package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unedtamps/gobackend/pkg/utils"
)

func JWT(j *utils.JWTGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token string from request
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(401, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		// handle Bearer
		token = token[7:]
		// parse token
		claims, err := j.ParseClaims(token)
		if err != nil {
			c.JSON(401, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
