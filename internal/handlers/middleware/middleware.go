package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

func DefaultMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return Chain(f, RequireAuth(), Cors())
}

func DefaultPublicMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return Chain(f, Cors())
}
