package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/util"
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
		} else if len(allowedOrigins) >= 1 && allowedOrigins[0] != "" {
			if util.Contains(&allowedOrigins, c.Request.Header.Get("Origin")) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
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
