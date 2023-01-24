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

const artistsURL = "https://groupietrackers.herokuapp.com/api/artists"

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
	}

	Relation struct {
		Index []RelationsDetail `json:"index"`
	}

	RelationsDetail struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	}
)

// struct for the Artist page
type artistPrint struct {
	A Artist
	R RelationsDetail
}

var artist []Artist

func main() {
	getArtists()

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
		fmt.Fprint(w, "Error 404, page not found")
		return
	}
	checkError(err)

	var toPrint artistPrint
	//artist detailed information
	toPrint.A = artist[id1-1]
	//concerts information. Taking the Relations URL from the concrete artist struct and using it in getRelations function
	toPrint.R = getRelations(artist[id1-1].Relations)

	tmpl := template.Must(template.ParseFiles("templates/artist.html"))
	//artist info to the artist page
	tmpl.Execute(w, toPrint)
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

func getArtists() {
	//takes the URL as the argument, and returns a response and error. For successful API calls, err will be non-nil
	req, err := http.NewRequest("GET", artistsURL, nil)
	checkError(err)

	//send the request to receive the response from the API.
	resp, err := http.DefaultClient.Do(req)
	checkError(err)
	//now we can display the data returned by the API. But we have to close the response body when finished with it.
	defer resp.Body.Close()

	// access the response body using the
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	//Convert response body to artist.
	json.Unmarshal(body, &artist)
}

func getRelations(url string) RelationsDetail {

	//takes Relations URL as the argument, which was in the concrete Artist struct
	req, err := http.NewRequest("GET", url, nil)
	checkError(err)

	resp, err := http.DefaultClient.Do(req)
	checkError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	//using RelationsDetail struct, which is inside the Relations URL
	var rel RelationsDetail

	json.Unmarshal(body, &rel)
	return rel
}

func checkError(err error) {
	if err != nil {
		fmt.Print(err.Error())
	}
}
