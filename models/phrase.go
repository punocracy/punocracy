// Data structures for our MongoDB data.

package models

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Display value type
type DisplayValue int

// Display value constants
const (
	// Unreviewed by a curator; do not display
	Unreviewed DisplayValue = iota
	// Curator is in the process of reviewing
	InReview
	// Accepted phrase to be displayed
	Accepted
	// Phrase was rejected -- do not discard, but don't save, either
	Rejected
)

// Rating maps the number of ratings of each star type. Allows computation of average rating
type Rating struct {
	OneStar   int `bson:"one"`
	TwoStar   int `bson:"two"`
	ThreeStar int `bson:"three"`
	FourStar  int `bson:"four"`
	FiveStar  int `bson:"five"`
}

// Pretty printing string method
func (r Rating) String() string {
	formatString := `{
	one: %v,
	two: %v,
	three: %v,
	four: %v,
	five: %v
}`
	return fmt.Sprintf(formatString, r.OneStar, r.TwoStar, r.ThreeStar, r.FourStar, r.FiveStar)
}

// Phrase data structure format in MongoDB according to diagram
type Phrase struct {
	PhraseID        primitive.ObjectID `bson:"_id"`
	SubmitterUserID int64              `bson:"submitterUserID"`
	SubmissionDate  time.Time          `bson:"submissionDate"`
	PhraseRatings   Rating             `bson:"ratings"`
	WordList        []int              `bson:"wordList"`
	ReviewedBy      int64              `bson:"reviewedBy"`
	ReviewDate      time.Time          `bson:"reviewDate"`
	PhraseText      string             `bson:"phraseText"`
	DisplayPublic   DisplayValue       `bson:"displayValue"`
}

// Pretty printing like a JSON document for Phrase
func (p Phrase) String() string {
	formatString := `{
	_id: ObjectId("%v"),
	submitterUserID: "%v",
	submissionDate: %v,
	ratings: {
		one: %v,
		two: %v,
		three: %v,
		four: %v,
		five: %v
	},
	wordList: %v,
	reviewedBy: %v,
	reviewDate: %v,
	phraseText: "%v",
	displayValue: %v
}`
	return fmt.Sprintf(formatString, p.PhraseID, p.SubmitterUserID, p.SubmissionDate, p.PhraseRatings.OneStar, p.PhraseRatings.TwoStar, p.PhraseRatings.ThreeStar, p.PhraseRatings.FourStar, p.PhraseRatings.FiveStar, p.WordList, p.ReviewedBy, p.ReviewDate, p.PhraseText, p.DisplayPublic)
}

// Type for sorting phrases.
type phraseSorter struct {
	phrase    Phrase
	avgRating float32
}

// Type for sorting a list of phrases. Implements sort.Interface
type phraseSorterList []phraseSorter

// Length of the data
func (p phraseSorterList) Len() int {
	return len(p)
}

// Less than
func (p phraseSorterList) Less(i, j int) bool {
	return p[i].avgRating < p[j].avgRating
}

// Swap
func (p phraseSorterList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Create a new instance of the phrase collection
func NewPhraseConnection(db *mongo.Database) *mongo.Collection {
	return db.Collection("phrases")
}

// Get wordID list from SQL database. Unsafe!
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
func InsertPhrase(phraseText string, creator UserRow, wordInstance *Word, phrasesCollection *mongo.Collection) error {
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
	wordIDs, err := wordInstance.GetWordIDList(nil, uniqueWords)
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
		PhraseRatings:   Rating{},
		WordList:        wordIDs,
		ReviewedBy:      0,
		ReviewDate:      time.Now(),
		PhraseText:      phraseText,
		DisplayPublic:   Unreviewed,
	}

	// Insert into collection
	_, err = phrasesCollection.InsertOne(context.Background(), candPhrase)
	if err != nil {
		return err
	}

	// Insert the record
	return nil
}

// Accept a reviewed phrase
func AcceptPhrase(phraseIDString string, reviewer UserRow, phrasesCollection *mongo.Collection) error {

	phraseID, _ := primitive.ObjectIDFromHex(phraseIDString)
	// Build update document filter (by _id)
	filter := bson.M{"_id": phraseID}

	// Update document
	updateDocument := bson.M{"$set": bson.M{"reviewedBy": reviewer.ID, "reviewDate": time.Now(), "displayValue": Accepted}}

	// Update the phrase in Mongo to set it to be accepted
	_, err := phrasesCollection.UpdateOne(context.Background(), filter, updateDocument)
	if err != nil {
		return err
	}

	return nil
}

