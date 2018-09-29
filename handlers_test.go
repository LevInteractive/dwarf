package main

import (
	"context"
	"log"
	"testing"

	"github.com/LevInteractive/dwarf/pb"
	"github.com/LevInteractive/dwarf/storage"
	"github.com/go-redis/redis"
)

func TestGRPCCreate(t *testing.T) {
	store := storage.Redis{
		CharFloor: 2,
		Conn: &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       2,
		},
	}

	store.Init()

	server := &CreateServer{
		Store: store,
	}

	tests := []struct {
		req   *pb.CreateRequest
		value int
	}{
		{
			req:   &pb.CreateRequest{Urls: []string{"http://hello.com"}},
			value: 1,
		},
		{
			req:   &pb.CreateRequest{Urls: []string{"http://a.com", "http://foo.com"}},
			value: 2,
		},
	}

	for _, tt := range tests {
		resp, err := server.Create(context.Background(), tt.req)
		if err != nil {
			log.Fatal(err)
		}

		if len(resp.Urls) != tt.value {
			t.Fatalf("Length of response should be 1. Got %v", len(tt.req.Urls))
		}
	}
}
