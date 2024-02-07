package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/richard-egeli/htd/views"
)

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

func main() {
	h1 := func(w http.ResponseWriter, _ *http.Request) {

	}

	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/", h1)
	http.HandleFunc("/events", createEventHandler())
	// http.HandleFunc("/login", LoginPage)
	// http.HandleFunc("/events/login", LoginEvent)
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := ":8080"
	log.Printf("Serving files on http://localhost %s/", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
