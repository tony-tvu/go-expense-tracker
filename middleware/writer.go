package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WriterAndCookies struct {
	Writer                http.ResponseWriter
	EncryptedRefreshToken string
	EncryptedAccessToken  string
}

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func (c *WriterAndCookies) SetToken(name, value string, expires time.Time) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func GetWriterAndCookies(ctx context.Context) *WriterAndCookies {
	return ctx.Value(userCtxKey).(*WriterAndCookies)
}

func setValInCtx(ctx *gin.Context, val interface{}) {
	newCtx := context.WithValue(ctx.Request.Context(), userCtxKey, val)
	ctx.Request = ctx.Request.WithContext(newCtx)
}

// Middleware allows the request's writer and cookies to be passed down to graphql resolvers
func CookieProvider() gin.HandlerFunc {
	return func(c *gin.Context) {
		wc := WriterAndCookies{
			Writer: c.Writer,
		}
		setValInCtx(c, &wc)

		encrpytedRT := ""
		encrpytedAT := ""

		refreshCookie, err := c.Request.Cookie("goexpense_refresh")
		if err == nil {
			encrpytedRT = refreshCookie.Value
		}
		accessCookie, err := c.Request.Cookie("goexpense_access")
		if err == nil {
			encrpytedAT = accessCookie.Value
		}

		wc.EncryptedRefreshToken = encrpytedRT
		wc.EncryptedAccessToken = encrpytedAT
		c.Next()
	}
}
