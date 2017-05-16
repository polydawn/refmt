package objLegacy

import (
	"reflect"
	"testing"

	. "github.com/polydawn/refmt/testutil"
	. "github.com/polydawn/refmt/tok"
)

func TestUnmarshalMachineLiteral(t *testing.T) {
	tt := []struct {
		title    string
		slotter  slotter
		tokenSeq []Token
		expect   interface{}
	}{{
		title:   "simple literal",
		slotter: &slotForInt{},
		tokenSeq: []Token{
			TokInt(4),
		},
		expect: 4,
	}, {
		title:   "simple literal into wildcard",
		slotter: &slotForIface{},
		tokenSeq: []Token{
			TokInt(4),
		},
		expect: 4,
	}}
	for _, tr := range tt {
		mach := &UnmarshalMachineLiteral{}
		mach.target = tr.slotter.Slot()

		// Run steps.
		var done bool
		var err error
		for n, tok := range tr.tokenSeq {
			done, err = mach.Step(nil, &tok)
			if err != nil {
				t.Errorf("step %d (inputting %s) errored: %s", n, tok, err)
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
