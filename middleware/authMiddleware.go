package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	helper "test-go-project/helpers"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprint("No authorization header provided")})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error"})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("email", claims.FirstName)
		c.Set("email", claims.LastName)
		c.Set("email", claims.Uid)
		c.Set("email", claims.UserType)
		c.Next()
	}
}
