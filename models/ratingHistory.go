// User rating history support functions

package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// A single User Rating; a log of the phrase, rating, and date it was rated
type UserRating struct {
	PhraseID    primitive.ObjectID `bson:"phraseID"`
	RatingValue int                `bson:"ratingValue"`
	RateDate    time.Time          `bson:"rateDate"`
}

// Pretty printing
func (u UserRating) String() string {
	formatString := `{
	phraseID: ObjectID("%v")
	rateDate: %v
}`
	return fmt.Sprintf(formatString, u.PhraseID, u.RatingValue, u.RateDate)
}

// Adds a rating for a specific user
func AddRating(user UserRow, rating int, ratedPhrase Phrase, ratingHistory *mongo.Collection) error {
	//
	return nil
}

// TODO: write ChangeRating function
func ChangeRating(user UserRow, rating int, ratedPhrase Phrase, ratingHistory *mongo.Collection) error {
	return nil
}

// TODO: write DeleteRating function
func DeleteRating(user UserRow, rating int, ratedPhrase Phrase, ratingHistory *mongo.Collection) error {
	return nil
}

// TODO: write GetRatingsByUserID function (sorted by date)
func GetRatingsByUserID(user UserRow, ratingHistory *mongo.Collection) ([]UserRating, error) {
	return nil, nil
}

// TODO: write DeleteUserRatings function
// NOTE: Make everything propagate to the phrases table
