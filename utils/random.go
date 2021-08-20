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

// usft8Chars is a list of useful chars to use for the random string. First char to use is 0x21 "!" until 0x7e "~"
// https://www.fileformat.info/info/charset/UTF-8/list.htm
var utf8Chars = [...]byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f, 0x40, 0x41,
	0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53,
	0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f, 0x60, 0x61, 0x62, 0x63, 0x64, 0x65,
	0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77,
	0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e}

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

//RandASCIIBytes - A helper function create and fill a slice of length n with characters from a-zA-Z0-9_-. It panics if there are any problems getting random bytes.
func RandASCIIBytes(n int) []byte {
	output := make([]byte, n)

	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)

	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}

	l := len(chars)
	// fill output
	for pos := range output {
		// get random item
		random := uint8(randomness[pos])

		// random % 64
		randomPos := random % uint8(l)

		// put into output
		output[pos] = chars[randomPos]
	}

	return output
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

func RandomUTF8(randomString []byte) error {
	if _, err := mrand.Read(randomString); err != nil {
		return err
	}

	for i, b := range randomString {
		randomString[i] = utf8Chars[b%byte(len(utf8Chars))]
	}

	return nil
}
