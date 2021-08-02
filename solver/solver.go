package solver

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

// Check if hex starts with dificulty number of 0s.
func hexStartsWith(hash [20]byte, dificulty int) bool {
	// Improve method to use bit manipulation for more optimal comparion.
	// E.g: Convert to int64 and compare it's content with expected value from dificulty.
	sha1_hash := hex.EncodeToString(hash[:])
	prefixDifficulty := strings.Repeat("0", dificulty)
	res := strings.HasPrefix(sha1_hash, prefixDifficulty)
	return res
}

// Check calculates the SHA1 of the authdata + suffix and that it starst with as meny 0s as the difficulty number.
func Check(authdata, suffix string, difficulty int) string {
	cksum_in_hex := sha1.Sum([]byte(authdata + suffix))
	// fmt.Printf("  SHA1: %x\n", cksum_in_hex)

	// check if the checksum has enough leading zeros
	// (length of leading zeros should be equal to the difficulty)
	if hexStartsWith(cksum_in_hex, difficulty) {
		return suffix
	}

	return ""
}
