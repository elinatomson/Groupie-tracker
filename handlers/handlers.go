package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
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
	tmpl.Execute(w, artist)
}

func ArtistPage(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/artist/")
	id1, err := strconv.Atoi(id)

	checkError(err)
	if id1 > 52 {
		fmt.Fprint(w, "Artist not found")
		return
	}
	checkError(err)

	tmpl := template.Must(template.ParseFiles("templates/artist.html"))
	tmpl.Execute(w, artist[id1-1])
}

func Search(w http.ResponseWriter, r *http.Request) {
	searched := r.FormValue("searched")
	data := artist

	if searched != "" {
		data = search(searched)
	}

	tmpl := template.Must(template.ParseFiles("templates/mainpage.html"))
	tmpl.Execute(w, data)
}

func Filter(w http.ResponseWriter, r *http.Request) {

	var filteredArtist []Artist
	filteredArtistInFilter := make(map[int]Artist)
	var isFilterActive bool
	isFilterActive = false

	//filter for creation year
	CreationDateFrom := r.FormValue("CreationDateFrom")
	CreationDateTo := r.FormValue("CreationDateTo")
	if CreationDateFrom != "" || CreationDateTo != "" {

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

		isFilterActive = true
		for _, value := range artist {
			if CDF <= value.CreationDate && value.CreationDate <= CDT {
				filteredArtistInFilter[value.ID] = value
			}
		}
	}

	//filter for the year of first album
	FirstAlbumFrom := r.FormValue("FirstAlbumFrom")
	FirstAlbumTo := r.FormValue("FirstAlbumTo")
	if FirstAlbumFrom != "" || FirstAlbumTo != "" {

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
		//if the previous filter was active (isFilterActive = true), then ranging over of these artists, which matched with the previous filter
		if isFilterActive {
			for artistID, value := range filteredArtistInFilter {
				FA := strings.Split(value.FirstAlbum, "-")
				fa, _ := separationArray(FA)
				//if the artist do not match with both filters, then the code is deleting this artist from the result
				if yearFrom > fa[2] || yearTo < fa[2] {
					delete(filteredArtistInFilter, artistID)

				}
			}
		} else {
			//if previous filter is not active, then else meaning ranging over all artists
			//setting it as true for the next filter
			isFilterActive = true
			for _, value := range artist {
				FA := strings.Split(value.FirstAlbum, "-")
				fa, _ := separationArray(FA)
				if yearFrom <= fa[2] && fa[2] <= yearTo {
					filteredArtistInFilter[value.ID] = value
				}
			}
		}
	}

	//filter for locations
	LocationOfConcertsValue := r.FormValue("LOC")

	if LocationOfConcertsValue != "" {
		if isFilterActive {
			for artistID, value := range filteredArtistInFilter {
				for j := range value.PlacesDates {
					if !caseIns(j, LocationOfConcertsValue) {
						delete(filteredArtistInFilter, artistID)
					}
				}
			}
		} else {
			isFilterActive = true
			for _, value := range artist {
				for j := range value.PlacesDates {
					if caseIns(j, LocationOfConcertsValue) {
						filteredArtistInFilter[value.ID] = value
					}
				}
			}
		}
	}

	//filter for number of members
	var NoM1, NoM2, NoM3, NoM4, NoM5, NoM6, NoM7, NoM8 int
	NoM1, _ = strconv.Atoi(r.FormValue("NoM1"))
	NoM2, _ = strconv.Atoi(r.FormValue("NoM2"))
	NoM3, _ = strconv.Atoi(r.FormValue("NoM3"))
	NoM4, _ = strconv.Atoi(r.FormValue("NoM4"))
	NoM5, _ = strconv.Atoi(r.FormValue("NoM5"))
	NoM6, _ = strconv.Atoi(r.FormValue("NoM6"))
	NoM7, _ = strconv.Atoi(r.FormValue("NoM7"))
	NoM8, _ = strconv.Atoi(r.FormValue("NoM8"))

	if NoM1 != 0 || NoM2 != 0 || NoM3 != 0 || NoM4 != 0 || NoM5 != 0 || NoM6 != 0 || NoM7 != 0 || NoM8 != 0 {
		if isFilterActive {
			for artistID, value := range filteredArtistInFilter {
				if NoM1 != len(value.Members) && NoM2 != len(value.Members) && NoM3 != len(value.Members) && NoM4 != len(value.Members) && NoM5 != len(value.Members) && NoM6 != len(value.Members) && NoM7 != len(value.Members) && NoM8 != len(value.Members) {
					delete(filteredArtistInFilter, artistID)
				}
			}
		} else {
			for _, value := range artist {
				if NoM1 == len(value.Members) || NoM2 == len(value.Members) || NoM3 == len(value.Members) || NoM4 == len(value.Members) || NoM5 == len(value.Members) || NoM6 == len(value.Members) || NoM7 == len(value.Members) || NoM8 == len(value.Members) {
					filteredArtistInFilter[value.ID] = value
				}
			}

		}
	}

	for _, value := range filteredArtistInFilter {
		filteredArtist = append(filteredArtist, value)
	}

	tmpl := template.Must(template.ParseFiles("templates/mainpage.html"))
	tmpl.Execute(w, filteredArtist)
}
