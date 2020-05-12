package handlers

import (
	"encoding/json"
	"net/http"

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
