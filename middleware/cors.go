package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/util"
)

func CorsHeaders(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if util.Contains(&allowedOrigins, origin) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Plaid-Verification")
		}
		c.Next()
	}
}
