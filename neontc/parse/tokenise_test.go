package parse

import (
	"bytes"
	"reflect"
	"testing"
)

func FuzzTokens(f *testing.F) {
	f.Add([]byte("hello world\n{{ func bananas }} \\{{ \\}} {{ endfunc }}"))
	f.Fuzz(func(t *testing.T, input []byte) {
		output, err := tokens(input)
		if err != nil {
			// invalid input for whatever reason
			return
		}

		joined := bytes.Join(output, nil)

		if !bytes.Equal(input, joined) {
			t.Errorf("output mangled: %v vs %v", input, joined)
		}
	})
}

func TestTokens(t *testing.T) {
	tests := []struct {
		name    string
		args    []byte
		want    [][]byte
		wantErr bool
	}{
		{"Normal", []byte("hello world\n{{ func bananas }} \\{{ \\}} {{ endfunc }} trailing"), [][]byte{[]byte("hello world\n"), []byte("{{ func bananas }}"), []byte(" \\{{ \\}} "), []byte("{{ endfunc }}"), []byte(" trailing")}, false},
		{"Really short", []byte{'0'}, [][]byte{{'0'}}, false},
		{"Starting with curlies", []byte("{{ func }} bananas"), [][]byte{[]byte("{{ func }}"), []byte(" bananas")}, false},
		{"Ending with curlies", []byte("bananas {{ func }}"), [][]byte{[]byte("bananas "), []byte("{{ func }}")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokens(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokens() got = %v, want %v", got, tt.want)
			}
		})
	}
}
