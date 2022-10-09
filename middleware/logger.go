package middleware

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Middleware logs request info
func Logger(env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if env == "development" {
			start := time.Now()
			defer func() {
				since := fmt.Sprintf("%sms", strconv.FormatInt(time.Since(start).Milliseconds(), 10))
				log.Printf("%-6s %-3d %-5s %s", c.Request.Method, c.Writer.Status(), since, c.Request.URL.Path)
			}()
		}
		c.Next()
	}
}
