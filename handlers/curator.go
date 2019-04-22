package handlers

import (
	"html/template"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/go-playground/form"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/mongo"
)

type curatorPageData struct {
	CurrentUser *models.UserRow
	IsCurator   bool
	Phrases     []curatePhrase
}

type curatePhrase struct {
	PhraseID   string
	PhraseText string
}

// TestData I was testing the "github.com/go-playground/form" library. This helped with parsing array/struct/map like input from html forms
type TestData struct {
	Status map[string]string
}

// GetCurator handles the loading of the curator page.
// It first checks if the user is a curator. If not it should redirect the user to the main page.
// A query to the phrases DB is made and it returns the phrases that are in review.
func GetCurator(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)

	if !ok {
		http.Redirect(w, r, "/now", 302)
		return
	}

	isCurator := currentUser.PermLevel <= models.Curator

	if !isCurator {
		http.Redirect(w, r, "/now", 302)
		return
	}

	mongdb := r.Context().Value("mongodb").(*mongo.Database)
	phrasesCollection := models.NewPhraseConnection(mongdb)
	phrases, _ := models.GetPhraseListForCurators(5, *currentUser, phrasesCollection)

	logrus.Infoln(phrases)

	pagePhrases := []curatePhrase{}

	for _, v := range phrases {
		pagePhrases = append(pagePhrases, curatePhrase{v.PhraseID.Hex(), v.PhraseText})
	}

	data := curatorPageData{CurrentUser: currentUser, IsCurator: isCurator, Phrases: pagePhrases}

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/curator.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}

// PostCurator handles POST requests to the system
func PostCurator(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/now", 302)
		return
	}

	isCurator := currentUser.PermLevel <= models.Curator

	if !isCurator {
		http.Redirect(w, r, "/now", 302)
		return
	}

	mongdb := r.Context().Value("mongodb").(*mongo.Database)
	phrasesCollection := models.NewPhraseConnection(mongdb)

	r.ParseForm()
	dec := form.NewDecoder()

	var res TestData

	dec.Decode(&res, r.Form)

	// TODO: Update DB based on the status of each of the reviewed phrases
	for k, v := range res.Status {
		if v == "accept" {
			logrus.Infoln(k, "was", v)
			models.AcceptPhrase(k, *currentUser, phrasesCollection)
		} else if v == "reject" {
			logrus.Infoln(k, "was", v)
			models.RejectPhrase(k, *currentUser, phrasesCollection)
		}
	}

	// TODO: Load more phrases from DB to put on the view
	phrases, _ := models.GetPhraseListForCurators(5, *currentUser, phrasesCollection)

	pagePhrases := []curatePhrase{}

	for _, v := range phrases {
		pagePhrases = append(pagePhrases, curatePhrase{v.PhraseID.Hex(), v.PhraseText})
	}

	data := curatorPageData{CurrentUser: currentUser, IsCurator: isCurator, Phrases: pagePhrases}

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/curator.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}
