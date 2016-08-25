package service

import (
	"net/http"

	"golang.org/x/net/context"
)

const (
	GET     = "GET"
	PUT     = "PUT"
	OPTIONS = "OPTIONS"
	HEAD    = "HEAD"
	POST    = "POST"
	DELETE  = "DELETE"
)

type contextHandler func(context.Context, http.ResponseWriter, *http.Request)

func (h contextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(nil, w, r)
}

func (h contextHandler) ServeWithContext(c context.Context, w http.ResponseWriter, r *http.Request) {
	h(c, w, r)
}

func toContextHandler(f interface{}) contextHandler {
	var h contextHandler

	switch f.(type) {
	case func(context.Context, http.ResponseWriter, *http.Request):
		h = contextHandler(f.(func(context.Context, http.ResponseWriter, *http.Request)))
	case contextHandler:
		h = f.(contextHandler)
	case func(http.ResponseWriter, *http.Request):
		h = func(c context.Context, w http.ResponseWriter, r *http.Request) {
			f.(func(http.ResponseWriter, *http.Request))(w, r)
		}
	default:
		if h, ok := f.(http.Handler); ok {
			return toContextHandler(h.ServeHTTP)
		}
		panic("Unsupported Handler")
	}

	return h
}

func Chain(stack ...interface{}) contextHandler {
	s := make([]contextHandler, 0, len(stack))

	for i := range stack {
		m := toContextHandler(s[i])
		s = append(s, m)
	}

	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		for _, m := range s {
			m(c, w, r)
		}
	}
}
