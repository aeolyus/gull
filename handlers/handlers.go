package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/aeolyus/gull/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type App struct {
	DB *gorm.DB
}

type URLEntry struct {
	URL   string `gorm:"unique" json:"url"`
	Alias string `gorm:"unique" json:"alias"`
}

func (a *App) Initialize(dbDriver string, dbURI string) {
	// Setup database
	db, err := gorm.Open(dbDriver, dbURI)
	if err != nil {
		panic("Failed to connect to database")
	}
	a.DB = db
	// Schema Migration
	db.AutoMigrate(&URLEntry{})
}

func (a *App) ListAll(w http.ResponseWriter, r *http.Request) {
	var res []URLEntry
	a.DB.Find(&res)
	resJson, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}

func (a *App) GetURL(w http.ResponseWriter, r *http.Request) {
	var urlEntry URLEntry
	args := mux.Vars(r)
	a.DB.Where("alias = ?", args["alias"]).First(&urlEntry)
	url := urlEntry.URL
	if url != "" {
		http.Redirect(w, r, string(url), http.StatusFound)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "No such link")
	}
}

func (a *App) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal("Could not parse JSON")
	}
	reqURL := r.PostFormValue("url")
	alias := r.PostFormValue("alias")
	// Verify URL is valid
	if !utils.IsValidUrl(reqURL) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid URL")
		return
	}
	// Check if URL entry already exists
	var urlEntry URLEntry
	a.DB.Where("url = ?", reqURL).Find(&urlEntry)
	if urlEntry.URL != "" {
		alias = urlEntry.Alias
	} else {
		// Check if alias is already taken
		for alias == "" || !a.DB.Where("alias = ?", alias).First(&urlEntry).RecordNotFound() {
			alias = utils.RandString(6)
		}
		newURLEntry := &URLEntry{URL: reqURL, Alias: alias}
		a.DB.Create(newURLEntry)
	}

	// Write HTTP Response
	shortlink := r.Host + "/s/" + alias
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", shortlink)
	fmt.Fprintf(w, shortlink)
}
