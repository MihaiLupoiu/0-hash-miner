package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/MihaiLupoiu/interview-exasol/miner"
	"github.com/MihaiLupoiu/interview-exasol/solver"
	"github.com/MihaiLupoiu/interview-exasol/utils"
	"github.com/MihaiLupoiu/interview-exasol/worker"
	"github.com/paulbellamy/ratecounter"
)

func runningtime(s string) (string, time.Time) {
	log.Println("Start:	", s)
	return s, time.Now()
}

func track(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End:	", s, "took", endTime.Sub(startTime))
}

func execute(authdata string, difficulty int, minStringLength, maxStringLength int) {
	defer track(runningtime("simple"))

	hashrateCounter := ratecounter.NewRateCounter(1 * time.Second)
	stop := make(chan bool, 1)
	go utils.HashRate(hashrateCounter, stop)

	for {
		length := rand.Intn(maxStringLength-minStringLength+1) + minStringLength
		suffix, _ := utils.RandStringRunes(length)
		hashrateCounter.Incr(1)
		if solver.CalculateAndCheckHash(authdata, suffix, difficulty) != "" {
			fmt.Printf("Authdata: %s\nSuffix: %s\nDifficulty: %d\n", authdata, suffix, difficulty)
			break
		}
	}
}

func testCorutine(authdata string, workers, difficulty, minStringlength, maxStringLength int) {
	defer track(runningtime("gorutines"))

	fmt.Println(authdata, difficulty)

	wp := worker.New(workers)
	hashrateCounter := ratecounter.NewRateCounter(1 * time.Second)

	context, cancelWorkerPool := context.WithCancel(context.Background())
	defer cancelWorkerPool()

	// Start workers
	go wp.Run(context)
	stop := make(chan bool, 1)
	go utils.HashRate(hashrateCounter, stop)

	jobs := miner.GenerateWorkerJobs(wp.GetWorkerCount(), difficulty, minStringlength, maxStringLength, authdata, hashrateCounter)
	go wp.SendBulkJobs(jobs)

	if suff, err := miner.GetResults(wp); err == nil && suff != "" {
		hash := sha1.Sum([]byte(authdata + suff))
		fmt.Println(hex.EncodeToString(hash[:]))
		stop <- true
	}
}

func main() {
	var difficulty int
	var minStringlength int
	var maxStringLength int
	var workers int
	flag.IntVar(&difficulty, "diff", 6, "Difficulty is the number of 0 in hex")
	flag.IntVar(&minStringlength, "minS", 15, "minStringlength is the minimum size of the random string to generate")
	flag.IntVar(&maxStringLength, "maxS", 30, "maxStringLength is the maximum size of the random string to generate")
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "number of workers to run in the pool")

	flag.Parse()

	// using authdata from a connection request.
	// authdata := "fSHmbbPePDavjmRTSdOITaUnTtBkbPcnIiYjWemfoBMUoGZNTmIPrnNEUAGtYrKn"
	authdata := "cQokBByiRKwFNFhsXUvtTuEwRPwXdFjBeLjelxqPXoQHhIZaXMucoBSBpKFRkDFR"

	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	rand.Seed(1) // Set random number to make calculate the same hash values.

	testCorutine(authdata, workers, difficulty, minStringlength, maxStringLength)

	f2, err := os.Create("mem.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f2.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f2); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	rand.Seed(1) // Set random number to make calculate the same hash values.

	execute(authdata, difficulty, minStringlength, maxStringLength)

}
