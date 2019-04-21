package handlers

import (
	"html/template"
	"net/http"

	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/gorilla/sessions"
)

type submitPageData struct {
	CurrentUser *models.UserRow
	IsCurator   bool
}

// GetSubmit generates a page for logged in users to submit their own phrases.
func GetSubmit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, _ := session.Values["user"].(*models.UserRow)

	pageData := submitPageData{currentUser, false}

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/submit.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, pageData)
}

// PostSubmit handles the submission of a phrase.
// A phrase will be stored in the phrase DB as a phrase that needs to be reviewed.
// It then redirects the user to the GetSubmit handler
func PostSubmit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, _ := session.Values["user"].(*models.UserRow)

	pageData := submitPageData{currentUser, false}

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/submit.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, pageData)
}
