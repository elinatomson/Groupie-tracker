package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	artistsURL  = "https://groupietrackers.herokuapp.com/api/artists"
	relationURL = "https://groupietrackers.herokuapp.com/api/relation"
)

type (
	Artist struct {
		ID           int      `json:"id"`
		Image        string   `json:"image"`
		Name         string   `json:"name"`
		Members      []string `json:"members"`
		CreationDate int      `json:"creationDate"`
		FirstAlbum   string   `json:"firstAlbum"`
		Locations    string   `json:"locations"`
		ConcertDates string   `json:"concertDates"`
		Relations    string   `json:"relations"`
		PlacesDates  map[string][]string
	}

	Relation struct {
		Index []RelationsDetail `json:"index"`
	}

	RelationsDetail struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	}
)

var artist []Artist
var relation Relation

func main() {
	getData(artistsURL, relationURL)

	fileServer := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/artist/", artistHandler)
	http.HandleFunc("/filters/", filterHandler)

	fmt.Printf("Starting server at port 8080\nOpen http://localhost:8080\nUse Ctrl+C to close the port\n")
	http.ListenAndServe(":8080", nil)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error 404, page not found")
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "Error 405, method not allowed")
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/mainpage.html"))
	//all artists icons to the main page
	tmpl.Execute(w, artist)
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/artist/")
	id1, err := strconv.Atoi(id)
	if id1 > 52 {
		fmt.Fprint(w, "Artist not found")
		return
	}
	checkError(err)

	tmpl := template.Must(template.ParseFiles("templates/artist.html"))
	tmpl.Execute(w, artist[id1-1])
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	var filteredArtist []Artist

	//for HTML
	CreationDateFrom := r.FormValue("CreationDateFrom")
	CreationDateTo := r.FormValue("CreationDateTo")

	//if the field is left empty
	if CreationDateFrom == "" {
		CreationDateFrom = "0"
	}
	if CreationDateTo == "" {
		CreationDateTo = "2111"
	}

	//strings to int
	dateFrom, _ := strconv.Atoi(CreationDateFrom)
	dateTo, _ := strconv.Atoi(CreationDateTo)

	//conditions for the years
	if dateFrom > dateTo {
		fmt.Fprint(w, "From year has to be earlier than to year")
		return
	}
	if dateFrom < 0 || dateTo < 0 {
		fmt.Fprint(w, "Years have to be positive numbers")
		return
	}

	//checking creation year of artists
	if len(artist) > 0 {
		for _, value := range artist {
			if dateFrom <= value.CreationDate && value.CreationDate <= dateTo {
				filteredArtist = append(filteredArtist, value)
			}
		}
	}
	tmpl := template.Must(template.ParseFiles("templates/mainpage.html"))
	tmpl.Execute(w, filteredArtist)
}

func getData(artistURL, relationURL string) ([]Artist, int) {

	reqArtists, err := http.NewRequest("GET", artistsURL, nil)
	checkError(err)

	respArtists, err := http.DefaultClient.Do(reqArtists)
	checkError(err)
	defer respArtists.Body.Close()

	body1, err := ioutil.ReadAll(respArtists.Body)
	checkError(err)

	json.Unmarshal(body1, &artist)

	reqRelations, err := http.NewRequest("GET", relationURL, nil)
	checkError(err)

	respRelations, err := http.DefaultClient.Do(reqRelations)
	checkError(err)
	defer respRelations.Body.Close()

	body2, err := ioutil.ReadAll(respRelations.Body)
	checkError(err)

	json.Unmarshal(body2, &relation)

	for i, v := range relation.Index {
		artist[i].PlacesDates = v.DatesLocations
	}

	return artist, http.StatusOK
}

func checkError(err error) {
	if err != nil {
		fmt.Print(err.Error())
	}
}
