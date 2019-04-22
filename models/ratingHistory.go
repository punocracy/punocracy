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

// Adds a rating for a specific user
func AddRating(user UserRow, ratedPhrase Phrase, ratingHistory *mongo.Collection) error {
	//
}

// TODO: write ChangeRating function
// TODO: write DeleteRating function
// TODO: write GetRatingsByUserID function (sorted by date)
// TODO: write DeleteUserRatings function
// NOTE: Make everything propagate to the phrases table