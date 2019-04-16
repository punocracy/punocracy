package handlers

import (
	"html/template"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/gorilla/sessions"
)

type homePageData struct {
	CurrentUser *models.UserRow
	Words       []string
	Phrases     []string
}

type resultPageData struct {
	CurrentUser *models.UserRow
	FoundWord   bool
	Puns        []string
}

// GetHome generates the home page of the system
func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	// TODO: Query DB for random words and top rated phrases

	pageData := homePageData{CurrentUser: currentUser, Words: nil, Phrases: nil}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/search.html.tmpl", "templates/home.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, pageData)
}

// PostHome posts home
func PostHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	logrus.Infoln(r.URL.Path)

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	// word := r.FormValue("queryWord")
	// TODO: Query DB for words in the same word group
	// TODO: Query DB for phrases and perform word replacement

	pageData := resultPageData{CurrentUser: currentUser, FoundWord: false, Puns: nil}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/search.html.tmpl", "templates/query.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, pageData)
}
