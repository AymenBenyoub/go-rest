package middleware 

import (
	"net/http")

// Chain is a type that represents a chain of middleware handlers. 

type Middleware func (http.Handler) http.Handler 


func Chain(middlewares ...Middleware) Middleware {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}