package alphwordlist

import (
    "log"
    "github.com/jmoiron/sqlx"
    "github.com/alvarosness/punocracy/models"
)

//This struct holds the words quered from the DB
//uses the db:tags to match in database, according to the specification in sqlx
/*
type Word struct {
    WordID int `db:"wordID"`
    Word string `db:"word"`
    HomophoneGroup int `db:"homophoneGroup"`
}
*/

//pointer to the DB object in use
var usedDBptr **sqlx.DB = nil

//A flag to test if the DB is valid
var validDB bool = false

//set the DB to be used for querying 
//input: a pointer of type DB
func SetDB( inputDBptr **sqlx.DB ) (reStatus bool){

    if inputDBptr == nil{
        validDB = false;
        return false
    }

    usedDBptr = inputDBptr
    validDB = true;
    return true
}

//Query the firstLetter for a lit of strucs
//input: a single "rune" representing the first letter
//this is case insensitive 
//output: a list of type Word
func QueryAlph( firstLetter rune ) (wordList []models.Word){
    words := []models.Word{}//this creates a Word array initiaizing it to empty it seems

    queryString := string(firstLetter) + "%"

    erri := (*usedDBptr).Select(&words,"SELECT * FROM Words_T WHERE word LIKE ? ORDER BY word;", queryString)

    if erri != nil{
        log.Fatalln(erri)
    }

    return words
}
