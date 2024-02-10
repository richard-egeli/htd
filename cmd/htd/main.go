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

func loginPOST(w http.ResponseWriter, r *http.Request) {
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

func dashboardHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	htdRouter := router.HtdRouter{}
	htdFileserver := router.HtdFileserver{Dir: "/static"}
	isDevelopment := true

	loginRoute := router.HtdRoute{
		Path: "/login",
		GET:  router.Page(pages.LoginPage),
		POST: router.Route(loginPOST),
	}

	defaultRoute := router.HtdRoute{
		Path:    "/",
		GET:     router.Redirect("/login"),
		DEFAULT: router.Page(pages.NotFoundPage),
	}

	if isDevelopment {
		loginRoute.GET.AddMiddleware(router.BrowserSSERefreshMiddleware)
		defaultRoute.GET.AddMiddleware(router.BrowserSSERefreshMiddleware)
		defaultRoute.DEFAULT.AddMiddleware(router.BrowserSSERefreshMiddleware)
		router.EnableBrowserSSEEvents("/events")
	}

	htdRouter.Routes = append(htdRouter.Routes, loginRoute, defaultRoute)
	htdRouter.Create()
	htdFileserver.Create()

	port := ":8080"
	log.Printf("Serving files on http://localhost %s/", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
