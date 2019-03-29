package model

import "time"

// Rating is the number of ratings per phrase
type Rating struct {
	OneStar   int
	TwoStar   int
	ThreeStar int
	FourStar  int
	FiveStar  int
}

// Phrase is another core of our project
type Phrase struct {
	PhraseID        string    `bson:"_id"`
	SubmitterUserID int       `bson:"submitterUserID"`
	SubmissionDate  time.Time `bson:"submissionDate"`
	Ratings         Rating    `bson:"ratings"`
	WordList        []int     `bson:"wordList"`
	ApprovedBy      int       `bson:"approvedBy"`
	ApprovalDate    time.Time `bson:"approvalDate"`
	PhraseText      string    `bson:"phraseText"`
}
