package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const defaultTimeout = 2 * time.Second

func TimeoutMiddleware(next http.Handler) http.Handler {
	timeout := getTimeout()
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()

		done := make(chan bool, 1)
		wrapped := &responseWriter{ResponseWriter: w, written: false}
		
		go func() {
			next.ServeHTTP(wrapped, r.WithContext(ctx))
			done <- true
		}()

		select {
		case <-ctx.Done():
			if !wrapped.written {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusGatewayTimeout)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "request timeout"})
			}
			return
		case <-done:
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	written bool
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.written = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.written = true
	return rw.ResponseWriter.Write(b)
}

func getTimeout() time.Duration {
	if timeoutStr := os.Getenv("REQUEST_TIMEOUT"); timeoutStr != "" {
		if timeout, err := strconv.Atoi(timeoutStr); err == nil && timeout > 0 {
			return time.Duration(timeout) * time.Millisecond
		}
	}
	return defaultTimeout
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("%s %s %v", r.Method, r.URL.Path, duration)
	})
}

