package handlers

import (
	"html/template"
	"net/http"

	"github.com/Sirupsen/logrus"
)

// GetNow renders the main query page
func GetNow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	t, _ := template.ParseFiles("templates/dashboard.html.tmpl", "templates/now.html.tmpl")
	t.Execute(w, nil)
}

// PostNow does that one thing
func PostNow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	r.ParseForm()

	logrus.Infoln(r.PostForm)

	t, _ := template.ParseFiles("templates/dashboard.html.tmpl", "templates/now.html.tmpl")
	t.Execute(w, nil)
}
