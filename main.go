package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {

	css := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", css))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/collection/character", CollectionCharacter)
	http.HandleFunc("/submit_sort", FilterCharacter)
	http.HandleFunc("/submit_sort_location", FilterLocation)
	http.HandleFunc("/collection/episodes", CollectionEpisodes)
	http.HandleFunc("/collection/location", CollectionLocation)
	http.HandleFunc("/ressource/character", RessourceCharacter)
	http.HandleFunc("/ressource/episode", RessourceEpisodes)
	http.HandleFunc("/ressource/location", RessourceLocation)
	http.HandleFunc("/ressource/favorite", RessourceFavorite)
	http.HandleFunc("/ressourceTemp/favoriteTemp/", FavoriteTemp)
	http.HandleFunc("/ressource/about", AboutTemp)
	http.HandleFunc("/search", SearchHandler)
	http.HandleFunc("/favorite/delete", DeleteFavorite)
	//empty et jaune

	fmt.Println("Serveur lancé sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "Home", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
