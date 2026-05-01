package json

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	. "github.com/polydawn/refmt/tok"
)

const (
	// fuzzTokenCountFloor handles small inputs where the per-byte
	// multiplier alone is too tight (e.g. "[]" yields TArrOpen +
	// TArrClose, two tokens for two wire bytes).
	fuzzTokenCountFloor = 64

	// fuzzTokenCountPerByte caps token emission relative to the input
	// length. JSON tokens consume at least one wire byte each (the
	// shortest scalar like a single digit, or a single bracket); 4
	// leaves generous headroom so legitimate inputs never trip the bound.
	fuzzTokenCountPerByte = 4

	// fuzzScalarMultiplier accounts for valid expansion paths from wire
	// bytes to decoded string bytes. The dominant case is invalid UTF-8
	// bytes inside a quoted string getting replaced with the Unicode
	// replacement character U+FFFD, which is three UTF-8 bytes per
	// substituted input byte. Escape sequences (\n, \uXXXX, surrogate
	// pairs) all shrink rather than expand, so 3x is the legitimate
	// upper bound on scalar/wire ratio.
	fuzzScalarMultiplier = 3

	// fuzzScalarSlackBytes covers the constant per-decode accounting
	// margin on top of the multiplier (e.g. quotes, structural chars
	// counted by the wire side but not by scalar accumulation).
	fuzzScalarSlackBytes = 32

	// fuzzTokenCountCap bounds drain-loop iterations within a single
	// fuzz sample regardless of the per-byte rate.
	fuzzTokenCountCap = 1 << 20
)

type decodeStats struct {
	tokenCount  int
	scalarBytes int
}

// drainDecoder pumps Step until done or error. It accepts no panic; any
// panic surfaces as a fuzz failure.
func drainDecoder(d *Decoder) (decodeStats, error) {
	var stats decodeStats
	for {
		var tk Token
		done, err := d.Step(&tk)
		stats.tokenCount++
		if tk.Type == TString {
			stats.scalarBytes += len(tk.Str)
		}
		if err != nil {
			return stats, err
		}
		if done {
			return stats, nil
		}
		if stats.tokenCount > fuzzTokenCountCap {
			return stats, nil
		}
	}
}

func assertDecodeNoPanic(t *testing.T, payload []byte) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("json decoder panicked for payload %x: %v", payload, r)
		}
	}()
	_, _ = drainDecoder(NewDecoder(bytes.NewReader(payload)))
}

func assertDecodeBounds(t *testing.T, payload []byte, stats decodeStats) {
	t.Helper()

	maxTokens := fuzzTokenCountFloor
	if n := len(payload)*fuzzTokenCountPerByte + fuzzTokenCountFloor; n > maxTokens {
		maxTokens = n
	}
	if stats.tokenCount > maxTokens && stats.tokenCount < fuzzTokenCountCap {
		t.Fatalf("decoder produced too many tokens for input size: tokens=%d max=%d len=%d", stats.tokenCount, maxTokens, len(payload))
	}

	maxScalarBytes := len(payload)*fuzzScalarMultiplier + fuzzScalarSlackBytes
	if stats.scalarBytes > maxScalarBytes {
		t.Fatalf("decoder expanded scalar payload beyond encoded size budget: scalar=%d max=%d len=%d", stats.scalarBytes, maxScalarBytes, len(payload))
	}
}

func TestZeroToTwoByteInputs(t *testing.T) {
	assertDecodeNoPanic(t, nil)
	for i := 0; i <= 0xff; i++ {
		assertDecodeNoPanic(t, []byte{byte(i)})
	}
	for i := 0; i <= 0xff; i++ {
		for j := 0; j <= 0xff; j++ {
			assertDecodeNoPanic(t, []byte{byte(i), byte(j)})
		}
	}
}

func FuzzJSONDecode(f *testing.F) {
	seedJSONCorpus(f)
	f.Fuzz(func(t *testing.T, payload []byte) {
		stats, _ := drainDecoder(NewDecoder(bytes.NewReader(payload)))
		assertDecodeBounds(t, payload, stats)
	})
}

func FuzzJSONUnmarshalInterface(f *testing.F) {
	seedJSONCorpus(f)
	f.Fuzz(func(t *testing.T, payload []byte) {
		var out interface{}
		_ = Unmarshal(payload, &out)
	})
}

func seedJSONCorpus(f *testing.F) {
	f.Helper()
	seen := make(map[string]struct{})
	add := func(payload []byte) {
		key := string(payload)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		f.Add(append([]byte(nil), payload...))
	}

	add(nil)
	for _, payload := range handcraftedJSONSeeds() {
		add(payload)
	}
	walkCorpus(f, "testdata/fuzz_corpus", add)
}

// handcraftedJSONSeeds covers small edge cases we want guaranteed
// in the seed set: bare scalars, malformed escapes, truncated structures,
// and a few inputs that historically tripped the JSON parser.
func handcraftedJSONSeeds() [][]byte {
	return [][]byte{
		[]byte("null"),
		[]byte("true"),
		[]byte("false"),
		[]byte(`""`),
		[]byte(`[]`),
		[]byte(`{}`),
		[]byte(`{"foo":1,"foo":2}`),
		[]byte(`{"foo":[[[[[null]]]]]}`),
		[]byte(`{"foo":"\uD800"}`),
		[]byte(`{"foo":"\uD800\uDC00"}`),
		[]byte(`{"foo":"\x"}`),
		[]byte(`{"foo":"\u00"}`),
		[]byte(`{"foo":"\uZZZZ"}`),
		[]byte(`{"foo":-}`),
		[]byte(`{"foo":1e}`),
		[]byte(`{"foo":01}`),
		[]byte(`{"foo":`),
		[]byte(`["`),
		[]byte("\xff"),
	}
}

func walkCorpus(f *testing.F, root string, add func([]byte)) {
	f.Helper()
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		payload, readErr := os.ReadFile(path)
		if readErr != nil {
			f.Fatalf("read corpus file %s: %v", path, readErr)
		}
		add(payload)
		return nil
	})
	if err != nil {
		f.Fatalf("walk corpus %s: %v", root, err)
	}
}
