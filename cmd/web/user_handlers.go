package main

import (
	"fmt"
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
		a.Logger.Error(err.Error())
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
		a.Logger.Error(err.Error())
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
	loginData, err := forms.Parse(r)
	if err != nil {
		a.Logger.Error(err.Error())
	}

	validator := loginData.Validator()
	validator.Require("email")
	validator.MatchEmail("email")
	validator.MaxLength("email", 280)
	validator.Require("password")
	validator.LengthRange("password", 8, 280)

	if validator.HasErrors() {
		tmpl := pongo2.Must(pongo2.FromFile("./templates/login.gohtml"))
		errorMap := validator.ErrorMap()
		var emailFieldErrors string
		if len(errorMap["email"]) > 0 {
			emailFieldErrors = "Your email seems invalid"
		} else {
			emailFieldErrors = ""
		}
		var passwordFieldErrors string
		if len(errorMap["password"]) > 0 {
			passwordFieldErrors = "Your password must be be more than 8 characters"
		} else {
			passwordFieldErrors = ""
		}
		err := tmpl.ExecuteWriter(pongo2.Context{"loggedin": a.IsAuthenticated(r), "emailFieldError": emailFieldErrors, "passwordFieldError": passwordFieldErrors, "emailFieldData": loginData.Get("email")}, w)
		if err != nil {
			http.Error(w, "could not display page", http.StatusInternalServerError)
			return
		}
		return
	}
	user_id, err := a.Users.Authenticate(r.Context(), models.AuthUserParams{
		Email:    loginData.Get("email"),
		Password: loginData.Get("password"),
	})

	if err != nil {
		tmpl := pongo2.Must(pongo2.FromFile("./templates/login.gohtml"))
		err := tmpl.ExecuteWriter(pongo2.Context{"loginError": "sorry, we could not log you in", "loggedin": a.IsAuthenticated(r)}, w)
		if err != nil {
			http.Error(w, "could not display page", http.StatusInternalServerError)
			return
		}
		return
	}
	a.SessionManager.RenewToken(r.Context())
	a.SessionManager.Put(r.Context(), "userID", user_id)
	a.SessionManager.RememberMe(r.Context(), true)
	http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
}

func (a *application) LogoutUser(w http.ResponseWriter, r *http.Request) {
	err := a.SessionManager.RenewToken(r.Context())
	if err != nil {
		a.Logger.Error(err.Error())
	}

	a.SessionManager.Remove(r.Context(), "userID")

	a.SessionManager.RenewToken(r.Context())

	a.SessionManager.Put(r.Context(), "flash", "logged out successfully")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (a *application) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("./templates/forgot_password.gohtml"))

	err := tmpl.ExecuteWriter(nil, w)
	if err != nil {
		http.Error(w, "could not display page", http.StatusInternalServerError)
		return
	}
}

func (a *application) ResetPassword(w http.ResponseWriter, r *http.Request) {
	resetData, err := forms.Parse(r)
	if err != nil {
		a.Logger.Error(err.Error())
	}

	validator := resetData.Validator()
	validator.Require("email")
	validator.MatchEmail("email")
	validator.Require("password")
	validator.LengthRange("password", 8, 280)
	validator.Require("confirm_password")
	validator.Equal("confirm_password", "password")

	if validator.HasErrors() {
		tmpl := pongo2.Must(pongo2.FromFile("./templates/forgot_password.gohtml"))
		errorMap := validator.ErrorMap()
		var emailFieldErrors string
		if len(errorMap["email"]) > 0 {
			emailFieldErrors = "Your email seems invalid"
		} else {
			emailFieldErrors = ""
		}
		var passwordFieldErrors string
		if len(errorMap["password"]) > 0 {
			passwordFieldErrors = "Your password must be be more than 8 characters"
		} else {
			passwordFieldErrors = ""
		}
		var confirmPasswordFieldErrors string
		if len(errorMap["confirm_password"]) > 0 {
			confirmPasswordFieldErrors = "They are not equal"
		} else {
			confirmPasswordFieldErrors = ""
		}

		err := tmpl.ExecuteWriter(pongo2.Context{"emailFieldErrors": emailFieldErrors, "passwordFieldErrors": passwordFieldErrors, "confirmPasswordFieldErrors": confirmPasswordFieldErrors, "emailFieldData": resetData.Get("email")}, w)
		if err != nil {
			http.Error(w, "could not display page", http.StatusInternalServerError)
			return
		}
		return
	}
	err = a.Users.ResetPassword(r.Context(), models.ResetPasswordParams{
		Email:    resetData.Get("email"),
		Password: resetData.Get("password"),
	})

	if err != nil {
		tmpl := pongo2.Must(pongo2.FromFile("./templates/forgot_password.gohtml"))
		err := tmpl.ExecuteWriter(pongo2.Context{"notUpdated": "password could not be reset"}, w)
		if err != nil {
			http.Error(w, "could not display page", http.StatusInternalServerError)
			return
		}
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
