package middleware

import (
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

type ipRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func newIPRateLimiter(rps int, burst int) *ipRateLimiter {
	limit := rate.Limit(rps)
	if limit <= 0 {
		limit = rate.Inf
	}
	return &ipRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     limit,
		burst:    burst,
	}
}

func (l *ipRateLimiter) get(ip string) *rate.Limiter {
	l.mu.RLock()
	limiter, ok := l.limiters[ip]
	l.mu.RUnlock()
	if ok {
		return limiter
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	limiter, ok = l.limiters[ip]
	if ok {
		return limiter
	}
	limiter = rate.NewLimiter(l.rate, l.burst)
	l.limiters[ip] = limiter
	return limiter
}

func RateLimit(rps int, burst int) func(http.Handler) http.Handler {
	limiter := newIPRateLimiter(rps, burst)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				parts := strings.Split(forwarded, ",")
				if len(parts) > 0 {
					ip = strings.TrimSpace(parts[0])
				}
			}
			if !limiter.get(ip).Allow() {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

