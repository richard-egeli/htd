package router

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
)

func CSRFMiddleware(next http.Handler) http.Handler {
	tokenKey := os.Getenv("CSRF_SECRET")
	if len(tokenKey) != 32 {
		log.Fatal("CSRF Token Key Is Invalid")
	}

	csrf := csrf.Protect([]byte(tokenKey), csrf.Secure(false))
	return csrf(next)
}
