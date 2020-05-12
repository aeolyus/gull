package main

import (
	"./handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	a := &handlers.App{}
	a.Initialize("sqlite3", "./data.db")
	defer a.DB.Close()

	r := mux.NewRouter()
	r.Handle("/", http.FileServer(http.Dir("./static"))).Methods("GET")
	r.HandleFunc("/all", a.ListAll).Methods("GET")
	http.Handle("/", r)

	log.Println("Server is listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
