package router

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type Router struct {
	mux         *http.ServeMux
	middlewares []Middleware
}

func NewRouter() *Router {
	return &Router{
		mux:         http.NewServeMux(),
		middlewares: []Middleware{},
	}
}

func (r *Router) Use(mw ...Middleware) {
	r.middlewares = append(r.middlewares, mw...)
}

func (r *Router) Handle(pattern string, handler http.HandlerFunc, routeMiddlewares ...Middleware) {
	var finalHandler http.Handler = handler

	for i := len(routeMiddlewares) - 1; i >= 0; i-- {
		finalHandler = routeMiddlewares[i](finalHandler)
	}

	for i := len(r.middlewares) - 1; i >= 0; i-- {
		finalHandler = r.middlewares[i](finalHandler)
	}

	r.mux.Handle(pattern, finalHandler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}