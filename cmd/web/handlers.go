package main

import (
	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/julienschmidt/httprouter"
)

// TODO: Setup the basic handlers
var tmpl = pongo2.Must(pongo2.FromFile("./templates/home.gohtml"))

func (a *application) Home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := tmpl.ExecuteWriter(nil, w)
	if err != nil {
		http.Error(w, "Error displaying page", http.StatusInternalServerError)
	}
}
