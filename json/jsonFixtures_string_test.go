package json

import (
	"testing"

	"github.com/polydawn/refmt/tok/fixtures"
)

func testStringEncoding(t *testing.T) {
	t.Run("encode string", func(t *testing.T) {
		t.Run("empty string", func(t *testing.T) {
			checkEncoding(t, fixtures.SequenceMap["empty string"],
				`""`, nil)
		})
		t.Run("flat string", func(t *testing.T) {
			checkEncoding(t, fixtures.SequenceMap["flat string"],
				`"value"`, nil)
		})
		t.Run("strings needing escape", func(t *testing.T) {
			checkEncoding(t, fixtures.SequenceMap["strings needing escape"],
				`"str\nbroken\ttabbed"`, nil)
		})
	})
}

func testStringDecoding(t *testing.T) {
	t.Run("decode string", func(t *testing.T) {
		t.Run("empty string", func(t *testing.T) {
			checkDecoding(t, fixtures.SequenceMap["empty string"],
				`""`, nil)
			t.Run("with extra whitespace", func(t *testing.T) {
				checkDecoding(t, fixtures.SequenceMap["empty string"],
					`  "" `, nil)
			})
		})
		t.Run("flat string", func(t *testing.T) {
			checkDecoding(t, fixtures.SequenceMap["flat string"],
				`"value"`, nil)
		})
		t.Run("strings needing escape", func(t *testing.T) {
			checkDecoding(t, fixtures.SequenceMap["strings needing escape"],
				`"str\nbroken\ttabbed"`, nil)
		})
	})
}
