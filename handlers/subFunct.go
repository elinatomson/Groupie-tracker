package handlers

import (
	"fmt"
	"strconv"
	"strings"
)

func search(search_input string) []Artist {
	result, _ := GetData(artistsURL, relationURL, locationURL)
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
