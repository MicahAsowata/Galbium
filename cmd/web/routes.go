package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *application) routes() http.Handler {
	mux := httprouter.New()
	//TODO: Setup proper routing system
	//TODO: Setup all the routes
	mux.HandlerFunc(http.MethodGet, "/", a.Home)
	mux.HandlerFunc(http.MethodGet, "/todo/new", a.NewTodo)
	mux.HandlerFunc(http.MethodPost, "/todo/create", a.CreateTodo)
	mux.GET("/todo/view/:id", a.GetTodo)
	mux.HandlerFunc(http.MethodGet, "/todo", a.Index)
	mux.GET("/todo/edit/:id", a.EditTodo)
	mux.POST("/todo/update/:id", a.UpdateTodo)
	return mux
}
