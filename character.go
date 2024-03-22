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

var charactersList []datasCharacter
var filteredCharacters []datasCharacter

type datasCharacter struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Gender  string `json:"gender"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Type    string `json:"type"`
	Url     string `json:"url"`
}

type InfoC struct {
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}

type DataPage struct {
	RecoverCharacters []datasCharacter
	PrevPage          string
	NextPage          string
	Species           string
}

type AllCharacters struct {
	InfoC   InfoC            `json:"info"`
	Results []datasCharacter `json:"results"`
}

func CollectionCharacter(w http.ResponseWriter, r *http.Request) {
	page, errPage := strconv.Atoi(r.FormValue("page"))
	if errPage != nil || page < 0 {
		page = 0
	}
	tmpl, err := template.ParseFiles("./templates/CollectionCharacter.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = http.Get("https://rickandmortyapi.com/api/character") // Effectuer la demande à l'API
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialisation de l'URL pour la première page
	url := "https://rickandmortyapi.com/api/character"

	var allCharacters AllCharacters

	// Boucle pour récupérer toutes les pages
	for url != "https://rickandmortyapi.com/api/character?page=42" {

		// Utiliser le client HTTP standard pour effectuer la requête
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close() // Fermer le corps de la réponse

		// Décoder la réponse JSON de l'API
		err = json.NewDecoder(resp.Body).Decode(&allCharacters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Mettre à jour l'URL pour la page suivante
		url = allCharacters.InfoC.Next

		charactersList = append(charactersList, allCharacters.Results...)

	}

	charactersList := TenCharacter(page)

	datas := DataPage{
		RecoverCharacters: charactersList,
		PrevPage:          fmt.Sprintf("/collection/character?page=%v", page-1),
		NextPage:          fmt.Sprintf("/collection/character?page=%v", page+1),
	}
	// Passer la structure de données à votre modèle HTML
	err = tmpl.ExecuteTemplate(w, "CollectionCharacter", datas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func TenCharacter(page int) []datasCharacter {
	var data []datasCharacter
	for i := 0; i < 10; i++ {
		data = append(data, charactersList[page*10+i])
	}
	return data
}

func RessourceCharacter(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("./templates/RessourceCharacter.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := r.URL.Query().Get("id")

	url1 := "https://rickandmortyapi.com/api/character/"
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

	var Characters datasCharacter

	if err := json.Unmarshal(body, &Characters); err != nil {
		log.Fatal(err)
	}
	Characters.Url = fmt.Sprintf("/ressource/character?id=%v", Characters.ID)
	err = tmpl.ExecuteTemplate(w, "RessourceCharacter", Characters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func FilterCharacter(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/CollectionCharacter.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Effacer les résultats précédents
	filteredCharacters = nil

	speciesFilter := strings.ToLower(r.URL.Query().Get("species"))

	// Verifier si la liste des personnages est initialisée
	if len(charactersList) == 0 {
		http.Error(w, "charactersList is not initialized", http.StatusInternalServerError)
		return
	}

	// Application du filtrage
	for _, character := range charactersList {
		if strings.ToLower(character.Species) == speciesFilter {
			filteredCharacters = append(filteredCharacters, character)
		}
	}

	// Récupérer la valeur de la page à partir de la requête
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page < 0 {
		page = 0
	}

	// Nombre d'éléments par page
	itemsPerPage := 10

	// Calculer l'index de début et de fin des personnages pour la page actuelle
	startIndex := page * itemsPerPage
	endIndex := (page + 1) * itemsPerPage
	if endIndex > len(filteredCharacters) {
		endIndex = len(filteredCharacters)
	}

	// Vérifier si la page suivante a des éléments
	nextPage := ""
	if endIndex < len(filteredCharacters) {
		nextPage = fmt.Sprintf("/submit_sort?species=%v&page=%v", speciesFilter, page+1)
	} else {
		nextPage = fmt.Sprintf("/submit_sort?species=%v&page=%v", speciesFilter, 0) // Retourner à la première page
	}

	// Extraire les personnages pour la page actuelle
	currentPageCharacters := filteredCharacters[startIndex:endIndex]

	datas := DataPage{
		RecoverCharacters: currentPageCharacters,
		PrevPage:          fmt.Sprintf("/submit_sort?species=%v&page=%v", speciesFilter, page-1),
		NextPage:          nextPage,
		Species:           speciesFilter,
	}

	err = tmpl.ExecuteTemplate(w, "CollectionCharacter", datas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
