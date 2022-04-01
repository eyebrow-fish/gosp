package gosp

import "net/http"

type PageHandlerOption[T any] func(p *PageHandler[T])

// DefaultErrorHandler writes the content of the error as the response and responds with a 500.
func DefaultErrorHandler(err error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Error is ignored because we won't be able to write anything.
		// Just do the 500 instead.
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	})
}

// DefaultEmptyHandler simply responds with a 404.
func DefaultEmptyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
