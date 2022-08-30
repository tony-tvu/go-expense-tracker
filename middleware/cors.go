package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var subdomain string

func init() {
	godotenv.Load(".env")
	subdomain = os.Getenv("SUB_DOMAIN")
}

func CORS(env *string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if *env == "development" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		if subdomain != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", subdomain)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Plaid-Public-Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
