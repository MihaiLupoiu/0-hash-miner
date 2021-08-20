package solver

import (
	"context"
	"reflect"
	"testing"

	"github.com/MihaiLupoiu/interview-exasol/utils"
	"github.com/google/uuid"
)

func TestCalculateAndCheckHash(t *testing.T) {
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
			if got := CalculateAndCheckHash(tt.args.authdata, tt.args.suffix, tt.args.difficulty); got != tt.want {
				t.Errorf("CalculateAndCheckHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCheckConstantString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CalculateAndCheckHash("", "aa", 9)
	}
}

func BenchmarkCheckWithUUID(b *testing.B) {
	for n := 0; n < b.N; n++ {
		suffix := uuid.New().String()
		CalculateAndCheckHash("", suffix, 9)
	}
}

func BenchmarkCheckWithRandomString10(b *testing.B) {
	for n := 0; n < b.N; n++ {
		suffix, _ := utils.RandStringRunes(10)
		CalculateAndCheckHash("", suffix, 9)
	}
}

func BenchmarkCheckWithRandomString20(b *testing.B) {
	for n := 0; n < b.N; n++ {
		suffix, _ := utils.RandStringRunes(20)
		CalculateAndCheckHash("", suffix, 9)
	}
}

func BenchmarkCheckWithRandomString30(b *testing.B) {
	for n := 0; n < b.N; n++ {
		suffix, _ := utils.RandStringRunes(30)
		CalculateAndCheckHash("", suffix, 9)
	}
}

func BenchmarkCheckWithRandomString32(b *testing.B) {
	for n := 0; n < b.N; n++ {
		suffix, _ := utils.RandStringRunes(32)
		CalculateAndCheckHash("", suffix, 9)
	}
}

func TestCalculateHash(t *testing.T) {
	type args struct {
		ctx  context.Context
		args interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "should work sha1",
			args: args{
				ctx:  context.TODO(),
				args: "l",
			},
			want:    [20]byte{7, 195, 66, 190, 110, 86, 14, 127, 67, 132, 46, 46, 33, 183, 116, 230, 29, 133, 240, 71},
			wantErr: false,
		},
		{
			name: "should fail with wrong type of argument",
			args: args{
				ctx:  context.TODO(),
				args: 18,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateHash(tt.args.ctx, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckDificulty(t *testing.T) {
	type args struct {
		hash      [20]byte
		dificulty int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should return true",
			args: args{
				hash:      [20]byte{7, 195, 66, 190, 110, 86, 14, 127, 67, 132, 46, 46, 33, 183, 116, 230, 29, 133, 240, 71},
				dificulty: 1,
			},
			want: true,
		},
		{
			name: "should return false",
			args: args{
				hash:      [20]byte{7, 195, 66, 190, 110, 86, 14, 127, 67, 132, 46, 46, 33, 183, 116, 230, 29, 133, 240, 71},
				dificulty: 2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckDificulty(tt.args.hash, tt.args.dificulty); got != tt.want {
				t.Errorf("CheckDificulty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchForHashWithDificulty(t *testing.T) {
	type args struct {
		authdata   []byte
		length     int
		difficulty int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Testing new solver",
			args: args{
				authdata:   []byte("cQokBByiRKwFNFhsXUvtTuEwRPwXdFjBeLjelxqPXoQHhIZaXMucoBSBpKFRkDFR"),
				length:     32,
				difficulty: 9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SearchForHashWithDificulty(tt.args.authdata, tt.args.length, tt.args.difficulty)
		})
	}
}

func TestHexStartsWith3(t *testing.T) {
	type args struct {
		hash      []byte
		dificulty int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Testing valid dificulty 6 hash",
			args: args{
				hash:      []byte{0, 0, 0, 190, 110, 86, 14, 127, 67, 132, 46, 46, 33, 183, 116, 230, 29, 133, 240, 71},
				dificulty: 6,
			},
			want: true,
		},
		{
			name: "Testing invalid dificulty 2 hash",
			args: args{
				hash:      []byte{7, 195, 66, 190, 110, 86, 14, 127, 67, 132, 46, 46, 33, 183, 116, 230, 29, 133, 240, 71},
				dificulty: 2,
			},
			want: false,
		},
		{
			name: "Testing invalid dificulty hash",
			args: args{
				hash:      []byte{48, 48, 173, 195, 152, 181, 227, 131, 214, 50, 135, 230, 158, 235, 173, 65, 253, 164, 140, 187},
				dificulty: 2,
			},
			want: false,
		},
		{
			name: "Testing valid dificulty 9 hash",
			args: args{
				hash:      []byte{0, 0, 0, 0, 7, 86, 14, 127, 67, 132, 46, 46, 33, 183, 116, 230, 29, 133, 240, 71},
				dificulty: 9,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexStartsWith3(tt.args.hash, tt.args.dificulty); got != tt.want {
				t.Errorf("HexStartsWith3() = %v, want %v", got, tt.want)
			}
		})
	}
}
