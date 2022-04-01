package gosp

import (
	"html/template"
	"net/http"
)

// Handler contains the main logic of a PageHandler.
//
// Three possible states occur in the return of a Handler:
//
//  The handler...
//    1. returns a pointer to a struct. This is our good case.
//    2. does not return a pointer, but nil instead. This is our empty case.
//    3. returns an error. This is our error case.
type Handler[T any] func() (*T, error)

// PageHandler represent the content to be served as the result of some logic being performed in the Handler.
// The template can be seen a structure for this data.
type PageHandler[T any] struct {
	// Handler determines what data is to be filled in Template.
	Handler Handler[T]

	// Template is the HTML template for our page.
	// The fields of this Template are filled by the result of Handler.
	Template *template.Template

	// ErrorHandler is used whenever Handler return an error.
	// The default value is DefaultErrorHandler.
	ErrorHandler func(error) http.Handler

	// EmptyHandler is used whenever Handler returns a nil.
	// The default value is DefaultEmptyHandler.
	EmptyHandler http.Handler
}

// NewPageHandler requires only the "Handler" and Template fields be provided.
// Everything else has a default.
func NewPageHandler[T any](handler Handler[T], template *template.Template, options ...PageHandlerOption[T]) *PageHandler[T] {
	p := &PageHandler[T]{handler, template, DefaultErrorHandler, http.HandlerFunc(DefaultEmptyHandler)}
	for _, o := range options {
		o(p)
	}
	return p
}

func (p *PageHandler[any]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := p.Handler()
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
