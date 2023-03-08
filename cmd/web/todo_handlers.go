package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/MicahAsowata/Galbium/internal/models"
	"github.com/albrow/forms"
	"github.com/alexedwards/flow"
	"github.com/dustin/go-humanize"
	"github.com/flosch/pongo2/v6"
)

func (a *application) Home(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/home.gohtml"))
	isAuthenticated := a.IsAuthenticated(r)
	err := tmpl.ExecuteWriter(pongo2.Context{"isAuthenticated": isAuthenticated}, w)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (a *application) NewTodo(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/create_todo.gohtml"))
	isAuthenticated := a.IsAuthenticated(r)
	err := tmpl.ExecuteWriter(pongo2.Context{"isAuthenticated": isAuthenticated}, w)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
func (a *application) CreateTodo(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := a.IsAuthenticated(r)
	todoData, err := forms.Parse(r)
	if err != nil {
		a.Logger.Error(err.Error())
	}
	validator := todoData.Validator()
	validator.Require("name")
	validator.MaxLength("name", 280)
	validator.Require("details")
	validator.MaxLength("details", 10000)
	if validator.HasErrors() {
		tmpl := pongo2.Must(pongo2.FromFile("./templates/create_todo.gohtml"))
		errorMap := validator.ErrorMap()
		var nameFieldError string
		var detailFieldError string
		if len(errorMap["name"]) > 0 {
			nameFieldError = "Invalid task name"
		} else {
			nameFieldError = ""
		}

		if len(errorMap["details"]) > 0 {
			detailFieldError = "Invalid task details"
		} else {
			detailFieldError = ""
		}
		nameFieldData := todoData.Get("name")
		detailFieldData := todoData.Get("detail")
		err := tmpl.ExecuteWriter(pongo2.Context{"isAuthenticated": isAuthenticated, "nameFieldError": nameFieldError, "detailFieldError": detailFieldError, "nameFieldData": nameFieldData, "detailFieldData": detailFieldData}, w)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}
	results, err := a.Queries.CreateTodo(r.Context(), models.CreateTodoParams{
		Name:      todoData.Get("name"),
		Details:   todoData.Get("details"),
		Completed: false,
		UserID:    a.SessionManager.GetInt(r.Context(), "userID"),
	})

	if err != nil {
		tmpl := pongo2.Must(pongo2.FromFile("./templates/create_todo.gohtml"))
		createTodoError := "Task not created"
		err := tmpl.ExecuteWriter(pongo2.Context{"isAuthenticated": isAuthenticated, "createTodoError": createTodoError}, w)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}

	insertedID, err := results.LastInsertId()
	if err != nil {
		a.Logger.Error(err.Error())
		return
	}
	path := fmt.Sprintf("/todo/view/%d", int(insertedID))
	a.SessionManager.Put(r.Context(), "flash", "Task created successfully")
	http.Redirect(w, r, path, http.StatusSeeOther)
}
func (a *application) GetTodo(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/view_todo.gohtml"))
	isAuthenticated := a.IsAuthenticated(r)
	userID := a.SessionManager.GetInt(r.Context(), "userID")
	id, err := strconv.Atoi(flow.Param(r.Context(), "id"))
	if err != nil {
		a.Logger.Error(err.Error())
		return
	}

	todo, err := a.Queries.GetTodo(r.Context(), id, userID)
	if err != nil {
		tmpl := pongo2.Must(pongo2.FromFile("./templates/not_found.gohtml"))
		err := tmpl.ExecuteWriter(pongo2.Context{"isAuthenticated": isAuthenticated}, w)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	}
	flash := a.SessionManager.PopString(r.Context(), "flash")
	err = tmpl.ExecuteWriter(pongo2.Context{"isAuthenticated": isAuthenticated, "todo": todo, "created": humanize.Time(todo.Created), "flash": flash}, w)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (a *application) Index(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/todo_index.gohtml"))
	userID := a.SessionManager.GetInt(r.Context(), "userID")

	todos, err := a.Queries.ListTodo(r.Context(), userID)
	if err != nil {
		a.Logger.Error(err.Error())
	}
	var username string
	if userID != 0 {
		getusername, err := a.Users.Get(r.Context(), userID)
		if err != nil {
			a.Logger.Error(err.Error())
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
	userID := a.SessionManager.GetInt(r.Context(), "userID")
	id, err := strconv.Atoi(flow.Param(r.Context(), "id"))
	if err != nil {
		a.Logger.Error(err.Error())
	}
	todo, err := a.Queries.GetTodo(r.Context(), id, userID)
	if errors.Is(err, sql.ErrNoRows) {
		a.Logger.Error(err.Error())
	}
	err = tmpl.ExecuteWriter(pongo2.Context{"todo": todo}, w)
	if err != nil {
		a.Logger.Error(err.Error())
	}
}
func (a *application) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(flow.Param(r.Context(), "id"))
	userID := a.SessionManager.GetInt(r.Context(), "userID")
	if err != nil {
		a.Logger.Error(err.Error())
	}

	updateTodoData, err := forms.Parse(r)
	if err != nil {
		a.Logger.Error(err.Error())
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
		ID:        id,
		UserID:    userID,
	})

	if err != nil {
		a.Logger.Error(err.Error())
	}
	path := fmt.Sprintf("/todo/view/%d", id)
	http.Redirect(w, r, path, http.StatusSeeOther)
}

func (a *application) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(flow.Param(r.Context(), "id"))
	userID := a.SessionManager.GetInt(r.Context(), "userID")
	if err != nil {
		a.Logger.Error(err.Error())
	}

	err = a.Queries.DeleteTodo(r.Context(), id, userID)
	if err != nil {
		a.Logger.Error(err.Error())
	}
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}
