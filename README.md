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

Examples of these handlers are in [example/](./example), but here is a truncated version of the page handler example:

```go
// Your main function or something.

pageHandler := gosp.NewPageHandler[testStruct](
    // The actual handler for this page. Normally you would make a query or two here.
    func(_ *http.Request) (*testStruct, error) {
        return &testStruct{"foobar"}, nil
    },
    // Using a file with go:embed is quite useful here.
    template.Must(template.New("test").Parse("<div>{{ .Content }}</div>")),
)
	
// Using "pageHandler" with net/http.
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
