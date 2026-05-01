package cbor

import (
	"bytes"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/polydawn/refmt/shared"
	. "github.com/polydawn/refmt/tok"
)

const (
	// fuzzTokenCountFloor is the smallest input size for which the
	// per-byte multiplier alone would be too tight: even a 1-byte
	// payload like 0x80 (empty array) emits TArrOpen + TArrClose.
	fuzzTokenCountFloor = 64

	// fuzzTokenCountPerByte caps the rate at which the decoder may emit
	// tokens relative to the input length. The smallest CBOR token
	// occupies one wire byte (e.g. 0x00, uint(0)), and pairs like
	// 0x80/0x9f-0xff produce at most 2 tokens per wire byte; 4 leaves
	// generous headroom so legitimate inputs never trip the bound.
	fuzzTokenCountPerByte = 4

	// fuzzScalarSlackBytes is the small margin allowed on scalar (string
	// and byte-string) payload size relative to the input length. CBOR
	// strings and bytes are stored 1:1 on the wire plus a header.
	fuzzScalarSlackBytes = 32

	// fuzzTokenCountCap bounds how many tokens we drain from a single
	// decode iteration regardless of the per-byte rate, so a fuzz
	// iteration is always bounded in time.
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
		done, stepErr := d.Step(&tk)
		stats.tokenCount++
		switch tk.Type {
		case TString:
			stats.scalarBytes += len(tk.Str)
		case TBytes:
			stats.scalarBytes += len(tk.Bytes)
		}
		if stepErr != nil {
			return stats, stepErr
		}
		if done {
			return stats, nil
		}
		if stats.tokenCount > fuzzTokenCountCap {
			return stats, nil
		}
	}
}

func FuzzCborDecode(f *testing.F) {
	seedCborCorpus(f)
	f.Fuzz(func(t *testing.T, payload []byte) {
		stats, _ := drainDecoder(NewDecoder(DecodeOptions{}, bytes.NewReader(payload)))
		assertDecodeBounds(t, payload, stats)

		stats, _ = drainDecoder(NewDecoder(DecodeOptions{RejectIndefinite: true}, bytes.NewReader(payload)))
		assertDecodeBounds(t, payload, stats)

		// Tight indefinite cap exercises the accumulator-bound path.
		stats, _ = drainDecoder(NewDecoder(DecodeOptions{MaxIndefiniteSize: 4096}, bytes.NewReader(payload)))
		assertDecodeBounds(t, payload, stats)
	})
}

// FuzzCborRoundtrip asserts the decode/re-encode pipeline is stable: if
// a payload decodes without error, encoding the resulting token stream
// and decoding+re-encoding that output must produce identical bytes.
// Catches asymmetries where the decoder accepts more than the encoder
// can emit (or vice versa).
func FuzzCborRoundtrip(f *testing.F) {
	seedCborCorpus(f)
	f.Fuzz(func(t *testing.T, payload []byte) {
		first, err := reencode(payload)
		if err != nil {
			return
		}
		second, err := reencode(first)
		if err != nil {
			t.Fatalf("re-decode of canonicalized output failed: %v\ncanonical: %x", err, first)
		}
		if !bytes.Equal(first, second) {
			t.Fatalf("roundtrip not stable:\nfirst:  %x\nsecond: %x", first, second)
		}
	})
}

func reencode(payload []byte) ([]byte, error) {
	var buf bytes.Buffer
	pump := shared.TokenPump{
		TokenSource: NewDecoder(DecodeOptions{}, bytes.NewReader(payload)),
		TokenSink:   NewEncoder(&buf),
	}
	if err := pump.Run(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func FuzzCborUnmarshalInterface(f *testing.F) {
	seedCborCorpus(f)
	f.Fuzz(func(t *testing.T, payload []byte) {
		var out interface{}
		_ = Unmarshal(DecodeOptions{}, payload, &out)

		out = nil
		_ = Unmarshal(DecodeOptions{RejectIndefinite: true}, payload, &out)
	})
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

	maxScalarBytes := len(payload) + fuzzScalarSlackBytes
	if stats.scalarBytes > maxScalarBytes {
		t.Fatalf("decoder expanded scalar payload beyond encoded size budget: scalar=%d max=%d len=%d", stats.scalarBytes, maxScalarBytes, len(payload))
	}
}

func seedCborCorpus(f *testing.F) {
	f.Helper()
	seen := make(map[string]struct{})
	add := func(payload []byte) {
		if len(payload) == 0 {
			return
		}
		key := string(payload)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		f.Add(append([]byte(nil), payload...))
	}

	for _, payload := range handcraftedCborSeeds() {
		add(payload)
	}
	walkCorpus(f, "testdata/fuzz_corpus", add)
}

// handcraftedCborSeeds covers small edge cases we want guaranteed in the
// seed set independent of the vendored corpus: bare sigils, indefinite
// frames, oversized lengths that should error out, and a few CBOR shapes
// that historically tripped the decoder.
func handcraftedCborSeeds() [][]byte {
	return [][]byte{
		{0x00},
		{0xf6},
		{0x80},
		{0xa0},
		{0x9f, 0xff},
		{0xbf, 0xff},
		{0x5f, 0xff},
		{0x7f, 0xff},
		{0x9f, 0x9f, 0x9f, 0xff, 0xff, 0xff},
		{0xbf, 0x61, 0x61, 0x9f, 0xff, 0xff},
		{0xd8, 0x18, 0x43, 0x01, 0x02, 0x03},
		{0xc1, 0xc1, 0xf6},
		{0xff},
		{0x81, 0xff},
		{0xa1, 0x61, 0x61, 0xff},
		{0xfb},
		{0xfa, 0x7f},
		{0x5a, 0xff, 0xff, 0xff, 0xff},
		{0x7a, 0xff, 0xff, 0xff, 0xff},
		{0x9a, 0xff, 0xff, 0xff, 0xff},
		{0xba, 0xff, 0xff, 0xff, 0xff},
		mustHex("a3636261720363666f6f0163666f6f02"),
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

func mustHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
