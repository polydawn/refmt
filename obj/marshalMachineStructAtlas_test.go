package obj

import (
	"testing"

	"github.com/polydawn/go-xlate/obj/atlas"
	. "github.com/polydawn/go-xlate/tok"
)

func TestMarshalMachineStructAtlas(t *testing.T) {
	tt := []struct {
		title     string
		targetFn  func() interface{} // func returns target, so test source looks like your call param
		atlas     atlas.Atlas
		expectSeq []Token
	}{{
		title: "struct of several primitives",
		targetFn: func() interface{} {
			return &struct {
				X int
				Y int
				Z string
			}{
				1, 4, "nine",
			}
		},
		atlas: atlas.Atlas{
			Fields: []atlas.Entry{
				{Name: "x", FieldRoute: []int{0}},
				{Name: "y", FieldRoute: []int{1}},
				{Name: "z", FieldRoute: []int{2}},
			},
		},
		expectSeq: []Token{
			Token_MapOpen,
			TokStr("x"), TokInt(1),
			TokStr("y"), TokInt(4),
			TokStr("z"), TokStr("nine"),
			Token_MapClose,
		},
	}}
	for _, tr := range tt {
		// Setup
		tgt := tr.targetFn()
		machine := NewMarshalMachineStructAtlas(
			tr.atlas,
		)
		// Placeholders required for recursing on.
		suite := &Suite{}
		suite.Add(tgt, machine)
		driver := NewMarshaler(suite, tgt)

		// Run steps.
		var done bool
		var err error
		var tok Token
		for n, expectTok := range tr.expectSeq {
			done, err = machine.Step(driver, suite, &tok)
			if !IsTokenEqual(expectTok, tok) {
				t.Errorf("step %d yielded wrong token: expected %s, got %s", n, TokenToString(expectTok), TokenToString(tok))
			}
			if err != nil {
				t.Errorf("step %d (yielded %#v) errored: %s", n, tok, err)
			}
			if done && n != len(tr.expectSeq)-1 {
				t.Errorf("done early! on step %d out of %d tokens", n, len(tr.expectSeq))
			}
		}
		if !done {
			t.Errorf("still not done after %d tokens!", len(tr.expectSeq))
		}
		t.Logf("test %q halted correctly and passed", tr.title)
	}
}
