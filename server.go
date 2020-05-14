package main

import (
	"github.com/aeolyus/gull/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	PUBLIC_DIR string = "./public"
	PORT       int    = 8081
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
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(PUBLIC_DIR))).Methods("GET")
	http.Handle("/", r)

	log.Println("Server is listening on port", PORT)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(PORT), nil))
}
