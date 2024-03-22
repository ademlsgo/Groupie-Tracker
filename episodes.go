package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var episodesList []datasEpisodes

type datasEpisodes struct {
	ID      int    `json:"id"`
	Episode string `json:"episode"`
	Name    string `json:"name"`
	Date    string `json:"air_date"`
	Url     string `json:"url"`
}

type InfoE struct {
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Next  string `json:"next"`
	Prev  string `json:"prev"`
}

type AllEpisodes struct {
	InfoE   InfoE           `json:"info"`
	Results []datasEpisodes `json:"results"`
}

func CollectionEpisodes(w http.ResponseWriter, r *http.Request) {
	page, errPage := strconv.Atoi(r.FormValue("page"))
	if errPage != nil || page < 0 {
		page = 0
	}
	tmpl, err := template.ParseFiles("./templates/CollectionEpisodes.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = http.Get("https://rickandmortyapi.com/api/episode") // Effectuer la demande à l'API
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialisation de l'URL pour la première page
	url := "https://rickandmortyapi.com/api/episode"

	var allEpisodes AllEpisodes

	// Boucle pour récupérer toutes les pages
	for url != "https://rickandmortyapi.com/api/episode?page=3" {

		// Utiliser le client HTTP standard pour effectuer la requête
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close() // Fermer le corps de la réponse

		// Décoder la réponse JSON de l'API
		err = json.NewDecoder(resp.Body).Decode(&allEpisodes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Mettre à jour l'URL pour la page suivante
		url = allEpisodes.InfoE.Next

		episodesList = append(episodesList, allEpisodes.Results...)
	}

	episodesList := TenEpisodesE(page)

	datas := DataPageE{
		RecoverEpisodes: episodesList,
		PrevPage:        page - 1,
		NextPage:        page + 1,
	}
	// Passer la structure de données à votre modèle HTML
	err = tmpl.ExecuteTemplate(w, "CollectionEpisodes", datas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type DataPageE struct {
	RecoverEpisodes []datasEpisodes
	PrevPage        int
	NextPage        int
}

func TenEpisodesE(page int) []datasEpisodes {
	var data []datasEpisodes
	for i := 0; i < 10; i++ {
		data = append(data, episodesList[page*10+i])
	}
	return data
}

func RessourceEpisodes(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/RessourceEpisodes.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := r.URL.Query().Get("id")

	url1 := "https://rickandmortyapi.com/api/episode/"
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

	var Episodes datasEpisodes

	if err := json.Unmarshal(body, &Episodes); err != nil {
		log.Fatal(err)
	}
	Episodes.Url = fmt.Sprintf("/ressource/episode?id=%v", Episodes.ID)

	err = tmpl.ExecuteTemplate(w, "RessourceEpisodes", Episodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
