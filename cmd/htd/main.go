package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	_ "github.com/mattn/go-sqlite3"

	"github.com/joho/godotenv"
	"github.com/richard-egeli/htd/pkg/router"
	"github.com/richard-egeli/htd/pkg/store"
	"github.com/richard-egeli/htd/views/layout"
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
		if r.URL.Path != "/" {
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

	if err := store.Open(); err != nil {
		log.Fatal("Failed to open sqlite database")
	}

	defer store.Close()

	loginData := pages.LoginData{
		GenerateCSRFToken: csrf.Token,
		Title:             "Login",
	}

	sidebarData := layout.SidebarData{
		GenerateCSRFToken: csrf.Token,
	}

	dashboardData := pages.DashboardData{
		SidebarData: &sidebarData,
	}

	settingsData := pages.SettingsData{
		SidebarData: &sidebarData,
	}

	ordersData := pages.OrdersData{
		SidebarData: &sidebarData,
	}

	productsData := pages.ProductsData{
		SidebarData: &sidebarData,
	}

	base := router.New()
	dash := base.Sub("dashboard")
	api := base.Sub("api")

	// base.Use(router.CSRFMiddleware)
	base.Use(router.CorsMiddleware)
	base.Use(router.GzipMiddleware)
	base.SetupBrowserRefreshEvent()

	base.Dir("/scripts/", "./web/src", []router.Middleware{router.TypescriptTranspilationMiddleware})
	base.Dir("/static/", "./static", nil)

	base.Post("/login", nil, router.Route(loginPost))
	base.Post("/logout", nil, router.Redirect("/login"))

	base.Get("/", []router.Middleware{DefaultRouteMiddleware}, router.Redirect("/login"))
	base.Get("/login", nil, router.Page(pages.LoginPage, &loginData))

	dash.Get("/", []router.Middleware{DashboardRouteMiddleware}, router.Page(pages.DashboardPage, &dashboardData))
	dash.Get("/settings", nil, router.Page(pages.SettingsPage, &settingsData))
	dash.Get("/products", nil, router.Page(pages.ProductsPage, &productsData))
	dash.Get("/orders", nil, router.Page(pages.OrdersPage, &ordersData))

	api.Post("/products/create", nil, router.Route(store.CreateProduct))
	api.Get("/products/fetch", nil, router.Route(store.FetchProductsAll))
	api.Post("/images", nil, router.Route(store.CreateImage))
	api.Delete("/images/{id}", nil, router.Route(store.DeleteImage))

	base.Listen("8080")
}
