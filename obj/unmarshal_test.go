package obj

import (
	"reflect"
	"testing"

	. "github.com/polydawn/go-xlate/testutil"
	"github.com/polydawn/go-xlate/tok"
)

type slotter interface {
	Slot() interface{} // Returns `&slot` -- the slot's type info is intact, but shadowed.
}

type slotForInt struct{ slot int }
type slotForString struct{ slot string }
type slotForIface struct{ slot interface{} }
type slotForSliceOfIface struct{ slot []interface{} }
type slotForMapOfStringToIface struct{ slot map[string]interface{} }

func (s *slotForInt) Slot() interface{}                { return &s.slot }
func (s *slotForString) Slot() interface{}             { return &s.slot }
func (s *slotForIface) Slot() interface{}              { return &s.slot }
func (s *slotForSliceOfIface) Slot() interface{}       { return &s.slot }
func (s *slotForMapOfStringToIface) Slot() interface{} { return &s.slot }

func TestUnmarshaller(t *testing.T) {
	tt := []struct {
		title    string
		slotter  slotter
		tokenSeq []tok.Token
		expect   interface{}
	}{
		{
			title:   "simple literal",
			slotter: &slotForInt{},
			tokenSeq: []tok.Token{
				4,
			},
			expect: 4,
		},
		{
			title:   "empty map into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []tok.Token{
				tok.Token_MapOpen,
				tok.Token_MapClose,
			},
			expect: map[string]interface{}{},
		},
		{
			title:   "simple flat map into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []tok.Token{
				tok.Token_MapOpen,
				"key", 6,
				tok.Token_MapClose,
			},
			expect: map[string]interface{}{"key": 6},
		},
		{
			title:   "map with nested map into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []tok.Token{
				tok.Token_MapOpen,
				"k1",
				tok.Token_MapOpen,
				"k2", "vvv",
				tok.Token_MapClose,
				tok.Token_MapClose,
			},
			expect: map[string]interface{}{"k1": map[string]interface{}{"k2": "vvv"}},
		},
		{
			title:   "array into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []tok.Token{
				tok.Token_ArrOpen,
				"v1",
				"v2",
				3,
				tok.Token_ArrClose,
			},
			expect: []interface{}{"v1", "v2", 3},
		},
		{
			title:   "map with nested array into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []tok.Token{
				tok.Token_MapOpen,
				"k1",
				tok.Token_ArrOpen,
				"v1",
				"v2",
				3,
				tok.Token_ArrClose,
				tok.Token_MapClose,
			},
			expect: map[string]interface{}{"k1": []interface{}{"v1", "v2", 3}},
		},
		{
			title:   "arrays with nested arrays into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []tok.Token{
				tok.Token_ArrOpen,
				"v1",
				tok.Token_ArrOpen,
				"v1",
				tok.Token_ArrOpen,
				"v1",
				"v2",
				3,
				tok.Token_ArrClose,
				3,
				tok.Token_ArrClose,
				3,
				tok.Token_ArrClose,
			},
			expect: []interface{}{"v1", []interface{}{"v1", []interface{}{"v1", "v2", 3}, 3}, 3},
		},
		{
			title:   "arrays with nested map into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []tok.Token{
				tok.Token_ArrOpen,
				"v1",
				tok.Token_MapOpen,
				"k2", "vvv",
				tok.Token_MapClose,
				3,
				tok.Token_ArrClose,
			},
			expect: []interface{}{"v1", map[string]interface{}{"k2": "vvv"}, 3},
		},
		{
			title:   "complex deeply nested structure into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []tok.Token{
				tok.Token_MapOpen,
				"k1",
				tok.Token_ArrOpen,
				"v1",
				tok.Token_MapOpen,
				tok.Token_MapClose,
				3,
				tok.Token_ArrOpen,
				14,
				15,
				tok.Token_MapOpen,
				"k9", "v10",
				tok.Token_MapClose,
				tok.Token_ArrOpen,
				tok.Token_ArrClose,
				tok.Token_ArrOpen,
				16,
				tok.Token_ArrClose,
				tok.Token_ArrClose,
				tok.Token_ArrClose,
				tok.Token_MapClose,
			},
			expect: map[string]interface{}{"k1": []interface{}{
				"v1",
				map[string]interface{}{},
				3,
				[]interface{}{
					14,
					15,
					map[string]interface{}{"k9": "v10"},
					[]interface{}(nil), // REVIEW: this behavior is questionable.  the type is right; a nil here may be... rude.
					[]interface{}{16},
				},
			}},
		},
	}
	for _, tr := range tt {
		// Create var receiver, aimed at the slotter.
		sink := NewUnmarshaler(tr.slotter.Slot())

		// Run steps.
		var done bool
		var err error
		for n, tok := range tr.tokenSeq {
			done, err = sink.Step(&tok)
			if err != nil {
				t.Errorf("step %d (inputting %#v) errored: %s", n, tok, err)
			}
			if done && n != len(tr.tokenSeq)-1 {
				t.Errorf("done early! on step %d out of %d tokens", n, len(tr.tokenSeq))
			}
		}
		if !done {
			t.Errorf("still not done after %d tokens!", len(tr.tokenSeq))
		}

		// Get value back out.  Some reflection required to get around pointers.
		v := reflect.ValueOf(tr.slotter.Slot()).Elem().Interface()
		Assert(t, tr.title, tr.expect, v)
	}
}
