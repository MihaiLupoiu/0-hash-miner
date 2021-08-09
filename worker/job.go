package worker

import "context"

// ExecutionFn is interface for what job will be executed.
type ExecutionFn func(ctx context.Context, args interface{}) (interface{}, error)

// Job is the definition of how work wil be passed and what funtion to execute.
type Job struct {
	ID     string
	ExecFn ExecutionFn
	Args   interface{}
}

// Result is the result of the job execution.
type Result struct {
	Value interface{}
	Err   error
	JobID string
}

// execute will execute the function and return the result.
func (j Job) execute(ctx context.Context) Result {
	value, err := j.ExecFn(ctx, j.Args)
	if err != nil {
		return Result{
			Err:   err,
			JobID: j.ID,
		}
	}

	return Result{
		Value: value,
		JobID: j.ID,
	}
}
