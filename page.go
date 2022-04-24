package gosp

import (
	"html/template"
	"net/http"
)

// PageHandler represent the content to be served as the result of some logic being performed in the MainPageHandler.
// The template can be seen a structure for this data.
type PageHandler[T any] struct {
	// Handler determines what data is to be filled in Template.
	Handler func(*http.Request) (*T, error)

	// Template is the HTML template for our page.
	// The fields of this Template are filled by the result of Handler.
	Template *template.Template

	// ErrorHandler is used whenever Handler returns an error.
	// The default value is DefaultErrorHandler.
	ErrorHandler func(error) http.Handler

	// EmptyHandler is used whenever Handler returns a nil.
	// The default value is DefaultEmptyHandler.
	EmptyHandler http.Handler
}

// NewPageHandler requires only PageHandler.Handler and PageHandler.Template fields be provided.
// Everything else has a default.
func NewPageHandler[T any](handler func(*http.Request) (*T, error), template *template.Template, options ...PageHandlerOption[T]) *PageHandler[T] {
	p := &PageHandler[T]{handler, template, DefaultErrorHandler, http.HandlerFunc(DefaultEmptyHandler)}
	for _, o := range options {
		o(p)
	}
	return p
}

// ServeHTTP delegates all logic to the handler methods of PageHandler.
func (p *PageHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := p.Handler(r)
	if err != nil {
		p.ErrorHandler(err).ServeHTTP(w, r)
		return
	}

	if t == nil {
		p.EmptyHandler.ServeHTTP(w, r)
		return
	}

	if err = p.Template.Execute(w, t); err != nil {
		p.ErrorHandler(err).ServeHTTP(w, r)
	}
}
