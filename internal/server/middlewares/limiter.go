package middlewares

import (
	"net/http"
	"time"

	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	limiter := ratelimit.NewBucket(1*time.Second, 1) // 1 request per second
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter.Take(1)
		next.ServeHTTP(w, r)
	})
}
