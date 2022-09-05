package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/util"
	"gorm.io/gorm"
)

// Middleware restricts access to logged in users only
func AuthRequired(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshCookie, err := c.Request.Cookie("goexpense_refresh")

		// refresh_token missing - make user log in (all requests must have a refresh token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// validate refresh_token from request
		refreshClaims, err := auth.ValidateTokenAndGetClaims(refreshCookie.Value)

		// refresh_token has expired
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		_, err = c.Request.Cookie("goexpense_access")

		// handle expired or missing access_token
		if err != nil {

			// find existing session
			var s *entity.Session
			if result := db.Where("username = ?", refreshClaims.Username).First(&s); result.Error != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// verify token from db session matches request's token
			if s.RefreshToken != refreshCookie.Value {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// validate refresh_token
			_, err := auth.ValidateTokenAndGetClaims(s.RefreshToken)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// renew access_token
			renewed, err := auth.GetEncryptedToken(auth.AccessToken, refreshClaims.Username, refreshClaims.UserType)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			util.SetCookie(c.Writer, "goexpense_access", renewed.Value, renewed.ExpiresAt)
		}

		c.Next()
	}
}
