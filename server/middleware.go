package server

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/time/rate"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func AdminUserMiddleware(f http.HandlerFunc, a *app.App) http.HandlerFunc {
	f = NoCache()(f)
	f = Admin(a)(f)
	f = LoggedIn(a)(f)
	f = RateLimit()(f)
	f = Logging(a)(f)
	return f
}

func RegularUserMiddleware(f http.HandlerFunc, a *app.App) http.HandlerFunc {
	f = NoCache()(f)
	f = LoggedIn(a)(f)
	f = RateLimit()(f)
	f = Logging(a)(f)
	return f
}

func GuestUserMiddleware(f http.HandlerFunc, a *app.App) http.HandlerFunc {
	f = NoCache()(f)
	f = RateLimit()(f)
	f = Logging(a)(f)
	return f
}

// Log all requests
func Logging(a *app.App) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if a.Env != "test" {
				start := time.Now()
				defer func() { log.Println(r.URL.Path, r.Method, time.Since(start)) }()
			}
			f(w, r)
		}
	}
}

// Method ensures that url can only be requested with a specific method
func Method(m string) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
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

var loginLimiter = rate.NewLimiter(1, 10)

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

// Restricts access to logged in users only
func LoggedIn(a *app.App) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			cookie, err := r.Cookie("goexpense_access")

			// no access token - make user log in
			if err != nil || cookie.Value == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// decrypt access token
			decrypted, err := auth.Decrypt(a.EncryptionKey, cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tkn, err := jwt.ParseWithClaims(decrypted, &auth.Claims{},
				func(token *jwt.Token) (interface{}, error) {
					return []byte(a.JwtKey), nil
				})
			claims := tkn.Claims.(*auth.Claims)

			// token is expired and missing correct claims - make user log in
			if err != nil && claims.UserID == "" && claims.UserType == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// if token is invalid/expired - check for existing session and renew access token
			if !tkn.Valid {

				// find existing session (refresh_token)
				var s *models.Session
				err := a.Sessions.FindOne(
					ctx, bson.D{{Key: "user_id", Value: claims.UserID}}).Decode(&s)

				// session not found - make user log in
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// decrypt refresh token
				decrypted, err = auth.Decrypt(a.EncryptionKey, s.RefreshToken)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// verify session is still valid
				_, err = jwt.ParseWithClaims(decrypted, &auth.Claims{},
					func(token *jwt.Token) (interface{}, error) {
						return []byte(a.JwtKey), nil
					})

				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				// renew access token
				renewed, err := auth.CreateAccessToken(ctx, a, claims.UserID, claims.UserType)
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

// Restricts access to admins only
func Admin(a *app.App) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("goexpense_access")

			// access token missing
			if err != nil || cookie.Value == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// decrypt token
			decrypted, err := auth.Decrypt(a.EncryptionKey, cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tkn, err := jwt.ParseWithClaims(decrypted, &auth.Claims{},
				func(token *jwt.Token) (interface{}, error) {
					return []byte(a.JwtKey), nil
				})
			claims := tkn.Claims.(*auth.Claims)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if claims.UserType != string(models.AdminUser) {
				w.WriteHeader(http.StatusUnauthorized)
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
