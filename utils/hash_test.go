package utils

import (
	"crypto/sha1"
	"testing"
)

func BenchmarkNormalHash(b *testing.B) {
	authdata := []byte("cQokBByiRKwFNFhsXUvtTuEwRPwXdFjBeLjelxqPXoQHhIZaXMucoBSBpKFRkDFR")
	suffix := []byte("sba(BE(p7`\"0]%>5X),n1$?n>~%(6G+j")
	bytes := make([]byte, len(authdata)+len(suffix))
	copy(bytes, authdata)
	copy(bytes[len(authdata):], suffix)
	for i := 0; i < b.N; i++ {
		sha1.Sum(bytes)
	}
}

func BenchmarkImprovedHash(b *testing.B) {
	authdata := []byte("cQokBByiRKwFNFhsXUvtTuEwRPwXdFjBeLjelxqPXoQHhIZaXMucoBSBpKFRkDFR")
	suffix := []byte("sba(BE(p7`\"0]%>5X),n1$?n>~%(6G+j")
	var ctx = NewHash(authdata)
	for i := 0; i < b.N; i++ {
		ctx.Sum(suffix)
	}
}

func TestCheckDificulty(t *testing.T) {
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
			name: "should return true",
			args: args{
				hash:      []byte{7, 195, 66, 190, 110, 86, 14, 127, 67, 132, 46, 46, 33, 183, 116, 230, 29, 133, 240, 71},
				dificulty: 1,
			},
			want: true,
		},
		{
			name: "should return false",
			args: args{
				hash:      []byte{7, 195, 66, 190, 110, 86, 14, 127, 67, 132, 46, 46, 33, 183, 116, 230, 29, 133, 240, 71},
				dificulty: 2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckDificulty1(tt.args.hash, tt.args.dificulty); got != tt.want {
				t.Errorf("CheckDificulty() = %v, want %v", got, tt.want)
			}
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
			if got := CheckDificulty(tt.args.hash, tt.args.dificulty); got != tt.want {
				t.Errorf("HexStartsWith3() = %v, want %v", got, tt.want)
			}
		})
	}
}
