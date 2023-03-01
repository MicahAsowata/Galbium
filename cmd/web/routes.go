package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *application) routes() http.Handler {
	mux := chi.NewRouter()
	//TODO: Setup proper routing system
	//TODO: Setup all the routes
	mux.Get("/", a.Home)
	mux.Get("/todo/new", a.NewTodo)
	mux.Post("/todo/create", a.CreateTodo)
	mux.Get("/todo/view/{id}", a.GetTodo)
	mux.Get("/todo", a.Index)
	mux.Get("/todo/edit/{id}", a.EditTodo)
	mux.Post("/todo/update/{id}", a.UpdateTodo)
	mux.Get("/todo/delete/{id}", a.DeleteTodo)
	return mux
}
