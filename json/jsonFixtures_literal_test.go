package json

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/polydawn/refmt/tok"
)

func TestMalformedLiteralsRejected(t *testing.T) {
	cases := []struct {
		payload []byte
		eof     bool
	}{
		{[]byte("n"), true},
		{[]byte("nu"), true},
		{[]byte("nul"), true},
		{[]byte("nall"), false},
		{[]byte("t"), true},
		{[]byte("tr"), true},
		{[]byte("tru"), true},
		{[]byte("tre"), false},
		{[]byte("f"), true},
		{[]byte("fa"), true},
		{[]byte("fal"), true},
		{[]byte("fals"), true},
		{[]byte("folse"), false},
	}
	for _, tc := range cases {
		t.Run(string(tc.payload), func(t *testing.T) {
			var token tok.Token
			_, err := NewDecoder(bytes.NewReader(tc.payload)).Step(&token)
			if err == nil {
				t.Fatalf("expected malformed literal %q to be rejected", tc.payload)
			}
			if tc.eof && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
				t.Fatalf("expected EOF-family error for %q, got %v", tc.payload, err)
			}
		})
	}
}
