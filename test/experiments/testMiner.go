// package is used to test the speed of the solver and the different corutines implementations
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
			stop <- true
			break
		}
	}
}

func execute2(authData string, difficulty, lenght int) {
	defer track(runningtime("simple 2"))

	hashrateCounter := ratecounter.NewRateCounter(1 * time.Second)
	stop := make(chan bool, 1)
	go utils.HashRate(hashrateCounter, stop)

	authdata := []byte(authData)

	bytes := make([]byte, len(authdata)+lenght)

	copy(bytes, authdata)
	randomGenerator := utils.InitRandomWithSeed(1)

	suffix := bytes[len(authdata):]
	for {

		utils.RandomUTF8(randomGenerator, suffix)
		hashrateCounter.Incr(1)
		if solver.CalculateHashAndCheckDifficulty(bytes, difficulty) {
			fmt.Printf("Authdata: %s\nSuffix: %s\nDifficulty: %d\n", authdata, suffix, difficulty)
			stop <- true
			break
		}
	}
}

func execute3(authData string, difficulty, lenght int) {
	defer track(runningtime("simple 3"))

	hashrateCounter := ratecounter.NewRateCounter(1 * time.Second)
	stop := make(chan bool, 1)
	go utils.HashRate(hashrateCounter, stop)

	authdata := []byte(authData)
	suffix := make([]byte, lenght)

	var ctx = utils.NewHash(authdata)
	randomGenerator := utils.InitRandomWithSeed(1)

	for {
		utils.RandomUTF8(randomGenerator, suffix)
		hashrateCounter.Incr(1)

		hash := ctx.Sum(suffix)

		if solver.HexStartsWith2(hash, difficulty) {
			fmt.Printf("Authdata: %s\nSuffix: %s\nDifficulty: %d\n", authdata, suffix, difficulty)
			fmt.Println("Hash:", hex.EncodeToString(hash))
			fmt.Printf("%v \n", hash)
			stop <- true
			break
		}
	}
}

func execute4(authData string, difficulty, minStringLength, maxStringLength int) {
	defer track(runningtime("simple 4"))

	hashrateCounter := ratecounter.NewRateCounter(1 * time.Second)
	stop := make(chan bool, 1)
	go utils.HashRate(hashrateCounter, stop)

	authdata := []byte(authData)

	var ctx = utils.NewHash(authdata)
	randomGenerator := utils.InitRandomWithSeed(1)

	for {
		suffix := make([]byte, rand.Intn(maxStringLength-minStringLength+1)+minStringLength)
		utils.RandomUTF8(randomGenerator, suffix)
		hashrateCounter.Incr(1)

		hash := ctx.Sum(suffix)

		if solver.HexStartsWith2(hash, difficulty) {
			fmt.Printf("Authdata: %s\nSuffix: %s\nDifficulty: %d\n", authdata, suffix, difficulty)
			fmt.Println("Hash:", hex.EncodeToString(hash))
			fmt.Printf("%v \n", hash)
			stop <- true
			break
		}
	}
}

func execute5(authData string, difficulty, lenght int) {
	defer track(runningtime("simple 5"))

	hashrateCounter := ratecounter.NewRateCounter(1 * time.Second)
	stop := make(chan bool, 1)
	go utils.HashRate(hashrateCounter, stop)

	authdata := []byte(authData)
	suffix := make([]byte, lenght)

	var ctx = utils.NewHash(authdata)
	randomGenerator := utils.InitRandomWithSeed(1)

	for {
		utils.RandomUTF8(randomGenerator, suffix)
		hashrateCounter.Incr(1)

		hash := ctx.Sum(suffix)

		if solver.HexStartsWith3(hash, difficulty) {
			fmt.Printf("Authdata: %s\nSuffix: %s\nDifficulty: %d\n", authdata, suffix, difficulty)
			fmt.Println("Hash:", hex.EncodeToString(hash))
			fmt.Printf("%v \n", hash)
			stop <- true
			break
		}
	}
}

