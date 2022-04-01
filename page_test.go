package gosp

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testStruct struct{ Content string }

func TestPage_happyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	NewPageHandler[testStruct](
		func() (*testStruct, error) {
			return &testStruct{"foobar"}, nil
		},
		template.Must(template.New("test").Parse("<div>{{ .Content }}</div>")),
	).ServeHTTP(rec, req)

	res := rec.Result()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "<div>foobar</div>", string(body))
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestPage_errorInHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	NewPageHandler[testStruct](
		func() (*testStruct, error) {
			return nil, errors.New("oops")
		},
		template.Must(template.New("test").Parse("<div>{{ .Content }}</div>")),
		func(handler *PageHandler[testStruct]) {
			handler.ErrorHandler = func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, err := w.Write([]byte("handler error"))
					assert.Nil(t, err)
				})
			}
		},
	).ServeHTTP(rec, req)

	res := rec.Result()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "handler error", string(body))
}

func TestPage_errorInHandler_defaultError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	NewPageHandler[testStruct](
		func() (*testStruct, error) {
			return nil, errors.New("handler error")
		},
		template.Must(template.New("test").Parse("<div>{{ .Content }}</div>")),
	).ServeHTTP(rec, req)

	res := rec.Result()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "handler error", string(body))
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestPage_errorInTemplateExecution(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	NewPageHandler[testStruct](
		func() (*testStruct, error) {
			return &testStruct{"foobar"}, nil
		},
		template.Must(template.New("test").Parse("<d")),
		func(handler *PageHandler[testStruct]) {
			handler.ErrorHandler = func(err error) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					_, err := w.Write([]byte("template error"))
					assert.Nil(t, err)
				})
			}
		},
	).ServeHTTP(rec, req)

	res := rec.Result()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "template error", string(body))
}

func TestPage_defaultEmpty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	NewPageHandler[testStruct](
		func() (*testStruct, error) {
			return nil, nil
		},
		template.Must(template.New("test").Parse("<d")),
	).ServeHTTP(rec, req)

	res := rec.Result()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "", string(body))
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}
