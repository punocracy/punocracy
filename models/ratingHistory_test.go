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
	_, err := phrasesCollection.DeleteOne(context.Background(), bson.M{"phraseID": p.PhraseID})
	return err
}

// Test checkIfPhraseExists function
func TestCheckIfPhraseExists(t *testing.T) {
	// Connect to MongoDB and get phrases collection
	mongoDB, err := connectToMongo("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}
	phrasesCollection := NewPhraseConnection(mongoDB)

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
		t.Error("Failed test. Expected phrase to be present in collection.")
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
	phrasesCollection := NewPhraseConnection(mongoDB)

	testUser := newTestUser()
	testPhrase := newTestPhrase(testUser)

	// Insert into the phrases collection. AddRating should work
	_, err = phrasesCollection.InsertOne(context.Background(), testPhrase)
	if err != nil {
		t.Fatal(err)
	}

	// Delete phrase
	err = deletePhraseFromPhrases(testPhrase, phrasesCollection)
	if err != nil {
		t.Fatal(err)
	}
}

// TestAddRating tests the AddRating function
func TestAddOrChangeRating(t *testing.T) {
	// Connect to MongoDB with default URL string
	mongoDB, err := connectToMongo("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	// Get userRatings collection and phrases collection
	userRatings := NewUserRatingsConnection(mongoDB)
	phrasesCollection := NewPhraseConnection(mongoDB)

	// Get test user and phrase
	testUser := newTestUser()
	testPhrase := newTestPhrase(testUser)

	// Add to user's history: should fail
	err = AddOrChangeRating(testUser, 5, testPhrase, phrasesCollection, userRatings)
	if err == nil {
		t.Error("Should not update phrase not in the phrases collection!")
	} else if err != ErrPhraseNotFound {
		t.Fatal(err)
	}

	// Insert into the phrases collection. AddRating should work
	_, err = phrasesCollection.InsertOne(context.Background(), testPhrase)
	if err != nil {
		t.Fatal(err)
	}

	// Add to the user's history. Should succeed
	err = AddOrChangeRating(testUser, 5, testPhrase, phrasesCollection, userRatings)
	if err != nil {
		t.Fatal(err)
	}

	// Check for the rating.
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

	// TODO: check if the ratings were updated
	// TODO: delete stuff
	_, err = userRatings.DeleteOne(context.Background(), bson.M{"userID": testUser.ID})
	if err != nil {
		t.Fatal(err)
	}
	_, err = phrasesCollection.DeleteOne(context.Background(), bson.M{"_id": testPhrase.PhraseID})
	if err != nil {
		t.Fatal(err)
	}

	// Add second rating
	//testRating.PhraseID = primitive.NewObjectID()
	//testRating.SubmissionDate = time.Now()
}

// TODO: write TestChangeRating function
// TODO: write TestDeleteRating function
// TODO: write TestGetRatingsByUserID function (sorted by date)
// TODO: write TestDeleteUserRatings function
