package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/MihaiLupoiu/interview-exasol/solver"
	"github.com/MihaiLupoiu/interview-exasol/utils"
)

func runningtime(s string) (string, time.Time) {
	log.Println("Start:	", s)
	return s, time.Now()
}

func track(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End:	", s, "took", endTime.Sub(startTime))
}

func execute(difficulty int) {
	defer track(runningtime("execute"))
	// using authdata from a connection request.
	authdata := "cQokBByiRKwFNFhsXUvtTuEwRPwXdFjBeLjelxqPXoQHhIZaXMucoBSBpKFRkDFR"
	for {
		suffix, _ := utils.RandStringRunes(30)
		if solver.Check(authdata, suffix, difficulty) != "" {
			fmt.Printf("Authdata: %s\nSuffix: %s\nDifficulty: %d\n", authdata, suffix, difficulty)
			break
		}
	}
}

func main() {
	rand.Seed(1) // Set random number to make calculate the same hash values.
	for i := 0; i < 10; i++ {
		execute(i)
	}
}
