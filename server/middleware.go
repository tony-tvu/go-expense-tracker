package server

import (
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

var Middlewares = []Middleware{
	Logging(),
	RateLimit(),
	NoCache(),
}

// Chain applies multiple middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	// loop in reverse to preserve middleware order
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f
}

func Logging() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() { log.Println(r.URL.Path, r.Method, time.Since(start)) }()
			f(w, r)
		}
	}
}

// A Limiter controls how frequently events are allowed to happen. 
// It implements a "token bucket" of size b, initially full and refilled at rate r tokens per second.
var refillRatePerSecond rate.Limit = 25
var bucketSize = 100
var limiter = rate.NewLimiter(refillRatePerSecond, bucketSize)

func RateLimit() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			f(w, r)
		}
	}

}

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

// NoCache sets:
// Expires: Thu, 01 Jan 1970 00:00:00 UTC
// Cache-Control: no-cache, private, max-age=0
// X-Accel-Expires: 0
// Pragma: no-cache (for HTTP/1.0 proxies/clients)
func NoCache() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			// Delete any ETag headers that may have been set
			for _, v := range etagHeaders {
				if r.Header.Get(v) != "" {
					r.Header.Del(v)
				}
			}

			// Set our NoCache headers
			for k, v := range noCacheHeaders {
				w.Header().Set(k, v)
			}

			f(w, r)
		}
	}
}
