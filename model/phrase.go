package model

import "time"

// Rating is the number of ratings per phrase
type Rating struct {
	oneStar   int
	twoStar   int
	threeStar int
	fourStar  int
	fiveStar  int
}

// Phrase is another core of our project
type Phrase struct {
	phraseID        string    `bson:"_id"`
	submitterUserID int       `bson:"submitterUserID"`
	submissionDate  time.Time `bson:"submissionDate"`
	ratings         Rating    `bson:"ratings"`
	wordList        []int     `bson:"wordList"`
	approvedBy      int       `bson:"approvedBy"`
	approvalDate    time.Time `bson:"approvalDate"`
	phraseText      string    `bson:"phraseText"`
}
