package router

import (
	"fmt"
	"net/http"
	"time"
)

var initialized = false

// Handles sending a <script> tag that enables automatic browser refresh over SSE (Server Sent Events)
func (router *Router) SetupBrowserRefreshEvent() {
	if !initialized {
		eventHandler := func() http.HandlerFunc {
			shouldReload := false

			return func(w http.ResponseWriter, r *http.Request) {
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

		router.mux.HandleFunc("GET /server/sent/event/browser/reload", eventHandler())
		initialized = true
	}
}
