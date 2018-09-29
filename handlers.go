package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/LevInteractive/dwarf/pb"
	"github.com/LevInteractive/dwarf/storage"
	"github.com/gorilla/mux"
)

// LoadResult represents what's returned from the Load method.
type LoadResult struct {
	LongURL string `json:"longUrls"`
}

// CreateServer for the gRPC server.
type CreateServer struct {
	Store storage.IStorage
}

// Create new short url given a set of long urls via gRPC.
func (s *CreateServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	res := &pb.CreateResponse{}

	// Validate incoming urls. Don't allow any non-legit URL's.
	for _, u := range req.Urls {
		_, err := url.ParseRequestURI(u)
		if err != nil {
			log.Printf("CreateHandler.error: recieved bad payload: %v", err)
			return res, err
		}
	}

	// Create/get the short code for the URL and build.
	for _, u := range req.Urls {
		code, err := s.Store.Save(u)
		if err != nil {
			log.Printf("CreateHandler.error: failed to shorten url: %v", err)
			return res, err
		}

		res.Urls = append(
			res.Urls,
			fmt.Sprintf("%s/%s", BaseURL, code),
		)
	}

	return res, nil
}

// LookupHandler for adding a new URL to the store via HTTP.
func LookupHandler(store storage.IStorage) func(http.ResponseWriter, *http.Request) {
	fnc := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		code := vars["hash"]
		fullURL, err := store.Load(code)

		if err == storage.ErrNotFound {
			log.Printf("LookupHander.warn: could not find url with code %s", code)
			http.Redirect(w, r, NotFoundURL, 301)
			return
		}

		if err != nil {
			log.Printf("LookupHandler.error: receieved error while looking up url: %v", err)
			http.Redirect(w, r, NotFoundURL, 301)
			return
		}

		http.Redirect(w, r, fullURL, 301)
	}
	return fnc
}
