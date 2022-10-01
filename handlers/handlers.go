package handlers

import "net/http"

type Middleware interface {
	Wrap(next http.Handler) http.Handler
}

func Wrap(handler http.Handler, mw ...Middleware) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		handler = mw[i].Wrap(handler)
	}

	return handler
}
