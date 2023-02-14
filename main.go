package main

import (
	"Groupie-tracker/functions"
	"fmt"
	"net/http"
)

const (
	artistsURL  = "https://groupietrackers.herokuapp.com/api/artists"
	relationURL = "https://groupietrackers.herokuapp.com/api/relation"
	locationURL = "https://groupietrackers.herokuapp.com/api/locations"
)

func main() {
	functions.GetData(artistsURL, relationURL, locationURL)

	fileServer := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	http.HandleFunc("/", functions.MainPageHandler)
	http.HandleFunc("/artist/", functions.ArtistHandler)
	http.HandleFunc("/filters/", functions.FilterHandler)
	http.HandleFunc("/search/", functions.SearchHandler)

	fmt.Printf("Starting server at port 8080\nOpen http://localhost:8080\nUse Ctrl+C to close the port\n")
	http.ListenAndServe(":8080", nil)
}
