package main

import (
	"context"
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

func execute(difficulty int, stringLength int) {
	defer track(runningtime("simple"))
	// using authdata from a connection request.
	authdata := "cQokBByiRKwFNFhsXUvtTuEwRPwXdFjBeLjelxqPXoQHhIZaXMucoBSBpKFRkDFR"
	for {
		suffix, _ := utils.RandStringRunes(stringLength)
		if solver.CalculateAndCheckHash(authdata, suffix, difficulty) != "" {
			fmt.Printf("Authdata: %s\nSuffix: %s\nDifficulty: %d\n", authdata, suffix, difficulty)
			break
		}
	}
}

func testCorutine(workers, difficulty, suffixLength int) {
	defer track(runningtime("gorutines"))

	rand.Seed(1) // Set random number to make calculate the same hash values.
	authdata := "fSHmbbPePDavjmRTSdOITaUnTtBkbPcnIiYjWemfoBMUoGZNTmIPrnNEUAGtYrKn"
	fmt.Println(authdata, difficulty)

	wp := worker.New(4)
	hashrateCounter := ratecounter.NewRateCounter(1 * time.Second)

	context, cancelWorkerPool := context.WithCancel(context.Background())
	defer cancelWorkerPool()

	// Start workers
	go wp.Run(context)
	stop := make(chan bool, 1)
	go utils.HashRate(hashrateCounter, stop)

	jobs := miner.GenerateWorkerJobs(wp.GetWorkerCount(), difficulty, suffixLength, authdata, hashrateCounter)
	go wp.SendBulkJobs(jobs)

	if suff, err := miner.GetResults(wp); err == nil && suff != "" {
		fmt.Println("Suff: ", suff)
		stop <- true
	}
}

func main() {
	var difficulty int
	var suffixLength int
	flag.IntVar(&difficulty, "diff", 6, "Difficulty is the number of 0 in hex")
	flag.IntVar(&suffixLength, "suff", 15, "suffixLength is the size of the random string to generate")
	flag.Parse()

	rand.Seed(1) // Set random number to make calculate the same hash values.

	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	testCorutine(runtime.NumCPU(), difficulty, suffixLength)

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

	for i := 0; i <= difficulty; i++ {
		execute(i, suffixLength)
	}
}
