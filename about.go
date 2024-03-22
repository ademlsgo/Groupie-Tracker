package main

import "net/http"
import "html/template"

func AboutTemp(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./templates/about.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "AboutTemp", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
