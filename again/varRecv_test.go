package again

import (
	"fmt"
	"reflect"
	"testing"
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

func TestWow(t *testing.T) {
	tt := []struct {
		title    string
		slotter  slotter
		tokenSeq []Token
		expect   interface{}
	}{
		{
			title:   "simple literal",
			slotter: &slotForInt{},
			tokenSeq: []Token{
				4,
			},
			expect: 4,
		},
		{
			title:   "empty map into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []Token{
				Token_MapOpen,
				Token_MapClose,
			},
			expect: map[string]interface{}{},
		},
		{
			title:   "simple flat map into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []Token{
				Token_MapOpen,
				"key", 6,
				Token_MapClose,
			},
			expect: map[string]interface{}{"key": 6},
		},
		{
			title:   "map with nested map into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []Token{
				Token_MapOpen,
				"k1",
				Token_MapOpen,
				"k2", "vvv",
				Token_MapClose,
				Token_MapClose,
			},
			expect: map[string]interface{}{"k1": map[string]interface{}{"k2": "vvv"}},
		},
		{
			title:   "array into wildcard",
			slotter: &slotForIface{},
			tokenSeq: []Token{
				Token_ArrOpen,
				"v1",
				"v2",
				3,
				Token_ArrClose,
			},
			expect: []interface{}{"v1", "v2", 3},
		},
	}
	for _, tr := range tt {
		// Create var receiver, aimed at the slotter.
		sink := NewVarReceiver(tr.slotter.Slot())

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
		if !stringyEquality(tr.expect, v) {
			t.Errorf("test %q FAILED:\n\texpected: %#v\n\tactual:   %#v",
				tr.title, tr.expect, v)
		}
	}
}

func stringyEquality(x, y interface{}) bool {
	return fmt.Sprintf("%#v", x) == fmt.Sprintf("%#v", y)
}

func assert(t *testing.T, title string, expect, actual interface{}) {
	if !stringyEquality(expect, actual) {
		t.Errorf("test %q FAILED:\n\texpected  %#v\n\tactual    %#v",
			title, expect, actual)
	}
}
