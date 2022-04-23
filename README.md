# gosp

[*Go*](https://go.dev/) *S*tatic *P*ages is a small library which wraps around [net/http](https://pkg.go.dev/net/http)
to help build static webpages. Using this library helps remove some necessary boilerplate that gets in the way of being
productive.

# goal

I don't think it's a big secret that many websites could probably be reduced down to a simple server-side rendered
service and actually gain big points in simplicity, maintainability, and speed. Which would mean that small projects
like this one could possibly help reduce on monstrously overcomplicated frontends move towards a simpler solution.

# use

There are two main "Handlers" that exist in this module which satisfy
[http.Handler](https://pkg.go.dev/net/http#Handler), these are `PageHandler` and `FormHandler`.

An example usage of `PageHandler` is as follows *(stolen from my tests)*:

```go
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
		func() (*testStruct, error) {
			return &testStruct{"foobar"}, nil
		},
		// Using a file with go:embed is quite useful here.
		template.Must(template.New("test").Parse("<div>{{ .Content }}</div>")),
	)

	http.Handle("/", pageHandler)

	http.ListenAndServe(":8080", nil) // Remember to handle your errors, kids!
}
```

Using `FormHandler` is slightly more simple because only one handler is required, and you don't have to
think about templates.

# some free "art"

```
     ()------()     ____________________
    | (.)  (.) |   |Thanks for reading! |
(( <|    uu    |>  <____________________|
    |          |      
    |.        .|      
     o--------o
```

Feel free to steal the above ASCII art.
