package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// Middleware applies to all routes
func RateLimit() gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted("300-M")
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	middleware := mgin.NewMiddleware(limiter.New(store, rate))
	return middleware
}

// Middleware applies to login route
func LoginRateLimit() gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted("5-M")
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	middleware := mgin.NewMiddleware(limiter.New(store, rate))
	return middleware
}
