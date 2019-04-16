package handlers

import (
	"html/template"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/go-playground/form"
	"github.com/gorilla/sessions"
)

type curatorPageData struct {
	CurrentUser *models.UserRow
	Phrases     []string
}

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
	r.ParseForm()
	dec := form.NewDecoder()
	vals := r.Form

	var res TestData

	dec.Decode(&res, vals)
	logrus.Infoln(res)
	logrus.Infoln(res.Status["Phrase0"])

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

	data := curatorPageData{CurrentUser: currentUser, Phrases: nil}

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/curator.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}
