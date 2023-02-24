package main

import (
	"log"
	"net/http"
)

type application struct{}

func main() {
	a := application{}
	log.Println("Starting server for http://localhost:4000")
	err := http.ListenAndServe(":4000", a.routes())
	if err != nil {
		log.Fatal(err)
	}
}
