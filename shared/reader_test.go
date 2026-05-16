package shared

import (
	"bytes"
	"errors"
	"io"
	"runtime"
	"testing"
)

func TestSlickReaderReadMethods(t *testing.T) {
	cases := []struct {
		name string
		new  func([]byte) SlickReader
	}{
		{
			name: "stream",
			new:  func(payload []byte) SlickReader { return NewReader(bytes.NewReader(payload)) },
		},
		{
			name: "bytes-buffer",
			new: func(payload []byte) SlickReader {
				return NewBytesReader(bytes.NewBuffer(payload))
			},
		},
		{
			name: "slice",
			new:  func(payload []byte) SlickReader { return NewSliceReader(payload) },
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.new([]byte("abcdef"))

			empty, err := r.Readn(0)
			if err != nil {
				t.Fatalf("zero-length Readn returned error: %v", err)
			}
			if len(empty) != 0 || cap(empty) != 0 || r.NumRead() != 0 {
				t.Fatalf("zero-length Readn = len %d cap %d NumRead %d, want all zero", len(empty), cap(empty), r.NumRead())
			}

			bs, err := r.Readnzc(2)
			if err != nil {
				t.Fatalf("Readnzc: %v", err)
			}
			if string(bs) != "ab" || r.NumRead() != 2 {
				t.Fatalf("Readnzc = %q NumRead %d, want ab/2", bs, r.NumRead())
			}

			copyOut, err := r.Readn(2)
			if err != nil {
				t.Fatalf("Readn: %v", err)
			}
			if string(copyOut) != "cd" || r.NumRead() != 4 {
				t.Fatalf("Readn = %q NumRead %d, want cd/4", copyOut, r.NumRead())
			}

			dst := []byte("xx")
			if err := r.Readb(dst); err != nil {
				t.Fatalf("Readb: %v", err)
			}
			if string(dst) != "ef" || r.NumRead() != 6 {
				t.Fatalf("Readb filled %q NumRead %d, want ef/6", dst, r.NumRead())
			}
		})
	}
}

func TestSlickReaderReadnReturnsStableCopy(t *testing.T) {
	payload := []byte("abcd")
	r := NewSliceReader(payload)

	bs, err := r.Readn(2)
	if err != nil {
		t.Fatalf("Readn: %v", err)
	}
	payload[0], payload[1] = 'x', 'y'
	if string(bs) != "ab" {
		t.Fatalf("Readn returned slice backed by input: got %q", bs)
	}
}

func TestSlickReaderUnreadAndTrack(t *testing.T) {
	cases := []struct {
		name string
		new  func([]byte) SlickReader
	}{
		{
			name: "stream",
			new:  func(payload []byte) SlickReader { return NewReader(bytes.NewReader(payload)) },
		},
		{
			name: "slice",
			new:  func(payload []byte) SlickReader { return NewSliceReader(payload) },
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.new([]byte("abcdef"))
			r.Track()

			b, err := r.Readn1()
			if err != nil {
				t.Fatalf("Readn1: %v", err)
			}
			if b != 'a' {
				t.Fatalf("Readn1 = %q, want a", b)
			}
			r.Unreadn1()
			if r.NumRead() != 0 {
				t.Fatalf("NumRead after unread = %d, want 0", r.NumRead())
			}

			bs, err := r.Readnzc(3)
			if err != nil {
				t.Fatalf("Readnzc after unread: %v", err)
			}
			if string(bs) != "abc" {
				t.Fatalf("Readnzc after unread = %q, want abc", bs)
			}
			b, err = r.Readn1()
			if err != nil {
				t.Fatalf("second Readn1: %v", err)
			}
			if b != 'd' {
				t.Fatalf("second Readn1 = %q, want d", b)
			}

			tracked := r.StopTrack()
			if string(tracked) != "abcd" {
				t.Fatalf("tracked bytes = %q, want abcd", tracked)
			}
		})
	}
}

func TestSlickReaderStreamPartialReadIsBounded(t *testing.T) {
	payload := bytes.Repeat([]byte{'x'}, 512)
	r := NewReader(bytes.NewReader(payload))

	runtime.GC()
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	bs, err := r.Readn(1 << 20)
	runtime.ReadMemStats(&after)

	if !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		t.Fatalf("Readn error = %v, want EOF-family error", err)
	}
	if len(bs) != len(payload) || string(bs) != string(payload) {
		t.Fatalf("Readn returned %d bytes, want original %d-byte payload", len(bs), len(payload))
	}
	if r.NumRead() != len(payload) {
		t.Fatalf("NumRead = %d, want %d", r.NumRead(), len(payload))
	}

	allocated := after.TotalAlloc - before.TotalAlloc
	const allowed = 64 << 10
	if allocated > allowed {
		t.Fatalf("Readn allocated %d bytes for truncated 1 MiB request, want <= %d", allocated, allowed)
	}
}

func TestSlickReaderStreamReadbTracksOnlyBytesRead(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte("ab")))
	r.Track()

	dst := []byte("xxxx")
	err := r.Readb(dst)
	if !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		t.Fatalf("Readb error = %v, want EOF-family error", err)
	}
	if string(dst[:2]) != "ab" {
		t.Fatalf("Readb prefix = %q, want ab", dst[:2])
	}
	if tracked := r.StopTrack(); string(tracked) != "ab" {
		t.Fatalf("tracked bytes after partial Readb = %q, want ab", tracked)
	}
	if r.NumRead() != 2 {
		t.Fatalf("NumRead after partial Readb = %d, want 2", r.NumRead())
	}
}

func TestSlickReaderSliceTruncatedReadDoesNotAdvance(t *testing.T) {
	r := NewSliceReader([]byte("abc"))

	bs, err := r.Readnzc(4)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("Readnzc error = %v, want unexpected EOF", err)
	}
	if len(bs) != 0 || r.NumRead() != 0 {
		t.Fatalf("truncated Readnzc returned len %d NumRead %d, want 0/0", len(bs), r.NumRead())
	}

	b, err := r.Readn1()
	if err != nil {
		t.Fatalf("Readn1 after truncated Readnzc: %v", err)
	}
	if b != 'a' {
		t.Fatalf("Readn1 after truncated Readnzc = %q, want a", b)
	}
}

func TestReaderToScannerUnreadAcrossBulkRead(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte("abcdef")))

	bs, err := r.Readnzc(3)
	if err != nil {
		t.Fatalf("Readnzc: %v", err)
	}
	if string(bs) != "abc" {
		t.Fatalf("Readnzc = %q, want abc", bs)
	}
	r.Unreadn1()

	b, err := r.Readn1()
	if err != nil {
		t.Fatalf("Readn1 after unread: %v", err)
	}
	if b != 'c' {
		t.Fatalf("Readn1 after unread = %q, want c", b)
	}
	bs, err = r.Readn(3)
	if err != nil {
		t.Fatalf("Readn after unread: %v", err)
	}
	if string(bs) != "def" {
		t.Fatalf("Readn after unread = %q, want def", bs)
	}
}
