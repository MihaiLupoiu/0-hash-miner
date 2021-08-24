package solver

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/MihaiLupoiu/interview-exasol/utils"
)

var prefixDifficultyMap = map[int]string{
	1:  "0",
	2:  "00",
	3:  "000",
	4:  "0000",
	5:  "00000",
	6:  "000000",
	7:  "0000000",
	8:  "00000000",
	9:  "000000000",
	10: "0000000000",
}

var prefixDifficultyBytesMap = map[int][]byte{
	1: {240, 0, 0, 0, 0, 0, 0, 0},
	2: {255, 0, 0, 0, 0, 0, 0, 0},
	3: {255, 240, 0, 0, 0, 0, 0, 0},
	4: {255, 255, 0, 0, 0, 0, 0, 0},
	5: {255, 255, 240, 0, 0, 0, 0, 0},
	6: {255, 255, 255, 0, 0, 0, 0, 0},
	7: {255, 255, 255, 240, 0, 0, 0, 0},
	8: {255, 255, 255, 255, 0, 0, 0, 0},
	9: {255, 255, 255, 255, 240, 0, 0, 0},
}

// Check if hex starts with dificulty number of 0s.
func hexStartsWith(hash [20]byte, dificulty int) bool {
	// Improve method to use bit manipulation for more optimal comparion.
	// E.g: Convert to uint64 and compare it's content with expected value from dificulty.
	sha1_hash := hex.EncodeToString(hash[:])
	res := strings.HasPrefix(sha1_hash, prefixDifficultyMap[dificulty])
	return res
}

func HexStartsWith2(hash []byte, dificulty int) bool {
	// Improve method to use bit manipulation for more optimal comparion.
	// E.g: Convert to uint64 and compare it's content with expected value from dificulty.
	sha1_hash := hex.EncodeToString(hash)
	res := strings.HasPrefix(sha1_hash, prefixDifficultyMap[dificulty])
	return res
}

func HexStartsWith3(hash []byte, dificulty int) bool {
	// return bytes.Equal(prefixDifficultyBytesMap[dificulty], hash[:len(prefixDifficultyBytesMap[dificulty])])

	for i := 0; i < len(prefixDifficultyBytesMap[dificulty]); i++ {
		res := hash[i] & prefixDifficultyBytesMap[dificulty][i]
		if res != 0 {
			return false
		}
	}
	return true
}

// CalculateAndCheckHash calculates the SHA1 of the authdata + suffix and that it starst with as meny 0s as the difficulty number.
func CalculateAndCheckHash(authdata, suffix string, difficulty int) string {
	cksum_in_hex := sha1.Sum([]byte(authdata + suffix))
	// fmt.Printf("  SHA1: %x\n", cksum_in_hex)

	// check if the checksum has enough leading zeros
	// (length of leading zeros should be equal to the difficulty)
	if hexStartsWith(cksum_in_hex, difficulty) {
		return suffix
	}

	return ""
}

// CalculateHash will generate the SHA1 of the arguments.
func CalculateHash(ctx context.Context, args interface{}) (interface{}, error) {
	argVal, ok := args.(string)
	if !ok {
		return nil, errors.New("wrong argument type")
	}

	return sha1.Sum([]byte(argVal)), nil
}

// CalculateHash will check if the hash has the dificulty ammount of 0 as a suffix.
func CheckDificulty(hash [20]byte, dificulty int) bool {
	return hexStartsWith(hash, dificulty)
}

func SearchForHashWithDificulty(authdata []byte, length int, dificulty int) {
	bytes := make([]byte, len(authdata)+length)

	copy(bytes, authdata)
	randomGenerator := utils.InitRandomWithRandomSeed()

	suffix := bytes[len(authdata):]
	utils.RandomUTF8(randomGenerator, suffix)
}

/*
func verifyHash(hash [20]byte, difficulty int) bool {
}

*/

func CalculateHashAndCheckDifficulty(bytes []byte, difficulty int) bool {
	cksum_in_hex := sha1.Sum(bytes)

	return CheckDificulty(cksum_in_hex, difficulty)
}
