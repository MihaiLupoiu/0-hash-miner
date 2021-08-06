package workerpool

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

var (
	errDefault = errors.New("wrong argument type")
	jobID      = "1"
)

func multiplyByTwo(ctx context.Context, args interface{}) (interface{}, error) {
	argVal, ok := args.(int)
	if !ok {
		return nil, errDefault
	}
	return argVal * 2, nil
}

func TestJob_execute(t *testing.T) {
	type fields struct {
		ID     string
		ExecFn ExecutionFn
		Args   interface{}
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Result
	}{
		// TODO: Add test cases.
		{
			name: "Succesful job execution",
			fields: fields{
				ID:     jobID,
				ExecFn: multiplyByTwo,
				Args:   10,
			},
			want: Result{
				Value: 20,
				JobID: jobID,
			},
		},
		{
			name: "Failed job execution",
			fields: fields{
				ID:     jobID,
				ExecFn: multiplyByTwo,
				Args:   "10",
			},
			want: Result{
				Err:   errDefault,
				JobID: jobID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := Job{
				ID:     tt.fields.ID,
				ExecFn: tt.fields.ExecFn,
				Args:   tt.fields.Args,
			}
			if got := j.execute(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Job.execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
