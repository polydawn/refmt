package json

import (
	"testing"

	"github.com/polydawn/refmt/tok/fixtures"
)

func testString(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		seq := fixtures.SequenceMap["empty string"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `""`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `""`, nil)
		})
		t.Run("decode with extra whitespace", func(t *testing.T) {
			checkDecoding(t, seq, `  "" `, nil)
		})
	})
	t.Run("flat string", func(t *testing.T) {
		seq := fixtures.SequenceMap["flat string"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `"value"`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `"value"`, nil)
		})
	})
	t.Run("strings needing escape", func(t *testing.T) {
		seq := fixtures.SequenceMap["strings needing escape"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `"str\nbroken\ttabbed"`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `"str\nbroken\ttabbed"`, nil)
		})
	})
}
