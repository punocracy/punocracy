package handlers

import (
	"html/template"
	"net/http"

	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/gorilla/sessions"
)

type historyPageData struct {
	CurrentUser      *models.UserRow
	RatedPhrases     []string
	SubmittedPhrases []string
}

// GetHistory generates a page showing the users' history of phrase ratings and phrase submissions
func GetHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	pageData := historyPageData{currentUser, nil, nil}

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/history.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, pageData)
}

// PostHistory handles the update of user ratings for phrases
func PostHistory(w http.ResponseWriter, r *http.Request) {

}
