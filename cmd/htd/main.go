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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file")
	}

	htdRouter := router.Create()
	htdSubRouter := htdRouter.Sub("/dashboard")

	htdFileserver := router.HtdFileserver{Dir: "/static"}

	loginData := pages.LoginData{
		GenerateCSRFToken: csrf.Token,
		Title:             "Hello, World",
	}

	dashboardData := pages.DashboardData{
		GenerateCSRFToken: csrf.Token,
	}

	htdRouter.Use(router.RefreshMiddleware)
	htdRouter.Use(router.CSRFMiddleware)
	htdRouter.Use(router.CorsMiddleware)

	htdRouter.Get("/login", nil, router.Page(pages.LoginPage, &loginData))
	htdRouter.Post("/login", nil, router.Route(loginPost))
	htdRouter.Post("/logout", nil, router.Redirect("/login"))
	htdRouter.Get("/", nil, router.Redirect("/login"))
	htdSubRouter.Get("/", nil, router.Page(pages.DashboardPage, &dashboardData))
	htdRouter.Get("*", nil, router.Page(pages.NotFoundPage, nil))

	htdFileserver.Create()

	port := 8080
	log.Printf("Serving files on http://localhost %d/", port)
	if err := htdRouter.Listen(port); err != nil {
		log.Fatal(err)
	}
}
