package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var allowedOrigins []string

func init() {
	godotenv.Load(".env")
	allowedOrigins = strings.Split(os.Getenv("ALLOWED_ORIGIN_DOMAINS"), ",")
}

func CORS(env *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if *env == "development" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else if len(allowedOrigins) != 0 {
			for _, origin := range allowedOrigins {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Plaid-Public-Token")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
