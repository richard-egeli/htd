package router

import "net/http"

type HtdMiddleware func(http.Handler) http.Handler
