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
	http.HandleFunc("/search/", searchHandler)

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

func searchHandler(w http.ResponseWriter, r *http.Request) {
	searched := r.FormValue("searched")
	data := artist

	if searched != "" {
		data = search(searched)
	}

	tmpl := template.Must(template.ParseFiles("templates/mainpage.html"))
	tmpl.Execute(w, data)
}

func filterHandler(w http.ResponseWriter, r *http.Request) {

	var filteredArtist []Artist

	//filter for creation year
	CreationDate := r.FormValue("CreationDate")
	if CreationDate == "on" {
		CreationDateFrom := r.FormValue("CreationDateFrom")
		CreationDateTo := r.FormValue("CreationDateTo")

		if CreationDateFrom == "" {
			CreationDateFrom = "1958"
		}
		if CreationDateTo == "" {
			CreationDateTo = "2015"
		}

		CDF, _ := strconv.Atoi(CreationDateFrom)
		CDT, _ := strconv.Atoi(CreationDateTo)

		if CDF > CDT {
			fmt.Fprint(w, "From year has to be earlier than to year")
			return
		}
		if CDF < 0 || CDT < 0 {
			fmt.Fprint(w, "Inserted year numbers have to be positive numbers")
			return
		}

		if len(artist) > 0 {
			for _, value := range artist {
				if CDF <= value.CreationDate && value.CreationDate <= CDT {
					filteredArtist = append(filteredArtist, value)
				}
			}
		}
	}

	//filter for the year of first album
	FirstAlbum := r.FormValue("FirstAlbum")
	if FirstAlbum == "on" {
		FirstAlbumFrom := r.FormValue("FirstAlbumFrom")
		FirstAlbumTo := r.FormValue("FirstAlbumTo")

		if FirstAlbumFrom == "" {
			FirstAlbumFrom = "1985"
		}
		if FirstAlbumTo == "" {
			FirstAlbumTo = "2025"
		}

		yearFrom, _ := strconv.Atoi(FirstAlbumFrom)
		yearTo, _ := strconv.Atoi(FirstAlbumTo)

		if yearFrom > yearTo {
			fmt.Fprint(w, "From year has to be earlier than to year")
			return
		}
		if yearFrom < 0 || yearTo < 0 {
			fmt.Fprint(w, "Inserted year numbers have to be positive numbers")
			return
		}

		if len(artist) > 0 {
			for _, value := range artist {
				FA := strings.Split(value.FirstAlbum, "-")
				fa, _ := separationArray(FA)
				if yearFrom <= fa[2] && fa[2] <= yearTo {
					filteredArtist = append(filteredArtist, value)
				}
			}
		}
	}

	//filter for number of members
	NumberOfMembers := r.FormValue("NumberOfMembers")
	if NumberOfMembers == "on" {
		NumberOfMembers2 := r.FormValue("NumberOfMembers2")

		NOM, _ := strconv.Atoi(NumberOfMembers2)

		if len(artist) > 0 {
			for _, value := range artist {
				if NOM == len(value.Members) {
					filteredArtist = append(filteredArtist, value)
				}
			}
		}
	}

	//filter for locations
	LocationOfConcerts := r.FormValue("LocationOfConcerts")
	if LocationOfConcerts == "on" {
		LocationOfConcertsValue := r.FormValue("LOC")

		if LocationOfConcertsValue != "" {
			if len(artist) > 0 {
				for _, value := range artist {
					for j := range value.PlacesDates {
						if caseIns(j, LocationOfConcertsValue) {
							filteredArtist = append(filteredArtist, value)
						}
					}
				}
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

func search(search_input string) []Artist {
	result, _ := getData(artistsURL, relationURL)
	var data []Artist

	for band := range result {
		if caseIns(result[band].Name, search_input) {
			data = append(data, result[band])
			continue
		}

		for j := range result[band].PlacesDates {
			if caseIns(j, search_input) {
				data = append(data, result[band])
				continue
			}
		}

		if search_input == result[band].FirstAlbum {
			data = append(data, result[band])
			continue
		}

		if search_input == strconv.Itoa(result[band].CreationDate) {
			data = append(data, result[band])
			continue
		}
		for j := range result[band].Members {
			if caseIns(result[band].Members[j], search_input) {
				data = append(data, result[band])
				continue
			}
		}
	}
	return data
}

// first album date into integer array
func separationArray(array []string) ([]int, bool) {
	arrayInt := []int{}
	for _, value := range array {
		Int, err := strconv.Atoi(value)
		if err != nil {
			return []int{}, false
		}
		arrayInt = append(arrayInt, Int)
	}
	return arrayInt, true
}

// handling search input as case-insensitive
func caseIns(search_input, result string) bool {
	return strings.Contains(
		strings.ToLower(search_input),
		strings.ToLower(result),
	)
}

func checkError(err error) {
	if err != nil {
		fmt.Print(err.Error())
	}
}
