package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Song struct {
	Title  string
	Artist string
}

func (s Song) String() string {
	return s.Title + " by " + s.Artist
}

func main() {
	fmt.Println("My MongoDB app!")

	// Connect to localhost
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Connect
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected!")

	// Get the songs collection from the cool_songs database
	songs := client.Database("cool_songs").Collection("songs")

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
