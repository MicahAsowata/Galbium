package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MicahAsowata/Galbium/internal/models"
	"github.com/albrow/forms"
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

func (a *application) NewTodo(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/createTodo.gohtml"))
	err := tmpl.ExecuteWriter(nil, w)
	if err != nil {
		http.Error(w, "Error displaying page", http.StatusInternalServerError)
	}

}
func (a *application) CreateTodo(w http.ResponseWriter, r *http.Request) {
	todoData, err := forms.Parse(r)
	if err != nil {
		log.Fatal("Could not parse form data")
	}
	validator := todoData.Validator()
	validator.Require("name")
	if validator.HasErrors() {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}

	var completedAsBool bool
	if todoData.Get("completed") == "on" {
		completedAsBool = true
	} else {
		completedAsBool = false
	}
	results, err := a.Queries.CreateTodo(r.Context(), models.CreateTodoParams{
		Name:      todoData.Get("name"),
		Details:   todoData.Get("details"),
		Completed: completedAsBool,
	})

	if err != nil {
		log.Fatal(err)
	}

	insertedID, err := results.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	// path := fmt.Sprintf("/todo/%d", int(insertedID))
	// http.Redirect(w, r, "/todo/view", http.StatusSeeOther)
	// fmt.Fprintln(w, path)
	fmt.Fprintf(w, "Created todo %d", int(insertedID))
}
func (a *application) GetTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Fatal(err)
	}
	todoId := int64(id)

	todo, err := a.Queries.GetTodo(r.Context(), todoId)
	if err != nil {
		http.NotFound(w, r)
		return
	}

}
