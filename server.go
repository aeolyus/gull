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
	publicDir          string = "./public"
	port               int    = 8081
	envFlagAllowCreate string = "GULL_ALLOW_CREATE"
)

var (
	allowCreateEnv      = os.Getenv(envFlagAllowCreate)
	allowCreate    bool = allowCreateEnv == "true" || allowCreateEnv == ""
)

func main() {
	a := &handlers.App{}
	os.MkdirAll("./data", os.ModePerm)
	a.Initialize("sqlite3", "./data/data.db")
	defer a.DB.Close()

	r := mux.NewRouter()
	if allowCreate {
		r.HandleFunc("/", a.CreateShortURL).Methods("POST")
	}
	r.HandleFunc("/all", a.ListAll).Methods("GET")
	r.HandleFunc("/s/{alias:.*}", a.GetURL).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(publicDir))).Methods("GET")
	http.Handle("/", r)

	log.Println("Server is listening on port", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
