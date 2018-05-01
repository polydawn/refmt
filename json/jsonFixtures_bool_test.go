package json

import (
	"testing"

	"github.com/polydawn/refmt/tok/fixtures"
)

func testBoolEncoding(t *testing.T) {
	t.Run("encode bool", func(t *testing.T) {
		checkEncoding(t, fixtures.SequenceMap["true"], `true`, nil)
		checkEncoding(t, fixtures.SequenceMap["false"], `false`, nil)
	})
}

func testBoolDecoding(t *testing.T) {
	t.Run("decode bool", func(t *testing.T) {
		checkDecoding(t, fixtures.SequenceMap["true"], `true`, nil)
		checkDecoding(t, fixtures.SequenceMap["false"], `false`, nil)
	})
}
