package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

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

	cookies := r.Cookies()

	for i, r := range cookies {
		fmt.Printf("Index %d, Value %s", i, r)
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	cookie := http.Cookie{
		Name:     "TestCookie",
		Value:    "Go Programming",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
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
	htdRouter := router.Create()
	htdFileserver := router.HtdFileserver{Dir: "/static"}

	htdRouter.EnableBrowserReload()
	htdRouter.Get("/login", nil, router.Page(pages.LoginPage))
	htdRouter.Post("/login", nil, router.Route(loginPost))
	htdRouter.Get("/", nil, router.Redirect("/login"))
	htdRouter.Get("*", nil, router.Page(pages.NotFoundPage))

	htdFileserver.Create()

	port := ":8080"
	log.Printf("Serving files on http://localhost %s/", port)

	if err := htdRouter.Listen(8080); err != nil {
		log.Fatal(err)
	}
}
