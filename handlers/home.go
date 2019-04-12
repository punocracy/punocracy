package handlers

import (
	"html/template"
	"net/http"

	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/gorilla/sessions"
)

// type HomeRow struct {
// 	Words   []string
// 	Phrases []string
// }

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	words := []string{"Quill", "Fun", "Pun"}
	// words := []string{}
	phrases := []string{}
	// phrases := []string{"All your Bayes are belong to us", "I was up kniting all night last night. I am sew tired."}

	data := struct {
		CurrentUser *models.UserRow
		Words       []string
		Phrases     []string
	}{
		currentUser,
		words,
		phrases,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}

func PostHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	data := struct {
		CurrentUser *models.UserRow
	}{
		currentUser,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}
