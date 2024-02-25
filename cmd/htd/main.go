package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/csrf"
	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
	"github.com/richard-egeli/htd/pkg/router"
	"github.com/richard-egeli/htd/views/pages"
)

func loginPost(w http.ResponseWriter, r *http.Request) {
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

func InitDB() {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bytes, err := os.ReadFile("./sql/create_table.sql")
	if err != nil {
		log.Fatal(err)
	}

	createTable := string(bytes)

	_, err = db.Exec(createTable)
	if err != nil {
		log.Printf("%q: %s\n", err, createTable)
		return
	}
}

func DashboardRouteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dashboard" && r.URL.Path != "/dashboard/" {
			w.Header().Add("Content-Type", "text/html; charset=utf8")
			pages.NotFoundPage(w, r, nil).Render(context.Background(), w)
			log.Println("Not Found Page")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func DefaultRouteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") && len(r.URL.Path) > 1 {
			w.Header().Add("Content-Type", "text/html; charset=utf8")
			pages.NotFoundPage(w, r, nil).Render(context.Background(), w)
			log.Println("Not Found Page")
			return
		}

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

	htdData := pages.HtdData{
		Name:              "htd",
		GenerateCSRFToken: csrf.Token,
	}

	base := router.New()
	dash := base.Sub("dashboard")

	base.Use(base.RefreshMiddleware)
	base.Use(router.CSRFMiddleware)
	base.Use(router.CorsMiddleware)
	base.Use(router.GzipMiddleware)

	base.Dir("/scripts/", "./web/src", []router.Middleware{router.TypescriptTranspilationMiddleware})
	base.Dir("/static/", "./static", nil)

	base.Post("/login", nil, router.Route(loginPost))

	dash.Get("/", []router.Middleware{DashboardRouteMiddleware}, router.Page(pages.HtdPage, &htdData))
	base.Get("/", []router.Middleware{DefaultRouteMiddleware}, router.Page(pages.LoginPage, &loginData))

	base.Listen("8080")
}
