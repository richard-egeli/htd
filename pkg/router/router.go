package router

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

type HtdRouter struct {
	globalMiddleware []HtdMiddleware
	routes           map[string]*HtdRoute
	parent           *HtdRouter
	subRouters       []*HtdRouter
	subPath          string
	fileDir          string
}

func Create() *HtdRouter {
	router := new(HtdRouter)
	router.routes = make(map[string]*HtdRoute)
	return router
}

func (h *HtdRouter) Dir(dir string) {
	h.fileDir = dir
}

func (h *HtdRouter) Use(m HtdMiddleware) {
	h.globalMiddleware = append(h.globalMiddleware, m)
}

func (h *HtdRouter) Sub(path string) *HtdRouter {
	sub := Create()
	sub.parent = h
	sub.subPath = path
	h.subRouters = append(h.subRouters, sub)
	return sub
}

func (h *HtdRouter) getAbsolutePath(path string) string {
	subPath := h.subPath
	parent := h.parent

	for parent != nil {
		subPath = parent.subPath + subPath
		parent = parent.parent
	}

	return subPath + path
}

func (h *HtdRouter) setMethod(method HtdMethod, path string, mw []HtdMiddleware, handler http.Handler) error {
	if route, ok := h.routes[path]; ok {
		methodFunc := route.GetMethodHandler(method)

		if methodFunc != nil {
			log.Printf("Failed setting " + string(method))
			return errors.New(string(method) + " already contains a route on path " + path)
		}

		route.SetMethodHandler(method, handler)
	} else {
		newRoute := HtdRoute{Path: path}
		newRoute.SetMethodHandler(method, handler)
		h.routes[path] = &newRoute
	}

	handle := h.routes[path].GetMethodHandler(method)
	for i := len(mw) - 1; i >= 0; i-- {
		m := mw[i]
		*handle = m(*handle)
	}

	return nil
}

func (h *HtdRouter) Post(path string, middlewares []HtdMiddleware, handler http.Handler) error {
	return h.setMethod(POST, path, middlewares, handler)
}

func (h *HtdRouter) Get(path string, middlewares []HtdMiddleware, handler http.Handler) error {
	return h.setMethod(GET, h.subPath+path, middlewares, handler)
}

func (router *HtdRouter) applyDefaultRoutesRecursive(defaultRoute *HtdRoute) {
	for i := range router.routes {
		route := router.routes[i]
		path := route.Path

		if len(path) > 1 {
			path = "/" + strings.Trim(path, "/")
		}

		http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			err := route.Handler(w, r)
			if err == nil {
				return
			}

			log.Printf("Accessing unknown path '%s' -- '%v'", route.Path, err)
			if defaultRoute != nil {
				method := defaultRoute.GetMethodHandler(HtdMethod(r.Method))
				if method != nil {
					(*method).ServeHTTP(w, r)
				}
			}
		})
	}

	for i := range router.subRouters {
		route := router.subRouters[i]
		route.applyDefaultRoutesRecursive(defaultRoute)
	}
}

func (router *HtdRouter) applyMiddlewareRecursive(parentMiddleware *[]HtdMiddleware) {
	if parentMiddleware == nil {
		return
	}

	for i := range router.routes {
		route := router.routes[i]

		for _, j := range route.GetMethodIterator() {
			if handle := route.GetMethodHandler(j); handle != nil {
				for _, mw := range *parentMiddleware {
					*handle = mw(*handle)
				}
			}
		}
	}

	for i := range router.subRouters {
		subRouter := router.subRouters[i]
		subRouter.applyMiddlewareRecursive(&subRouter.globalMiddleware)

		if parentMiddleware != nil {
			subRouter.applyMiddlewareRecursive(parentMiddleware)
		}
	}
}
