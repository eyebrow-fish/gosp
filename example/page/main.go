package main

import (
	"github.com/eyebrow-fish/gosp"
	"html/template"
	"net/http"
)

type testStruct struct{ Content string }

func main() {
	pageHandler := gosp.NewPageHandler[testStruct](
		// The actual handler for this page. Normally you would make a query or two here.
		func(_ *http.Request) (*testStruct, error) {
			return &testStruct{"foobar"}, nil
		},
		// Using a file with go:embed is quite useful here.
		template.Must(template.New("test").Parse("<div>{{ .Content }}</div>")),
	)

	http.Handle("/", pageHandler)

	// Remember to handle your errors, kids!
	_ = http.ListenAndServe(":8080", nil)
}
