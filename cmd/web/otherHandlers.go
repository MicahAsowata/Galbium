package main

import (
	"net/http"

	"github.com/flosch/pongo2/v6"
)

func (a *application) NotFound(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/notFound.gohtml"))
	err := tmpl.ExecuteWriter(nil, w)
	if err != nil {
		http.Error(w, "Error displaying page", http.StatusInternalServerError)
	}
}
