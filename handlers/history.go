package handlers

import (
	"html/template"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/mongo"
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

	// Getting submitted phrases
	mongdb := r.Context().Value("mongodb").(*mongo.Database)
	phrasesCollection := models.NewPhraseConnection(mongdb)

	phrases, err := models.GetPhraseHistory(*currentUser, phrasesCollection)
	if err != nil {
		logrus.Error(err.Error())
	}

	submittedPhrases := []string{}

	for _, phrase := range phrases {
		submittedPhrases = append(submittedPhrases, phrase.PhraseText)
	}

	pageData := historyPageData{CurrentUser: currentUser, IsCurator: isCurator, RatedPhrases: nil, SubmittedPhrases: submittedPhrases}

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
