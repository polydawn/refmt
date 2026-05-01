package cbor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"runtime"
	"testing"

	. "github.com/polydawn/refmt/tok"
)

// nextToken pulls one token via the decoder's Step method.
func nextToken(t *testing.T, opts DecodeOptions, payload []byte) (Token, bool, error) {
	t.Helper()
	d := NewDecoder(opts, bytes.NewReader(payload))
	var tk Token
	done, err := d.Step(&tk)
	return tk, done, err
}

func assertDecodeNoPanic(t *testing.T, opts DecodeOptions, payload []byte) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("decoder panicked for payload %x with opts %+v: %v", payload, opts, r)
		}
	}()
	_, _ = drainDecoder(NewDecoder(opts, bytes.NewReader(payload)))
}

func TestZeroToTwoByteInputs(t *testing.T) {
	cases := []struct {
		name    string
		payload []byte
	}{
		{name: "empty", payload: nil},
	}
	for i := 0; i <= 0xff; i++ {
		cases = append(cases, struct {
			name    string
			payload []byte
		}{
			name:    fmt.Sprintf("%02x", i),
			payload: []byte{byte(i)},
		})
	}
	for i := 0; i <= 0xff; i++ {
		for j := 0; j <= 0xff; j++ {
			cases = append(cases, struct {
				name    string
				payload []byte
			}{
				name:    fmt.Sprintf("%02x%02x", i, j),
				payload: []byte{byte(i), byte(j)},
			})
		}
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Run("default", func(t *testing.T) {
				assertDecodeNoPanic(t, DecodeOptions{}, tc.payload)
			})
			t.Run("reject_indefinite", func(t *testing.T) {
				assertDecodeNoPanic(t, DecodeOptions{RejectIndefinite: true}, tc.payload)
			})
		})
	}
}

func TestRejectIndefinite(t *testing.T) {
	cases := []struct {
		name    string
		payload []byte
	}{
		{"indef-bytes", []byte{0x5F, 0x41, 0x00, 0xFF}},
		{"indef-string", []byte{0x7F, 0x61, 'a', 0xFF}},
		{"indef-array", []byte{0x9F, 0xFF}},
		{"indef-map", []byte{0xBF, 0xFF}},
	}
	for _, tc := range cases {
		t.Run(tc.name+" rejected when option set", func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{RejectIndefinite: true}, tc.payload)
			if !errors.Is(err, ErrIndefiniteLength) {
				t.Fatalf("expected ErrIndefiniteLength, got %v", err)
			}
		})
		t.Run(tc.name+" accepted by default", func(t *testing.T) {
			_, _, err := nextToken(t, DecodeOptions{}, tc.payload)
			if err != nil {
				t.Fatalf("unexpected error decoding indefinite: %v", err)
			}
		})
	}
}

func TestMaxIndefiniteSize(t *testing.T) {
	// Build an indefinite-bytes payload of N chunks of 1 KiB each.
	build := func(chunks int) []byte {
		const chunkSize = 1024
		var buf bytes.Buffer
		buf.WriteByte(0x5F)
		chunk := append([]byte{0x59, 0x04, 0x00}, bytes.Repeat([]byte{0x00}, chunkSize)...) // bytes(1024)
		for i := 0; i < chunks; i++ {
			buf.Write(chunk)
		}
		buf.WriteByte(0xFF)
		return buf.Bytes()
	}

	t.Run("custom cap rejects when exceeded", func(t *testing.T) {
		payload := build(10) // 10 KiB declared total
		_, _, err := nextToken(t, DecodeOptions{MaxIndefiniteSize: 5 * 1024}, payload)
		if !errors.Is(err, ErrIndefiniteSizeExceeded) {
			t.Fatalf("expected ErrIndefiniteSizeExceeded, got %v", err)
		}
	})

	t.Run("custom cap accepts when within limit", func(t *testing.T) {
		payload := build(2) // 2 KiB total
		_, _, err := nextToken(t, DecodeOptions{MaxIndefiniteSize: 4 * 1024}, payload)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("default cap bounds peak allocation for over-size payload", func(t *testing.T) {
		// Build a payload that aggregates beyond the 32 MiB default total cap.
		// 33 KiB of 1 MiB chunks = 33 MiB total, just over the default.
		const chunkSize = 1 << 20
		const chunks = 33
		var buf bytes.Buffer
		buf.WriteByte(0x5F)
		chunkHeader := []byte{0x5A, 0x00, 0x10, 0x00, 0x00} // bytes(1 MiB) header
		for i := 0; i < chunks; i++ {
			buf.Write(chunkHeader)
			buf.Write(bytes.Repeat([]byte{0x00}, chunkSize))
		}
		buf.WriteByte(0xFF)

		runtime.GC()
		var before, after runtime.MemStats
		runtime.ReadMemStats(&before)
		_, _, err := nextToken(t, DecodeOptions{}, buf.Bytes())
		runtime.ReadMemStats(&after)

		if !errors.Is(err, ErrIndefiniteSizeExceeded) {
			t.Fatalf("expected ErrIndefiniteSizeExceeded, got %v", err)
		}
		// TotalAlloc accumulates every realloc, and refmt grows the
		// accumulator by doubling. The decoder also stages each chunk before
		// copying it into the aggregate buffer so truncated reads can return
		// partial data without preallocating the declared size. Bound at 7x
		// to allow for that extra copy while still asserting that we stay
		// tightly bounded relative to the cap (not the payload size).
		allocated := after.TotalAlloc - before.TotalAlloc
		const allowed = uint64(defaultMaxIndefiniteSize) * 7
		if allocated > allowed {
			t.Fatalf("allocation %d exceeded bound %d", allocated, allowed)
		}
		t.Logf("allocated %d bytes (cap %d, payload %d, allowed %d)",
			allocated, defaultMaxIndefiniteSize, buf.Len(), allowed)
	})
}

func TestTruncatedDefiniteBytes(t *testing.T) {
	payload := []byte{0x5a, 0x00, 0x30, 0x30, 0x30} // bytes(0x00303030) with no body

	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_, _, err := nextToken(t, DecodeOptions{}, payload)
	runtime.ReadMemStats(&after)

	if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected EOF-family error, got %v", err)
	}

	allocated := after.TotalAlloc - before.TotalAlloc
	const allowed = 256 << 10
	if allocated > allowed {
		t.Fatalf("allocation %d exceeded bound %d for truncated definite bytes", allocated, allowed)
	}
}

func TestTruncatedIndefiniteBytes(t *testing.T) {
	payload := []byte{0x5f, 0x58, 0x30, 0x30} // indefinite bytes, first chunk claims 48 bytes but has 1

	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_, _, err := nextToken(t, DecodeOptions{}, payload)
	runtime.ReadMemStats(&after)

	if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected EOF-family error, got %v", err)
	}

	allocated := after.TotalAlloc - before.TotalAlloc
	const allowed = 256 << 10
	if allocated > allowed {
		t.Fatalf("allocation %d exceeded bound %d for truncated indefinite bytes", allocated, allowed)
	}
}
