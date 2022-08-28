package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

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

/*
Middleware sets following in response headers:
Expires: Thu, 01 Jan 1970 00:00:00 UTC
Cache-Control: no-cache, private, max-age=0
X-Accel-Expires: 0
Pragma: no-cache (for HTTP/1.0 proxies/clients)
*/
func NoCache(c *gin.Context) {
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
