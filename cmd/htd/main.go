package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/richard-egeli/htd/pkg/router"
	"github.com/richard-egeli/htd/views"
	"github.com/richard-egeli/htd/views/pages"
)

func createEventHandler() func(http.ResponseWriter, *http.Request) {
	shouldReload := false

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		if !shouldReload {
			flusher, ok := w.(http.Flusher)
			if !ok {
				http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
				return
			}

			shouldReload = true
			fmt.Fprintf(w, "data: %s\n\n", time.Now().Format("2006-01-02T15:04:05Z07:00"))
			flusher.Flush() // Ensure the message is sent immediately
		}
	}
}

func loginPOST(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form input", http.StatusInternalServerError)
		return
	}

	cookies := r.Cookies()

	for i, r := range cookies {
		fmt.Printf("Index %d, Value %s", i, r)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	cookie := http.Cookie{
		Name:     "TestCookie",
		Value:    "Go Programming",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
		MaxAge:   3600,
	}

	log.Printf("Username %s", username)
	log.Printf("Password %s", password)

	if len(username) <= 0 || len(password) <= 0 {
		pages.LoginErrorPage().Render(context.Background(), w)
		return
	}

	http.SetCookie(w, &cookie)
	w.Header().Add("HX-Redirect", "/dashboard")
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

}

func middlewareRefreshEvent(next router.HtdHandler) router.HtdHandler {
	return router.HtdHandler{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		// Add a script tag after handling the next event, in order to insert a <script> tag at the
		// Very end of the HTML page
		if r.Method == "GET" {
			log.Printf("Sending reload script!")
			err := views.ReloadScript().Render(context.Background(), w)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	})}
}

func setupFileServer() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
}

func main() {
	htdRouter := router.HtdRouter{}
	isDevelopment := true

	loginRoute := router.HtdRoute{
		Path: "/login",
		GET:  router.Page(pages.LoginPage),
		POST: router.Route(loginPOST),
	}

	defaultRoute := router.HtdRoute{
		Path:    "/",
		GET:     router.Redirect("/login"),
		DEFAULT: router.Page(pages.NotFoundPage),
	}

	if isDevelopment {
		loginRoute.GET.AddMiddleware(middlewareRefreshEvent)
		defaultRoute.GET.AddMiddleware(middlewareRefreshEvent)
		defaultRoute.DEFAULT.AddMiddleware(middlewareRefreshEvent)
		http.HandleFunc("/events", createEventHandler())
	}

	htdRouter.Routes = append(htdRouter.Routes, loginRoute, defaultRoute)
	htdRouter.Create()
	setupFileServer()

	port := ":8080"
	log.Printf("Serving files on http://localhost %s/", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
