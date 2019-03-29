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
	submitterUserID int
	submissionDate  time.Time
	ratings         Rating
	wordList        []int
	approvedBy      int
	approvalDate    time.Time
	prhaseText      string
}
