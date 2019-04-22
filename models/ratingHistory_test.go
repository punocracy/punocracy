package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

// TODO: write TestAddRating function
// TestAddRating tests the AddRating function
func TestAddRating(t *testing.T) {
	// Connect to MongoDB with default URL string
	mongoDB, err := connectToMongo("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	// Get userRatings collection
	userRatings := NewUserRatingsConnection(mongoDB)

	// Get test user
	testUser := newTestUser()

	// Create the phrase
	testPhrase := Phrase{
		PhraseID:        primitive.NewObjectID(),
		SubmitterUserID: testUser.ID,
		SubmissionDate:  time.Now(),
		PhraseRatings:   Rating{},
		WordList:        []int{1454, 518, 588, 189, 71},
		ReviewedBy:      0,
		ReviewDate:      time.Now(),
		PhraseText:      "All your base are belong to us.",
		DisplayPublic:   Unreviewed,
	}

	// Add to user's history
	err = AddRating(testUser, 5, testPhrase, userRatings)
	if err != nil {
		t.Fatal(err)
	}

	// Check for the rating
	var testUserHist UserHistory
	err = userRatings.FindOne(context.Background(), bson.M{"userID": testUser.ID}).Decode(&testUserHist)
	if err != nil {
		t.Fatal(err)
	}

	// Look for the rating in user's history
	found := false
	for _, r := range testUserHist.RatingHistory {
		if r.PhraseID == testPhrase.PhraseID {
			found = true
			break
		}
	}
	if !found {
		t.Error("Could not find testUser rating. testUserHist:", testUserHist)
	}

	// Add second rating
	//testRating.PhraseID = primitive.NewObjectID()
	//testRating.SubmissionDate = time.Now()

}

// TODO: write TestChangeRating function
// TODO: write TestDeleteRating function
// TODO: write TestGetRatingsByUserID function (sorted by date)
// TODO: write TestDeleteUserRatings function
