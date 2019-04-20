package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

// Get the testing user
func newTestUser() UserRow {
	// Example UserRow
	testUser := UserRow{
		ID:        2,
		Username:  "testerUser",
		Email:     "test@testerson.com",
		Password:  "asdf",
		PermLevel: 0,
	}
	return testUser
}

// Connect to MongoDB instance
func connectToMongo(urlString string) (*mongo.Database, error) {
	// Connect to localhost
	client, err := mongo.NewClient(options.Client().ApplyURI(urlString))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Connect
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	// Check connection with ping
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	return client.Database("punocracy"), nil
}

// Test phrase insertion function
//func TestInsertCandidatePhrase(t *testing.T) {
//	// Connect to MongoDB with default URL string
//	db, err := connectToMongo("mongodb://localhost:27017")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Get the phrases collection from the cool_songs database
//	phrasesCollection := NewInReviewConnection(db)
//
//	// Example UserRow
//	testUser := newTestUser()
//
//	// Test cases
//	var testPhrases = []struct {
//		input  string
//		output bool
//	}{
//		{"All your base are belong to us.", true},      // Homophones: all, are, base, to, your
//		{"This has zero homophones within it.", false}, // No homophones
//	}
//
//	// Insert each phrase
//	for _, phrase := range testPhrases {
//		// Try to insert the phrase
//		var successVal bool
//		err := InsertPhrase(phrase.input, testUser, phrasesCollection)
//		successVal = (err == nil)
//
//		// Check the value
//		if successVal != phrase.output {
//			t.Error("Phrase: ", phrase.input, " not inserted successfully.")
//		}
//
//	}
//}

// Test InsertPhrase directly
func TestInsertPhrase(t *testing.T) {
	// Connect to MongoDB with default URL string
	db, err := connectToMongo("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	// Get the phrases collection from the cool_songs database
	phrasesCollection := NewPhraseConnection(db)

	// Get test user
	testUser := newTestUser()

	// Empty rating
	var emptyRating Rating

	// Phrase example
	testPhrase := Phrase{
		PhraseID:        primitive.NewObjectID(),
		SubmitterUserID: testUser.ID,
		SubmissionDate:  time.Now(),
		Ratings:         emptyRating,
		WordList:        []int{588, 817},
		PhraseText:      "The base of the project.",
	}

	// Log PhraseID
	t.Log(testPhrase.PhraseID)

	// Insert test phrase
	err = InsertPhrase(testPhrase, testUser, phrasesCollection)
	if err != nil {
		t.Error(err)
	}

	// Try to delete the phrase
	_, err = phrasesCollection.DeleteOne(context.Background(), bson.M{"_id": testPhrase.PhraseID})
	if err != nil {
		t.Error(err)
	}

}

// Test the GetPhraseList object
//func TestGetPhraseList(t *testing.T) {
//	// Connect to MongoDB with default URL string
//	db, err := connectToMongo("mongodb://localhost:27017")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Get the phrases collection from the cool_songs database
//	phrases := NewPhraseConnection(db)
//
//	// List of words for phrase query
//	//	var wordList = []Word{
//	//		{
//
//	// Query for songs
//	songList, err := ArrayQuery(songs)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Print songlist
//	for _, s := range songList {
//		fmt.Println(s)
//	}
//}
