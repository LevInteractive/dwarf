package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/psaia/dwarf/storage"
)

// CreatePayload represents both the incoming and outgoing payload of URL's.
type CreatePayload struct {
	Urls []string `json:"urls"`
}

// LoadResult represents what's returned from the Load method.
type LoadResult struct {
	LongURL string `json:"longUrl"`
}

// CreateHandler for adding a new URL to the store.
func CreateHandler(store storage.IStorage, baseURL string) func(http.ResponseWriter, *http.Request) {
	fnc := func(w http.ResponseWriter, r *http.Request) {
		var outgoing CreatePayload
		var incoming CreatePayload

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&incoming)

		if err != nil || len(incoming.Urls) < 1 {
			log.Printf("CreateHandler.error: recieved bad payload: %v", err)
			http.Error(w, "Malformed payload", http.StatusBadRequest)
			return
		}

		// Validate incoming urls.
		for _, u := range incoming.Urls {
			_, err := url.ParseRequestURI(u)
			if err != nil {
				log.Printf("CreateHandler.error: recieved bad payload: %v", err)
				http.Error(w, "Invalid URL passed - no urls shortened.", http.StatusBadRequest)
				return
			}
		}

		// Shorten them.
		for _, u := range incoming.Urls {
			code, err := store.Save(u)
			if err != nil {
				log.Printf("CreateHandler.error: failed to shorten url: %v", err)
				http.Error(w, "Unknown error occurred. Please try again.", http.StatusBadRequest)
				return
			}
			outgoing.Urls = append(
				outgoing.Urls,
				fmt.Sprintf("%s/%s", baseURL, code),
			)
		}

		json.NewEncoder(w).Encode(outgoing)
	}
	return fnc
}

// LookupHandler for adding a new URL to the store.
func LookupHandler(store storage.IStorage, baseURL string) func(http.ResponseWriter, *http.Request) {
	fnc := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		code := vars["hash"]
		fullURL, err := store.Load(code)

		if err == storage.ErrNotFound {
			log.Printf("LookupHander.warn: could not find url with code %s", code)
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		if err != nil {
			log.Printf("LookupHandler.error: receieved error while looking up url: %v", err)
			http.Error(w, "Unknown error while looking up full url", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, fullURL, 301)
	}
	return fnc
}
