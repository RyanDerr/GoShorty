package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/RyanDerr/GoShorty/pkg/response"
	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	tokens     int
	maxTokens  int
	refillRate int
	dur        time.Duration
	mutex      sync.Mutex
}

func NewRateLimiter(maxTokens, refillRate int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:     maxTokens,
		maxTokens:  maxTokens,
		refillRate: refillRate,
		dur:        duration,
	}
}

// RefillTokens is a method that refills the tokens of the RateLimiter set duration with a mutex lock for thread safety.
func (rl *RateLimiter) RefillTokens() {
	for {
		time.Sleep(rl.dur)
		rl.mutex.Lock()
		if rl.tokens < rl.maxTokens {
			rl.tokens += rl.refillRate
		}

		if rl.tokens > rl.maxTokens {
			rl.tokens = rl.maxTokens
		}
		rl.mutex.Unlock()
	}
}

// acquireToken is a method that acquires a token from the RateLimiter with a mutex lock for thread safety.
func (rl *RateLimiter) acquireToken() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

// IsRateLimited is a middleware function that checks if the RateLimiter is rate limited and returns a 429 status code if it is.
func (rl *RateLimiter) IsRateLimited() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.acquireToken() {
			response.ResponseError(c, "Rate limit exceeded", http.StatusTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}
