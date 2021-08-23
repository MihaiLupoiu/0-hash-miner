package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"reflect"
	"strings"
	"time"

	"github.com/paulbellamy/ratecounter"
)

type Hash struct {
	origv   reflect.Value
	hasherv reflect.Value
	hasher  hash.Hash

	result [sha1.Size]byte
}

func NewHash(fixed []byte) *Hash {
	h := sha1.New()
	h.Write(fixed)

	c := &Hash{origv: reflect.ValueOf(h).Elem()}
	hasherv := reflect.New(c.origv.Type())
	c.hasher = hasherv.Interface().(hash.Hash)
	c.hasherv = hasherv.Elem()

	return c
}

func (c *Hash) Sum(data []byte) []byte {
	// Set state of the fixed hash:
	c.hasherv.Set(c.origv)

	c.hasher.Write(data)
	return c.hasher.Sum(c.result[:0])
}

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
func CheckDificultyOriginal(hash [20]byte, dificulty int) bool {
	// Improve method to use bit manipulation for more optimal comparion.
	// E.g: Convert to uint64 and compare it's content with expected value from dificulty.
	sha1_hash := hex.EncodeToString(hash[:])
	res := strings.HasPrefix(sha1_hash, prefixDifficultyMap[dificulty])
	return res
}

func CheckDificulty1(hash []byte, dificulty int) bool {
	// Improve method to use bit manipulation for more optimal comparion.
	// E.g: Convert to uint64 and compare it's content with expected value from dificulty.
	sha1_hash := hex.EncodeToString(hash)
	res := strings.HasPrefix(sha1_hash, prefixDifficultyMap[dificulty])
	return res
}

func CheckDificulty(hash []byte, dificulty int) bool {
	// return bytes.Equal(prefixDifficultyBytesMap[dificulty], hash[:len(prefixDifficultyBytesMap[dificulty])])

	for i := 0; i < len(prefixDifficultyBytesMap[dificulty]); i++ {
		res := hash[i] & prefixDifficultyBytesMap[dificulty][i]
		if res != 0 {
			return false
		}
	}
	return true
}

// HashRate will print the mega hash rate per second.
func HashRate(counter *ratecounter.RateCounter, stop chan bool) {
	t := time.NewTimer(time.Second)
	interval := time.Second * time.Duration(1)
	fmt.Print("\033[s") // save the cursor position

	for {
		select {
		case <-stop:
			fmt.Println("Closing HashRate gorutine")
			return
		case <-t.C:
			fmt.Print("\033[u\033[K")
			fmt.Printf("%f MH/s", float64(counter.Rate())/float64(1000000))
		}
		t.Reset(interval)
	}
}
