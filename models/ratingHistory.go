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

// UserRating stores a single User Rating; a log of the phrase, rating, and date it was rated
type UserRating struct {
	ratingID    primitive.ObjectID `bson:"_id"`
	UserID      int64              `bson:"userID"`
	PhraseID    primitive.ObjectID `bson:"phraseID"`
	RatingValue int                `bson:"ratingValue"`
	RateDate    time.Time          `bson:"rateDate"`
}

// String implements the Stringer interface
func (u UserRating) String() string {
	formatString := `{
	"userID": ObjectID("%v"),
	"phraseID": ObjectID("%v"),
	"ratingValue": %v,
	"rateDate": %v
}`
	return fmt.Sprintf(formatString, u.UserID, u.PhraseID, u.RatingValue, u.RateDate)
}

// NewUserRatingsConnection creats a reference to the userRatings collection from DB pointer
func NewUserRatingsConnection(db *mongo.Database) *mongo.Collection {
	return db.Collection("userRatings")
}

// AddOrChangeRating adds or modifies a rating value given a user, phrase, and rating value
func AddOrChangeRating(user UserRow, rating int, thePhrase Phrase, phrases *mongo.Collection, userRatings *mongo.Collection) error {
	// Attempt to change the rating, and if an ErrNoDocuments error is encountered, add it afresh
	err := changeRating(user, rating, thePhrase, phrases, userRatings)
	if err == mongo.ErrNoDocuments {
		return addRating(user, rating, thePhrase, phrases, userRatings)
	}
	return err
}

// addRating adds a rating given a Phrase, rating value, UserRow, and the collection pointers. It places the rating in the userRatings collection without checking if one exists
// Returns ErrPhraseNotFound if the phrase does not exist in the phrases collection
func addRating(user UserRow, rating int, thePhrase Phrase, phrases *mongo.Collection, userRatings *mongo.Collection) error {
	// Check if the phrase exists and raise error if not
	ok, err := checkIfPhraseExists(thePhrase, phrases)
	if err != nil {
		return err
	} else if !ok {
		return ErrPhraseNotFound
	}

	// Construct the rating struct
	ratingEntry := UserRating{
		ratingID:    primitive.NewObjectID(),
		UserID:      user.ID,
		PhraseID:    thePhrase.PhraseID,
		RatingValue: rating,
		RateDate:    time.Now(),
	}

	// Insert the rating into the userRatings collection
	_, err = phrases.InsertOne(context.Background(), ratingEntry)
	if err != nil {
		return err
	}

	// Add the rating value to the phrase
	err = addRatingToPhrase(thePhrase, rating, phrases)
	return err
}

// changeRating changes the rating for a user and phrase pair given a new rating value. Assumes it exists in the userRatings collection
func changeRating(user UserRow, rating int, thePhrase Phrase, phrases *mongo.Collection, userRatings *mongo.Collection) error {
	// Get the old rating value
	oldRating, err := getRating(user, thePhrase, userRatings)
	if err != nil {
		return err
	}

	// Update the userRatings entry
	filterDoc := bson.M{"userID": user.ID, "phraseID": thePhrase.PhraseID}
	updateDoc := bson.M{"$set": bson.M{"ratingValue": rating, "rateDate": time.Now()}}
	_, err = userRatings.UpdateOne(context.Background(), filterDoc, updateDoc)
	if err != nil {
		return err
	}

	// Change rating for phrase and return the error
	return changeRatingForPhrase(thePhrase, oldRating.RatingValue, rating, phrases)
}

// getRating retrieves a rating given a user and phrase
func getRating(user UserRow, thePhrase Phrase, userRatings *mongo.Collection) (UserRating, error) {
	var theRating UserRating
	err := userRatings.FindOne(context.Background(), bson.M{"userID": user.ID, "phraseID": thePhrase.PhraseID}).Decode(&theRating)
	return theRating, err
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

// changeRatingForPhrases changes a rating in the phrases collection.
func changeRatingForPhrase(thePhrase Phrase, oldRating int, newRating int, phrases *mongo.Collection) error {
	// Update in the phrases collection
	filterDoc := bson.M{"_id": thePhrase.PhraseID}
	updateDoc := bson.M{"$inc": bson.M{"ratings." + ratingToRatingString(oldRating): -1, "ratings." + ratingToRatingString(newRating): 1}}
	_, err := phrases.UpdateOne(context.Background(), filterDoc, updateDoc)
	return err
}

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
