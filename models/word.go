package models

// Word is the core of our project
type Word struct {
	WordID         int    `db:"wordID"`
	Word           string `db:"word"`
	HomophoneGroup int    `db:"homophoneGroup"`
}
