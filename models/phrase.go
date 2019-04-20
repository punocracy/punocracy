// Data structures for our MongoDB data.

package models

import (
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

// Query for phrases from a list of words
func GetPhraseList(wordlist []Word) []Phrase {
	// Query for the word
}
