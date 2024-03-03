package router

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func Route(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(f)
}

func Component(component func(any ...any) templ.Component) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf8")

		if component := component().Render(context.Background(), w); component != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})
}

func Page[T any](page func(http.ResponseWriter, *http.Request, *T) templ.Component, data *T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf8")

		if component := page(w, r, data).Render(context.Background(), w); component != nil {
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
