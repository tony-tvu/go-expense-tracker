package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
)

var encryptionKey string
var jwtKey string

func init() {
	godotenv.Load(".env")
	encryptionKey = os.Getenv("ENCRYPTION_KEY")
	jwtKey = os.Getenv("JWT_KEY")
	if util.ContainsEmpty(encryptionKey, jwtKey) {
		log.Fatal("auth keys are missing")
	}
}

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

		// decrypt access token
		decrypted, err := auth.Decrypt(cookie.Value)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		accessTkn, err := jwt.ParseWithClaims(decrypted, &auth.Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtKey), nil
			})
		claims := accessTkn.Claims.(*auth.Claims)

		// token is expired and missing correct claims - make user log in
		if err != nil && claims.Email == "" && claims.UserType == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// if token is invalid/expired - check for existing session and renew access token
		if !accessTkn.Valid {

			// find existing session (refresh_token)
			var s *models.Session
			err := db.Sessions.FindOne(
				ctx, bson.D{{Key: "email", Value: claims.Email}}).Decode(&s)

			// session not found - make user log in
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// decrypt refresh token
			decrypted, err = auth.Decrypt(s.RefreshToken)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// verify session is still valid
			_, err = jwt.ParseWithClaims(decrypted, &auth.Claims{},
				func(token *jwt.Token) (interface{}, error) {
					return []byte(jwtKey), nil
				})

			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// renew access token
			renewed, err := auth.CreateAccessToken(ctx, claims.Email, claims.UserType)
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
