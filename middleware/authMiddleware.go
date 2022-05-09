package middleware

import (
	"net/http"
	"restaurantProject/helpers"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from header
		tokenHeader := c.Request.Header.Get("token")
		if tokenHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			c.Abort()
		}
		// Token is valid
		claim, err := helpers.ValidateToken(tokenHeader)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Unauthorized",
			})
			c.Abort()
		}
		c.Set("email", claim.Email)
		c.Set("firstName", claim.FirstName)
		c.Set("lastName", claim.LastName)
		c.Set("uid", claim.Uid)
		c.Next()
	}
}
