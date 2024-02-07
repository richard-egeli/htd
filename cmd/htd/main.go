package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

import "github.com/richard-egeli/htd/views/pages"

func createEventHandler() func(http.ResponseWriter, *http.Request) {
	shouldReload := false

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		if !shouldReload {
			shouldReload = true
			fmt.Fprintf(w, "data: %s\n\n", time.Now().Format("2006-01-02T15:04:05Z07:00"))
			flusher.Flush() // Ensure the message is sent immediately
		}
	}
}

func loginGET(w http.ResponseWriter, r *http.Request) {
	component := pages.LoginLayout()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := component.Render(context.Background(), w)
	if err != nil {
		http.Error(w, "Failed to render component", http.StatusInternalServerError)
		return
	}
}

func loginPOST(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form input", http.StatusInternalServerError)
		return
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
		http.Error(w, "Invalid username / password", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &cookie)
	w.Header().Add("HX-Redirect", "/dashboard")
}

func pageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	component := pages.PageNotFound()
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		loginGET(w, r)
	case "POST":
		loginPOST(w, r)
	}
}

func main() {
	fs := http.FileServer(http.Dir("static"))

	http.HandleFunc("/events", createEventHandler())
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/", pageNotFoundHandler)
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := ":8080"
	log.Printf("Serving files on http://localhost %s/", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
