package cbor

import (
	"errors"
	"testing"
)

func TestRejectNonMinimalInteger_NonMinimal(t *testing.T) {
	cases := []struct {
		name    string
		payload []byte
	}{
		{"uint zero in 1-byte form", []byte{0x18, 0x00}},
		{"uint 23 in 1-byte form", []byte{0x18, 0x17}},
		{"uint 255 in 2-byte form", []byte{0x19, 0x00, 0xff}},
		{"uint 0xffff in 4-byte form", []byte{0x1a, 0x00, 0x00, 0xff, 0xff}},
		{"uint 0xffffffff in 8-byte form", []byte{0x1b, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}},
		{"negint -1 in 1-byte form", []byte{0x38, 0x00}},
		{"bytes len 0 in 1-byte form", []byte{0x58, 0x00}},
		{"string len 0 in 1-byte form", []byte{0x78, 0x00}},
		{"array len 0 in 1-byte form", []byte{0x98, 0x00}},
		{"map len 0 in 1-byte form", []byte{0xb8, 0x00}},
		{"tag 42 in 2-byte form", []byte{0xd9, 0x00, 0x2a, 0x40}},
		{"tag 42 in 4-byte form", []byte{0xda, 0x00, 0x00, 0x00, 0x2a, 0x40}},
		{"tag 42 in 8-byte form", []byte{0xdb, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2a, 0x40}},
	}
	for _, tc := range cases {
		t.Run(tc.name+" rejected when option set", func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{RejectNonMinimalInteger: true}, tc.payload)
			if !errors.Is(err, ErrNonMinimalInteger) {
				t.Fatalf("expected ErrNonMinimalInteger, got %v", err)
			}
		})
		t.Run(tc.name+" accepted by default", func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{}, tc.payload)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRejectNonMinimalInteger_MinimalRoundtrip(t *testing.T) {
	cases := []struct {
		name    string
		payload []byte
	}{
		{"uint inline 0", []byte{0x00}},
		{"uint inline 23", []byte{0x17}},
		{"uint 1-byte 24", []byte{0x18, 0x18}},
		{"uint 1-byte 255", []byte{0x18, 0xff}},
		{"uint 2-byte 256", []byte{0x19, 0x01, 0x00}},
		{"uint 2-byte 0xffff", []byte{0x19, 0xff, 0xff}},
		{"uint 4-byte 0x10000", []byte{0x1a, 0x00, 0x01, 0x00, 0x00}},
		{"uint 4-byte 0xffffffff", []byte{0x1a, 0xff, 0xff, 0xff, 0xff}},
		{"uint 8-byte 0x100000000", []byte{0x1b, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}},
		{"tag 42 minimal", []byte{0xd8, 0x2a, 0x40}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{RejectNonMinimalInteger: true}, tc.payload)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
