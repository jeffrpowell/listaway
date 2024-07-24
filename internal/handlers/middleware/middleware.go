package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

var DefaultPublicMiddlewareSlice []Middleware = []Middleware{Cors()}
var DefaultMiddlewareSlice []Middleware = []Middleware{RequireAuth(), Cors()}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func DefaultMiddlewareChain(f http.HandlerFunc) http.HandlerFunc {
	return Chain(f, DefaultMiddlewareSlice...)
}

func DefaultPublicMiddlewareChain(f http.HandlerFunc) http.HandlerFunc {
	return Chain(f, DefaultPublicMiddlewareSlice...)
}
