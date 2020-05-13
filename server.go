package main

import (
	"github.com/aeolyus/gull/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	a := &handlers.App{}
	os.MkdirAll("./data", os.ModePerm)
	a.Initialize("sqlite3", "./data/data.db")
	defer a.DB.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", a.CreateShortURL).Methods("POST")
	r.HandleFunc("/all", a.ListAll).Methods("GET")
	r.HandleFunc("/s/{alias:.*}", a.GetURL).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/"))).Methods("GET")
	http.Handle("/", r)

	log.Println("Server is listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
