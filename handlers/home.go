package handlers

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/form"

	"go.mongodb.org/mongo-driver/bson/primitive"

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
	Phrases     []phraseDisplay
}

type phraseDisplay struct {
	PhraseID            string
	PhraseText          string
	Author              string
	TimeSinceSubmission string
	IsOneStar           bool
	IsTwoStar           bool
	IsThreeStar         bool
	IsFourStar          bool
	IsFiveStar          bool
}

type resultPageData struct {
	CurrentUser *models.UserRow
	QueryWord   string
	IsCurator   bool
	NoPhrases   bool
	NoWords     bool
	Puns        []string
	Phrases     []phraseDisplay
}

type phraseRatings struct {
	Ratings map[string]string
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

	currentUser, isCurator := getUser(session)

	db := r.Context().Value("db").(*sqlx.DB)
	wordTable := models.NewWord(db)

	words, _ := wordTable.RandWordsList(nil, 5)

	samplephrase := models.Phrase{
		PhraseID:        primitive.NewObjectID(),
		SubmitterUserID: 5,
		SubmissionDate:  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		PhraseRatings:   models.Rating{OneStar: 0, TwoStar: 0, ThreeStar: 1, FourStar: 4, FiveStar: 1},
		WordList:        []int{},
		ReviewedBy:      5,
		ReviewDate:      time.Now(),
		PhraseText:      "This is a test",
		DisplayPublic:   models.Accepted,
	}

	userTable := models.NewUser(db)
	sampleUser, _ := userTable.GetByID(nil, 5)
	now := time.Now()
	sampleTime := now.Sub(samplephrase.SubmissionDate)
	avgRating := math.Round(models.AverageRating(samplephrase.PhraseRatings))
	phraseList := []phraseDisplay{}

	logrus.Infoln(avgRating)
	phraseList = append(phraseList, phraseDisplay{
		PhraseText:          samplephrase.PhraseText,
		Author:              sampleUser.Username,
		TimeSinceSubmission: sampleTime.String(),
		IsOneStar:           avgRating == 1,
		IsTwoStar:           avgRating == 2,
		IsThreeStar:         avgRating == 3,
		IsFourStar:          avgRating == 4,
		IsFiveStar:          avgRating == 5,
	})

	pageData := homePageData{CurrentUser: currentUser, IsCurator: isCurator, Words: words, Phrases: phraseList}

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

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")

	currentUser, isCurator := getUser(session)

	r.ParseForm()

	queryWord := r.FormValue("queryWord")

	if queryWord == "" {
		decoder := form.NewDecoder()
		var ratings phraseRatings
		decoder.Decode(&ratings, r.Form)
		logrus.Infoln(ratings)

		mongdb := r.Context().Value("mongodb").(*mongo.Database)
		phrasesCollection := models.NewPhraseConnection(mongdb)
		ratingsCollection := models.NewUserRatingsConnection(mongdb)

		for k, v := range ratings.Ratings {
			phrID, _ := primitive.ObjectIDFromHex(k)
			rating, _ := strconv.Atoi(v)

			phr, _ := models.GetPhraseByID(phrID, phrasesCollection)
			models.AddOrChangeRating(*currentUser, rating, phr, phrasesCollection, ratingsCollection)
		}

		http.Redirect(w, r, "/now", 302)
		return
	} else {
		db := r.Context().Value("db").(*sqlx.DB)
		wordTable := models.NewWord(db)

		var noPhrases bool
		var noWords bool

		words, wordErr := wordTable.QueryHlistString(nil, strings.ToLower(queryWord))

		if wordErr != nil {
			noWords = true
		}

		mongdb := r.Context().Value("mongodb").(*mongo.Database)
		phrasesCollection := models.NewPhraseConnection(mongdb)
		phrases, phraseErr := models.GetPhraseList(words, phrasesCollection)

		if phraseErr != nil {
			noPhrases = true
		}

		if len(phrases) == 0 {
			noPhrases = true
		}

		userTable := models.NewUser(db)
		puns := models.GeneratePuns(queryWord, words, phrases)
		phraseList := []phraseDisplay{}

		for _, phrase := range phrases {
			submitter, _ := userTable.GetByID(nil, phrase.SubmitterUserID)
			now := time.Now()
			timeSinceSubmission := now.Sub(phrase.SubmissionDate)
			avgRating := math.Round(models.AverageRating(phrase.PhraseRatings))

			phraseList = append(phraseList, phraseDisplay{
				PhraseID:            phrase.PhraseID.Hex(),
				PhraseText:          phrase.PhraseText,
				Author:              submitter.Username,
				TimeSinceSubmission: timeSinceSubmission.String(),
				IsOneStar:           avgRating == 1,
				IsTwoStar:           avgRating == 2,
				IsThreeStar:         avgRating == 3,
				IsFourStar:          avgRating == 4,
				IsFiveStar:          avgRating == 5,
			})
		}
		pageData := resultPageData{CurrentUser: currentUser, QueryWord: queryWord, IsCurator: isCurator, NoPhrases: noPhrases, NoWords: noWords, Puns: puns, Phrases: phraseList}

		tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/search.html.tmpl", "templates/query.html.tmpl")
		if err != nil {
			libhttp.HandleErrorJson(w, err)
			return
		}

		tmpl.Execute(w, pageData)
	}

}

func getUser(session *sessions.Session) (*models.UserRow, bool) {
	currentUser, ok := session.Values["user"].(*models.UserRow)

	var isCurator bool

	if !ok {
		currentUser = nil
		isCurator = false
	} else {
		isCurator = currentUser.PermLevel <= models.Curator
	}

	return currentUser, isCurator
}