// Set the specified phrase as rejected after review
func RejectPhrase(phraseIDString string, reviewer UserRow, phrasesCollection *mongo.Collection) error {

	phraseID, _ := primitive.ObjectIDFromHex(phraseIDString)
	// Build update document filter (by _id)
	filter := bson.M{"_id": phraseID}

	// Update document
	updateDocument := bson.M{"$set": bson.M{"reviewedBy": reviewer.ID, "reviewDate": time.Now(), "displayValue": Rejected}}

	// Update the phrase in Mongo to set it to be accepted
	_, err := phrasesCollection.UpdateOne(context.Background(), filter, updateDocument)
	if err != nil {
		return err
	}

	return nil
}

// Retrieve phrases in review for curators up to a specified number
func GetPhraseListForCurators(maxPhrases int64, phrasesCollection *mongo.Collection) ([]Phrase, error) {
	// Build the query document
	queryDocument := bson.M{"displayValue": Unreviewed}
	queryOptions := &options.FindOptions{Limit: &maxPhrases}

	// Get a cursor pointing to the list of phrases as a result of the query
	cur, err := phrasesCollection.Find(context.Background(), queryDocument, queryOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	// List of phrases and ObjectIDs for update
	var phraseList []Phrase
	var phraseObjectIDs []primitive.ObjectID

	// Get query result and print
	//for i := 0; i < maxPhrases && cur.Next(context.Background()); i++ {
	for cur.Next(context.Background()) {
		// Decode into struct
		var onePhrase Phrase
		err = cur.Decode(&onePhrase)
		if err != nil {
			return nil, err
		}

		// Append result to phraseList and append ObjectID
		phraseList = append(phraseList, onePhrase)
		phraseObjectIDs = append(phraseObjectIDs, onePhrase.PhraseID)
	}

	// Check for cursor errors
	if err := cur.Err(); err != nil {
		return nil, err
	}

	// Set all phrases to be in review
	filter := bson.M{"_id": bson.M{"$in": phraseObjectIDs}}
	update := bson.M{"$set": bson.M{"displayValue": InReview}}
	_, err = phrasesCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	// Sort the phrases
	sortPhrases(phraseList)

	// Return the result
	return phraseList, nil
}

// TODO: add delete all by userID function
// TODO: add function to anonymize by userID
// TODO: add function to get phrases by userID

// TODO: write this
// Query for phrases from a list of words
func GetPhraseList(wordList []WordRow, phrasesCollection *mongo.Collection) ([]Phrase, error) {
	// Get list of word IDS from wordList
	var wordIDs []int
	for _, w := range wordList {
		wordIDs = append(wordIDs, w.WordID)
	}

	// Build the query document
	queryDocument := bson.M{"wordList": bson.M{"$in": wordIDs}, "displayValue": Accepted}
	//queryDocument := bson.M{"wordList": bson.M{"$in": wordIDs}}

	// Get a cursor pointing to the list of phrases as a result of the query
	cur, err := phrasesCollection.Find(context.Background(), queryDocument)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	// list of phrases
	var phraseList []Phrase

	// get query result and print
	for cur.Next(context.Background()) {
		// Decode into struct
		var onePhrase Phrase
		//var onePhrase bson.d
		err = cur.Decode(&onePhrase)
		if err != nil {
			return nil, err
		}
		// append to phraseList
		phraseList = append(phraseList, onePhrase)
	}

	// check for cursor errors
	if err := cur.Err(); err != nil {
		return nil, err
	}

	// Sort phraseList by rating
	//sortPhrases(phraseList)

	// return the result
	return phraseList, nil
}

// Sort phrases by average rating
func sortPhrases(phraseList []Phrase) {
	// Put list of phrases into phraseSorter array
	var phrs []phraseSorter
	for _, p := range phraseList {
		phrs = append(phrs, phraseSorter{phrase: p, avgRating: AverageRating(p.PhraseRatings)})
	}

	// Sort the list
	sort.Sort(phraseSorterList(phrs))

	// Copy back into phraseList
	for i := range phrs {
		phraseList[i] = phrs[i].phrase
	}
}

// Get average rating from rating struct
func AverageRating(r Rating) float32 {
	totalRatings := r.OneStar + r.TwoStar + r.ThreeStar + r.FourStar + r.FiveStar
	if totalRatings == 0 {
		return float32(0)
	}
	weightedRatings := 1*r.OneStar + 2*r.TwoStar + 3*r.ThreeStar + 4*r.FourStar + 5*r.FiveStar
	return 5.0 * float32(weightedRatings) / float32(5*totalRatings)
}
