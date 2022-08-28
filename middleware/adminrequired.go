package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/models"
)

// Middleware restricts access to admin users only
func AdminRequired(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("goexpense_access")
		if err != nil || cookie.Value == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// decrypt token
		decrypted, err := auth.Decrypt(encryptionKey, cookie.Value)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tkn, err := jwt.ParseWithClaims(decrypted, &auth.Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtKey), nil
			})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := tkn.Claims.(*auth.Claims)
		if claims.UserType != string(models.AdminUser) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
