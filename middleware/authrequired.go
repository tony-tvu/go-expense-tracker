package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

// Middleware restricts access to logged in users only
func AuthRequired(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		cookie, err := c.Request.Cookie("goexpense_access")

		// no access token - make user log in
		if err != nil || cookie.Value == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// decrypt access token
		decrypted, err := auth.Decrypt(a.EncryptionKey, cookie.Value)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tkn, err := jwt.ParseWithClaims(decrypted, &auth.Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(a.JwtKey), nil
			})
		claims := tkn.Claims.(*auth.Claims)

		// token is expired and missing correct claims - make user log in
		if err != nil && claims.Email == "" && claims.UserType == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// if token is invalid/expired - check for existing session and renew access token
		if !tkn.Valid {

			// find existing session (refresh_token)
			var s *models.Session
			err := a.Sessions.FindOne(
				ctx, bson.D{{Key: "email", Value: claims.Email}}).Decode(&s)

			// session not found - make user log in
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// decrypt refresh token
			decrypted, err = auth.Decrypt(a.EncryptionKey, s.RefreshToken)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// verify session is still valid
			_, err = jwt.ParseWithClaims(decrypted, &auth.Claims{},
				func(token *jwt.Token) (interface{}, error) {
					return []byte(a.JwtKey), nil
				})

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// renew access token
			renewed, err := auth.CreateAccessToken(ctx, a, claims.Email, claims.UserType)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "goexpense_access",
				Value:    renewed.Value,
				Expires:  renewed.ExpiresAt,
				HttpOnly: true,
			})
		}

		c.Next()
	}
}
