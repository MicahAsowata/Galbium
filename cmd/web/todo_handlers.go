package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MicahAsowata/Galbium/internal/models"
	"github.com/albrow/forms"
	"github.com/alexedwards/flow"
	"github.com/dustin/go-humanize"
	"github.com/flosch/pongo2/v6"
)

// TODO: Setup the basic handlers
func (a *application) Home(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/home.gohtml"))
	err := tmpl.ExecuteWriter(nil, w)
	if err != nil {
		http.Error(w, "Error displaying page", http.StatusInternalServerError)
	}
}

func (a *application) NewTodo(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/create_todo.gohtml"))
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
	results, err := a.Queries.CreateTodo(r.Context(), models.CreateTodoParams{
		Name:      todoData.Get("name"),
		Details:   todoData.Get("details"),
		Completed: true,
	})

	if err != nil {
		log.Fatal(err)
	}

	insertedID, err := results.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	path := fmt.Sprintf("/todo/view/%d", int(insertedID))
	a.SessionManager.Put(r.Context(), "flash", "todo created successfully")
	http.Redirect(w, r, path, http.StatusSeeOther)
}
func (a *application) GetTodo(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/view_todo.gohtml"))
	id, err := strconv.Atoi(flow.Param(r.Context(), "id"))
	if err != nil {
		log.Fatal(err)
	}
	todoId := int64(id)

	todo, err := a.Queries.GetTodo(r.Context(), todoId)
	if err != nil {
		log.Fatal(err)
	}
	flash := a.SessionManager.PopString(r.Context(), "flash")
	err = tmpl.ExecuteWriter(pongo2.Context{"todo": todo, "created": humanize.Time(todo.Created.UTC()), "flash": flash}, w)
	if err != nil {
		http.Error(w, "Error displaying page", http.StatusInternalServerError)
		return
	}
}

func (a *application) Index(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/todo_index.gohtml"))

	todos, err := a.Queries.ListTodo(r.Context())
	if err != nil {
		log.Fatal(err)
	}
	userID := a.SessionManager.GetInt(r.Context(), "userID")
	var username string
	if userID != 0 {
		getusername, err := a.Users.Get(r.Context(), userID)
		if err != nil {
			log.Fatal(err)
		}

		username = getusername
	}
	err = tmpl.ExecuteWriter(pongo2.Context{"todos": todos, "username": username}, w)
	if err != nil {
		http.Error(w, "Error displaying page", http.StatusInternalServerError)
		return
	}
}
func (a *application) EditTodo(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/edit_todo.gohtml"))
	id, err := strconv.Atoi(flow.Param(r.Context(), "id"))
	if err != nil {
		log.Fatal(err)
	}
	todoID := int64(id)
	todo, err := a.Queries.GetTodo(r.Context(), todoID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Fatal("User not found")
	}
	err = tmpl.ExecuteWriter(pongo2.Context{"todo": todo}, w)
	if err != nil {
		log.Fatal(err)
	}
}
func (a *application) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(flow.Param(r.Context(), "id"))
	if err != nil {
		log.Fatal(err)
	}

	todoId := int64(id)

	updateTodoData, err := forms.Parse(r)
	if err != nil {
		log.Fatal(err)
	}
	validator := updateTodoData.Validator()
	validator.Require("name")

	var completedAsBool bool
	if updateTodoData.Get("completed") == "on" {
		completedAsBool = true
	} else {
		completedAsBool = false
	}

	err = a.Queries.UpdateTodo(r.Context(), models.UpdateTodoParams{
		Name:      updateTodoData.Get("name"),
		Details:   updateTodoData.Get("details"),
		Completed: completedAsBool,
		ID:        todoId,
	})

	if err != nil {
		log.Fatal(err)
	}
	path := fmt.Sprintf("/todo/view/%d", id)
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func (a *application) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(flow.Param(r.Context(), "id"))
	if err != nil {
		log.Fatal(err)
	}

	todoID := int64(id)

	err = a.Queries.DeleteTodo(r.Context(), todoID)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}
