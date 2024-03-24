package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var locationList []datasLocation
var filteredLocation []datasLocation

type datasLocation struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Dimension string `json:"dimension"`
	Url       string `json:"url"`
}

type InfoL struct {
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}

type DataPageL struct {
	RecoverLocation []datasLocation
	PrevPage        string
	NextPage        string
	Type            string
}

type AllLocation struct {
	InfoL   InfoL           `json:"info"`
	Results []datasLocation `json:"results"`
}

func CollectionLocation(w http.ResponseWriter, r *http.Request) {
	page, errPage := strconv.Atoi(r.FormValue("page"))
	if errPage != nil || page < 0 {
		page = 0
	}
	tmpl, err := template.ParseFiles("./templates/CollectionLocation.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = http.Get("https://rickandmortyapi.com/api/location")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := "https://rickandmortyapi.com/api/location"

	var allLocation AllLocation

	for url != "https://rickandmortyapi.com/api/location?page=7" {

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&allLocation)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Mettre à jour l'URL pour la page suivante
		url = allLocation.InfoL.Next

		locationList = append(locationList, allLocation.Results...)
	}

	locationList := TenLocation(page)

	datas := DataPageL{
		RecoverLocation: locationList,
		PrevPage:        fmt.Sprintf("/collection/location?page=%v", page-1),
		NextPage:        fmt.Sprintf("/collection/location?page=%v", page+1),
	}
	err = tmpl.ExecuteTemplate(w, "CollectionLocation", datas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func TenLocation(page int) []datasLocation {
	var data []datasLocation
	for i := 0; i < 10; i++ {
		data = append(data, locationList[page*10+i])
	}
	return data
}

func RessourceLocation(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/RessourceLocation.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := r.URL.Query().Get("id")

	url1 := "https://rickandmortyapi.com/api/location/"
	url := url1 + id

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var Location datasLocation

	if err := json.Unmarshal(body, &Location); err != nil {
		log.Fatal(err)
	}

	Location.Url = fmt.Sprintf("/ressource/location?id=%v", Location.ID)

	err = tmpl.ExecuteTemplate(w, "RessourceLocation", Location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func FilterLocation(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/CollectionLocation.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filteredLocation = nil

	typeFilter := strings.ToLower(r.URL.Query().Get("type"))

	if len(locationList) == 0 {
		http.Error(w, "locationList is not initialized", http.StatusInternalServerError)
		return
	}

	for _, location := range locationList {
		if strings.ToLower(location.Type) == typeFilter {
			filteredLocation = append(filteredLocation, location)
		}
	}

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page < 0 {
		page = 0
	}

	itemsPerPage := 10

	startIndex := page * itemsPerPage
	endIndex := (page + 1) * itemsPerPage
	if endIndex > len(filteredLocation) {
		endIndex = len(filteredLocation)
	}

	nextPage := ""
	if endIndex < len(filteredLocation) {
		nextPage = fmt.Sprintf("/submit_sort_location?type=%v&page=%v", typeFilter, page+1)
	} else {
		nextPage = fmt.Sprintf("/submit_sort_location?type=%v&page=%v", typeFilter, 0) // Retourner à la première page
	}

	currentPageLocation := filteredLocation[startIndex:endIndex]

	datas := DataPageL{
		RecoverLocation: currentPageLocation,
		PrevPage:        fmt.Sprintf("/submit_sort_location?type=%v&page=%v", typeFilter, page-1),
		NextPage:        nextPage,
		Type:            typeFilter,
	}

	err = tmpl.ExecuteTemplate(w, "CollectionLocation", datas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
