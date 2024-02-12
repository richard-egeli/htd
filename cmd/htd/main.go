package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/richard-egeli/htd/pkg/router"
	"github.com/richard-egeli/htd/views/pages"
)

func loginPost(w http.ResponseWriter, r *http.Request) {
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
		SameSite: http.SameSiteStrictMode,
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

func testMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Test Middleware called from %s", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file")
	}

	loginData := pages.LoginData{
		GenerateCSRFToken: csrf.Token,
		Title:             "Login",
	}

	dashboardData := pages.DashboardData{
		GenerateCSRFToken: csrf.Token,
	}

	port := 8080
	base := router.Create()
	dashboard := base.Sub("/dashboard")

	base.Use(router.RefreshMiddleware)
	base.Use(router.CSRFMiddleware)
	base.Use(router.CorsMiddleware)
	base.Dir("/static")

	dashboard.Use(testMiddleware)

	// Main routes
	base.Get("/login", nil, router.Page(pages.LoginPage, &loginData))
	base.Post("/login", nil, router.Route(loginPost))
	base.Post("/logout", nil, router.Redirect("/login"))
	base.Get("/", nil, router.Redirect("/login"))
	base.Get("*", nil, router.Page(pages.NotFoundPage, nil))

	// Dashboard sub routes
	dashboard.Get("/", nil, router.Page(pages.DashboardPage, &dashboardData))

	log.Printf("Serving files on http://localhost %d/", port)
	if err := base.Listen(port); err != nil {
		log.Fatal(err)
	}
}
