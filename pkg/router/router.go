package router

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

type HtdMethod string

const (
	GET     HtdMethod = "GET"
	PUT     HtdMethod = "PUT"
	HEAD    HtdMethod = "HEAD"
	POST    HtdMethod = "POST"
	PATCH   HtdMethod = "PATCH"
	DELETE  HtdMethod = "DELETE"
	OPTIONS HtdMethod = "OPTIONS"
)

type HtdHandler struct {
	http.Handler
}

func (h *HtdHandler) IsValid() bool {
	return h.Handler != nil
}

func (h *HtdHandler) AddMiddleware(m HtdMiddleware) {
	h.Handler = m(*h)
}

type HtdMiddleware func(HtdHandler) HtdHandler

type HtdRoute struct {
	Path    string
	GET     HtdHandler
	POST    HtdHandler
	DELETE  HtdHandler
	PUT     HtdHandler
	PATCH   HtdHandler
	HEAD    HtdHandler
	OPTIONS HtdHandler
	DEFAULT HtdHandler
}

func (r *HtdRoute) ApplyMiddlewareAll(m HtdMiddleware) {
	if r.GET.IsValid() {
		r.GET = m(r.GET)
	}

	if r.PUT.IsValid() {
		r.PUT = m(r.PUT)
	}

	if r.HEAD.IsValid() {
		r.HEAD = m(r.HEAD)
	}

	if r.POST.IsValid() {
		r.POST = m(r.POST)
	}

	if r.PATCH.IsValid() {
		r.PATCH = m(r.PATCH)
	}

	if r.DELETE.IsValid() {
		r.DELETE = m(r.DELETE)
	}

	if r.OPTIONS.IsValid() {
		r.OPTIONS = m(r.OPTIONS)
	}

	if r.DEFAULT.IsValid() {
		r.DEFAULT = m(r.DEFAULT)
	}
}

type htdRouteHandler struct {
	method http.Handler
	error  string
	code   int
}

func (r *htdRouteHandler) Init(method http.Handler, error string, code int) {
	r.method = method
	r.error = error
	r.code = code
}

type HtdRouter struct {
	Routes []HtdRoute
}

func (r *HtdRouter) Middleware(m HtdMiddleware) {
	for i := range r.Routes {
		r.Routes[i].ApplyMiddlewareAll(m)
	}
}

func (r *HtdRouter) Create() {
	for i := range r.Routes {
		r.Routes[i].Create()
	}
}

func (route *HtdRoute) Create() {
	routeFunc := func(w http.ResponseWriter, r *http.Request) {
		var handle htdRouteHandler

		switch HtdMethod(r.Method) {
		case GET:
			handle.Init(route.GET, "Not Found", http.StatusNotFound)
		case POST:
			handle.Init(route.POST, "Method Not Allowed", http.StatusMethodNotAllowed)
		case DELETE:
			handle.Init(route.DELETE, "Method Not Allowed", http.StatusMethodNotAllowed)
		case PUT:
			handle.Init(route.DELETE, "Method Not Allowed", http.StatusMethodNotAllowed)
		case PATCH:
			handle.Init(route.DELETE, "Method Not Allowed", http.StatusMethodNotAllowed)
		case HEAD:
			handle.Init(route.DELETE, "Method Not Allowed", http.StatusMethodNotAllowed)
		case OPTIONS:
			handle.Init(route.DELETE, "Method Not Allowed", http.StatusMethodNotAllowed)
		default:
			handle.Init(nil, "Not Found", http.StatusNotFound)
		}

		if handle.method != nil && route.Path == r.URL.Path {
			handle.method.ServeHTTP(w, r)
		} else if route.DEFAULT.IsValid() {
			route.DEFAULT.ServeHTTP(w, r)
		} else {
			http.Error(w, handle.error, handle.code)
		}
	}

	http.HandleFunc(route.Path, routeFunc)
}

func Route(f http.HandlerFunc) HtdHandler {
	return HtdHandler{Handler: http.HandlerFunc(f)}
}

func Page(f func() templ.Component) HtdHandler {
	return HtdHandler{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf8")

		if component := f().Render(context.Background(), w); component != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})}
}

func Redirect(path string) HtdHandler {
	return HtdHandler{http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusFound)
	})}
}
