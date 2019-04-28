// Data structures for our MongoDB data.

package models

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
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
	avgRating float64
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
	logrus.Infoln(phraseID)
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

/*
TODO:
When curator deletes their account, a function is needed to anonimize their data and release inreview phrases
*/

/*
This function fully manages the curator phrases by:
> Assigning curators to phrases
> First collecting phrases already assigned to a curator before including new ones
> returning that bundle of curator phrases
> updates the database with the curator assignments
input:  maxPhrases (amount of phrases to generate)
        curatingUser (curator UserRow)
        phrasesCollection (the mongo collection to work with)
ouput:  Phrase slice and error

*/
func GetPhraseListForCurators(maxPhrases int64, curatingUser UserRow, phrasesCollection *mongo.Collection) ([]Phrase, error) {
	//get up to maxPhrases in review phrases
	inReviewPhrases, err := GetInReviewPhraseList(maxPhrases, curatingUser, phrasesCollection)
	if err != nil {
		return nil, err
	}

	//get the rest of the phrases from "new" phrases
	if int64(len(inReviewPhrases)) < maxPhrases {
		newPhrases, err2 := GetNewPhraseListForCurators((maxPhrases - int64(len(inReviewPhrases))), curatingUser, phrasesCollection)

		if err2 != nil {
			return nil, err2
		}
		//append inreview phrases and new phrases
		allPhrases := append(inReviewPhrases, newPhrases...)
		return allPhrases, nil
	} else {

		return inReviewPhrases, nil
	}

}

/*
This function will retireve phrases that are in review by a curator up to maxPhrases
*/
func GetInReviewPhraseList(maxPhrases int64, curatingUser UserRow, phrasesCollection *mongo.Collection) ([]Phrase, error) {

	// Build the query document
	queryDocument := bson.M{"displayValue": InReview, "reviewedBy": curatingUser.ID}
	queryOptions := &options.FindOptions{Limit: &maxPhrases}

	// Get a cursor pointing to the list of phrases as a result of the query
	cur, err := phrasesCollection.Find(context.Background(), queryDocument, queryOptions)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	// List of phrases
	var phraseList []Phrase

	// Get query result and print
	//for i := 0; i < maxPhrases && cur.Next(context.Background()); i++ {}
	for cur.Next(context.Background()) {
		// Decode into struct
		var onePhrase Phrase
		err = cur.Decode(&onePhrase)
		if err != nil {
			return nil, err
		}

		// Append result to phraseList and append ObjectID
		phraseList = append(phraseList, onePhrase)
	}

	// Check for cursor errors
	if err := cur.Err(); err != nil {
		return nil, err
	}

	// Sort the phrases
	sortPhrases(phraseList)

	// Return the result
	return phraseList, nil
}

/*
TODO:

Keep track of reviewer assigned to phrases

New functionality: first query the phrases that are already inreview by that curator and append them before appending new phrases to the slice

New input: UserRow

Function signature change, previous name:GetPhraseListForCurators
                           NEW NAME: GetNewPhraseListForCurators

Plan is to change this to a helper function to a new overall GetPhraseListForCurators function.
that will first query for exsisting curator assigned phrases then append the results of this function.
*/
// Retrieve phrases in review for curators up to a specified number
func GetNewPhraseListForCurators(maxPhrases int64, curatingUser UserRow, phrasesCollection *mongo.Collection) ([]Phrase, error) {
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

		//BEFORE APPENDING: edit these results then update them in the DB
		// Append result to phraseList and append ObjectID
		onePhrase.ReviewedBy = curatingUser.ID
		onePhrase.DisplayPublic = InReview
		phraseList = append(phraseList, onePhrase)
		phraseObjectIDs = append(phraseObjectIDs, onePhrase.PhraseID)
	}

	// Check for cursor errors
	if err := cur.Err(); err != nil {
		return nil, err
	}

	if len(phraseList) != 0 {
		// Set phrases to be in review
		// and assigned to curator
		filter := bson.M{"_id": bson.M{"$in": phraseObjectIDs}}
		update := bson.M{"$set": bson.M{"reviewedBy": curatingUser.ID, "displayValue": InReview}}
		_, err = phrasesCollection.UpdateMany(context.Background(), filter, update)
		if err != nil {
			return nil, err
		}

		// Sort the phrases
		sortPhrases(phraseList)

	}

	// Return the result
	return phraseList, nil
}

// Delete all phrases by a single userID
func DeleteByUserID(user UserRow, phrasesCollection *mongo.Collection) error {
	// Build query document
	filterDocument := bson.M{"submitterUserID": user.ID}

	// Execute delete statement
	_, err := phrasesCollection.DeleteMany(context.Background(), filterDocument)
	return err
}

// TODO: slice of phrases, reset their flags
// TODO: add function to anonymize by userID
// Anonimize user data for phrases
func AnonimizeUserData(user UserRow, phrasesCollection *mongo.Collection) error {
	//build query document
	filter := bson.M{"submitterUserID": user.ID}
	// Update all
	//define anon user
	//user 0 Should be const?
	update := bson.M{"$set": bson.M{"submitterUserID": 0}}
	_, err := phrasesCollection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// TODO: add function to get phrases by userID

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

// GetPhraseHistory for phrases from a list of words
func GetPhraseHistory(user UserRow, phrasesCollection *mongo.Collection) ([]Phrase, error) {
	// Build the query document
	queryDocument := bson.M{"submitterUserID": user.ID}

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

	// return the result
	return phraseList, nil
}

// Get average rating from rating struct
func AverageRating(r Rating) float64 {
	totalRatings := r.OneStar + r.TwoStar + r.ThreeStar + r.FourStar + r.FiveStar
	if totalRatings == 0 {
		return float64(0)
	}
	weightedRatings := 1*r.OneStar + 2*r.TwoStar + 3*r.ThreeStar + 4*r.FourStar + 5*r.FiveStar
	return 5.0 * float64(weightedRatings) / float64(5*totalRatings)
}

// GetTopPhrases gets a sorted list of the top phrases, limited by a number
func GetTopPhrases(limit int, phrases *mongo.Collection) ([]Phrase, error) {
	// Aggregation pipeline
	pipeline := bson.D{
		bson.M{
			"$addFields": bson.M{
				"numRatings": bson.M{
					"$sum": bson.D{
						"$ratings.one",
						"$ratings.two",
						"$ratings.three",
						"$ratings.four",
						"$ratings.five",
					},
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"avgRating": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": bson.D{"$numRatings", 0}},
						"then": 0,
						"else": bson.M{
							"$divide": bson.D{
								bson.M{
									"$sum": bson.D{
										"$ratings.one",
										bson.M{"$multiply": bson.D{"$ratings.two", 2}},
										bson.M{"$multiply": bson.D{"$ratings.three", 3}},
										bson.M{"$multiply": bson.D{"$ratings.four", 4}},
										bson.M{"$multiply": bson.D{"$ratings.five", 5}},
									},
								},
								"$numRatings",
							},
						},
					},
				},
			},
		},
		bson.M{
			"$sort": bson.M{"avgRating": -1, "numRatings": -1},
		},
		bson.M{
			"$limit": limit,
		},
	}

	// Execute aggregation
	cur, err := phrases.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	// Decode results into list
	var topPhrases []Phrase
	for cur.Next() {
		var thePhrase Phrase
		err = cur.Decode(&thePhrase)
		if err != nil {
			return nil, err
		}
		topPhrases = append(topPhrases, thePhrase)
	}
	return topPhrases, nil
}
