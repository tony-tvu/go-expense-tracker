package server

import (
	"log"
	"net/http"
	"time"

	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/time/rate"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

var SharedMiddlewares = []Middleware{
	Logging(),
	RateLimit(),
	NoCache(),
}

// Append additional middlewares along with SharedMiddlewares
func UseMiddlewares(additional ...Middleware) []Middleware {
	m := SharedMiddlewares
	m = append(m, additional...)
	return m
}

// Chain applies multiple middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	// loop in reverse to preserve middleware order
	for i := len(middlewares) - 1; i >= 0; i-- {
		f = middlewares[i](f)
	}
	return f
}

// Log all requests
func Logging() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() { log.Println(r.URL.Path, r.Method, time.Since(start)) }()
			f(w, r)
		}
	}
}

var refillRatePerSecond rate.Limit = 10
var bucketSize = 100
var limiter = rate.NewLimiter(refillRatePerSecond, bucketSize)

// A Limiter controls how frequently events are allowed to happen.
// It implements a "token bucket" of size b, initially full and refilled at rate r tokens per second.
func RateLimit() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			f(w, r)
		}
	}
}

var loginLimiter = rate.NewLimiter(1, 2)

// Limit number of login attemps per second
func LoginRateLimit() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !loginLimiter.Allow() {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			f(w, r)
		}
	}
}

// Restrict access to logged in users only
func LoginProtected(a *app.App) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			cookie, err := r.Cookie("goexpense_access")

			// no access token - make user log in
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			accessToken := cookie.Value

			// if token is invalid/expired - check for existing session and renew access token
			isAccessValid, claims := auth.IsTokenValid(a, accessToken)
			if !isAccessValid {

				// find existing session (refresh_token)
				var s *auth.Session
				err := a.Collections.Sessions.FindOne(
					ctx, bson.D{{Key: "user_id", Value: claims.UserId}}).Decode(&s)

				// session not found - make user log in
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// verify session is still valid
				isRefreshValid, _ := auth.IsTokenValid(a, s.RefreshToken)
				if !isRefreshValid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// renew access token
				renewed, err := auth.CreateAccessToken(ctx, a, claims.UserId, claims.Role)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "goexpense_access",
					Value:    renewed.Value,
					Expires:  renewed.ExpiresAt,
					HttpOnly: true,
				})
			}

			f(w, r)
		}
	}
}

// func AdminOnly(a *app.App) Middleware {
// 	return func(f http.HandlerFunc) http.HandlerFunc {
// 		return func(w http.ResponseWriter, r *http.Request) {
// 			accessCookie, err := r.Cookie("goexpense_access")
// 			if err != nil {
// 				log.Fatalf("Error occured while reading cookie")
// 			}

// 			// verify access token valid
// 			isValid, accessClaims := auth.IsTokenValid(a, accessCookie.Value)
// 			if !isValid {

// 			}

// 			if !loginLimiter.Allow() {
// 				w.WriteHeader(http.StatusTooManyRequests)
// 				return
// 			}
// 			f(w, r)
// 		}
// 	}
// }

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
NoCache sets:
Expires: Thu, 01 Jan 1970 00:00:00 UTC
Cache-Control: no-cache, private, max-age=0
X-Accel-Expires: 0
Pragma: no-cache (for HTTP/1.0 proxies/clients)
*/
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
