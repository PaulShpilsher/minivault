package server

import (
	"fmt"
	"net/http"
	"minivault/domain"
)

// BodyLimitMiddleware limits incoming request body size to 4KB (4096 bytes).
func BodyLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 4096)
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware returns a middleware that recovers from panics and logs them using the provided logger.
func RecoveryMiddleware(logger domain.LoggerPort, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				reqID := r.Header.Get("X-Request-ID")
				logger.LogError(fmt.Sprintf("panic recovered [reqID: %s]", reqID), fmt.Errorf("%v", rec))
				http.Error(w, "Internal Server Error [panic]", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
