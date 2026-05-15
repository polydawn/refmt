package cbor

import (
	"errors"
	"testing"
)

func TestRejectNaN(t *testing.T) {
	cases := []struct {
		name    string
		payload []byte
	}{
		{"f16 canonical NaN", []byte{0xf9, 0x7e, 0x00}},
		{"f16 alt NaN", []byte{0xf9, 0x7f, 0xf8}},
		{"f32 NaN", []byte{0xfa, 0x7f, 0xc0, 0x00, 0x00}},
		{"f64 NaN", []byte{0xfb, 0x7f, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"f64 alt NaN payload", []byte{0xfb, 0x7f, 0xf8, 0xca, 0xfe, 0xde, 0xad, 0xbe, 0xef}},
	}
	for _, tc := range cases {
		t.Run(tc.name+" rejected when option set", func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{RejectNaN: true}, tc.payload)
			if !errors.Is(err, ErrFloatNaN) {
				t.Fatalf("expected ErrFloatNaN, got %v", err)
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

func TestRejectInfinity(t *testing.T) {
	cases := []struct {
		name    string
		payload []byte
	}{
		{"f16 +Inf", []byte{0xf9, 0x7c, 0x00}},
		{"f16 -Inf", []byte{0xf9, 0xfc, 0x00}},
		{"f32 +Inf", []byte{0xfa, 0x7f, 0x80, 0x00, 0x00}},
		{"f32 -Inf", []byte{0xfa, 0xff, 0x80, 0x00, 0x00}},
		{"f64 +Inf", []byte{0xfb, 0x7f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"f64 -Inf", []byte{0xfb, 0xff, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}
	for _, tc := range cases {
		t.Run(tc.name+" rejected when option set", func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{RejectInfinity: true}, tc.payload)
			if !errors.Is(err, ErrFloatInfinity) {
				t.Fatalf("expected ErrFloatInfinity, got %v", err)
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

func TestRejectFloat_FiniteValuesUnaffected(t *testing.T) {
	cases := []struct {
		name    string
		payload []byte
	}{
		{"f16 1.0", []byte{0xf9, 0x3c, 0x00}},
		{"f16 -1.0", []byte{0xf9, 0xbc, 0x00}},
		{"f32 1.5", []byte{0xfa, 0x3f, 0xc0, 0x00, 0x00}},
		{"f64 1.5", []byte{0xfb, 0x3f, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"f64 zero", []byte{0xfb, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{RejectNaN: true, RejectInfinity: true}, tc.payload)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
