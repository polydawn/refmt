package json

import (
	"testing"

	"github.com/polydawn/refmt/tok/fixtures"
)

func testArray(t *testing.T) {
	t.Run("empty array", func(t *testing.T) {
		seq := fixtures.SequenceMap["empty array"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `[]`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `[]`, nil)
		})
		t.Run("decode with extra whitespace", func(t *testing.T) {
			checkDecoding(t, seq, `  [ ] `, nil)
		})
	})
	t.Run("single entry array", func(t *testing.T) {
		seq := fixtures.SequenceMap["single entry array"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `["value"]`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `["value"]`, nil)
		})
		t.Run("decode with extra whitespace", func(t *testing.T) {
			checkDecoding(t, seq, `  [ "value" ] `, nil)
		})
	})
	t.Run("duo entry array", func(t *testing.T) {
		seq := fixtures.SequenceMap["duo entry array"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `["value","v2"]`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `["value",  "v2"]`, nil)
		})
	})
}
