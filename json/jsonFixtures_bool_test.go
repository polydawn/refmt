package json

import (
	"testing"

	"github.com/polydawn/refmt/tok/fixtures"
)

func testBoolEncoding(t *testing.T) {
	t.Run("encode bool", func(t *testing.T) {
		checkEncoding(t, fixtures.SequenceMap["true"], `true`, nil)
	})
}
