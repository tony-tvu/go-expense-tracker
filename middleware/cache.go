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

func FrontendCache(c *gin.Context) {
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

func NoCache(c *gin.Context) {
	var epoch = time.Unix(0, 0).Format(time.RFC1123)
	var noCacheHeaders = map[string]string{
		"Expires":         epoch,
		"Cache-Control":   "no-cache, no-store, no-transform, must-revalidate, private, max-age=0",
		"Pragma":          "no-cache",
		"X-Accel-Expires": "0",
	}

	var etagHeaders = []string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}

	// Delete any ETag headers that may have been set
	for _, v := range etagHeaders {
		if c.Request.Header.Get(v) != "" {
			c.Request.Header.Del(v)
		}
	}
	// Set NoCache headers
	for k, v := range noCacheHeaders {
		c.Writer.Header().Set(k, v)
	}

	c.Next()
}
