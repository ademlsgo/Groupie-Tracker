package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type DatasCharacter struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	Gender  string `json:"gender"`
	Status  string `json:"status"`
	Species string `json:"species"`
	Type    string `json:"type"`
	Url     string `json:"url"`
}

type CharacterList struct {
	Info    Info             `json:"info"`
	Results []DatasCharacter `json:"results"`
}

type Info struct {
	Next string `json:"next"`
}

// VERIF si le nom du personnage contient la requête
func containsName(characterName, query string) bool {
	return strings.Contains(strings.ToLower(characterName), strings.ToLower(query))
}

// For obtenir le personnage par ID
func getCharacterByID(charactersList []DatasCharacter, id int) (*DatasCharacter, error) {
	for _, character := range charactersList {
		if character.ID == id {
			return &character, nil
		}
	}
	return nil, fmt.Errorf("Personnage introuvable avec l'ID %d", id)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Étape 1: Récupérer le nom du personnage depuis le formulaire
	query := r.FormValue("q")

	// Étape 2: Boucle pour récupérer toutes les pages
	page := 1
	var matchingCharacters []DatasCharacter
	for {
		// Construire l'URL avec le numéro de la page
		url := "https://rickandmortyapi.com/api/character?page=" + strconv.Itoa(page)

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var charactersList CharacterList
		err = json.NewDecoder(resp.Body).Decode(&charactersList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Étape 3: Vérifier si le nom est dans la liste
		for _, character := range charactersList.Results {
			if containsName(character.Name, query) {
				matchingCharacters = append(matchingCharacters, character)
			}
		}

		// Vérifier si la page suivante existe
		if charactersList.Info.Next == "" {
			break
		}

		// Passer à la page suivante
		page++
	}

	// Étape 4: Récupérer l'ID du premier personnage correspondant au nom recherché
	var characterID int
	var err error

	if len(matchingCharacters) > 0 {
		characterID = matchingCharacters[0].ID
	}

	// Étape 5: Récupérer le personnage correspondant à l'ID
	character, err := getCharacterByID(matchingCharacters, characterID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./templates/Search.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "Search", character)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
