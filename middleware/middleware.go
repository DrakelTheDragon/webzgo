package middleware

import "net/http"

type Interface interface {
	Wrap(next http.Handler) http.Handler
}

type Func func(next http.Handler) http.HandlerFunc

func (f Func) Wrap(next http.Handler) http.Handler {
	return f(next)
}
