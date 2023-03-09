package main

import (
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *application) routes() http.Handler {
	router := flow.New()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(a.SecureHeaders)
	router.Use(a.SessionManager.LoadAndSave)
	router.HandleFunc("/", a.Home, http.MethodGet)
	router.HandleFunc("/user/signup", a.SignUpUser, http.MethodGet)
	router.HandleFunc("/user/signup", a.SignUpUserPost, http.MethodPost)
	router.HandleFunc("/user/login", a.LoginUser, http.MethodGet)
	router.HandleFunc("/user/login", a.LoginUserPost, http.MethodPost)
	router.HandleFunc("/user/forgot_password", a.ForgotPassword, http.MethodGet)
	router.HandleFunc("/user/reset_password", a.ResetPassword, http.MethodPost)
	router.Group(func(m *flow.Mux) {
		router.Use(a.RequireAuth)
		router.HandleFunc("/todo", a.Index, http.MethodGet)
		router.HandleFunc("/todo/new", a.NewTodo, http.MethodGet)
		router.HandleFunc("/todo/create", a.CreateTodo, http.MethodPost)
		router.HandleFunc("/todo/view/:id", a.GetTodo, http.MethodGet)
		router.HandleFunc("/todo/edit/:id", a.EditTodo, http.MethodGet)
		router.HandleFunc("/todo/update/:id", a.UpdateTodo, http.MethodPost)
		router.HandleFunc("/todo/delete/:id", a.DeleteTodo, http.MethodGet)
		router.HandleFunc("/user/logout", a.LogoutUser, http.MethodGet)
	})
	return router
}
