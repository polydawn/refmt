package json

import (
	"testing"

	"github.com/polydawn/refmt/tok/fixtures"
)

func testMap(t *testing.T) {
	t.Run("empty map", func(t *testing.T) {
		seq := fixtures.SequenceMap["empty map"]
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, seq, `{}`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkEncoding(t, seq, `{}`, nil)
		})
		t.Run("decode with extra interior whitespace", func(t *testing.T) {
			checkDecoding(t, seq, `{  }`, nil)
		})
		t.Run("decode with extra flanking whitespace", func(t *testing.T) {
			checkDecoding(t, seq, `  {  }  `, nil)
		})
	})
}
