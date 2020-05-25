package handlers

import (
	"encoding/json"
	"github.com/aeolyus/gull/utils"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	// Use the sqlite db driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// App represents the appplication itself
type App struct {
	DB *gorm.DB
}

type urlEntry struct {
	URL   string `gorm:"unique" json:"url"`
	Alias string `gorm:"unique" json:"alias"`
}

type jsonRes struct {
	ShortURL string `json:"shorturl"`
}

// Initialize initializes the app, connects to database, and does auto migration
func (a *App) Initialize(dbDriver string, dbURI string) {
	// Setup database
	db, err := gorm.Open(dbDriver, dbURI)
	if err != nil {
		panic("Failed to connect to database")
	}
	a.DB = db
	// Schema Migration
	db.AutoMigrate(&urlEntry{})
}

// ListAll returns a JSON response of all the url entries in the database
func (a *App) ListAll(w http.ResponseWriter, r *http.Request) {
	res := &[]urlEntry{}
	a.DB.Find(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// GetURL takes a shortlink and redirects the user to the relevant link if it exists
func (a *App) GetURL(w http.ResponseWriter, r *http.Request) {
	u := &urlEntry{}
	args := mux.Vars(r)
	a.DB.Where("alias = ?", args["alias"]).First(u)
	if u.URL != "" {
		http.Redirect(w, r, string(u.URL), http.StatusFound)
	} else {
		http.Error(w, "No such link :(", http.StatusNotFound)
	}
}

// CreateShortURL will take a url and shorten it with an (custom) alias
func (a *App) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	u := &urlEntry{}
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Verify URL is valid
	if !utils.IsValidURL(u.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	// Verify alias is valid
	if u.Alias != "" && !utils.IsValidAlias(u.Alias) {
		http.Error(w, "Invalid Alias", http.StatusBadRequest)
		return
	}
	// Check if URL entry already exists
	existingURL := &urlEntry{}
	a.DB.Where("url = ?", u.URL).Find(existingURL)
	if existingURL.URL != "" {
		u.Alias = existingURL.Alias
	} else {
		// Verify alias is unique
		for u.Alias == "" || !a.DB.Where("alias = ?", u.Alias).First(existingURL).RecordNotFound() {
			u.Alias = utils.RandString(6)
		}
		a.DB.Create(u)
	}
	// Write HTTP Response
	shortlink := r.Host + "/s/" + u.Alias
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", shortlink)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&jsonRes{ShortURL: shortlink})
}