func execute6(authData string, difficulty, minStringLength, maxStringLength int) {
	defer track(runningtime("simple 6"))

	hashrateCounter := ratecounter.NewRateCounter(1 * time.Second)
	stop := make(chan bool, 1)
	go utils.HashRate(hashrateCounter, stop)

	authdata := []byte(authData)

	var ctx = utils.NewHash(authdata)
	randomGenerator := utils.InitRandomWithSeed(1)

	for {
		suffix := make([]byte, rand.Intn(maxStringLength-minStringLength+1)+minStringLength)
		utils.RandomUTF8(randomGenerator, suffix)
		hashrateCounter.Incr(1)

		hash := ctx.Sum(suffix)

		if solver.HexStartsWith3(hash, difficulty) {
			fmt.Printf("Authdata: %s\nSuffix: %s\nDifficulty: %d\n", authdata, suffix, difficulty)
			fmt.Println("Hash:", hex.EncodeToString(hash))
			fmt.Printf("%v \n", hash)
			stop <- true
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
	// go utils.HashRate(hashrateCounter, stop)

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
	var multi bool
	var benchSimple bool
	var benchMulti bool
	// var noMetrics bool
	flag.IntVar(&difficulty, "diff", 6, "Difficulty is the number of 0 in hex")
	flag.IntVar(&minStringlength, "minS", 2, "minStringlength is the minimum size of the random string to generate")
	flag.IntVar(&maxStringLength, "maxS", 5, "maxStringLength is the maximum size of the random string to generate")
	flag.IntVar(&workers, "workers", runtime.NumCPU(), "number of workers to run in the pool")
	flag.BoolVar(&multi, "multi", false, "run on corutines.")
	flag.BoolVar(&benchSimple, "benchSimple", false, "benchmark only simple funtion.")
	flag.BoolVar(&benchMulti, "benchMulti", false, "benchmark only corutines funtion.")
	// flag.BoolVar(&noMetrics, "noMetrics", false, "no hashrate metrics.")

	flag.Parse()

	// using authdata from a connection request.
	// authdata := "fSHmbbPePDavjmRTSdOITaUnTtBkbPcnIiYjWemfoBMUoGZNTmIPrnNEUAGtYrKn"
	authdata := "cQokBByiRKwFNFhsXUvtTuEwRPwXdFjBeLjelxqPXoQHhIZaXMucoBSBpKFRkDFR"

	if benchSimple {
		for i := 1; i < 10; i++ {
			fmt.Printf("Starting Difficulty: %d\n", i)

			// rand.Seed(1) // Set random number to make calculate the same hash values.
			// execute(authdata, i, minStringlength, maxStringLength)

			// rand.Seed(1) // Set random number to make calculate the same hash values.
			// execute2(authdata, i, 64)

			// rand.Seed(1) // Set random number to make calculate the same hash values.
			// execute3(authdata, i, 32)

			// rand.Seed(1) // Set random number to make calculate the same hash values.
			// execute4(authdata, i, minStringlength, maxStringLength)

			rand.Seed(1) // Set random number to make calculate the same hash values.
			execute5(authdata, i, 32)

			// rand.Seed(1) // Set random number to make calculate the same hash values.
			// execute6(authdata, i, minStringlength, maxStringLength)

		}
		return
	}

	if benchMulti {
		for i := 1; i < 10; i++ {
			rand.Seed(1) // Set random number to make calculate the same hash values.
			testCorutine(authdata, workers, i, minStringlength, maxStringLength)
		}
		return
	}

	// Generate CPU pprof
	f, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	/*
		var wg sync.WaitGroup

		go makeRandomNumbers(&wg, 1)
		go makeRandomNumbers(&wg, 2)
		go makeRandomNumbers(&wg, 3)
		go makeRandomNumbers(&wg, 4)
		go makeRandomNumbers(&wg, 5)
		wg.Add(5)
		wg.Wait()
	*/

	if multi {
		// rand.Seed(1) // Set random number to make calculate the same hash values.
		testCorutine(authdata, workers, difficulty, minStringlength, maxStringLength)
	} else {
		length := rand.Intn(maxStringLength-minStringlength+1) + minStringlength
		execute5(authdata, difficulty, length)

	}

	// Generate memory pprof
	f2, err := os.Create("mem.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f2.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f2); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	/*
		// Generate CPU pprof
		f, err := os.Create("cpuSimple.pprof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()

		rand.Seed(1) // Set random number to make calculate the same hash values.

		execute5(authdata, difficulty, 35)
	*/
}

/*
func makeRandomNumbers(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	source := rand.NewSource(1)
	generator := rand.New(source)
	sum := 0
	for i := 0; i < 1000000; i++ {
		sum += generator.Intn(100)
		// fmt.Println("id: ", id, "Random: ", generator.Intn(100))
		// time.Sleep(2 * time.Millisecond)
	}

	fmt.Println("id: ", id, "SUM:  ", sum)
}

*/
