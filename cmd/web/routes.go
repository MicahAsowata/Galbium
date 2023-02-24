package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *application) routes() http.Handler {
	mux := httprouter.New()
	//TODO: Setup proper routing system
	//TODO: Setup all the routes
	mux.GET("/", a.Home)

	return mux
}
