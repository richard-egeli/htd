package router

import (
	"net/http"
	"strings"
)

type Router struct {
	root        string
	parent      *Router
	mux         *http.ServeMux
	middlewares []Middleware
}

func New() *Router {
	return &Router{
		root:        "",
		mux:         http.NewServeMux(),
		middlewares: nil,
	}
}

func (router *Router) getMiddlewares() []Middleware {
	if router.parent != nil {
		return append(router.parent.getMiddlewares(), router.middlewares...)
	} else {
		return router.middlewares
	}
}

func (router *Router) set(path string, method Method, middlewares []Middleware, handler http.Handler) {
	path = router.root + path
	length := len(path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if len(method) > 0 {
		path = string(method) + " " + path
	}

	for _, mw := range append(router.getMiddlewares(), middlewares...) {
		handler = mw(handler)
	}

	if strings.HasSuffix(path, "/") && length > 1 {
		absPath := strings.TrimSuffix(path, "/")
		router.mux.Handle(absPath, handler)
	}

	router.mux.Handle(path, handler)
}

func (router *Router) route() string {
	if router.parent != nil {
		return router.parent.route() + router.root
	}

	return router.root
}

func (router *Router) Use(middleware Middleware) {
	router.middlewares = append(router.middlewares, middleware)
}

func (router *Router) Sub(subroute string) *Router {
	return &Router{
		root:        router.route() + subroute,
		mux:         router.mux,
		parent:      router,
		middlewares: nil,
	}
}

func (router *Router) Get(path string, middlewares []Middleware, handler http.Handler) {
	router.set(path, GET, middlewares, handler)
}

func (router *Router) Post(path string, middlewares []Middleware, handler http.Handler) {
	router.set(path, POST, middlewares, handler)
}

func (router *Router) Delete(path string, middlewares []Middleware, handler http.Handler) {
	router.set(path, DELETE, middlewares, handler)
}

func (router *Router) Put(path string, middlewares []Middleware, handler http.Handler) {
	router.set(path, PUT, middlewares, handler)
}

func (router *Router) Patch(path string, middlewares []Middleware, handler http.Handler) {
	router.set(path, PATCH, middlewares, handler)
}

func (router *Router) Head(path string, middlewares []Middleware, handler http.Handler) {
	router.set(path, HEAD, middlewares, handler)
}

func (router *Router) Options(path string, middlewares []Middleware, handler http.Handler) {
	router.set(path, OPTIONS, middlewares, handler)
}

func (router *Router) Any(path string, middlewares []Middleware, handler http.Handler) {
	router.set(path, "", middlewares, handler)
}

func (router *Router) Dir(path string, dir string, middlewares []Middleware) {
	fs := http.FileServer(http.Dir(dir))

	for _, mw := range middlewares {
		fs = mw(fs)
	}

	router.mux.Handle("GET "+path, http.StripPrefix(path, fs))
}

func (router *Router) Listen(port string) error {

	return http.ListenAndServe(":"+port, router.mux)
}
