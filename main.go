package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/LevInteractive/dwarf/pb"
	"github.com/LevInteractive/dwarf/storage"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

// BaseURL is the base URL of the application. Short codes are appended on to it
// to ultimately construct the full short url.
var BaseURL string

// NotFoundURL is the URL used when Dwarf needs to redirect if a short code
// doesn't exist. We also redirect to this if the user lands on the root (/) of
// Dwarf.
var NotFoundURL string

func init() {
	BaseURL = GetEnv("APP_BASE_URL", "https://example.com")
	NotFoundURL = GetEnv("NOTFOUND_REDIRECT_URL", "https://google.com")
}

// GetEnv grabs env with a fallback
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
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

	go setupHTTPDiscovery(store)
	setupGrpcDiscovery(store)
}

func setupGrpcDiscovery(store storage.IStorage) {
	lis, err := net.Listen("tcp", GetEnv("GRPC_PORT", ":8001"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterDwarfServer(s, &CreateServer{
		Store: store,
	})

	log.Printf(
		"Dwarf's gRPC server is listening here -> %s",
		GetEnv("GRPC_PORT", ":8001"),
	)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func setupHTTPDiscovery(store storage.IStorage) {
	appPort := GetEnv("APP_PORT", ":8000")

	r := mux.NewRouter()
	r.HandleFunc("/{hash:.*}", LookupHandler(store)).Methods("GET")
	http.Handle("/", r)
	log.Printf(
		"Dwarf's public HTTP server is listening here -> %s and exposed here -> %s",
		appPort,
		BaseURL,
	)
	log.Fatal(http.ListenAndServe(appPort, r))
}
