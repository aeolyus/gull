package handlers

import (
	"encoding/json"
	"github.com/aeolyus/gull/utils"
	"net/http"

	valid "github.com/asaskevich/govalidator"
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

type JSONRes struct {
	ShortURL string `json:"shorturl"`
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
	res := &[]URLEntry{}
	a.DB.Find(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (a *App) GetURL(w http.ResponseWriter, r *http.Request) {
	u := &URLEntry{}
	args := mux.Vars(r)
	a.DB.Where("alias = ?", args["alias"]).First(u)
	if u.URL != "" {
		http.Redirect(w, r, string(u.URL), http.StatusFound)
	} else {
		http.Error(w, "No such link :(", http.StatusNotFound)
	}
}

func (a *App) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	u := &URLEntry{}
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Verify URL is valid
	if !valid.IsRequestURL(u.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	// Verify alias is valid
	if u.Alias != "" && !utils.IsValidAlias(u.Alias) {
		http.Error(w, "Invalid Alias", http.StatusBadRequest)
		return
	}
	// Check if URL entry already exists
	existingURL := &URLEntry{}
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
	json.NewEncoder(w).Encode(&JSONRes{ShortURL: shortlink})
}
