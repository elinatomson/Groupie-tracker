package main

import (
	"Groupie-tracker/handlers"
	"fmt"
	"log"
	"net/http"
)

const (
	artistsURL  = "https://groupietrackers.herokuapp.com/api/artists"
	relationURL = "https://groupietrackers.herokuapp.com/api/relation"
	locationURL = "https://groupietrackers.herokuapp.com/api/locations"
)

func main() {
	handlers.GetData(artistsURL, relationURL, locationURL)

	fileServer := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	http.HandleFunc("/", handlers.MainPage)
	http.HandleFunc("/artist/", handlers.ArtistPage)
	http.HandleFunc("/filters/", handlers.Filter)
	http.HandleFunc("/search/", handlers.Search)

	fmt.Printf("Starting server at port 8080\nOpen http://localhost:8080\nUse Ctrl+C to close the port\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
