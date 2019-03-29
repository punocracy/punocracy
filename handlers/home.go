package handlers

import (
	"html/template"
	"net/http"
)

// HomeHandler handles the main page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	t, _ := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")
	t.Execute(w, nil)
}
