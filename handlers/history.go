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
	IsCurator        bool
	RatedPhrases     []string
	SubmittedPhrases []string
}

// GetHistory generates a page showing the users' history of phrase ratings and phrase submissions
func GetHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)

	var isCurator bool

	if !ok {
		currentUser = nil
		isCurator = false
	} else {
		isCurator = currentUser.PermLevel <= models.Curator
	}

	pageData := historyPageData{CurrentUser: currentUser, IsCurator: isCurator, RatedPhrases: nil, SubmittedPhrases: nil}

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
