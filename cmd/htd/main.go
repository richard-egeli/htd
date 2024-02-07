package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// type IndexData struct {
// 	Title string
// }
//
// func LoginPage(w http.ResponseWriter, _ *http.Request) {
// 	loginButtonData := button.HtdButton{
// 		Type: button.Submit,
// 		Id:   "loginButton",
// 		Text: "Login",
// 	}
//
// 	resetButtonData := button.HtdButton{
// 		Type: button.Submit,
// 		Text: "Reset Password",
// 	}
//
// 	usernameInputData := input.HtdInput{
// 		Type:        input.Text,
// 		Name:        "username",
// 		Placeholder: "Username",
// 	}
//
// 	passwordInputData := input.HtdInput{
// 		Type:        input.Password,
// 		Name:        "password",
// 		Placeholder: "Password",
// 	}
//
// 	files, err := template.ParseFiles("src/views/index.html", "src/views/login.html", "src/views/button.html", "src/views/input.html")
//
// 	templ := template.Must(files, err)
// 	templ.ExecuteTemplate(w, "index.html", map[string]interface{}{
// 		"Title":         "Main Page",
// 		"LoginButton":   loginButtonData,
// 		"ResetPassword": resetButtonData,
// 		"UsernameInput": usernameInputData,
// 		"PasswordInput": passwordInputData,
// 	})
// }
//
// func LoginEvent(w http.ResponseWriter, r *http.Request) {
// 	if err := r.ParseForm(); err != nil {
// 		fmt.Fprintf(w, "ParseForm() err: %v", err)
// 		return
// 	}
//
// 	username := r.FormValue("username")
// 	password := r.FormValue("password")
//
// 	log.Printf("Username %v Password %v", username, password)
//
// 	time.Sleep(5 * time.Second)
//
// 	if len(username) > 0 {
// 		w.Header().Add("HX-Location", "/dashboard")
// 		w.Header().Add("Location", "/dashboard")
// 	} else {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		fmt.Fprintf(w, "Unauthorized")
// 	}
// }

func createEventHandler() func(http.ResponseWriter, *http.Request) {
	shouldReload := false

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		if !shouldReload {
			shouldReload = true
			fmt.Fprintf(w, "data: %s\n\n", time.Now().Format("2006-01-02T15:04:05Z07:00"))
			flusher.Flush() // Ensure the message is sent immediately
		}
	}
}

func main() {
	h1 := func(w http.ResponseWriter, _ *http.Request) {
		// tmpl := template.Must(template.ParseFiles("./src/views/index.html"))

	}

	fs := http.FileServer(http.Dir("static"))
	http.HandleFunc("/", h1)
	http.HandleFunc("/events", createEventHandler())
	// http.HandleFunc("/login", LoginPage)
	// http.HandleFunc("/events/login", LoginEvent)
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := ":8080"
	log.Printf("Serving files on http://localhost %s/", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
