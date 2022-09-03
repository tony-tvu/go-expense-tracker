package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/entity"
	"gorm.io/gorm"
)

// Middleware restricts access to admin users only
func AdminRequired(db *gorm.DB) gin.HandlerFunc {
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

		if claims.UserType != string(entity.AdminUser) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
