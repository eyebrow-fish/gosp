package gosp

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestForm_happyPath(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("foo=bar"))
	req.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	req.Header["Referer"] = []string{"/prev"}

	type testStruct struct {
		Foo string `schema:"foo"`
	}

	NewFormHandler[testStruct](func(_ *http.Request, s *testStruct) error {
		assert.Equal(t, "bar", s.Foo)
		return nil
	}).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Equal(t, []string{"/prev"}, rec.Header()["Location"])
}

func TestForm_customRedirectHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("foo=bar"))
	req.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	req.Header["Referer"] = []string{"/prev"}

	type testStruct struct {
		Foo string `schema:"foo"`
	}

	NewFormHandler[testStruct](
		func(_ *http.Request, s *testStruct) error {
			assert.Equal(t, "bar", s.Foo)
			return nil
		},
		func(h *FormHandler[testStruct]) {
			h.RedirectHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/foo", http.StatusAccepted)
			})
		},
	).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.Equal(t, []string{"/foo"}, rec.Header()["Location"])
}

func TestForm_errorHandle(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	type testStruct struct {
		Foo string `schema:"foo"`
	}

	NewFormHandler[testStruct](
		func(_ *http.Request, s *testStruct) error { return errors.New("oof") },
	).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "oof", rec.Body.String())
}
