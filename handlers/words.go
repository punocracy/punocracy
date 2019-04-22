package handlers

import (
	"html/template"
	"net/http"

	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

type wordPageData struct {
	CurrentUser *models.UserRow
	IsCurator   bool
	Words       []string
}

// GetWords loads the page listing all of the words in our system
func GetWords(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	letter := rune(vars["letter"][0])

	db := r.Context().Value("db").(*sqlx.DB)

	wordTable := models.NewWord(db)
	wordsRows, _ := wordTable.QueryAlph(nil, letter)

	words := []string{}

	for _, v := range wordsRows {
		words = append(words, v.Word)
	}

	pageData := wordPageData{CurrentUser: currentUser, IsCurator: isCurator, Words: words}

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/word.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, pageData)
}
