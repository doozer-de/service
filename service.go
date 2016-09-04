package service

import (
	"context"
	"net/http"
	"path"
	"strings"
)

type Service struct {
	preHandler contextHandler

	routes map[string]*node

	postHandler contextHandler

	notFoundHandler    contextHandler
	serverErrorHandler contextHandler
	unauthorizedHander contextHandler

	trimSlash bool
	baseURI   string
}

func New(baseURI string) *Service {
	return &Service{
		baseURI:   path.Join("/", baseURI, "/"),
		routes:    map[string]*node{},
		trimSlash: true,
	}
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.ServeWithContext(context.Background(), w, r)
}

func (s *Service) ServeWithContext(c context.Context, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if s.preHandler != nil {
		s.preHandler(c, w, r)
	}

	if r.URL.Path != "/" && s.trimSlash {
		r.URL.Path = strings.TrimRight(r.URL.Path, "/")
	}

	var h contextHandler
	var ps KVs

	n, ok := s.routes[r.Method]

	if ok {
		h, ps, _ = n.getValue(r.URL.Path)
		// First place where we could attach information to the context
	}

	if h == nil {
		if s.notFoundHandler != nil {
			s.notFoundHandler(c, w, r)
		} else {
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	} else {
		for _, p := range ps {
			r.Form.Set(p.Key, p.Value)
		}

		h(c, w, r)
	}

	if s.postHandler != nil {
		s.postHandler(c, w, r)
	}
}

// TODO(cs): Add a Swagger information here and register it with the route
func (s *Service) Route(method, uri string, handler interface{}) {
	//TODO(cs): Check if the http is valid (one of the constants above)
	h := toContextHandler(handler)

	if n := s.routes[method]; n == nil {
		s.routes[method] = &node{}
	}

	s.routes[method].addRoute(path.Join(s.baseURI, strings.TrimRight(uri, "/")), h)
}

func (s *Service) SetNotFound(f interface{}) {
	if f == nil {
		s.notFoundHandler = nil
		return
	}

	h := toContextHandler(f)
	s.notFoundHandler = h
}
