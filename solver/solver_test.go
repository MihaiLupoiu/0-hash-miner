package solver

import (
	"testing"

	"github.com/google/uuid"
)

func TestCheck(t *testing.T) {
	type args struct {
		authdata   string
		suffix     string
		difficulty int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "No zeros in the beginning of the hash",
			args: args{
				authdata:   "",
				suffix:     "",
				difficulty: 1,
			},
			want: "",
		},
		{
			name: "1 zero in the beginning of the hash",
			args: args{
				authdata:   "",
				suffix:     "l",
				difficulty: 1,
			},
			want: "l",
		},
		{
			name: "2 zeros in the beginning of the hash",
			args: args{
				authdata:   "",
				suffix:     "543935c4-59e1-4b85-b062-4f9b01914336",
				difficulty: 2,
			},
			want: "543935c4-59e1-4b85-b062-4f9b01914336",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Check(tt.args.authdata, tt.args.suffix, tt.args.difficulty); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCheckConstantString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Check("", "aa", 9)
	}
}

func BenchmarkCheckWithUUID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		suffix := uuid.New().String()
		Check("", suffix, 9)
	}
}
