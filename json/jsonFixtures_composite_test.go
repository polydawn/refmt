package json

import (
	"testing"

	"github.com/polydawn/refmt/tok/fixtures"
)

func testComposite(t *testing.T) {
	t.Run("array nested in map as non-first and final entry", func(t *testing.T) {
		seq := fixtures.SequenceMap["array nested in map as non-first and final entry"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `{"k1":"v1","ke":["oh","whee","wow"]}`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `{"k1":"v1","ke":["oh","whee","wow"]}`, nil)
		})
	})
	t.Run("array nested in map as first and non-final entry", func(t *testing.T) {
		seq := fixtures.SequenceMap["array nested in map as first and non-final entry"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `{"ke":["oh","whee","wow"],"k1":"v1"}`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `{"ke":["oh","whee","wow"],"k1":"v1"}`, nil)
		})
	})
	t.Run("maps nested in array", func(t *testing.T) {
		seq := fixtures.SequenceMap["maps nested in array"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `[{"k":"v"},"whee",{"k1":"v1"}]`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `[{"k":"v"},"whee",{"k1":"v1"}]`, nil)
		})
	})
	t.Run("arrays in arrays in arrays", func(t *testing.T) {
		seq := fixtures.SequenceMap["arrays in arrays in arrays"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `[[[]]]`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `[[[]]]`, nil)
		})
	})
	t.Run("maps nested in maps", func(t *testing.T) {
		seq := fixtures.SequenceMap["maps nested in maps"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `{"k":{"k2":"v2"}}`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `{"k":{"k2":"v2"}}`, nil)
		})
	})
}
