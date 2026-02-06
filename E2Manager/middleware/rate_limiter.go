package middleware

import (
	"e2mgr/logger"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	logger   *logger.Logger
	limit    rate.Limit
	burst    int
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(logger *logger.Logger, requestsPerSecond float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		logger:   logger,
		limit:    rate.Limit(requestsPerSecond),
		burst:    burst,
	}
	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.limit, rl.burst)
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(5 * time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		limiter := rl.getVisitor(ip)
		rl.logger.Infof("RateLimiter hit by: %s", ip)
		if !limiter.Allow() {
			rl.logger.Warnf("Rate limit exceeded for IP: %s", ip)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
