// TODO: Refactor this package and eliminate utils.
package utils

import (
	"bytes"
	"math/rand"
	"testing"
)

func Test_RandomUTF8(t *testing.T) {
	type args struct {
		randomString []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    string
	}{
		{
			name:    "Random string",
			args:    args{randomString: make([]byte, 32)},
			wantErr: false,
			want:    "sba(BE(p7`\"0]%>5X),n1$?n>~%(6G+j",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rand.Seed(1)
			if err := RandomUTF8(tt.args.randomString); (err != nil) != tt.wantErr {
				t.Errorf("RandomUTF8() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !bytes.Equal(tt.args.randomString, []byte(tt.want)) {
				t.Errorf("RandomUTF8() = `%v`, want `%v`", string(tt.args.randomString), tt.want)
			}
		})
	}
}

func BenchmarkSecureRandomStrings32(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SecureRandomString(32)
	}
}

func BenchmarkRandStringRunes32(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandStringRunes(32)
	}
}

func BenchmarkRandASCIIBytes32(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandASCIIBytes(32)
	}
}

func BenchmarkRandomUTF8_32(b *testing.B) {
	suffix := make([]byte, 32)
	for n := 0; n < b.N; n++ {
		RandomUTF8(suffix)
	}
}

func BenchmarkSecureRandomStrings256(b *testing.B) {
	for n := 0; n < b.N; n++ {
		SecureRandomString(256)
	}
}

func BenchmarkRandStringRunes256(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandStringRunes(256)
	}
}

func BenchmarkRandASCIIBytes256(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandASCIIBytes(256)
	}
}

func BenchmarkRandomUTF8_256(b *testing.B) {
	suffix := make([]byte, 256)
	for n := 0; n < b.N; n++ {
		RandomUTF8(suffix)
	}
}
