// Data structures for our MongoDB data.

package models

import (
	"errors"
	_ "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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
	PhraseID        bson.ObjectId `bson:"_id"`
	SubmitterUserID int           `bson:"submitterUserID"`
	SubmissionDate  time.Time     `bson:"submissionDate"`
	Ratings         Rating        `bson:"ratings"`
	WordList        []int         `bson:"wordList"`
	ApprovedBy      int           `bson:"approvedBy"`
	ApprovalDate    time.Time     `bson:"approvalDate"`
	PhraseText      string        `bson:"phraseText"`
}

// Create a new instance of the phrase collection
func NewPhraseConnection(db *mongo.Database) *mongo.Collection {
	return db.Collection("phrases")
}

// Create a new instance of the phrase collection
func NewInReviewConnection(db *mongo.Database) *mongo.Collection {
	return db.Collection("inReview")
}

// Insert a candidate phrase submitted by a user
func InsertCandidatePhrase(phraseText string, creator UserRow, inReviewCollection *mongo.Collection) error {

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

// Query for phrases from a list of words
func GetPhraseList(wordlist []Word, phrasesCollection *mongo.Collection) ([]Phrase, error) {
	// Build the query document
	var queryDocument bson.D

	// Get a cursor pointing to the list of phrases as a result of the query
	cur, err := phrasesCollection.Find(context.Background(), queryDocument)
	if err != nil {
		log.fatal(err)
	}
	defer cur.close(context.background())

	// list of phrases
	var phraseList []Phrase

	// get query result and print
	for cur.next(context.background()) {
		// decode into struct
		var onePhrase Phrase
		//var onePhrase bson.d
		err = cur.decode(&onePhrase)
		if err != nil {
			return nil, err
		}
		// append to phraseList
		phraseList = append(phraseList, onePhrase)
	}

	// check for cursor errors
	if err := cur.err(); err != nil {
		return nil, err
	}

	// return the result
	return phraseList, nil
}
