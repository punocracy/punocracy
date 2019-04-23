package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

type homePageData struct {
	CurrentUser *models.UserRow
	IsCurator   bool
	Words       []string
	Phrases     []string
}

type resultPageData struct {
	CurrentUser *models.UserRow
	QueryWord   string
	IsCurator   bool
	NoPhrases   bool
	NoWords     bool
	Puns        []string
	Phrases     []models.Phrase
}

// HandleRoot redirects to now
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/now", 302)
}

// HandleNotFound shows a 404 page
// func HandleNotFound(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

// 	session, _ := sessionStore.Get(r, "punocracy-session")
// 	currentUser, ok := session.Values["user"].(*models.UserRow)

// 	var isCurator bool

// 	if !ok {
// 		currentUser = nil
// 		isCurator = false
// 	} else {
// 		isCurator = currentUser.PermLevel <= models.Curator
// 	}
// 	pageData := homePageData{CurrentUser: currentUser, IsCurator: isCurator, Words: nil, Phrases: nil}

// 	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/not-found.html.tmpl")
// 	if err != nil {
// 		libhttp.HandleErrorJson(w, err)
// 		return
// 	}

// 	tmpl.Execute(w, pageData)
// }

// GetHome generates the home page of the system
func GetHome(w http.ResponseWriter, r *http.Request) {
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

	db := r.Context().Value("db").(*sqlx.DB)
	wordTable := models.NewWord(db)

	words, _ := wordTable.RandWordsList(nil, 5)

	pageData := homePageData{CurrentUser: currentUser, IsCurator: isCurator, Words: words, Phrases: nil}

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

	var isCurator bool

	if !ok {
		currentUser = nil
		isCurator = false
	} else {
		isCurator = currentUser.PermLevel <= models.Curator
	}

	queryWord := r.FormValue("queryWord")

	// TODO: Query DB for words in the same word group
	db := r.Context().Value("db").(*sqlx.DB)
	wordTable := models.NewWord(db)

	var noPhrases bool
	var noWords bool

	words, wordErr := wordTable.QueryHlistString(nil, strings.ToLower(queryWord))

	if wordErr != nil {
		noWords = true
	}

	// TODO: Query DB for phrases
	mongdb := r.Context().Value("mongodb").(*mongo.Database)
	phrasesCollection := models.NewPhraseConnection(mongdb)
	phrases, phraseErr := models.GetPhraseList(words, phrasesCollection)

	if phraseErr != nil {
		noPhrases = true
	}

	if len(phrases) == 0 {
		noPhrases = true
	}

	// TODO: perform word replacement
	puns := models.GeneratePuns(queryWord, words, phrases)
	pageData := resultPageData{CurrentUser: currentUser, QueryWord: queryWord, IsCurator: isCurator, NoPhrases: noPhrases, NoWords: noWords, Puns: puns, Phrases: phrases}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/search.html.tmpl", "templates/query.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, pageData)
}
