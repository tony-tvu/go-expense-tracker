package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/util"
)

func CorsHeaders(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if util.Contains(&allowedOrigins, c.Request.Header.Get("Origin")) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "webhooksURL")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		c.Next()
	}
}
