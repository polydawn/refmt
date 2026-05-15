package cbor

import (
	"errors"
	"testing"
)

func TestRejectNarrowFloat(t *testing.T) {
	rejected := []struct {
		name    string
		payload []byte
	}{
		{"f16 1.0", []byte{0xf9, 0x3c, 0x00}},
		{"f16 +Inf", []byte{0xf9, 0x7c, 0x00}},
		{"f16 NaN", []byte{0xf9, 0x7e, 0x00}},
		{"f32 1.5", []byte{0xfa, 0x3f, 0xc0, 0x00, 0x00}},
		{"f32 +Inf", []byte{0xfa, 0x7f, 0x80, 0x00, 0x00}},
	}
	for _, tc := range rejected {
		t.Run(tc.name+" rejected when option set", func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{RejectNarrowFloat: true}, tc.payload)
			if !errors.Is(err, ErrNarrowFloat) {
				t.Fatalf("expected ErrNarrowFloat, got %v", err)
			}
		})
		t.Run(tc.name+" accepted by default", func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{}, tc.payload)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}

	accepted := []struct {
		name    string
		payload []byte
	}{
		{"f64 1.5", []byte{0xfb, 0x3f, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"f64 zero", []byte{0xfb, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}
	for _, tc := range accepted {
		t.Run(tc.name+" still accepted", func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{RejectNarrowFloat: true}, tc.payload)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// Sigil byte alone, no payload: option must reject before any read so a
// truncated stream surfaces ErrNarrowFloat, not io.EOF.
func TestRejectNarrowFloat_RejectsBeforeRead(t *testing.T) {
	cases := []struct {
		name    string
		payload []byte
	}{
		{"f16 sigil only", []byte{0xf9}},
		{"f32 sigil only", []byte{0xfa}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{RejectNarrowFloat: true}, tc.payload)
			if !errors.Is(err, ErrNarrowFloat) {
				t.Fatalf("expected ErrNarrowFloat, got %v", err)
			}
		})
	}
}
