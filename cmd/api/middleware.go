package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")
					app.serverErrorResponse(w, fmt.Errorf("%s", err), nil)
				}
			}()

			next.ServeHTTP(w, r)
		})
}
