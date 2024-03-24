package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type datasCharacterFav struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Gender  string `json:"gender"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Type    string `json:"type"`
}

type datasEpisodesFav struct {
	ID      int    `json:"id"`
	Episode string `json:"episode"`
	Name    string `json:"name"`
	Date    string `json:"air_date"`
}

type datasLocationFav struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Dimension string `json:"dimension"`
}

var favorites []struct {
	ID           int    `json:"id"`
	TypeEndpoint string `json:"typeEndpoint"`
}

func FavoriteTemp(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/FavoriteTemplate.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := os.ReadFile("Favorite.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &favorites)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var datasFavCharacter []datasCharacterFav
	var datasFavEpisodes []datasEpisodesFav
	var datasFavLocation []datasLocationFav

	for _, fav := range favorites {
		if fav.TypeEndpoint == "character" {
			characterData := fetchDataCharacter(fav.ID) // Récupère les données du personnage, location, épisode
			datasFavCharacter = append(datasFavCharacter, characterData)
		} else if fav.TypeEndpoint == "episode" {
			episodeData := fetchDataEpisode(fav.ID)
			datasFavEpisodes = append(datasFavEpisodes, episodeData)
		} else if fav.TypeEndpoint == "location" {
			locationData := fetchDataLocation(fav.ID)
			datasFavLocation = append(datasFavLocation, locationData)
		}
	}

	dataForTemplate := struct {
		FavoritesCharacter []datasCharacterFav
		FavoritesEpisodes  []datasEpisodesFav
		FavoritesLocation  []datasLocationFav
	}{
		FavoritesCharacter: datasFavCharacter,
		FavoritesEpisodes:  datasFavEpisodes,
		FavoritesLocation:  datasFavLocation,
	}

	err = tmpl.ExecuteTemplate(w, "FavoriteTemp", dataForTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func fetchData(url string, d interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("La requête HTTP a échoué avec le code d'état : %d", response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(d)
	if err != nil {
		return err
	}

	return nil
}

func fetchDataCharacter(id int) datasCharacterFav {
	url := "https://rickandmortyapi.com/api/character/" + strconv.Itoa(id)
	characterData := datasCharacterFav{}
	err := fetchData(url, &characterData)
	if err != nil {
		return datasCharacterFav{}
	}
	return characterData
}

func fetchDataEpisode(id int) datasEpisodesFav {
	url := "https://rickandmortyapi.com/api/episode/" + strconv.Itoa(id)
	episodeData := datasEpisodesFav{}
	err := fetchData(url, &episodeData)
	if err != nil {
		return datasEpisodesFav{}
	}
	return episodeData
}

func fetchDataLocation(id int) datasLocationFav {
	url := "https://rickandmortyapi.com/api/location/" + strconv.Itoa(id)
	locationData := datasLocationFav{}
	err := fetchData(url, &locationData)
	if err != nil {
		return datasLocationFav{}
	}
	return locationData
}
