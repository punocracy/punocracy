package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

func newTestPhrase(testUser UserRow) Phrase {
	// Create the phrase
	return Phrase{
		PhraseID:        primitive.NewObjectID(),
		SubmitterUserID: testUser.ID,
		SubmissionDate:  time.Now(),
		PhraseRatings:   Rating{},
		WordList:        []int{1454, 518, 588, 189, 71},
		ReviewedBy:      0,
		ReviewDate:      time.Now(),
		PhraseText:      "All your base are belong to us.",
		DisplayPublic:   Accepted,
	}
}

// Delete a phrase from the phrases collection
func deletePhraseFromPhrases(p Phrase, phrasesCollection *mongo.Collection) error {
	_, err := phrasesCollection.DeleteOne(context.Background(), bson.M{"_id": p.PhraseID})
	return err
}

// Return phrases_test collection connection
func newTestPhraseConnection(db *mongo.Database) *mongo.Collection {
	return db.Collection("phrases_test")
}

// Return userRatings_test collection connection
func newTestUserRatingsConnection(db *mongo.Database) *mongo.Collection {
	return db.Collection("userRatings_test")
}

// Test checkIfPhraseExists function
func TestCheckIfPhraseExists(t *testing.T) {
	// Connect to MongoDB and get phrases collection
	mongoDB, err := connectToMongo("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}
	phrasesCollection := newTestPhraseConnection(mongoDB)

	// Test user and phrase
	testUser := newTestUser()
	testPhrase := newTestPhrase(testUser)

	// Check if it exists. Should be false
	result, err := checkIfPhraseExists(testPhrase, phrasesCollection)
	if err != nil {
		t.Fatal(err)
	}
	if result == true {
		t.Error("Failed test. Expected phrase not to be in collection.")
	}

	// Insert into the phrases collection. AddRating should work
	_, err = phrasesCollection.InsertOne(context.Background(), testPhrase)
	if err != nil {
		t.Fatal(err)
	}

	// Check if it exists. Should be true
	result, err = checkIfPhraseExists(testPhrase, phrasesCollection)
	if err != nil {
		t.Fatal(err)
	}
	if result == false {
		t.Error("Result:", result, "expected: true. Expected phrase to be present in collection with ID: ", testPhrase.PhraseID)
	}

	// Delete phrase
	err = deletePhraseFromPhrases(testPhrase, phrasesCollection)
	if err != nil {
		t.Fatal(err)
	}
}

// Test addRatingToPhrase and removeRatingFromPhrase functions
func TestAddRemoveRatingToPhrase(t *testing.T) {
	// Connect to MongoDB and get phrases collection
	mongoDB, err := connectToMongo("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}
	phrasesCollection := newTestPhraseConnection(mongoDB)

	testUser := newTestUser()
	testPhrase := newTestPhrase(testUser)

	// Insert into the phrases collection.
	_, err = phrasesCollection.InsertOne(context.Background(), testPhrase)
	if err != nil {
		t.Fatal(err)
	}

	// Add rating to the phrase
	err = addRatingToPhrase(testPhrase, 5, phrasesCollection)
	if err != nil {
		t.Fatal(err)
	}

	// Check the phrase
	testPhrase, err = getPhraseByID(testPhrase.PhraseID, phrasesCollection)
	if err != nil {
		t.Fatal(err)
	}
	if testPhrase.PhraseRatings.FiveStar != 1 {
		t.Error("Phrase rating not changed when added!")
	}

	// Remove the rating and test if ErrNegativeRatings works
	err = removeRatingFromPhrase(testPhrase, 4, phrasesCollection)
	if err == nil {
		t.Error("Rating remove successfully when it wasn't. Shouldn't be any ratings")
	} else if err != ErrNegativeRatings {
		t.Fatal(err)
	}

	// Remove the 5 star rating
	err = removeRatingFromPhrase(testPhrase, 5, phrasesCollection)
	if err == ErrNegativeRatings {
		t.Error("Rating was not successfully removed when it should have been.")
	} else if err != nil {
		t.Fatal(err)
	}

	// Delete phrase
	err = deletePhraseFromPhrases(testPhrase, phrasesCollection)
	if err != nil {
		t.Fatal(err)
	}
}

// TestAddRating tests the AddRating function
//func TestAddOrChangeRating(t *testing.T) {
//	// Connect to MongoDB with default URL string
//	mongoDB, err := connectToMongo("mongodb://localhost:27017")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Get userRatings collection and phrases collection
//	userRatings := newTestUserRatingsConnection(mongoDB)
//	phrasesCollection := newTestPhraseConnection(mongoDB)
//
//	// Get test user and phrase
//	testUser := newTestUser()
//	testPhrase := newTestPhrase(testUser)
//
//}

// TODO: write TestChangeRating function
// TODO: write TestDeleteRating function
// TODO: write TestGetRatingsByUserID function (sorted by date)
// TODO: write TestDeleteUserRatings function
