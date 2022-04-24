package gosp

import (
	"github.com/gorilla/schema"
	"net/http"
)

// FormDecoder is the schema.Decoder used for every FormHandler.
var FormDecoder = schema.NewDecoder() // Exposed to allow users to modify behaviour.

// FormHandler is a http.Handler compatible struct which handles http.MethodPost requests with form bodies.
type FormHandler[T any] struct {
	// Handler is the main handler of a FormHandler.
	// The parameter given is the deserialized form sent.
	Handler func(*T) error

	// ErrorHandler is used whenever Handler returns an error.
	// The default value is DefaultErrorHandler.
	ErrorHandler func(error) http.Handler

	// RedirectHandler is invoked once Handler succeeds, and is meant to provide functionality for redirecting
	// after a successful POST request.
	// The default value is DefaultRedirectHandler.
	RedirectHandler http.Handler
}

// NewFormHandler initializes a new FormHandler, and only the FormHandler.Handler is required.
// All other fields have default values.
func NewFormHandler[T any](handler func(*T) error, options ...FormHandlerOption[T]) *FormHandler[T] {
	h := &FormHandler[T]{handler, DefaultErrorHandler, http.HandlerFunc(DefaultRedirectHandler)}
	for _, o := range options {
		o(h)
	}
	return h
}

// ServeHTTP delegates all logic to the handler methods of FormHandler.
func (f FormHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	form := new(T)

	if err := r.ParseForm(); err != nil {
		f.ErrorHandler(err).ServeHTTP(w, r)
		return
	}

	if err := FormDecoder.Decode(form, r.PostForm); err != nil {
		f.ErrorHandler(err).ServeHTTP(w, r)
		return
	}

	if err := f.Handler(form); err != nil {
		f.ErrorHandler(err).ServeHTTP(w, r)
		return
	}

	f.RedirectHandler.ServeHTTP(w, r)
}
