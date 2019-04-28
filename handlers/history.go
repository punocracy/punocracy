package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/alvarosness/punocracy/libhttp"
	"github.com/alvarosness/punocracy/models"
	"github.com/go-playground/form"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type historyPageData struct {
	CurrentUser      *models.UserRow
	IsCurator        bool
	RatedPhrases     []ratedPhraseDisplay
	SubmittedPhrases []string
}

type submittedPhraseDisplay struct {
}

type ratedPhraseDisplay struct {
	PhraseID            string
	PhraseText          string
	TimeSinceSubmission string
	IsOneStar           bool
	IsTwoStar           bool
	IsThreeStar         bool
	IsFourStar          bool
	IsFiveStar          bool
}

// GetHistory generates a page showing the users' history of phrase ratings and phrase submissions
func GetHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")

	currentUser, isCurator := getUser(session)

	// Getting submitted phrases
	mongdb := r.Context().Value("mongodb").(*mongo.Database)
	phrasesCollection := models.NewPhraseConnection(mongdb)
	ratingsCollection := models.NewUserRatingsConnection(mongdb)

	ratings, _ := models.GetRatingsByUserID(*currentUser, ratingsCollection)

	ratedPhrases := []ratedPhraseDisplay{}

	for _, rating := range ratings {
		now := time.Now()
		timeSinceRating := now.Sub(rating.RateDate)
		phrase, _ := models.GetPhraseByID(rating.PhraseID, phrasesCollection)

		ratedPhrases = append(ratedPhrases, ratedPhraseDisplay{
			PhraseID:            rating.PhraseID.Hex(),
			PhraseText:          phrase.PhraseText,
			TimeSinceSubmission: timeSinceRating.String(),
			IsOneStar:           rating.RatingValue == 1,
			IsTwoStar:           rating.RatingValue == 2,
			IsThreeStar:         rating.RatingValue == 3,
			IsFourStar:          rating.RatingValue == 4,
			IsFiveStar:          rating.RatingValue == 5,
		})
	}

	phrases, err := models.GetPhraseHistory(*currentUser, phrasesCollection)
	if err != nil {
		logrus.Error(err.Error())
	}

	submittedPhrases := []string{}

	for _, phrase := range phrases {
		submittedPhrases = append(submittedPhrases, phrase.PhraseText)
	}

	pageData := historyPageData{CurrentUser: currentUser, IsCurator: isCurator, RatedPhrases: ratedPhrases, SubmittedPhrases: submittedPhrases}

	tmpl, err := template.ParseFiles("templates/dashboard-nosearch.html.tmpl", "templates/history.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, pageData)
}

// PostHistory handles the update of user ratings for phrases
func PostHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "punocracy-session")

	currentUser, _ := getUser(session)

	r.ParseForm()

	decoder := form.NewDecoder()

	var ratings phraseRatings

	decoder.Decode(&ratings, r.Form)

	mongdb := r.Context().Value("mongodb").(*mongo.Database)
	phrasesCollection := models.NewPhraseConnection(mongdb)
	ratingsCollection := models.NewUserRatingsConnection(mongdb)

	for k, v := range ratings.Ratings {
		phrID, _ := primitive.ObjectIDFromHex(k)
		rating, _ := strconv.Atoi(v)

		phr, _ := models.GetPhraseByID(phrID, phrasesCollection)
		models.AddOrChangeRating(*currentUser, rating, phr, phrasesCollection, ratingsCollection)
	}

	http.Redirect(w, r, "/history", 302)
}
