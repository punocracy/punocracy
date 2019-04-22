package handlers

import (
	"html/template"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

type submitPageData struct {
	CurrentUser *models.UserRow
	IsCurator   bool
}

// TODO: Address illegal access to these views

// GetSubmit generates a page for logged in users to submit their own phrases.
func GetSubmit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)

	var isCurator bool

	if !ok {
		http.Redirect(w, r, "/now", http.StatusBadRequest)
	} else {
		isCurator = currentUser.PermLevel <= models.Curator
	}

	pageData := submitPageData{currentUser, isCurator}

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
	currentUser, ok := session.Values["user"].(*models.UserRow)

	var isCurator bool

	if !ok {
		http.Redirect(w, r, "/now", http.StatusBadRequest)
	} else {
		isCurator = currentUser.PermLevel <= models.Curator
	}

	db := r.Context().Value("db").(*sqlx.DB)

	mongdb := r.Context().Value("mongodb").(*mongo.Database)
	phrase := r.FormValue("phraseText")

	phrasesCollection := models.NewPhraseConnection(mongdb)
	word := models.NewWord(db)

	err := models.InsertPhrase(phrase, *currentUser, word, phrasesCollection)
	logrus.Infoln("Before")
	if err != nil {
		logrus.Errorln(err.Error())
		logrus.Infoln("After")
		// TODO: Handle multiple types of errors
		http.Redirect(w, r, "/now", 302)
		return
	}

	pageData := submitPageData{currentUser, isCurator}

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/submit.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, pageData)
}
