package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
)

func main() {
	fmt.Println("Hello, world!")

	// Connect to localhost, and get reference to collection
	Session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	Database := "cool_songs"
	Collection := "songs"
	songs := Session.DB(Database).C(Collection)

	// Map data structure for BSON stuff
	SongMap := new(bson.M)

	// Query for some song
	//err = songs.Find(nil).One(&SongMap)
	err = songs.Find(bson.M{"title": "The Year 3000"}).One(&SongMap)
	if err != nil {
		log.Fatal(err)
	}

	// Print query result
	fmt.Println(SongMap)
}
