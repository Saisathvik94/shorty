package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
)

func RedisRateLimiter(limiter *redis_rate.Limiter, rateLimit int, category string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		ipAddress := c.ClientIP()
		redisKey := fmt.Sprintf("ratelimit:%s:%s", category, ipAddress)

		limitRule := redis_rate.PerMinute(rateLimit)

		res, err := limiter.Allow(ctx, redisKey, limitRule)
		if err != nil {
			// Fail-open strategy: log error and let the request proceed
			fmt.Printf("Redis rate limiter error: %v\n", err)
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", res.Limit.Burst))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", res.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", int(res.ResetAfter.Seconds())))

		if res.Allowed == 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too Many Requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()

	}
}
