package json

import (
	"strconv"
	"testing"

	. "github.com/polydawn/refmt/tok"
	"github.com/polydawn/refmt/tok/fixtures"
)

func testNumber(t *testing.T) {
	t.Run("integer zero", func(t *testing.T) {
		seq := fixtures.Sequence{Tokens: fixtures.Tokens{{Type: TInt, Int: 0}}}
		checkCanonical(t, seq, "0")
	})
	t.Run("integer one", func(t *testing.T) {
		seq := fixtures.Sequence{Tokens: fixtures.Tokens{{Type: TInt, Int: 1}}}
		checkCanonical(t, seq, "1")
	})
	t.Run("integer neg 1", func(t *testing.T) {
		seq := fixtures.Sequence{Tokens: fixtures.Tokens{{Type: TInt, Int: -1}}}
		checkCanonical(t, seq, "-1")
	})
	t.Run("integer neg 100", func(t *testing.T) {
		seq := fixtures.Sequence{Tokens: fixtures.Tokens{{Type: TInt, Int: -100}}}
		checkCanonical(t, seq, "-100")
	})
	t.Run("integer 1000000", func(t *testing.T) {
		seq := fixtures.Sequence{Tokens: fixtures.Tokens{{Type: TInt, Int: 1000000}}}
		checkCanonical(t, seq, "1000000")
	})
	t.Run("float 1 e+100", func(t *testing.T) {
		seq := fixtures.Sequence{Tokens: fixtures.Tokens{{Type: TFloat64, Float64: 1.0e+300}}}
		checkCanonical(t, seq, `1e+300`)
	})
	t.Run("integer too big to parse", func(t *testing.T) {
		seq := fixtures.Sequence{Tokens: fixtures.Tokens{{Type: TInt, Int: 2<<62 - 1}}}
		t.Run("decode", func(t *testing.T) {
			checkDecoding(t, seq, `18446744073709551617`, &strconv.NumError{"ParseInt", "18446744073709551617", strconv.ErrRange})
		})
	})
}
