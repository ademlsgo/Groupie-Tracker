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

	_, err = http.Get("https://rickandmortyapi.com/api/location") // Effectuer la demande à l'API
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialisation de l'URL pour la première page
	url := "https://rickandmortyapi.com/api/location"

	var allLocation AllLocation

	// Boucle pour récupérer toutes les pages
	for url != "https://rickandmortyapi.com/api/location?page=7" {

		// Utiliser le client HTTP standard pour effectuer la requête
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close() // Fermer le corps de la réponse

		// Décoder la réponse JSON de l'API
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
	// Passer la structure de données à votre modèle HTML
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

	// Effacer les résultats précédents
	filteredLocation = nil

	typeFilter := strings.ToLower(r.URL.Query().Get("type"))

	// Verifier si la liste des lieux est initialisée
	if len(locationList) == 0 {
		http.Error(w, "locationList is not initialized", http.StatusInternalServerError)
		return
	}

	// Application du filtrage
	for _, location := range locationList {
		if strings.ToLower(location.Type) == typeFilter {
			filteredLocation = append(filteredLocation, location)
		}
	}

	// Récupérer la valeur de la page à partir de la requête
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page < 0 {
		page = 0
	}

	// Nombre d'éléments par page
	itemsPerPage := 10

	// Calculer l'index de début et de fin des lieux pour la page actuelle
	startIndex := page * itemsPerPage
	endIndex := (page + 1) * itemsPerPage
	if endIndex > len(filteredLocation) {
		endIndex = len(filteredLocation)
	}

	// Vérifier si la page suivante a des éléments
	nextPage := ""
	if endIndex < len(filteredLocation) {
		nextPage = fmt.Sprintf("/submit_sort_location?type=%v&page=%v", typeFilter, page+1)
	} else {
		nextPage = fmt.Sprintf("/submit_sort_location?type=%v&page=%v", typeFilter, 0) // Retourner à la première page
	}

	// Extraire les lieux pour la page actuelle
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
