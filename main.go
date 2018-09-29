package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/psaia/dwarf/storage"
)

// GetEnv grabs env with a fallback
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	// Require a base url.
	base := GetEnv("APP_BASE_URL", "")
	if base == "" {
		panic("APP_BASE_URL must be set to a valid url. E.g. https://short.in")
	}

	// Configure the Redis store. Redis is all we have now so eh.
	db, err := strconv.Atoi(GetEnv("REDIS_DB", "0"))
	if err != nil {
		log.Fatal(err)
	}

	charFloor, err := strconv.Atoi(GetEnv("CHAR_FLOOR", "3"))
	if err != nil {
		log.Fatal(err)
	}

	store := storage.Redis{
		CharFloor: charFloor,
		Conn: &redis.Options{
			Addr:     GetEnv("REDIS_SERVER", "localhost:6379"),
			Password: GetEnv("REDIS_PASS", ""),
			DB:       db,
		},
	}

	store.Init()

	// App routing.
	appPort := GetEnv("APP_PORT", ":8000")
	r := mux.NewRouter()
	r.HandleFunc("/create", CreateHandler(store, base)).Methods("POST")
	r.HandleFunc("/{hash}", LookupHandler(store, base)).Methods("GET")
	http.Handle("/", r)
	log.Printf(
		"Dwarf is listening here -> %s and exposed here -> %s",
		appPort,
		GetEnv("APP_BASE_URL", "!!BASE DOMAIN NOT SET!!"),
	)
	log.Fatal(http.ListenAndServe(appPort, r))
}
