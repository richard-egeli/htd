package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/richard-egeli/htd/views"
)

// Handles sending a <script> tag that enables automatic browser refresh over SSE (Server Sent Events)
func BrowserSSERefreshMiddleware(next HtdHandler) HtdHandler {
	return HtdHandler{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		// Add a script tag after handling the next event, in order to insert a <script> tag at the
		// Very end of the HTML page
		if r.Method == "GET" {
			err := views.ReloadScript().Render(context.Background(), w)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	})}
}

func EnableBrowserSSEEvents(path string) {
	eventHandler := func() http.HandlerFunc {
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

	http.HandleFunc(path, eventHandler())
}
