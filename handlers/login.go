package handlers

import (
	"html/template"
	"net/http"
)

// LoginHandler is responsible for logins
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	t, _ := template.ParseFiles("templates/dashboard.html.tmpl", "templates/login.html.tmpl")
	t.Execute(w, nil)
}
