package handlers

import (
	"html/template"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/gorilla/sessions"
)

type aboutPageData struct {
	CurrentUser *models.UserRow
	IsCurator   bool
}

// GetAbout generates the about page
func GetAbout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, _ := session.Values["user"].(*models.UserRow)

	pageData := aboutPageData{CurrentUser: currentUser, IsCurator: false}

	logrus.Infoln(pageData)

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/about.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, pageData)
}
