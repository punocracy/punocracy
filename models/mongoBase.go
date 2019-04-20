// Baseline package for implementing MongoDB simple query operations
package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	// Find a song
	cur, err := songs.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	var songList []Song

	// Get query result and print
	for cur.Next(context.Background()) {
		// Decode into struct
		var oneSong Song
		//var oneSong bson.D
		err = cur.Decode(&oneSong)
		if err != nil {
			log.Fatal(err)
		}
		// Append to songlist
		songList = append(songList, oneSong)
	}

	// Check for cursor errors
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Print songlist
	for _, s := range songList {
		fmt.Println(s)
	}
}
