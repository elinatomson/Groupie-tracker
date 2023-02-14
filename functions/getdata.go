package functions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	artistsURL  = "https://groupietrackers.herokuapp.com/api/artists"
	relationURL = "https://groupietrackers.herokuapp.com/api/relation"
	locationURL = "https://groupietrackers.herokuapp.com/api/locations"
)

type (
	Artist struct {
		ID           int      `json:"id"`
		Image        string   `json:"image"`
		Name         string   `json:"name"`
		Members      []string `json:"members"`
		CreationDate int      `json:"creationDate"`
		FirstAlbum   string   `json:"firstAlbum"`
		Locations    []string `json:"locations"`
		ConcertDates []string `json:"concertDates"`
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

	Location struct {
		Index []LocationsDetail `json:"index"`
	}

	LocationsDetail struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	}
)

var artist []Artist
var relation Relation
var locations Location

func GetData(artistURL, relationURL, locationURL string) ([]Artist, int) {

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
	//geolocalizationi jaoks eraldi Locations andmete võtmine juurde lisatud
	reqLocations, err := http.NewRequest("GET", locationURL, nil)
	checkError(err)

	respLocations, err := http.DefaultClient.Do(reqLocations)
	checkError(err)
	defer respLocations.Body.Close()

	body3, err := ioutil.ReadAll(respLocations.Body)
	checkError(err)

	json.Unmarshal(body3, &locations)

	for i, value := range locations.Index {
		artist[i].Locations = value.Locations
	}

	return artist, http.StatusOK
}
