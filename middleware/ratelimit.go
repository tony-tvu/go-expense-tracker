package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var globalRate string
var loginRate string

func init() {
	godotenv.Load(".env")
	globalRate = os.Getenv("GLOBAL_RATE_LIMIT")
	loginRate = os.Getenv("LOGIN_RATE_LIMIT")
}

// Middleware applies to all routes
func RateLimit() gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(globalRate)
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	middleware := mgin.NewMiddleware(limiter.New(store, rate))
	return middleware
}

// Middleware applies to login route
func LoginRateLimit() gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(loginRate)
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	middleware := mgin.NewMiddleware(limiter.New(store, rate))
	return middleware
}
