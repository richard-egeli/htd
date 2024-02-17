package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

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
	w.Header().Add("HX-Redirect", "/htd")
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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file")
	}

	InitDB()

	loginData := pages.LoginData{
		GenerateCSRFToken: csrf.Token,
		Title:             "Login",
	}

	htdData := pages.HtdData{
		Name:              "htd",
		GenerateCSRFToken: csrf.Token,
	}

	port := 8080
	base := router.Create()
	htd := base.Sub("/" + htdData.Name)

	base.Use(router.RefreshMiddleware)
	base.Use(router.CSRFMiddleware)
	base.Use(router.CorsMiddleware)
	base.Dir("/static")

	// Main routes
	base.Get("/login", nil, router.Page(pages.LoginPage, &loginData))
	base.Post("/login", nil, router.Route(loginPost))
	base.Post("/logout", nil, router.Redirect("/login"))
	base.Get("/", nil, router.Redirect("/login"))
	base.Get("*", nil, router.Page(pages.NotFoundPage, nil))

	// Dashboard sub routes
	htd.Get("/", nil, router.Page(pages.HtdPage, &htdData))
	htd.Get("/pages", nil, router.Page(pages.Pages, &pages.PagesData{HtdData: &htdData}))

	log.Printf("Serving files on http://localhost %d/", port)
	if err := base.Listen(os.Getenv("SERVER_PORT")); err != nil {
		log.Fatal(err)
	}
}
