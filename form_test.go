package gosp

import (
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

	type testStruct struct {
		Foo string `schema:"foo"`
	}

	NewFormHandler[testStruct](func(s *testStruct) error {
		assert.Equal(t, "bar", s.Foo)
		return nil
	}).ServeHTTP(rec, req)

	assert.Equal(t, 301, rec.Code)
}
