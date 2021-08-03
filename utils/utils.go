// TODO: Refactor this package and eliminate utils.
package utils

import (
	"crypto/rand"
	mrand "math/rand"
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
