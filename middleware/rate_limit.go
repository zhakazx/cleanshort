package middleware

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/zhakazx/cleanshort/models"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Allow(key string) (bool, int, time.Time) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	requests := rl.requests[key]

	// Filter out requests outside the window
	var validRequests []time.Time
	for _, req := range requests {
		if req.After(windowStart) {
			validRequests = append(validRequests, req)
		}
	}

	// Check if limit is exceeded
	if len(validRequests) >= rl.limit {
		// Calculate reset time (when the oldest request expires)
		resetTime := validRequests[0].Add(rl.window)
		return false, rl.limit - len(validRequests), resetTime
	}

	validRequests = append(validRequests, now)
	rl.requests[key] = validRequests

	return true, rl.limit - len(validRequests), now.Add(rl.window)
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)

		for key, requests := range rl.requests {
			var validRequests []time.Time
			for _, req := range requests {
				if req.After(windowStart) {
					validRequests = append(validRequests, req)
				}
			}

			if len(validRequests) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validRequests
			}
		}
		rl.mutex.Unlock()
	}
}

func RateLimitMiddleware(limit int, window time.Duration) fiber.Handler {
	limiter := NewRateLimiter(limit, window)

	return func(c *fiber.Ctx) error {
		// Use IP address as the key
		key := c.IP()

		allowed, remaining, resetTime := limiter.Allow(key)

		c.Set("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if !allowed {
			return c.Status(fiber.StatusTooManyRequests).JSON(models.ErrorResponse{
				Error: models.ErrorDetail{
					Code:      "TOO_MANY_REQUESTS",
					Message:   fmt.Sprintf("Rate limit exceeded. Try again after %v", time.Until(resetTime).Round(time.Second)),
					RequestID: c.Locals("requestid").(string),
				},
			})
		}

		return c.Next()
	}
}

func AuthRateLimitMiddleware(limit int) fiber.Handler {
	return RateLimitMiddleware(limit, time.Minute)
}

func RedirectRateLimitMiddleware(limit int) fiber.Handler {
	return RateLimitMiddleware(limit, time.Minute)
}
