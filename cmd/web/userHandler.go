package main

import (
	"fmt"
	"log"
	"net/http"

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

}

func (a *application) LoginUser(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/login.gohtml"))
	err := tmpl.ExecuteWriter(nil, w)
	if err != nil {
		http.Error(w, "could not display page", http.StatusInternalServerError)
		return
	}
}

func (a *application) LoginUserPost(w http.ResponseWriter, r *http.Request) {

}

func (a *application) LogoutUser(w http.ResponseWriter, r *http.Request) {

}
