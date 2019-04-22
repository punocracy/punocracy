package models

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"testing"
	"time"
)

// Connect to SQL database
func newDBConnection() (*sqlx.DB, error) {
	// DSN string
	defaultDSN := strings.Replace("nathaniel:nathaniel@tcp(localhost:3306)/punocracy?parseTime=true", "-", "_", -1)

	// Connect to DB
	db, err := sqlx.Connect("mysql", defaultDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Get the testing user
func newTestUser() UserRow {
	// Example UserRow
	testUser := UserRow{
		ID:           2,
		Username:     "testerUser",
		Email:        "test@testerson.com",
		PasswordHash: "asdf",
		PermLevel:    0,
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

// Test GetPhraseList
func TestGetPhraseList(t *testing.T) {
	// Connect to MongoDB with default URL string
	mongoDB, err := connectToMongo("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	// Get the phrases collection from the cool_songs database
	phrasesCollection := NewPhraseConnection(mongoDB)

	// List of words to search for
	wordList := []WordRow{
		{1414, "two", 625},
		{189, "to", 625},
		//{831, "too", 625},
	}

	// Get a list of phrases
	phraseList, err := GetPhraseList(wordList, phrasesCollection)
	if err != nil {
		t.Fatal(err)
	}

	// Print all phrases
	for _, p := range phraseList {
		t.Log(p.String())
	}
}

// Test GetPhraseListForCurators
func TestGetPhrasesForCurators(t *testing.T) {
	// Connect to MongoDB with default URL string
	mongoDB, err := connectToMongo("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	// Get the phrases collection from the cool_songs database
	phrasesCollection := NewPhraseConnection(mongoDB)

	// Connect to MySQL database
	sqlDB, err := newDBConnection()
	if err != nil {
		t.Fatal(err)
	}
	wordInstance := NewWord(sqlDB)

	// Example UserRow
	testUser := newTestUser()

	// Test cases
	var testPhrases = []string{
		"All your base are belong to us.",
		"To live is to dream.",
		"Live free or die hard.",
		"This has no homophones in it.",
	}
	maxPhrases := 3

	// Insert each phrase
	for _, phrase := range testPhrases {
		// Try to insert the phrase
		err := InsertPhrase(phrase, testUser, wordInstance, phrasesCollection)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Get phrases for curator list
	phrases, err := GetPhraseListForCurators(int64(maxPhrases), phrasesCollection)
	if err != nil {
		t.Fatal(err)
	}

	// Check length
	if len(phrases) != maxPhrases {
		t.Error("Got too many phrases! Expected", maxPhrases, "got", len(phrases))
	}

	// Check the fields for all and print the phrases
	for _, p := range phrases {
		var result Phrase
		err = phrasesCollection.FindOne(context.Background(), bson.M{"_id": p.PhraseID}).Decode(&result)
		if err != nil {
			t.Fatal(err)
		}

		if result.DisplayPublic != InReview {
			t.Error("Display value is not InReveiw! Expected", InReview, "got value:", result.DisplayPublic)
		}
		t.Log("PhraseText:", p.PhraseText)
	}

	// Try to delete the phrases
	for _, phrase := range testPhrases {
		_, err = phrasesCollection.DeleteOne(context.Background(), bson.M{"phraseText": phrase})
		if err != nil {
			t.Error(err)
		}
	}
}

// Test average rating function
func TestAverageRating(t *testing.T) {
	// List of ratings
	tests := []struct {
		input    Rating
		expected float32
	}{
		{Rating{0, 0, 0, 4, 0}, float32(4)},
		{Rating{0, 0, 0, 0, 0}, float32(0)},
		{Rating{0, 0, 0, 2, 2}, 4.5},
		{Rating{1, 2, 5, 2, 1}, float32(3)},
	}

	// Test
	for _, test := range tests {
		if output := AverageRating(test.input); output != test.expected {
			t.Error("Test failed: {} inputted, {} expected, {} received", test.input, test.expected, output)
		}
	}
}

// Test fake query for word IDs
func TestFakeGetWordIDList(t *testing.T) {
	// Connect to database
	db, err := newDBConnection()
	if err != nil {
		t.Fatal(err)
	}

	// Get word IDs
	wordIds, err := fakeGetWordIDList([]string{"base", "two", "mom"}, db)
	if err != nil {
		t.Fatal(err)
	}

	// Print the result
	t.Log(wordIds)
}

// Test phrase insertion function
func TestInsertPhrase(t *testing.T) {
	// Connect to MongoDB with default URL string
	mongoDB, err := connectToMongo("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	// Get the phrases collection from the cool_songs database
	phrasesCollection := NewPhraseConnection(mongoDB)

	// Connect to MySQL database
	sqlDB, err := newDBConnection()
	if err != nil {
		t.Fatal(err)
	}

	// Create instance of Word type
	wordInstance := NewWord(sqlDB)

	// Example UserRow
	testUser := newTestUser()

	// Test cases
	var testPhrases = []struct {
		input  string
		output bool
	}{
		{"All your base are belong to us.", true}, // Homophones: all, are, base, to, your
		{"To live is to dream.", true},
		{"Live free or die hard.", true},
		{"This has zero homophones within it.", false}, // No homophones
	}

	// Insert each phrase
	for _, phrase := range testPhrases {
		// Try to insert the phrase
		var successVal bool
		err := InsertPhrase(phrase.input, testUser, wordInstance, phrasesCollection)
		successVal = (err == nil)

		// Check the value
		if successVal != phrase.output {
			t.Error("Phrase: ", phrase.input, " not inserted successfully.")
		}

	}

	// Try to delete the phrases
	for _, phrase := range testPhrases {
		_, err = phrasesCollection.DeleteOne(context.Background(), bson.M{"phraseText": phrase.input})
		if err != nil {
			t.Error(err)
		}
	}
}

// Test InsertPhrase directly
func TestAcceptRejectPhrase(t *testing.T) {
	// Connect to MongoDB with default URL string
	mongoDB, err := connectToMongo("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	// Get the phrases collection from the cool_songs database
	phrasesCollection := NewPhraseConnection(mongoDB)

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

	// Insert into collection
	_, err = phrasesCollection.InsertOne(context.Background(), testPhrase)
	if err != nil {
		t.Fatal(err)
	}

	// Test accept
	err = AcceptPhrase(testPhrase, testUser, phrasesCollection)
	if err != nil {
		t.Fatal(err)
	}

	// Find the phrase by ID and see if it's accepted
	var queryPhrase Phrase
	err = phrasesCollection.FindOne(context.Background(), bson.M{"_id": testPhrase.PhraseID}).Decode(&queryPhrase)
	if err != nil {
		t.Fatal(err)
	}

	if queryPhrase.DisplayPublic != Accepted {
		t.Error("Phrase was not accepted! PhraseID: ", queryPhrase.PhraseID)
	}

	// Set phrase as rejected
	err = RejectPhrase(testPhrase, testUser, phrasesCollection)
	if err != nil {
		t.Fatal(err)
	}

	// Check that out
	err = phrasesCollection.FindOne(context.Background(), bson.M{"_id": testPhrase.PhraseID}).Decode(&queryPhrase)
	if err != nil {
		t.Fatal(err)
	}

	if queryPhrase.DisplayPublic != Rejected {
		t.Error("Phrase was not rejected! PhraseID: ", queryPhrase.PhraseID)
	}

	// Try to delete the phrase
	_, err = phrasesCollection.DeleteOne(context.Background(), bson.M{"_id": testPhrase.PhraseID})
	if err != nil {
		t.Fatal(err)
	}

}
