package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
)

// Middleware restricts access to admin users only
func AdminRequired(db *database.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("goexpense_refresh")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateTokenAndGetClaims(cookie.Value)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserType != string(models.AdminUser) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
