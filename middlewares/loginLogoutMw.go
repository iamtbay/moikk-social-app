package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("accessToken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Please login first.",
			})
			c.Abort()
		}
		c.Next()
	}
}

//CHECK USER LOGGED OUT OR NOT.

func LogoutMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("accessToken")
		if err != nil {
			c.Next()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User already logged in!",
		})
		c.Abort()
	}
}
