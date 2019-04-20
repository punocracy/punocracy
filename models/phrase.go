// Data structures for our MongoDB data.

package models

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

// Rating maps the number of ratings of each star type. Allows computation of average rating
type Rating struct {
	OneStar   int `bson:"one"`
	TwoStar   int `bson:"two"`
	ThreeStar int `bson:"three"`
	FourStar  int `bson:"four"`
	FiveStar  int `bson:"five"`
}

// Phrase data structure format in MongoDB according to diagram
type Phrase struct {
	PhraseID        primitive.ObjectID `bson:"_id"`
	SubmitterUserID int64              `bson:"submitterUserID"`
	SubmissionDate  time.Time          `bson:"submissionDate"`
	Ratings         Rating             `bson:"ratings"`
	WordList        []int              `bson:"wordList"`
	ApprovedBy      int64              `bson:"approvedBy"`
	ApprovalDate    time.Time          `bson:"approvalDate"`
	PhraseText      string             `bson:"phraseText"`
}

// Create a new instance of the phrase collection
func NewPhraseConnection(db *mongo.Database) *mongo.Collection {
	return db.Collection("phrases")
}

// Create a new instance of the phrase collection
func NewInReviewConnection(db *mongo.Database) *mongo.Collection {
	return db.Collection("inReview")
}

func fakeGetWordIDList(words []string, db *sqlx.DB) ([]int, error) {
	// Build query string
	queryString := "SELECT wordID FROM Words_T WHERE word IN ("
	for i := 0; i < len(words)-1; i++ {
		queryString += "'" + words[i] + "',"
	}
	queryString += "'" + words[len(words)-1] + "')"

	// Execute the query on the database
	var wordIDs []int
	err := db.Select(&wordIDs, queryString)
	if err != nil {
		return []int{}, err
	}

	// Query the database
	return wordIDs, nil
}

// Insert a candidate phrase submitted by a user
func InsertCandidatePhrase(phraseText string, creator UserRow, sqlDB *sqlx.DB, inReviewCollection *mongo.Collection) error {
	// Split into lowercase words by space character
	allWords := strings.Split(strings.ToLower(phraseText), " ")

	// Get all unique words
	wordMap := make(map[string]bool)
	for _, word := range allWords {
		if _, ok := wordMap[word]; !ok {
			wordMap[word] = true
		}
	}

	var uniqueWords []string
	for word := range wordMap {
		uniqueWords = append(uniqueWords, word)
	}

	// Query the database to check if any of the words are homophones
	// TODO: replace with safe function
	wordIDs, err := fakeGetWordIDList(uniqueWords, sqlDB)
	if err != nil {
		return err
	}

	// Check if the list is empty and return error
	if len(wordIDs) == 0 {
		return errors.New("Error: no homophones in candidate phrase.")
	}

	// Create the full record
	candPhrase := Phrase{
		PhraseID:        primitive.NewObjectID(),
		SubmitterUserID: creator.ID,
		SubmissionDate:  time.Now(),
		Ratings:         Rating{},
		WordList:        wordIDs,
		ApprovedBy:      0,
		ApprovalDate:    nil,
		PhraseText:      phraseText,
	}

	// Insert into collection
	_, err = inReviewCollection.InsertOne(context.Background(), candPhrase)
	if err != nil {
		return err
	}

	// Insert the record
	return nil
}

// Insert a phrase into the phrases from candidate phrase
func InsertPhrase(phrase Phrase, approver UserRow, phrasesCollection *mongo.Collection) error {
	// Set approver
	phrase.ApprovedBy = approver.ID
	phrase.ApprovalDate = time.Now()

	// Insert into phrases collection and propagate error
	// TODO: check first return value???
	_, err := phrasesCollection.InsertOne(context.Background(), phrase)
	if err != nil {
		return err
	}

	return nil
}

// TODO list:
//  - Get phrases for curators, take in max number of phrases
//  - Get phrases for display from the homophone list, ranked by rating, take in max number of phrases

// Query for phrases from a list of words
//func GetPhraseList(wordlist []Word, phrasesCollection *mongo.Collection) ([]Phrase, error) {
//	// Build the query document
//	var queryDocument bson.D
//
//	// Get a cursor pointing to the list of phrases as a result of the query
//	cur, err := phrasesCollection.Find(context.Background(), queryDocument)
//	if err != nil {
//		return nil, err
//	}
//	defer cur.Close(context.background())
//
//	// list of phrases
//	var phraseList []Phrase
//
//	// get query result and print
//	for cur.Next(context.background()) {
//		// decode into struct
//		var onePhrase Phrase
//		//var onePhrase bson.d
//		err = cur.decode(&onePhrase)
//		if err != nil {
//			return nil, err
//		}
//		// append to phraseList
//		phraseList = append(phraseList, onePhrase)
//	}
//
//	// check for cursor errors
//	if err := cur.err(); err != nil {
//		return nil, err
//	}
//
//	// return the result
//	return phraseList, nil
//}
