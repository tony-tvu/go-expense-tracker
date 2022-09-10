package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/util"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// Middleware applies rate limiting to all routes
func RateLimit() gin.HandlerFunc {
	godotenv.Load(".env")
	rateLimit := os.Getenv("RATE_LIMIT")
	if util.ContainsEmpty(rateLimit) {
		rateLimit = "30-S"
	}

	rate, err := limiter.NewRateFromFormatted(rateLimit)
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	middleware := mgin.NewMiddleware(limiter.New(store, rate))
	return middleware
}
