package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MicahAsowata/Galbium/internal/models"
	"github.com/flosch/pongo2/v6"
)

// TODO: Setup the basic handlers
var tmpl = pongo2.Must(pongo2.FromFile("./templates/home.gohtml"))

func (a *application) Home(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteWriter(nil, w)
	if err != nil {
		http.Error(w, "Error displaying page", http.StatusInternalServerError)
	}
}

func (a *application) Todo(w http.ResponseWriter, r *http.Request) {
	name := "Eat"
	details := "Eat whatever is cooked"
	completed := true

	results, err := a.Queries.CreateTodo(r.Context(), models.CreateTodoParams{
		Name:      name,
		Details:   details,
		Completed: completed,
	})

	if err != nil {
		log.Fatal(err)
	}

	insertedID, err := results.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Created todo %d", int(insertedID))
}
