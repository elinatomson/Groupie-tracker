package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	//HTML-is kasutatavad
	CreationDate := r.FormValue("CreationDate")
	if CreationDate == "on" {
		CreationDateFrom := r.FormValue("CreationDateFrom")
		CreationDateTo := r.FormValue("CreationDateTo")

		//kui väli jäetakse tühjaks
		if CreationDateFrom == "" {
			CreationDateFrom = "0"
		}
		if CreationDateTo == "" {
			CreationDateTo = "2111"
		}

		//kirjutatud stringid peab intiks tegema
		dateFrom, _ := strconv.Atoi(CreationDateFrom)
		dateTo, _ := strconv.Atoi(CreationDateTo)

		//aastate kirjapildi nõuded
		if dateFrom > dateTo {
			fmt.Fprint(w, "From year has to be earlier than to year")
			return
		}
		if dateFrom < 0 || dateTo < 0 {
			fmt.Fprint(w, "Years have to be positive numbers")
			return
		}

		//artisti aastaarvu kontrollimine
		if len(artist) > 0 {
			for _, value := range artist {
				// kui sisestatud fromDatest on väiksem või võrdne CreationDate-st ja CreationDate on omakorda väiksem võrdne toDatest
				if dateFrom <= value.CreationDate && value.CreationDate <= dateTo {
					filteredArtist = append(filteredArtist, value)
				}
			}
		}
	}

	//HTML-is kasutatavad
	FirstAlbumDate := r.FormValue("FirstAlbumDate")
	if FirstAlbumDate == "on" {
		FirstAlbumDateFrom := r.FormValue("FirstFrom")
		FirstAlbumDateTo := r.FormValue("FirstTo")

		//kui väli jäetakse tühjaks
		if FirstAlbumDateFrom == "" {
			FirstAlbumDateFrom = "20-01-1000"
		}
		if FirstAlbumDateTo == "" {
			FirstAlbumDateTo = "20-01-2111"
		}

		//kuupäevad sidekriipsuga splittimine
		from := strings.Split(FirstAlbumDateFrom, "-")
		to := strings.Split(FirstAlbumDateTo, "-")

		//kuupäevade kirjapildi nõudeks 3 argumenti
		if len(from) != 3 || len(to) != 3 {
			fmt.Fprint(w, "The date is not in a correct form")
			return
		}

		//artisti albumite kuupäevade kontrollimine
		fr, _ := separationArray(from)
		t, _ := separationArray(to)

		if len(artist) > 0 {
			for _, value := range artist {

				FAD := strings.Split(value.FirstAlbum, "-") //artisti andmetest splitib kuupäeva sidekriipsudega
				fad, _ := separationArray(FAD)

				if comparison(fad, fr, t) {
					filteredArtist = append(filteredArtist, value)
				}
			}
		}

	}

	tmpl := template.Must(template.ParseFiles("templates/mainpage.html"))
	tmpl.Execute(w, filteredArtist)
}

func separationArray(array []string) ([]int, bool) { //see func on natuke segane mulle
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

func comparison(fad, fr, t []int) bool { //see func on natuke segane mulle
	Fad := time.Month(fad[1])
	from := time.Month(fr[1])
	to := time.Month(t[1])

	FAD := time.Date(fad[2], Fad, fad[0], 0, 0, 0, 0, time.UTC)
	FROM := time.Date(fr[2], from, fr[0], 0, 0, 0, 0, time.UTC)
	TO := time.Date(t[2], to, t[0], 0, 0, 0, 0, time.UTC)

	From := FROM.Before(FAD)
	To := TO.After(FAD)
	EqualFrom := FROM.Equal(FAD)
	EqualTo := TO.Equal(FAD)

	if EqualFrom == true || EqualTo == true {
		return true
	}

	if From == true && To == true {
		return true
	}

	return false
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
