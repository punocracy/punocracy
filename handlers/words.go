package handlers

import (
	"html/template"
	"net/http"

	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type wordPageData struct {
	CurrentUser *models.UserRow
	IsCurator   bool
	Words       []string
}

// GetWords loads the page listing all of the words in our system
func GetWords(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, _ := session.Values["user"].(*models.UserRow)
	// if !ok {
	// 	http.Redirect(w, r, "/logout", 302)
	// 	return
	// }

	vars := mux.Vars(r)
	_ = vars["letter"]

	//
	// TODO: Query DB for words that start with the letter

	pageData := wordPageData{CurrentUser: currentUser, IsCurator: false, Words: nil}

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/word.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, pageData)
}
