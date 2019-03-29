package model

// Word is the core of our project
type Word struct {
	WordID         int    `bson:"wordID"`
	Word           string `bson:"word"`
	HomophoneGroup int    `bson:"homophoneGroup"`
}
