// User rating history support functions

package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var ErrPhraseNotFound = errors.New("models: no phrases with that PhraseID found in phrases collection")

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

// AddRating adds a rating for a specific user
func AddRating(user UserRow, rating int, ratedPhrase Phrase, userRatings *mongo.Collection) error {
	// Build UserRating document
	ratingEntry := UserRating{PhraseID: ratedPhrase.PhraseID, RatingValue: rating, RateDate: time.Now()}

	// Check if user has a rating history entry. If not, create user rating history
	var userHist UserHistory
	err := userRatings.FindOne(context.Background(), bson.M{"userID": user.ID}).Decode(&userHist)

	if err == mongo.ErrNoDocuments {
		userHist.UserID = user.ID
		userHist.RatingHistory = []UserRating{ratingEntry}
		_, err := userRatings.InsertOne(context.Background(), userHist)
		if err != nil {
			return err
		}
	} else if err == nil {
		// Add document to the user's rating history
		filterDoc := bson.M{"userID": user.ID}
		updateDoc := bson.M{"$addToSet": bson.M{"ratingHistory": ratingEntry}}
		_, err = userRatings.UpdateOne(context.Background(), filterDoc, updateDoc)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	// Add the user rating to the array for this user's rating history
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
