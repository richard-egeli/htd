package router

import (
	"errors"
	"net/http"
)

type HtdRoute struct {
	Path    string
	Get     *http.Handler
	Post    *http.Handler
	Delete  *http.Handler
	Put     *http.Handler
	Patch   *http.Handler
	Head    *http.Handler
	Options *http.Handler
}

func (r *HtdRoute) GetMethodIterator() []HtdMethod {
	return []HtdMethod{
		GET,
		POST,
		DELETE,
		PUT,
		PATCH,
		HEAD,
		OPTIONS,
	}
}

func (route *HtdRoute) SetMethodHandler(method HtdMethod, handler http.Handler) {
	switch method {
	case GET:
		route.Get = &handler
	case POST:
		route.Post = &handler
	case DELETE:
		route.Delete = &handler
	case PUT:
		route.Put = &handler
	case PATCH:
		route.Patch = &handler
	case HEAD:
		route.Head = &handler
	case OPTIONS:
		route.Options = &handler
	}
}

func (route *HtdRoute) GetMethodHandler(method HtdMethod) *http.Handler {
	switch method {
	case GET:
		return route.Get
	case POST:
		return route.Post
	case DELETE:
		return route.Delete
	case PUT:
		return route.Put
	case PATCH:
		return route.Patch
	case HEAD:
		return route.Head
	case OPTIONS:
		return route.Options
	default:
		return nil
	}
}

func (route *HtdRoute) Handler(w http.ResponseWriter, r *http.Request) error {
	handler := route.GetMethodHandler(HtdMethod(r.Method))

	if handler != nil && r.URL.Path == route.Path {
		(*handler).ServeHTTP(w, r)
		return nil
	}

	return errors.New("Failed path checks")
}
