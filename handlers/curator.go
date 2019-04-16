package handlers

import (
	"html/template"
	"net/http"

	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/go-playground/form"
	"github.com/gorilla/sessions"
)

type curatorPageData struct {
	CurrentUser *models.UserRow
	Phrases     []string
}

// I was testing the "github.com/go-playground/form" library. This helped with parsing array/struct/map like input from html forms
type TestData struct {
	Status map[string]string
}

// GetCurator handles requests for the curator. I need to better document these functions
func GetCurator(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	// TODO: Check if user is curator
	// if not curator redirect to main page

	// TODO: Query DB for a max number of phrases to be reviewed

	data := curatorPageData{CurrentUser: currentUser, Phrases: nil}

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
		http.Redirect(w, r, "/logout", 302)
		return
	}

	// TODO: Check if user is curator
	// if not curator redirect to main page

	r.ParseForm()
	dec := form.NewDecoder()

	var res TestData

	dec.Decode(&res, r.Form)
	// TODO: Update DB based on the status of each of the reviewed phrases

	// TODO: Load more phrases from DB to put on the view

	data := curatorPageData{CurrentUser: currentUser, Phrases: nil}

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/curator.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}
