package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery middleware handles panics and returns 500 Internal Server Error
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				log.Printf("panic: %v\n%s", err, debug.Stack())

				// Return 500 Internal Server Error
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}