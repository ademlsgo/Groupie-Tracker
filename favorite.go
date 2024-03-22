package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

var favorite []datasFavorite

type datasFavorite struct {
	ID           int    `json:"id"`
	TypeEndpoint string `json:"typeEndpoint"`
}

func RessourceFavorite(w http.ResponseWriter, r *http.Request) {

	ID := r.URL.Query().Get("id")
	TypeEndpoint := r.URL.Query().Get("typeEndpoint")
	urlRedir := r.URL.Query().Get("url")

	// Instance
	newFavorite := datasFavorite{
		ID:           convertStringToInt(ID),
		TypeEndpoint: TypeEndpoint,
	}

	data, err := os.ReadFile("Favorite.json")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(data, &favorite) // Décode les données Json en string
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Vérifier si le même favori et typeEndpoint existent déjà, si oui, rediriger
	for _, i := range favorite {
		if i.ID == newFavorite.ID && i.TypeEndpoint == newFavorite.TypeEndpoint {
			http.Redirect(w, r, urlRedir, http.StatusSeeOther)
			return
		}
	}

	// Ajouter le nouveau favori au tableau
	favorite = append(favorite, newFavorite)

	// Traduire les données string en JSON
	data, err = json.Marshal(favorite)
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Écrire les données dans le fichier JSON
	err = os.WriteFile("Favorite.json", data, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, urlRedir, http.StatusSeeOther)
}

func convertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres de l'URL
	ID := r.URL.Query().Get("id")
	TypeEndpoint := r.URL.Query().Get("typeEndpoint")

	// Vérifier si les paramètres sont vides
	if ID == "" || TypeEndpoint == "" {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	// Lire le fichier JSON
	data, err := os.ReadFile("Favorite.json")
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Décoder les données JSON en slice de favoris
	var favorites []datasFavorite
	err = json.Unmarshal(data, &favorites)
	if err != nil {
		fmt.Println("Erreur lors du décodage du JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Rechercher et supprimer le favori correspondant
	for i, fav := range favorites {
		if strconv.Itoa(fav.ID) == ID && fav.TypeEndpoint == TypeEndpoint {
			favorites = append(favorites[:i], favorites[i+1:]...)
			break
		}
	}

	// Réécrire les favoris dans le fichier JSON
	data, err = json.Marshal(favorites)
	if err != nil {
		fmt.Println("Erreur lors de la conversion en JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = os.WriteFile("Favorite.json", data, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture dans le fichier:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Rediriger vers la page des favoris
	http.Redirect(w, r, "/ressourceTemp/favoriteTemp/", http.StatusSeeOther)
}
