package worker

import "context"

type ExecutionFn func(ctx context.Context, args interface{}) (interface{}, error)

type Job struct {
	ID     string
	ExecFn ExecutionFn
	Args   interface{}
}

type Result struct {
	Value interface{}
	Err   error
	JobID string
}

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
