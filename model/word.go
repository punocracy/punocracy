package model

// Word is the core of our project
type Word struct {
	wordID         int    `bson:"wordID"`
	word           string `bson:"word"`
	homophoneGroup int    `bson:"homophoneGroup"`
}
