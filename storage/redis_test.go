package storage

import (
	"fmt"
	"testing"

	"github.com/LevInteractive/dwarf/storage"
	"github.com/go-redis/redis"
)

func inArray(val string, array []string) bool {
	for i := range array {
		if array[i] == val {
			return true
		}
	}
	return false
}

func TestCreation(t *testing.T) {
	store := storage.Redis{
		CharFloor: 2,
		Conn: &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       2,
		},
		BaseURL: "https://dwarf.short",
	}

	store.Init()

	var codes []string

	// Should create new without dupes.
	for i := 1; i <= 1000; i++ {
		code, err := store.Save(fmt.Sprintf("http://google.com/%d", i))
		if err != nil {
			t.Fatal(err)
		}
		if inArray(code, codes) != false {
			t.Fatalf("there was a duplicate code made: %s", code)
		}
		codes = append(codes, code)
	}

	// Should retrieve existing codes.
	for j := 1; j <= 1000; j++ {
		code, err := store.Save(fmt.Sprintf("http://google.com/%d", j))
		if err != nil {
			t.Fatal(err)
		}
		if inArray(code, codes) == false {
			t.Fatalf("code should have already existed: %s", code)
		}
	}
	store.Client.FlushDB()
}

func TestLoad(t *testing.T) {
	store := storage.Redis{
		CharFloor: 2,
		Conn: &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       2,
		},
		BaseURL: "https://dwarf.short",
	}

	store.Init()

	u := "http://google.com/"

	code, err := store.Save(u)
	if err != nil {
		t.Fatal(err)
	}

	ru, err := store.Load(code)
	store.Client.FlushDB()

	if err != nil {
		t.Fatal(err)
	}

	if ru != u {
		t.Fatalf("the url saved was not the same as the one returned %s / %s", u, ru)
	}
}
