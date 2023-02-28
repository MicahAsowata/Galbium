package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MicahAsowata/Galbium/internal/models"
	"github.com/albrow/forms"
	"github.com/dustin/go-humanize"
	"github.com/flosch/pongo2/v6"
	"github.com/julienschmidt/httprouter"
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
	path := fmt.Sprintf("/todo/view/%d", int(insertedID))
	http.Redirect(w, r, path, http.StatusSeeOther)
}
func (a *application) GetTodo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/view.gohtml"))
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		log.Fatal(err)
	}
	todoId := int64(id)

	todo, err := a.Queries.GetTodo(r.Context(), todoId)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.ExecuteWriter(pongo2.Context{"todo": todo, "created": humanize.Time(todo.Created)}, w)
	if err != nil {
		http.Error(w, "Error displaying page", http.StatusInternalServerError)
		return
	}
}

func (a *application) Index(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/todoIndex.gohtml"))

	todos, err := a.Queries.ListTodo(r.Context())
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.ExecuteWriter(pongo2.Context{"todos": todos}, w)
	if err != nil {
		http.Error(w, "Error displaying page", http.StatusInternalServerError)
		return
	}
}
