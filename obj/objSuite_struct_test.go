package obj

import (
	"reflect"
	"testing"

	"github.com/polydawn/refmt/obj/atlas"
	"github.com/polydawn/refmt/tok/fixtures"
)

func TestStructHandling(t *testing.T) {
	t.Run("tokens for map with one string field", func(t *testing.T) {
		seq := fixtures.SequenceMap["single row map"].Tokens
		t.Run("prism to object with explicit atlas", func(t *testing.T) {
			type tObjStr struct {
				X string
			}
			atlas := atlas.MustBuild(
				atlas.BuildEntry(tObjStr{}).StructMap().
					AddField("X", atlas.StructMapEntry{SerialName: "key"}).
					Complete(),
			)
			t.Run("marshal", func(t *testing.T) {
				value := tObjStr{"value"}
				checkMarshalling(t, atlas, value, seq, nil)
				checkMarshalling(t, atlas, &value, seq, nil)
			})
			t.Run("unmarshal", func(t *testing.T) {
				slot := &tObjStr{}
				expect := &tObjStr{"value"}
				checkUnmarshalling(t, atlas, slot, seq, expect, nil)
			})
			t.Run("unmarshal overwriting", func(t *testing.T) {
				slot := &tObjStr{"should be overruled"}
				expect := &tObjStr{"value"}
				checkUnmarshalling(t, atlas, slot, seq, expect, nil)
			})
		})
		t.Run("prism to object with autogen atlasentry", func(t *testing.T) {
			type tObjStr struct {
				Key string // these key downcased by default in autogen
			}
			atlas := atlas.MustBuild(
				atlas.BuildEntry(tObjStr{}).StructMap().Autogenerate().Complete(),
			)
			t.Run("marshal", func(t *testing.T) {
				value := tObjStr{"value"}
				checkMarshalling(t, atlas, value, seq, nil)
				checkMarshalling(t, atlas, &value, seq, nil)
			})
			t.Run("unmarshal", func(t *testing.T) {
				slot := &tObjStr{}
				expect := &tObjStr{"value"}
				checkUnmarshalling(t, atlas, slot, seq, expect, nil)
			})
		})
		t.Run("prism to object with additional fields", func(t *testing.T) {
			type tObjStr struct {
				Key   string
				Spare string
			}
			atlas := atlas.MustBuild(
				atlas.BuildEntry(tObjStr{}).StructMap().Autogenerate().Complete(),
			)
			t.Run("unmarshal", func(t *testing.T) {
				slot := &tObjStr{}
				expect := &tObjStr{"value", ""}
				checkUnmarshalling(t, atlas, slot, seq, expect, nil)
			})
			t.Run("unmarshal overwriting", func(t *testing.T) {
				slot := &tObjStr{"should be overruled", "untouched"}
				expect := &tObjStr{"value", "untouched"}
				checkUnmarshalling(t, atlas, slot, seq, expect, nil)
			})
		})
		t.Run("prism to object with no matching fields", func(t *testing.T) {
			type tObjStr struct {
				Spare string
			}
			atlas := atlas.MustBuild(
				atlas.BuildEntry(tObjStr{}).StructMap().Autogenerate().Complete(),
			)
			t.Run("unmarshal rejected", func(t *testing.T) {
				slot := &tObjStr{}
				expect := &tObjStr{}
				seq := seq[:2]
				checkUnmarshalling(t, atlas, slot, seq, expect, ErrNoSuchField{"key", reflect.TypeOf(tObjStr{}).String()})
			})
		})
	})
}
