package json

import (
	"testing"

	"github.com/polydawn/refmt/tok/fixtures"
)

func testMapEncoding(t *testing.T) {
	t.Run("encode map", func(t *testing.T) {
		t.Run("empty map", func(t *testing.T) {
			checkEncoding(t, fixtures.SequenceMap["empty map"],
				`{}`, nil)
		})
	})
}

func testMapDecoding(t *testing.T) {
	t.Run("decode map", func(t *testing.T) {
		t.Run("empty map", func(t *testing.T) {
			fix := fixtures.SequenceMap["empty map"]
			checkEncoding(t, fix,
				`{}`, nil)
			t.Run("with extra interior whitespace", func(t *testing.T) {
				checkDecoding(t, fix,
					`{  }`, nil)
			})
			t.Run("with extra flanking whitespace", func(t *testing.T) {
				checkDecoding(t, fix,
					`  {  }  `, nil)
			})
		})
	})
}
