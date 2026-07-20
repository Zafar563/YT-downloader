package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	ips sync.Map // maps IP (string) -> *client
	r   rate.Limit
	b   int
}

// NewIPRateLimiter initializes an IP-based rate limiter
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		r: r,
		b: b,
	}

	// Periodically run cleanup to prevent memory growth from idle IPs
	go limiter.cleanup()

	return limiter
}

// GetLimiter retrieves or initializes a token bucket limiter for a client IP
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	val, exists := i.ips.Load(ip)
	if exists {
		c := val.(*client)
		c.lastSeen = time.Now()
		return c.limiter
	}

	limiter := rate.NewLimiter(i.r, i.b)
	c := &client{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	i.ips.Store(ip, c)
	return limiter
}

func (i *IPRateLimiter) cleanup() {
	for {
		time.Sleep(3 * time.Minute)
		i.ips.Range(func(key, value interface{}) bool {
			c := value.(*client)
			// Purge client limiter if inactive for over 15 minutes
			if time.Since(c.lastSeen) > 15*time.Minute {
				i.ips.Delete(key)
			}
			return true
		})
	}
}

// RateLimitMiddleware blocks IPs exceeding limits with HTTP 429
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		lim := limiter.GetLimiter(ip)
		if !lim.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please wait a moment before trying again.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
