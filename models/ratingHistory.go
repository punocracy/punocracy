// User rating history support functions

package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var ErrPhraseNotFound = errors.New("models: no phrases with that PhraseID found in phrases collection")
var ErrNegativeRatings = errors.New("models: negative rating value.")
var ErrInvalidRating = errors.New("models: invalid rating value.")

// UserHistory stores the user's rating history
type UserHistory struct {
	UserID        int64        `bson:"userID"`
	RatingHistory []UserRating `bson:"ratingHistory"`
}

// UserRating stores a single User Rating; a log of the phrase, rating, and date it was rated
type UserRating struct {
	PhraseID    primitive.ObjectID `bson:"phraseID"`
	RatingValue int                `bson:"ratingValue"`
	RateDate    time.Time          `bson:"rateDate"`
}

// String implements the Stringer interface
func (u UserRating) String() string {
	formatString := `{
	phraseID: ObjectID("%v")
	rateDate: %v
}`
	return fmt.Sprintf(formatString, u.PhraseID, u.RatingValue, u.RateDate)
}

// NewUserRatingsConnection creats a reference to the userRatings collection from DB pointer
func NewUserRatingsConnection(db *mongo.Database) *mongo.Collection {
	return db.Collection("userRatings")
}

// Get a phrase by ID
func getPhraseByID(phraseID primitive.ObjectID, phrasesCollection *mongo.Collection) (Phrase, error) {
	var returnPhrase Phrase
	err := phrasesCollection.FindOne(context.Background(), bson.M{"phraseID": phraseID}).Decode(&returnPhrase)
	return returnPhrase, err
}

// Check if a phrase exists in the phrases collection
func checkIfPhraseExists(p Phrase, phrasesCollection *mongo.Collection) (bool, error) {
	// Check if the phrase exists in the phrases collection
	_, err := getPhraseByID(p.PhraseID, phrasesCollection)
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Removes a rating from a phrase in
func removeRatingFromPhrase(p Phrase, r int, phrasesCollection *mongo.Collection) error {
	ratingField := "ratings."
	// Check which rating it is and set
	switch r {
	case 1:
		if p.PhraseRatings.OneStar-1 < 0 {
			return ErrNegativeRatings
		}
		ratingField += "one"
	case 2:
		if p.PhraseRatings.TwoStar-1 < 0 {
			return ErrNegativeRatings
		}
		ratingField += "two"
	case 3:
		if p.PhraseRatings.ThreeStar-1 < 0 {
			return ErrNegativeRatings
		}
		ratingField += "three"
	case 4:
		if p.PhraseRatings.FourStar-1 < 0 {
			return ErrNegativeRatings
		}
		ratingField += "four"
	case 5:
		if p.PhraseRatings.FiveStar-1 < 0 {
			return ErrNegativeRatings
		}
		ratingField += "five"
	default:
		return ErrInvalidRating
	}

	// Update the rating in the database (decrement)
	phraseFilterDoc := bson.M{"_id": p.PhraseID}
	phraseUpdateDoc := bson.M{"$inc": bson.M{ratingField: -1}}
	_, err := phrasesCollection.UpdateOne(context.Background(), phraseFilterDoc, phraseUpdateDoc)
	return err
}

// Add a rating to the phrase
func addRatingToPhrase(p Phrase, rating int, phrasesCollection *mongo.Collection) error {
	// Update the phrase to include the rating
	phraseFilterDoc := bson.M{"_id": p.PhraseID}
	phraseUpdateDoc := bson.M{"$inc": bson.M{"ratings." + ratingToRatingString(rating): 1}}
	_, err := phrasesCollection.UpdateOne(context.Background(), phraseFilterDoc, phraseUpdateDoc)
	return err
}

// AddRating adds a rating for a specific user
func AddOrChangeRating(user UserRow, rating int, ratedPhrase Phrase, phrasesCollection *mongo.Collection, userRatings *mongo.Collection) error {
	// Ensure the phrase exists in th ecollection
	phraseExists, err := checkIfPhraseExists(ratedPhrase, phrasesCollection)
	if err != nil {
		return err
	} else if !phraseExists {
		return ErrPhraseNotFound
	}

	// Build UserRating document
	ratingEntry := UserRating{PhraseID: ratedPhrase.PhraseID, RatingValue: rating, RateDate: time.Now()}

	// Check if user has a rating history entry. If not, create user rating history
	var userHist UserHistory
	err = userRatings.FindOne(context.Background(), bson.M{"userID": user.ID}).Decode(&userHist)

	// If user has history, check if the rating history
	if err == nil {
		// Check if the user has rated this document before
		newRatingFlag := true
		var oldRating UserRating
		for _, r := range userHist.RatingHistory {
			if r.PhraseID == ratedPhrase.PhraseID {
				newRatingFlag = false
				oldRating = r
				break
			}
		}

		// Add to set if it's a new rating or update the rating if it's already been rated
		if newRatingFlag {
			// Add document to the user's rating history
			filterDoc := bson.M{"userID": user.ID}
			updateDoc := bson.M{"$addToSet": bson.M{"ratingHistory": ratingEntry}}
			_, err = userRatings.UpdateOne(context.Background(), filterDoc, updateDoc)
			if err != nil {
				return err
			}
		} else {
			// Change entry in user's rating history
			//filterDoc := bson.M{"userID": user.ID, "ratingHistory.phraseID": ratedPhrase.PhraseID}

			// Decrement the rating in the thingy
			removeRatingFromPhrase(ratedPhrase, oldRating.RatingValue, phrasesCollection)
		}
	} else if err == mongo.ErrNoDocuments {
		userHist.UserID = user.ID
		userHist.RatingHistory = []UserRating{ratingEntry}
		_, err := userRatings.InsertOne(context.Background(), userHist)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	// Update the phrase to include the rating
	addRatingToPhrase(ratedPhrase, rating, phrasesCollection)

	return nil
}

// TODO: write ChangeRating function
func ChangeRating(user UserRow, rating int, ratedPhrase Phrase, userRatings *mongo.Collection) error {
	return nil
}

// TODO: write DeleteRating function
func DeleteRating(user UserRow, rating int, ratedPhrase Phrase, userRatings *mongo.Collection) error {
	return nil
}

// TODO: write GetRatingsByUserID function (sorted by date)
func GetRatingsByUserID(user UserRow, userRatings *mongo.Collection) ([]UserRating, error) {
	return nil, nil
}

// TODO: write DeleteUserRatings function
// TODO: write updateRatingByUser to update the rating in the phrases collection
// NOTE: Make everything propagate to the phrases table

// Convert an integer rating to its corresponding document string
func ratingToRatingString(r int) string {
	var s string
	switch r {
	case 1:
		s = "one"
	case 2:
		s = "two"
	case 3:
		s = "three"
	case 4:
		s = "four"
	case 5:
		s = "five"
	default:
		s = ""
	}
	return s
}
