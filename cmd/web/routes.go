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
	mux.HandlerFunc(http.MethodGet, "/todo", a.Todo)

	return mux
}
