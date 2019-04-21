package models

import (
        "fmt"    
        "testing"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func newWordForTest(t *testing.T) *Word{
    //edit the Connect string with own username and password
    db, err := sqlx.Connect("mysql","root:@/proj_testing")
    if err != nil {
        t.Errorf(" database didnt connect. Error: %v", err)
    }
    return NewWord(db)
}

func TestQueryAlph(t *testing.T) {
    w := newWordForTest(t)
    wordList, _:= w.QueryAlph(nil,'a')
    
    fmt.Printf("checking for letter: a\n")
    for _,v := range wordList{
        fmt.Printf("%#v\n", v)
    }
}

func TestQueryHlistString(t *testing.T) {
    w := newWordForTest(t)
    wordList, _ := w.QueryHlistString(nil,"brono")
    
    fmt.Printf("checking for word: brono's homophones\n")
    for _,v := range wordList{
        fmt.Printf("%#v\n", v)
    }
}

func TestNILQueryHlistString(t *testing.T) {
    w := newWordForTest(t)
    //somthing that does not exist
    wordList, erri := w.QueryHlistString(nil,"fakedude")
    
    fmt.Printf("checking what happens on nil\n")
    if erri != nil{
        fmt.Printf("an error happend (as expected)\n")
        return
    }
    for _,v := range wordList{
        fmt.Printf("%#v\n", v)
    }
}

/*
test the wordId list generator
*/
func TestNormalIDList(t *testing.T){
    w := newWordForTest(t)
    wordSliceTest := []string{"brono","Arono"}
    intSlice, err := w.GetWordIDList(nil, wordSliceTest)
    
    if(err != nil){
        t.Errorf("Id list failed Error: %v",err)
    }
    fmt.Printf("ids: %v , %v \n",intSlice[0],intSlice[1])
}

func TestEmptyIDList(t *testing.T){
    w := newWordForTest(t)
    wordSliceTest := []string{}
    _, err := w.GetWordIDList(nil, wordSliceTest)

    if(err == nil){
        t.Error("An error was supposed to happen")
    }
}
