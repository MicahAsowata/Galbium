package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MicahAsowata/Galbium/internal/models"
	"github.com/albrow/forms"
	"github.com/flosch/pongo2/v6"
)

func (a *application) SignUpUser(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/signup.gohtml"))
	err := tmpl.ExecuteWriter(nil, w)
	if err != nil {
		http.Error(w, "could not display page", http.StatusInternalServerError)
		return
	}
}

func (a *application) SignUpUserPost(w http.ResponseWriter, r *http.Request) {
	todoData, err := forms.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	validator := todoData.Validator()
	validator.Require("name")
	validator.MaxLength("name", 280)
	validator.Require("email")
	validator.MaxLength("email", 280)
	validator.MatchEmail("email")
	validator.Require("username")
	validator.MaxLength("username", 280)
	validator.Require("password")
	validator.LengthRange("password", 8, 280)

	if validator.HasErrors() {
		fmt.Fprintln(w, "Invalid data")
		return
	}

	err = a.Users.Insert(r.Context(), models.InsertUsersParams{
		Name:     todoData.Get("name"),
		Email:    todoData.Get("email"),
		Username: todoData.Get("username"),
		Password: todoData.Get("password"),
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func (a *application) LoginUser(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/login.gohtml"))
	err := tmpl.ExecuteWriter(pongo2.Context{"loggedin": a.IsAuthenticated(r)}, w)
	if err != nil {
		http.Error(w, "could not display page", http.StatusInternalServerError)
		return
	}
}

func (a *application) LoginUserPost(w http.ResponseWriter, r *http.Request) {
	todoData, err := forms.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	validator := todoData.Validator()
	validator.Require("email")
	validator.MatchEmail("email")
	validator.MaxLength("email", 280)
	validator.Require("password")
	validator.LengthRange("password", 8, 280)

	if validator.HasErrors() {
		fmt.Fprintln(w, "Invalid data")
		return
	}
	id, err := a.Users.Authenticate(r.Context(), models.AuthUserParams{
		Email:    todoData.Get("email"),
		Password: todoData.Get("password"),
	})

	if err != nil {
		log.Fatal(err)
	}
	a.SessionManager.RenewToken(r.Context())
	a.SessionManager.Put(r.Context(), "userID", id)
	a.SessionManager.RememberMe(r.Context(), true)
	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}

func (a *application) LogoutUser(w http.ResponseWriter, r *http.Request) {
	err := a.SessionManager.RenewToken(r.Context())
	if err != nil {
		log.Fatal(err)
	}

	a.SessionManager.Remove(r.Context(), "userID")

	a.SessionManager.RenewToken(r.Context())

	a.SessionManager.Put(r.Context(), "flash", "logged out successfully")

	http.Redirect(w, r, "/todo", http.StatusSeeOther)
}
