package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func (app *app) recoverPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {

			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}

		}()
		next.ServeHTTP(w, r)
	})

}

func (app *app) rateLimiter(next http.Handler) http.Handler {
	rateLimit := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rateLimit.Allow() {
			app.tooManyRequest(w, r)
			return
		}
		next.ServeHTTP(w, r)

	})
}
