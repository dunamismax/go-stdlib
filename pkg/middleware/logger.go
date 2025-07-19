package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"duration", time.Since(start),
			"user_agent", r.UserAgent(),
		)
	})
}

func CORS(next http.Handler) http.Handler {
	return CORSWithOrigins(nil)(next)
}

func CORSWithOrigins(allowedOrigins []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if len(allowedOrigins) == 0 {
				// Default allowed origins for development
				defaultOrigins := []string{
					"http://localhost:3000",
					"http://localhost:3001",
					"http://127.0.0.1:3000",
					"http://127.0.0.1:3001",
				}
				allowedOrigins = defaultOrigins
			}

			// Check if origin is allowed
			originAllowed := false
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					originAllowed = true
					break
				}
			}

			if originAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, "Request timeout")
	}
}

type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Get existing requests for this IP
	requests := rl.requests[clientIP]

	// Remove old requests outside the window
	var validRequests []time.Time
	for _, req := range requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}

	// Check if we're under the limit
	if len(validRequests) >= rl.limit {
		rl.requests[clientIP] = validRequests
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[clientIP] = validRequests

	return true
}

func (rl *RateLimiter) cleanup() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	for ip, requests := range rl.requests {
		var validRequests []time.Time
		for _, req := range requests {
			if req.After(cutoff) {
				validRequests = append(validRequests, req)
			}
		}

		if len(validRequests) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = validRequests
		}
	}
}

func RateLimit(limit int, window time.Duration) func(next http.Handler) http.Handler {
	limiter := NewRateLimiter(limit, window)

	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for range ticker.C {
			limiter.cleanup()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := getClientIP(r)

			if !limiter.Allow(clientIP) {
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for load balancers/proxies)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, use the first one
		if firstIP := net.ParseIP(xff); firstIP != nil {
			return firstIP.String()
		}
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		if ip := net.ParseIP(xri); ip != nil {
			return ip.String()
		}
	}

	// Fall back to remote address
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

// LoginRateLimit provides stricter rate limiting for login endpoints
func LoginRateLimit() func(next http.Handler) http.Handler {
	return RateLimit(5, 5*time.Minute) // 5 attempts per 5 minutes
}

// APIRateLimit provides general rate limiting for API endpoints
func APIRateLimit() func(next http.Handler) http.Handler {
	return RateLimit(100, time.Minute) // 100 requests per minute
}

// GeneralRateLimit provides basic rate limiting for general endpoints
func GeneralRateLimit(next http.Handler) http.Handler {
	return RateLimit(60, time.Minute)(next) // 60 requests per minute
}
