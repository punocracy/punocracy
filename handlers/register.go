package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/alvarosness/punocracy/model"
)

// GetRegister creates a new user and stores the info in a database
func GetRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	t, _ := template.ParseFiles("templates/dashboard.html.tmpl", "templates/register.html.tmpl")
	t.Execute(w, nil)
}

// PostRegister registers the new user
func PostRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// create user object
	username := r.Form["username"][0]
	email := r.Form["email"][0]
	password := r.Form["password"][0]

	user := &model.User{UserID: 0, Username: username, Email: email, Password: password}

	fmt.Println(user)
	// Store User in database

	// Redirect to Home page
	http.Redirect(w, r, "/", http.StatusFound)
}
