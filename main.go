package main

import (
	"log"
	"net/http"

	"github.com/alvarosness/punocracy/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
