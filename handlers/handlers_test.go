package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// Initialize an in-memory database for testing
func setup() *App {
	app := &App{}
	app.Initialize("sqlite3", ":memory:")
	return app
}

// Discard the in-memory database
func teardown(app *App) {
	app.DB.Close()
}

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		url    string
		alias  string
		status int
	}{
		"https":                    {url: "https://google.com", alias: "ggls", status: http.StatusCreated},
		"http":                     {url: "http://google.com", alias: "ggl", status: http.StatusCreated},
		"url already exists":       {url: "http://google.com", alias: "", status: http.StatusCreated},
		"alias already exists":     {url: "http://asdf.com", alias: "ggl", status: http.StatusCreated},
		"both url and alias exist": {url: "http://asdf.com", alias: "ggl", status: http.StatusCreated},
	}

	app := setup()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			origURL := &URLEntry{
				URL:   tc.url,
				Alias: tc.alias,
			}
			JSONData, _ := json.Marshal(origURL)
			// Set up a new request.
			req, err := http.NewRequest("POST", "/", bytes.NewBuffer(JSONData))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			http.HandlerFunc(app.CreateShortURL).ServeHTTP(rr, req)
			// Test that the status code is correct.
			if status := rr.Code; status != tc.status {
				t.Errorf("Status code is invalid. Expected %d. Got %d instead", tc.status, status)
			}
			// Alias might have been randomly generated so get from result instead
			temp := &JSONRes{}
			json.Unmarshal(rr.Body.Bytes(), temp)
			origURL.Alias = strings.TrimPrefix(temp.ShortURL, "/s/")
			// Test that the created url entry is correct.
			createdURL := URLEntry{}
			app.DB.Where("url = ?", origURL.URL).First(&createdURL)
			if createdURL != *origURL {
				t.Errorf("Created entry is invalid. Expected %+v. Got %+v instead", origURL, createdURL)
			}
		})
	}
	teardown(app)
}

func TestInvalidCreate(t *testing.T) {
	tests := map[string]struct {
		url    string
		alias  string
		status int
	}{
		"bad":           {url: "https/agoogle.com", alias: "ggls", status: http.StatusBadRequest},
		"no colon":      {url: "http//google.com", alias: "ggl", status: http.StatusBadRequest},
		"empty":         {url: "", alias: "", status: http.StatusBadRequest},
		"asdf":          {url: "asdf", alias: "", status: http.StatusBadRequest},
		"spaces":        {url: "https://google.com", alias: "oh no spaces", status: http.StatusBadRequest},
		"question mark": {url: "https://google.com", alias: "huh?", status: http.StatusBadRequest},
		"percent":       {url: "https://google.com", alias: "test%20stuff", status: http.StatusBadRequest},
	}

	app := setup()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			origURL := &URLEntry{
				URL:   tc.url,
				Alias: tc.alias,
			}
			JSONData, _ := json.Marshal(origURL)
			// Set up a new request.
			req, err := http.NewRequest("POST", "application/json", bytes.NewBuffer(JSONData))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			http.HandlerFunc(app.CreateShortURL).ServeHTTP(rr, req)
			// Test that the status code is correct.
			if status := rr.Code; status != tc.status {
				t.Errorf("Status code is invalid. Expected %d. Got %d instead", tc.status, status)
			}
			// Make sure no entries are created
			createdEntry := &URLEntry{}
			if !app.DB.First(createdEntry).RecordNotFound() {
				t.Errorf("Should not have created an entry")
			}
		})
	}
	teardown(app)
}

func TestCorruptCreate(t *testing.T) {
	app := setup()
	JSONData, _ := json.Marshal([]byte("garbage data"))
	// Set up a new request.
	req, err := http.NewRequest("POST", "application/json", bytes.NewBuffer(JSONData))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(app.CreateShortURL).ServeHTTP(rr, req)
	// Test that the status code is correct.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Status code is invalid. Expected %d. Got %d instead", http.StatusBadRequest, status)
	}
	// Make sure no entries are created
	createdEntry := &URLEntry{}
	if !app.DB.First(createdEntry).RecordNotFound() {
		t.Errorf("Should not have created an entry")
	}
	teardown(app)
}

func TestGetURL(t *testing.T) {
	seed := []struct {
		url   string
		alias string
	}{
		{url: "https://google.com", alias: "ggl"},
	}

	tests := map[string]struct {
		url    string
		alias  string
		status int
	}{
		"exists":      {url: "https://google.com", alias: "ggl", status: http.StatusFound},
		"nonexistent": {url: "", alias: "whoami", status: http.StatusNotFound},
	}

	app := setup()
	for _, u := range seed {
		origURL := &URLEntry{
			URL:   u.url,
			Alias: u.alias,
		}
		app.DB.Create(origURL)
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Set up a new request.
			req, err := http.NewRequest("GET", "/s/"+tc.alias, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			r := mux.NewRouter()
			r.HandleFunc("/s/{alias:.*}", app.GetURL).Methods("GET")
			r.ServeHTTP(rr, req)
			// Test that the status code is correct.
			if status := rr.Code; status != tc.status {
				t.Errorf("Status code is invalid. Expected %d. Got %d instead", tc.status, status)
			}
			if rr.HeaderMap.Get("Location") != tc.url {
				t.Errorf("Created entry is invalid. Expected %+v. Got %+v instead", tc.url, rr.HeaderMap.Get("Location"))
			}
		})
	}
	teardown(app)
}

func TestListAll(t *testing.T) {
	tests := []struct {
		url   string
		alias string
	}{
		{url: "https://google1.com", alias: "ggl1"},
		{url: "https://google2.com", alias: "ggl2"},
		{url: "https://google3.com", alias: "ggl3"},
	}

	app := setup()
	origList := make([]URLEntry, 0)
	for _, tc := range tests {
		origURL := &URLEntry{
			URL:   tc.url,
			Alias: tc.alias,
		}
		origList = append(origList, *origURL)
		app.DB.Create(origURL)
	}
	// Set up a new request.
	req, err := http.NewRequest("GET", "/all", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/all", app.ListAll).Methods("GET")
	r.ServeHTTP(rr, req)
	lst := &[]URLEntry{}
	json.Unmarshal(rr.Body.Bytes(), lst)
	// Test that the status code is correct.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code is invalid. Expected %d. Got %d instead", http.StatusOK, status)
	}
	// Test that the created url entries are listed correctly.
	for i, createdURL := range *lst {
		if createdURL != origList[i] {
			t.Errorf("Created entry is invalid. Expected %+v. Got %+v instead", origList[i], createdURL)
		}
	}
	teardown(app)
}

func TestFailDBBadDriver(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	app := &App{}
	app.Initialize("baddriver", ":memory:")
}

func TestFailDBBadURI(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	app := &App{}
	app.Initialize("sqlite3", "./garbagepath/data.db")
}
