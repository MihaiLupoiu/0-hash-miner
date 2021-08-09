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

// HashRate will print the mega hash rate per second.
func HashRate(counter *ratecounter.RateCounter, stop chan bool) {
	t := time.NewTimer(time.Second)
	interval := time.Second * time.Duration(1)

	for {
		select {
		case <-stop:
			fmt.Println("Closing HashRate gorutine")
			return
		case <-t.C:
			fmt.Println(float64(counter.Rate())/float64(1000000), "Mh/s")
		}
		t.Reset(interval)
	}
}
