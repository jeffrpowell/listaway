package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

func Cors() Middleware {
	var corsImpl = cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			corsImpl.HandlerFunc(w, r)
			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
