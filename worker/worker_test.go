package worker

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

const (
	workerCount = 2
)

func TestWorkerPool(t *testing.T) {
	wp := New(workerCount)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go wp.Run(ctx)

	numJobs := 5
	for i := 0; i < numJobs; i++ {
		wp.SendJob(Job{
			ID:     strconv.Itoa(i),
			ExecFn: multiplyByTwo, // same function from job_test.
			Args:   i,
		})
	}

	for {
		select {
		case r, ok := <-wp.Results():
			if !ok {
				continue
			}

			if r.Err == nil {
				i, err := strconv.ParseInt(r.JobID, 10, 64)
				if err != nil {
					fmt.Println("->", r.JobID)
					t.Fatalf("unexpected error: %v", err)
				}

				val := r.Value.(int)
				if val != int(i)*2 {
					t.Fatalf("wrong value %v; expected %v", val, int(i)*2)
				} else {
					numJobs -= 1
				}
			} else {
				if r.Err != context.Canceled { // Context error do to context cancellation to stop gorutines.
					t.Fatalf("unexpected error: %v", r.Err)
				}
			}

		case <-wp.Done:
			return
		default:
			if numJobs == 0 {
				cancel()
			}
		}
	}
}

func TestWorkerPool_TimeOut(t *testing.T) {
	wp := New(workerCount)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Nanosecond*10)
	defer cancel()

	go wp.Run(ctx)

	for {
		select {
		case r := <-wp.Results():
			if r.Err != nil && r.Err != context.DeadlineExceeded {
				t.Fatalf("expected error: %v; got: %v", context.DeadlineExceeded, r.Err)
			}
		case <-wp.Done:
			return
		default:
		}
	}
}

func TestWorkerPool_Cancel(t *testing.T) {
	wp := New(workerCount)

	ctx, cancel := context.WithCancel(context.TODO())

	go wp.Run(ctx)
	cancel()

	for {
		select {
		case r := <-wp.Results():
			if r.Err != nil && r.Err != context.Canceled {
				t.Fatalf("expected error: %v; got: %v", context.Canceled, r.Err)
			}
		case <-wp.Done:
			return
		default:
		}
	}
}
