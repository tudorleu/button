package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/tudorleu/button/api"
	"github.com/tudorleu/button/models"
	"log"
	"net/http"
	"os"
)

var dbUrl = flag.String("dbUrl",
	"postgres://tudor@localhost/button?sslmode=disable",
	"Url for the postgres db.")

func main() {
	flag.Parse()
	models.InitDb(*dbUrl, false)
	defer models.CloseDb()

	r := mux.NewRouter()

	post := r.Methods("POST").Subrouter()
	post.HandleFunc("/user", api.WithContext(api.NewUser))
	post.HandleFunc("/user/{userId:[0-9]+}/transfer",
		api.WithContext(api.NewTransfer))

	get := r.Methods("GET").Subrouter()
	get.HandleFunc("/user/{userId:[0-9]+}", api.WithContext(api.GetUser))
	get.HandleFunc("/user/{userId:[0-9]+}/transfers",
		api.WithContext(api.GetTransfers))

	err := http.ListenAndServe(os.Getenv("PORT"), r)
	if err != nil {
		log.Fatal(err)
	}
}
