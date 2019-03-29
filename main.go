package main

import (
	"log"
	"net/http"

	"github.com/alvarosness/punocracy/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
