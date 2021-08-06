package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/MihaiLupoiu/interview-exasol/solver"
	"github.com/MihaiLupoiu/interview-exasol/utils"
	"github.com/MihaiLupoiu/interview-exasol/worker"
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
	defer track(runningtime("execute"))
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

///

func generateWork(authdata string, suffixLength int, stop chan bool, WPool worker.Pool) {
	for {
		select {
		case <-stop:
			fmt.Println("Closing HashRate gorutine")
			return
		default:
			// generate short random string, server accepts all utf-8 characters,
			// except [\n\r\t ], it means that the suffix should not contain the
			// characters: newline, carriege return, tab and space
			suffix, _ := utils.RandStringRunes(suffixLength)
			// fmt.Println("Suffix:", suffix)

			WPool.SendJob(worker.Job{
				ID:     suffix,
				ExecFn: solver.CalculateHash,
				Args:   authdata + suffix,
			})
		}
	}
}

func checkResults(WPool worker.Pool, difficulty int) (string, error) {
	select {
	case r, ok := <-WPool.Results():
		if !ok {
			return "", errors.New("could not read results")
		}

		if r.Err == nil {
			suffix := r.JobID
			hash := r.Value.([20]byte)

			if solver.CheckDificulty(hash, difficulty) {
				return suffix, nil
			}
		} else {
			if r.Err != context.Canceled { // Context error do to context cancellation to stop gorutines.
				fmt.Printf("unexpected error: %v", r.Err)
				return "", r.Err
			}
		}

	case <-WPool.Done:
		return "", nil
	}

	return "", nil
}

func testCorutine(workers, difficulty int) {
	rand.Seed(1) // Set random number to make calculate the same hash values.
	authdata := "cQokBByiRKwFNFhsXUvtTuEwRPwXdFjBeLjelxqPXoQHhIZaXMucoBSBpKFRkDFR"
	fmt.Println(authdata, difficulty)

	wp := worker.New(4)

	context, cancelWorkerPool := context.WithCancel(context.Background())
	defer cancelWorkerPool()

	// Start workers
	go wp.Run(context)
	stop := make(chan bool, 1)

	go generateWork(authdata, 15, stop, wp)

	for {
		if suff, err := checkResults(wp, difficulty); err == nil && suff != "" {
			fmt.Println(suff)
			stop <- true
			close(stop)
			break
		}
	}
}

///

func main() {
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

	testCorutine(4, 5)

	f2, err := os.Create("mem.pprof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f2.Close()
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f2); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	for i := 0; i <= 7; i++ {
		execute(i, 30)
	}
}
