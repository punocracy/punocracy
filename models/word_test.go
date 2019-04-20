package models

import (
        "fmt"    
        "testing"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func newWordForTest(t *testing.T) *Word{
    //edit the Connect string with own username and password
    db, err := sqlx.Connect("mysql","root:*password_here*@/proj_testing")
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
