package router

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
)

type HtdRouter struct {
	globalMiddleware []HtdMiddleware
	routes           map[string]*HtdRoute
}

func Create() HtdRouter {
	return HtdRouter{
		routes: make(map[string]*HtdRoute),
	}
}

func (h *HtdRouter) Use(m HtdMiddleware) {
	h.globalMiddleware = append(h.globalMiddleware, m)
}

func (h *HtdRouter) setMethod(method HtdMethod, path string, mw []HtdMiddleware, handler http.Handler) error {
	if route, ok := h.routes[path]; ok {
		methodFunc := route.GetMethodHandler(method)
		if methodFunc != nil {
			log.Printf("Failed setting " + string(method))
			return errors.New(string(method) + " already contains a route on path " + path)
		}
		route.SetMethodHandler(method, handler)
		log.Printf("Setting POST request " + string(method))

	} else {
		log.Printf("Creating new " + string(method) + " On path " + path)
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
	return h.setMethod(GET, path, middlewares, handler)
}

func (router *HtdRouter) Listen(port int) error {
	for i := range router.routes {
		route := router.routes[i]

		for _, method := range route.GetMethodIterator() {
			handler := route.GetMethodHandler(method)

			if handler != nil {
				for _, middleware := range router.globalMiddleware {
					*handler = middleware(*handler)
				}
			}
		}

		http.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			err := route.Handler(w, r)
			if err == nil || HtdMethod(r.Method) != GET {
				return
			}

			if handler, ok := router.routes["*"]; ok {
				if handler.Get != nil {
					(*handler.Get).ServeHTTP(w, r)
				}
			}
		})
	}

	if err := http.ListenAndServe(":"+fmt.Sprint(port), nil); err != nil {
		return err
	}

	return nil
}

func (h *HtdRouter) EnableBrowserReload() {
	h.Use(BrowserSSERefreshMiddleware)
	EnableBrowserSSEEvents()
}

func Route(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(f)
}

func Page(f func() templ.Component) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf8")
		if component := f().Render(context.Background(), w); component != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

func Redirect(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("HX-Redirect", path)
		http.Redirect(w, r, path, http.StatusFound)
	})
}
