package gosp

import (
	"github.com/gorilla/schema"
	"net/http"
)

var decoder = schema.NewDecoder()

type FormHandler[T any] struct {
	Handler func(*T) error
}

func NewFormHandler[T any](handler func(*T) error) *FormHandler[T] {
	return &FormHandler[T]{handler}
}

func (f FormHandler[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	form := new(T)

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := decoder.Decode(form, r.PostForm); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := f.Handler(form); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusMovedPermanently)
}
