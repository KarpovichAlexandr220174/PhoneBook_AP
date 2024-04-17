package main

import (
	"context"
	"net/http"
	"time"
)

func TimeoutHandler(timeout time.Duration, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
