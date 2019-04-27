package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"testing"
)

func newWordForTest(t *testing.T) *Word {
	//edit the Connect string with own username and password
	db, err := sqlx.Connect("mysql", "root:root@/punocracy")
	if err != nil {
		t.Errorf(" database didnt connect. Error: %v", err)
	}
	return NewWord(db)
}

func TestQueryAlph(t *testing.T) {
	w := newWordForTest(t)
	wordList, _ := w.QueryAlph(nil, 'a')

	t.Logf("checking for letter: a\n")
	for _, v := range wordList {
		t.Logf("%#v\n", v)
	}
}

func TestQueryHlistString(t *testing.T) {
	w := newWordForTest(t)
	wordList, _ := w.QueryHlistString(nil, "brono")

	t.Logf("checking for word: brono's homophones\n")
	for _, v := range wordList {
		t.Logf("%#v\n", v)
	}
}

func TestNILQueryHlistString(t *testing.T) {
	w := newWordForTest(t)
	//somthing that does not exist
	wordList, erri := w.QueryHlistString(nil, "fakedude")

	t.Logf("checking what happens on nil\n")
	if erri != nil {
		t.Logf("an error happend (as expected)\n")
		return
	}
	for _, v := range wordList {
		t.Logf("%#v\n", v)
	}
}

/*
test the wordId list generator
*/
func TestNormalIDList(t *testing.T) {
	w := newWordForTest(t)
	wordSliceTest := []string{"brono", "Arono"}
	intSlice, err := w.GetWordIDList(nil, wordSliceTest)

	if err != nil {
		t.Errorf("Id list failed Error: %v", err)
	}
	t.Logf("ids: %v , %v \n", intSlice[0], intSlice[1])
}

func TestEmptyIDList(t *testing.T) {
	w := newWordForTest(t)
	wordSliceTest := []string{}
	_, err := w.GetWordIDList(nil, wordSliceTest)

	if err == nil {
		t.Errorf("An error was supposed to happen")
	}
}

func TestRandWordsList(t *testing.T) {
	w := newWordForTest(t)
	words, err := w.RandWordsList(nil, 3)
	if err != nil {
		t.Errorf("an error occured getting random words %v", err)
	}

	for _, v := range words {
		t.Logf("ids: %v \n", v)
	}
}
