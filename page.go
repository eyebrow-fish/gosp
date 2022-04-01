package gosp

import (
	"html/template"
	"net/http"
)

type Handler[T any] func() (*T, error)

// PageHandler represent the content to be served as the result of some logic being performed in the Handler.
// The template can be seen a structure for this data.
type PageHandler[T any] struct {
	Handler      Handler[T]
	Template     *template.Template
	ErrorHandler func(error) http.Handler
}

// NewPageHandler requires only the "Handler" and Template fields be provided.
// Everything else has a default.
func NewPageHandler[T any](handler Handler[T], template *template.Template, options ...PageHandlerOption[T]) *PageHandler[T] {
	p := &PageHandler[T]{handler, template, DefaultErrorHandler}
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

	if err = p.Template.Execute(w, t); err != nil {
		p.ErrorHandler(err).ServeHTTP(w, r)
	}
}
