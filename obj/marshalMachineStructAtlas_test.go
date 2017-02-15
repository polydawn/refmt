package obj

import (
	"reflect"
	"testing"

	"github.com/polydawn/go-xlate/obj/atlas"
	. "github.com/polydawn/go-xlate/testutil"
	. "github.com/polydawn/go-xlate/tok"
)

func TestMarshalMachineStructAtlas(t *testing.T) {
	tt := []struct {
		title       string
		targetFn    func() interface{} // func returns target, so test source looks like your call param
		atlas       atlas.Atlas
		expectSeq   []Token
		expectErr   error
		expectPanic error
		errString   string
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
	}, {
		title: "struct containing nils",
		// This is not exactly a well-confined test: it covers the ptr machine,
		// not just this struct atlas machine.
		targetFn: func() interface{} {
			return &struct {
				Z *string
			}{}
		},
		atlas: atlas.Atlas{
			Fields: []atlas.Entry{
				{Name: "x", FieldRoute: []int{0}},
			},
		},
		expectSeq: []Token{
			Token_MapOpen,
			TokStr("x"), nil,
			Token_MapClose,
		},
	}, {
		title: "struct containing ptr to primitve",
		// This is not exactly a well-confined test: it covers the ptr machine,
		// not just this struct atlas machine.
		targetFn: func() interface{} {
			s := "asdf"
			return &struct {
				Z *string
			}{
				&s,
			}
		},
		atlas: atlas.Atlas{
			Fields: []atlas.Entry{
				{Name: "x", FieldRoute: []int{0}},
			},
		},
		expectSeq: []Token{
			Token_MapOpen,
			TokStr("x"), TokStr("asdf"),
			Token_MapClose,
		},
	}}
	for _, tr := range tt {
		// Setup
		tgt := tr.targetFn()
		// Placeholders required for recursing on.
		suite := &Suite{}
		suite.Add(tgt, Morphism{Atlas: tr.atlas})

		err := CapturePanics(func() {
			marshaller := NewMarshaler(suite)
			marshaller.Bind(tgt)

			// Run steps.
			var done bool
			var err error
			var tok Token
			for n, expectTok := range tr.expectSeq {
				done, err = marshaller.Step(&tok)
				if !IsTokenEqual(expectTok, tok) {
					t.Errorf("test %q failed: step %d yielded wrong token: expected %s, got %s",
						tr.title, n, TokenToString(expectTok), TokenToString(tok))
				}
				if err != nil {
					t.Errorf("test %q failed: step %d (expecting %#v) errored: %s",
						tr.title, n, expectTok, err)
				}
				if done && n != len(tr.expectSeq)-1 {
					t.Errorf("test %q failed: done early! on step %d out of %d tokens",
						tr.title, n, len(tr.expectSeq))
				}
			}
			if !done {
				t.Errorf("test %q failed: still not done after %d tokens!",
					tr.title, len(tr.expectSeq))
			}
		})
		if tr.expectPanic == nil && err == nil {
			t.Logf("test %q halted correctly and passed", tr.title)
		} else if err == nil {
			t.Errorf("test %q failed: expected panic of %T, but got nil",
				tr.title, tr.expectPanic)
		} else {
			ok := true
			if reflect.TypeOf(tr.expectPanic) != reflect.TypeOf(err) {
				t.Errorf("test %q failed: expected panic of type %T, but got %T",
					tr.title, tr.expectPanic, err)
				ok = false
			}
			if tr.errString != err.Error() {
				t.Errorf("test %q failed: expected panic of string of %q, but got %q",
					tr.title, tr.errString, err)
				ok = false
			}
			if ok {
				t.Logf("test %q panicked correctly and passed", tr.title)
			}
		}
	}
}
