package json

import (
	"testing"

	"github.com/polydawn/refmt/tok/fixtures"
)

func testBool(t *testing.T) {
	t.Run("bool true", func(t *testing.T) {
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, fixtures.SequenceMap["true"], `true`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, fixtures.SequenceMap["true"], `true`, nil)
		})
	})
	t.Run("bool false", func(t *testing.T) {
		t.Run("encode", func(t *testing.T) {
			checkEncoding(t, fixtures.SequenceMap["false"], `false`, nil)
		})
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, fixtures.SequenceMap["false"], `false`, nil)
		})
	})
}
