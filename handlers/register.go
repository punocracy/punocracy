package handlers

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/alvarosness/punocracy/model"
)

// RegisterHandler creates a new user and stores the info in a database
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html")

		t, _ := template.ParseFiles("templates/dashboard.html.tmpl", "templates/register.html.tmpl")
		t.Execute(w, nil)
	}

	if r.Method == "POST" {
		r.ParseForm()

		// create user object
		h := sha256.New()
		username := strings.Join(r.Form["username"], "")
		email := strings.Join(r.Form["email"], "")
		password := strings.Join(r.Form["password"], "")

		h.Write([]byte(password))

		user := &model.User{UserID: 0, Username: username, Email: email, Password: h.Sum(nil)}

		fmt.Println(user)
		// Store User in database

		// Redirect to Home page
	}

}
