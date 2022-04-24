package main

import (
	"fmt"
	"github.com/eyebrow-fish/gosp"
	"net/http"
)

type testStruct struct{ Content string }

func main() {
	formHandler := gosp.NewFormHandler[testStruct](
		// Form handling is done in this function. The form has already been deserialized.
		func(t *testStruct) error {
			fmt.Println(t.Content)
			return nil
		},
	)

	http.Handle("/", formHandler)

	// Remember to handle your errors, kids!
	_ = http.ListenAndServe(":8080", nil)
}
