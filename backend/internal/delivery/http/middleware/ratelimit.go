package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/org/itsm/pkg/apperror"
)

// RateLimiter implements token bucket algorithm for rate limiting
type RateLimiter struct {
	mu              sync.RWMutex
	buckets         map[string]*tokenBucket
	capacity        int
	refillRate      time.Duration
	cleanupInterval time.Duration
}

type tokenBucket struct {
	tokens     float64
	lastRefill time.Time
}

// NewRateLimiter creates a new rate limiter with specified capacity and refill rate
func NewRateLimiter(capacity int, refillRate time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets:         make(map[string]*tokenBucket),
		capacity:        capacity,
		refillRate:      refillRate,
		cleanupInterval: 5 * time.Minute,
	}
	go rl.cleanup()
	return rl
}

// Allow checks if a request from the given key is allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, exists := rl.buckets[key]
	now := time.Now()

	if !exists {
		rl.buckets[key] = &tokenBucket{
			tokens:     float64(rl.capacity) - 1,
			lastRefill: now,
		}
		return true
	}

	elapsed := now.Sub(bucket.lastRefill)
	tokensToAdd := float64(rl.capacity) * (elapsed.Seconds() / rl.refillRate.Seconds())
	bucket.tokens = minFloat(float64(rl.capacity), bucket.tokens+tokensToAdd)
	bucket.lastRefill = now

	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

// cleanup removes old buckets periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, bucket := range rl.buckets {
			if now.Sub(bucket.lastRefill) > 10*time.Minute {
				delete(rl.buckets, key)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit middleware returns a middleware that rate limits requests per authenticated user
func RateLimit(capacity int) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(capacity, time.Minute)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r)
			if !ok {
				apperror.WriteError(w, apperror.ErrTokenInvalid)
				return
			}

			key := claims.UserID.String()

			if !limiter.Allow(key) {
				w.Header().Set("Retry-After", "60")
				apperror.WriteError(w, apperror.ErrRateLimitExceeded)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
