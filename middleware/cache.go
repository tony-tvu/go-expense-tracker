package middleware

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var maxAgeStr string
var maxAge int

func init() {
	godotenv.Load(".env")
	maxAgeStr = os.Getenv("MAX_AGE")

	maxAgeInt, err := strconv.Atoi(maxAgeStr)
	if err != nil {
		// 24 hour default
		maxAge = 86400
	} else {
		maxAge = maxAgeInt
	}
}

func CacheControl(c *gin.Context) {
	epoch := time.Now().Add(time.Duration(maxAge) * time.Second).Format(time.RFC1123)
	cacheHeaders := map[string]string{
		"Expires":       epoch,
		"Cache-Control": fmt.Sprintf("must-revalidate, private, max-age=%d", maxAge),
	}

	// Set Cache headers
	for k, v := range cacheHeaders {
		c.Writer.Header().Set(k, v)
	}

	c.Next()
}
