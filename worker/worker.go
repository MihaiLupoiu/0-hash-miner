package worker

import (
	"context"
	"fmt"
	"sync"
)

type Pool struct {
	workersCount int
	jobs         chan Job
	results      chan Result
	Done         chan struct{}
}

func New(wcount int) Pool {
	return Pool{
		workersCount: wcount,
		jobs:         make(chan Job, wcount),
		results:      make(chan Result, wcount),
		Done:         make(chan struct{}),
	}
}

// GetWorkerCount returns the number of workers configured.
func (wp Pool) GetWorkerCount() int {
	return wp.workersCount
}

// Run will start the gorutines that will wait for jobs and
// the function will wait as well for all the jobs to finish.
func (wp Pool) Run(ctx context.Context) {
	var wg sync.WaitGroup

	for i := 0; i < wp.workersCount; i++ {
		wg.Add(1)
		go worker(ctx, &wg, wp.jobs, wp.results)
	}

	wg.Wait()
	close(wp.Done)
	close(wp.results)
}

// SendJob sends one job to be executed by the worker pool
func (wp Pool) SendJob(job Job) {
	wp.jobs <- job
}

// SendBulkJobs sends multiple jobs to be executed by the worker pool
func (wp Pool) SendBulkJobs(jobsBulk []Job) {
	for _, job := range jobsBulk {
		wp.jobs <- job
	}
	close(wp.jobs)
}

// Results returns the result of the job execution
func (wp Pool) Results() <-chan Result {
	return wp.results
}

// worker is the function that executes the job.
func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			results <- job.execute(ctx)
		case <-ctx.Done():
			fmt.Printf("cancelled worker. Error detail: %v\n", ctx.Err())
			results <- Result{
				Err: ctx.Err(),
			}
			return
		}
	}
}
