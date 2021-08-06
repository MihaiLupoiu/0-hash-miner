package miner

import (
	"context"
	"errors"
	"fmt"

	"github.com/MihaiLupoiu/interview-exasol/solver"
	"github.com/MihaiLupoiu/interview-exasol/utils"
	"github.com/MihaiLupoiu/interview-exasol/worker"
	"github.com/paulbellamy/ratecounter"
)

type Args struct {
	Authdata        string
	Difficulty      int
	SuffixLength    int
	Seed            int64
	HashrateCounter *ratecounter.RateCounter
}

func FindHash(ctx context.Context, args interface{}) (interface{}, error) {
	argVal, ok := args.(Args)
	if !ok {
		return nil, errors.New("wrong argument type")
	}

	for {
		// generate short random string, server accepts all utf-8 characters,
		// except [\n\r\t ], it means that the suffix should not contain the
		// characters: newline, carriege return, tab and space
		suffix, _ := utils.RandStringRunes(argVal.SuffixLength)
		argVal.HashrateCounter.Incr(1)

		if solver.CalculateAndCheckHash(argVal.Authdata, suffix, argVal.Difficulty) != "" {
			fmt.Printf("Authdata: %s\n Suffix: %s\n", argVal.Authdata, suffix)
			return suffix, nil
		}

		if ctx.Err() == context.Canceled {
			return "", nil
		}
	}
}

func GenerateWorkerJobs(jobsCount, difficulty, suffixLength int, authdata string, counter *ratecounter.RateCounter) []worker.Job {
	jobs := make([]worker.Job, jobsCount)
	for i := 0; i < jobsCount; i++ {
		jobs[i] = worker.Job{
			ID:     fmt.Sprintf("%v", i),
			ExecFn: FindHash,
			Args: Args{
				Authdata:        authdata,
				Difficulty:      difficulty,
				SuffixLength:    suffixLength,
				Seed:            int64(i),
				HashrateCounter: counter,
			},
		}
	}
	return jobs
}

func GetResults(wPool worker.Pool) (string, error) {
	select {
	case r, ok := <-wPool.Results():
		if !ok {
			return "", errors.New("could not read results")
		}

		if r.Err == nil {
			suffix := r.Value.(string)
			if suffix != "" {
				return suffix, nil
			}
		} else {
			if r.Err != context.Canceled { // Context error do to context cancellation to stop gorutines.
				fmt.Printf("unexpected error: %v", r.Err)
				return "", r.Err
			}
		}

	case <-wPool.Done:
		return "", nil
	}

	return "", nil
}
