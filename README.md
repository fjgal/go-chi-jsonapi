# go-chi-jsonapi

`go-chi-jsonapi` is a collection of utilities to facilitate the implementation of [JSON:API](https://jsonapi.org/) HTTP APIs with `go-chi`

## `render` package

The `render` package provides a JSON API responder and decoder that can be easily integrated into [go-chi/render](https://github.com/go-chi/render). It provides `DefaultResponder` and `DefaultDecoder` functions to be set in `go-chi/render` 

To make integration simpler it re-implements some of the `go-chi/render` functions and delegates to the orginal one, the implementation is heavily inspired by [go-chi/render](https://github.com/go-chi/render).

This package uses [google/jsonapi](https://github.com/google/jsonapi). Structs must have proper `jsonapi` tags, see [jsonapi Tag Reference](https://github.com/google/jsonapi#jsonapi-tag-reference). Payload must be given as a struct pointer or a slice of struct pointers, this is how `google/jsonapi` works.

Supported features:

* `DefaultResponder` and `DefaultDecoder` functions for easy integration with `go-chi/render`
* Detects Content-Type from request Header and decodes accordingly (can be overriden using `render.SetConentType` middleware)
* Detects Accept type from request Header and encodes accordingly (can be overriden using `render.SetConentType` middleware)
* Automatically encodes errors as JSON API Error Objects (when Accept is set to JSON API)
* Can be used by calling `render.Respond`/`render.Decode` or `render.Render`/`render.Bind`

## Examples

### Renderer

Using `render.Respond` and `render.Decode`

```
    import (
        jsonapi_render "github.com/fjgal/go-chi-jsonapi/render"
        "github.com/go-chi/render"
    )

	// plug jsonapi responder and decoder
	render.Respond = jsonapi_render.DefaultResponder
	render.Decode = jsonapi_render.DefaultDecoder

    // implement a handler function
    router.Get("/", func(w http.ResponseWriter, r *http.Request) {
        var blog Blog
        // decode the request body into a struct
        _ := render.Decode(r, &blog)
        // set response status
        render.Status(r, http.StatusAccepted)
        // render the response
        render.Respond(w, r, &blog)
    })
```

Responding with an error

```
    import (
        jsonapi_render "github.com/fjgal/go-chi-jsonapi/render"
        "github.com/go-chi/render"
    )

	// plug jsonapi responder and decoder
	render.Respond = jsonapi_render.DefaultResponder
	render.Decode = jsonapi_render.DefaultDecoder

    // implement a handler function
    router.Get("/", func(w http.ResponseWriter, r *http.Request) {
        render.Status(r, http.StatusUnprocessableEntity)
        render.Respond(w, r, err)
    })
```

Using `render.Render` and `render.Bind` (model structs must implement `render.Renderer` and  `renderBinder` interfaces)

```
    import (
        jsonapi_render "github.com/fjgal/go-chi-jsonapi/render"
        "github.com/go-chi/render"
    )

    router.Get("/", func(w http.ResponseWriter, r *http.Request) {
        var blog Blog
        _ := render.Bind(r, &blog)
        render.Status(r, http.StatusAccepted)
        err = render.Render(w, r, &blog)
        if err != nil {
            t.Log(err)
            return
        }
    })
```


## TODO

- [] implement `queryparams` package
- [] support json api streaming ?