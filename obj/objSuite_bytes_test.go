package obj

import (
	"testing"

	"github.com/polydawn/refmt/obj/atlas"
	"github.com/polydawn/refmt/tok/fixtures"
)

func TestBytes(t *testing.T) {
	t.Run("tokens for bytes", func(t *testing.T) {
		seq := fixtures.SequenceMap["short byte array"].Tokens
		t.Run("prism to []byte", func(t *testing.T) {
			atlas := atlas.MustBuild()
			t.Run("marshal", func(t *testing.T) {
				value := []byte(`value`)
				checkMarshalling(t, atlas, value, seq, nil)
				checkMarshalling(t, atlas, &value, seq, nil)
			})
			t.Run("unmarshal", func(t *testing.T) {
				slot := []byte{}
				expect := []byte(`value`)
				checkUnmarshalling(t, atlas, &slot, seq, &expect, nil)
			})
		})
		t.Run("prism to [5]byte", func(t *testing.T) {
			atlas := atlas.MustBuild()
			t.Run("marshal", func(t *testing.T) {
				value := [5]byte{'v', 'a', 'l', 'u', 'e'}
				checkMarshalling(t, atlas, value, seq, nil)
				checkMarshalling(t, atlas, &value, seq, nil)
			})
			t.Run("unmarshal", func(t *testing.T) {
				//	slot := [5]byte{}
				//	expect := [5]byte{'v', 'a', 'l', 'u', 'e'}
				//	checkUnmarshalling(t, atlas, &slot, seq, &expect, nil)
			})
		})
	})
}
