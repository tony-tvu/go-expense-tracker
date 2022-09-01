package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

// Middleware restricts access to logged in users only
func AuthRequired(db *database.Db) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		cookie, err := c.Request.Cookie("goexpense_access")

		// no access token - make user log in
		if err != nil || cookie.Value == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		isExpired, claims, err := auth.GetClaimsWithValidation(cookie.Value)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// if access token has expired, check for existing session and renew token
		if *isExpired {

			// find existing session (refresh token)
			var s *models.Session
			err := db.Sessions.FindOne(
				ctx, bson.D{{Key: "email", Value: claims.Email}}).Decode(&s)

			// session not found - make user log in
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// validate refresh token
			isExpired, claims, err := auth.GetClaimsWithValidation(s.RefreshToken)
			if err != nil || *isExpired {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// renew access token
			renewed, err := auth.GetEncryptedAccessToken(ctx, claims.Email, claims.UserType)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "goexpense_access",
				Value:    renewed.Value,
				Expires:  renewed.ExpiresAt,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
		}
		c.Next()
	}
}
