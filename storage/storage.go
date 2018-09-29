package storage

import (
	"errors"
	"math/rand"
	"time"
)

// IStorage defines the interface for storages to implement. Any store that is
// implemented must conform this.
type IStorage interface {
	Save(string) (string, error)
	Load(string) (string, error)
}

// ErrNotFound is returned when a url can't be found with a given code.
var ErrNotFound = errors.New("not found")

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// GenCode for a short URL at a specific length.
// Barrowed from: https://stackoverflow.com/a/31832326
func GenCode(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
