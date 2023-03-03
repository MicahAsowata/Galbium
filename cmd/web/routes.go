package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *application) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.CleanPath)
	router.Use(middleware.Recoverer)
	router.Use(a.SessionManager.LoadAndSave)
	router.NotFound(a.NotFound)
	router.Get("/", a.Home)
	router.Get("/todo/new", a.NewTodo)
	router.Post("/todo/create", a.CreateTodo)
	router.Get("/todo/view/{id}", a.GetTodo)
	router.Get("/todo", a.Index)
	router.Get("/todo/edit/{id}", a.EditTodo)
	router.Post("/todo/update/{id}", a.UpdateTodo)
	router.Get("/todo/delete/{id}", a.DeleteTodo)
	router.Get("/user/signup", a.SignUpUser)
	router.Post("/user/signup", a.SignUpUserPost)
	router.Get("/user/login", a.LoginUser)
	router.Post("/user/login", a.LoginUserPost)
	router.Post("/user/logout", a.LogoutUser)
	return router
}
