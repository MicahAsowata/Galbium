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
	router.Use(a.SessionManager.LoadAndSave)
	router.HandleFunc("/", a.Home, "GET")
	router.HandleFunc("/user/signup", a.SignUpUser, "GET")
	router.HandleFunc("/user/signup", a.SignUpUserPost, "POST")
	router.HandleFunc("/user/login", a.LoginUser, "GET")
	router.HandleFunc("/user/login", a.LoginUserPost, "POST")
	router.Group(func(m *flow.Mux) {
		router.Use(a.RequireAuth)
		router.HandleFunc("/todo", a.Index, "GET")
		router.HandleFunc("/todo/new", a.NewTodo, "GET")
		router.HandleFunc("/todo/create", a.CreateTodo, "POST")
		router.HandleFunc("/todo/view/:id", a.GetTodo, "GET")
		router.HandleFunc("/todo/edit/:id", a.EditTodo, "GET")
		router.HandleFunc("/todo/update/:id", a.UpdateTodo, "POST")
		router.HandleFunc("/todo/delete/:id", a.DeleteTodo, "GET")
		router.HandleFunc("/user/logout", a.LogoutUser, "POST")
	})

	return router
}
