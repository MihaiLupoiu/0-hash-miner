// TODO: Refactor this package and eliminate utils.
package utils

import (
	"crypto/rand"
	"fmt"
	mrand "math/rand"
	"time"

	"github.com/paulbellamy/ratecounter"
)

const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// SecureRandomString generates a cryptography secure random string but at cost of performance
func SecureRandomString(length int) (string, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

// RandStringRunes generates a non-secure secure random string but with better performance
func RandStringRunes(length int) (string, error) {
	bytes := make([]byte, length)

	if _, err := mrand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

func HashRate(counter *ratecounter.RateCounter) {
	for range time.Tick(time.Second * 1) {
		fmt.Println("Hash rate:", counter.Rate())
	}
}
