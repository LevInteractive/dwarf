// Redis Store Strategy:
//
// Will store 2 keys per URL:
//
// [prefix]:code:[fullURL] = [code]
// [prefix]:url:[code] = [fullURL]
//
// Using this we can lookup both ways very efficiently.

package storage

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

const prefix = "dwarf:"

// Redis storage.
type Redis struct {
	CharFloor int
	BaseURL   string
	Client    *redis.Client
	Conn      *redis.Options
}

// Init the connection to redis.
func (s *Redis) Init() {
	s.Client = redis.NewClient(s.Conn)
	_, err := s.Client.Ping().Result()

	if err != nil {
		log.Fatal(err)
	}
}

// Save the url in store with a new code if it isn't already there.
func (s Redis) Save(u string) (string, error) {
	codehash := fmt.Sprintf("%s:code:%s", prefix, u)
	existingCode, err := s.Client.Do("get", codehash).String()
	if existingCode != "" {
		log.Printf("redis.Save.info: didn't need to create new - existed: %s with code %s", u, existingCode)
		return existingCode, nil
	}

	code, err := discover(s.Client, s.CharFloor)
	if err != nil {
		return "", err
	}

	err = set(s.Client, u, code)

	if err != nil {
		return "", err
	}

	log.Printf("redis.Save.info: created new short url for %s / %s", u, code)
	return code, nil
}

// Load will lookup the code and return the full URL.
func (s Redis) Load(code string) (string, error) {
	hash := fmt.Sprintf("%s:url:%s", prefix, code)
	fullURL, err := s.Client.Do("get", hash).String()

	if err == redis.Nil {
		return "", ErrNotFound
	} else if err != nil {
		log.Printf("redis.Load.error: had redis error %v", err)
		return "", err
	}

	log.Printf("redis.Load.info: loaded url %s", fullURL)

	return fullURL, nil
}

// Discover a truly unqiue key.
func discover(c *redis.Client, n int) (string, error) {
	code := GenCode(n)
	hash := fmt.Sprintf("%s:url:%s", prefix, code)
	_, err := c.Do("get", hash).String()

	if err == redis.Nil {
		return code, nil
	} else if err != nil {
		log.Printf("redis.discover.error: had redis error %v", err)
		return "", err
	}

	log.Printf("redis.discover.info: had key collision, incrementing 1 char and looking again")
	return discover(c, n+1)
}

// Set each of our key hashes.
func set(c *redis.Client, fullURL string, code string) error {
	codehash := fmt.Sprintf("%s:code:%s", prefix, fullURL)
	urlhash := fmt.Sprintf("%s:url:%s", prefix, code)
	if err := c.Set(codehash, code, 0).Err(); err != nil {
		return err
	}
	if err := c.Set(urlhash, fullURL, 0).Err(); err != nil {
		log.Printf("redis.set.error: hit error with setting url hash. rolling back previous codehash. %v", err)
		c.Del(codehash)
		return err
	}

	return nil
}
