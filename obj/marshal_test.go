package obj

import (
	"reflect"
	"testing"

	//. "github.com/polydawn/go-xlate/testutil"
	"github.com/polydawn/go-xlate/obj/atlas"
	. "github.com/polydawn/go-xlate/tok"
)

func TestMarshaller(t *testing.T) {
	type NN struct {
		F int
		X string
	}
	type BB struct {
		Z string
	}
	type AA struct {
		X string
		Y BB
	}

	tt := []struct {
		title     string
		targetFn  func() interface{} // func returns target, so test source looks like your call param
		suite     *Suite
		expectSeq []Token
		expectErr error
	}{
		{
			title:    "simple literal",
			targetFn: func() interface{} { i := 4; return &i },
			expectSeq: []Token{
				TokInt(4),
			},
		},
		{
			title: "simple struct of literals",
			targetFn: func() interface{} {
				return &NN{
					7, "s",
				}
			},
			suite: (&Suite{}).
				Add(NN{}, NewMarshalMachineStructAtlas(atlas.Atlas{
					Type: reflect.TypeOf(NN{}),
					Fields: []atlas.Entry{
						{Name: "F", FieldName: atlas.FieldName{"F"}},
						{Name: "X", FieldName: atlas.FieldName{"X"}},
					},
				})),
			expectSeq: []Token{
				Token_MapOpen,
				TokStr("F"), TokInt(7),
				TokStr("X"), TokStr("s"),
				Token_MapClose,
			},
		},
		{
			title: "nested structs and literals",
			targetFn: func() interface{} {
				return &AA{
					"s",
					BB{},
				}
			},
			suite: (&Suite{}).
				Add(NN{}, NewMarshalMachineStructAtlas(atlas.Atlas{
					Fields: []atlas.Entry{ /* this should be extraneous */ },
				})).
				Add(AA{}, NewMarshalMachineStructAtlas(atlas.Atlas{
					Type: reflect.TypeOf(AA{}),
					Fields: []atlas.Entry{
						{Name: "a.y", FieldName: atlas.FieldName{"Y"}},
						{Name: "a.x", FieldName: atlas.FieldName{"X"}},
					},
				})).
				Add(BB{}, NewMarshalMachineStructAtlas(atlas.Atlas{
					Type: reflect.TypeOf(BB{}),
					Fields: []atlas.Entry{
						{Name: "zee", FieldName: atlas.FieldName{"Z"}},
					},
				})),
			expectSeq: []Token{
				Token_MapOpen,
				TokStr("a.y"), Token_MapOpen,
				TokStr("zee"), TokStr(""),
				Token_MapClose,
				TokStr("a.x"), TokStr("s"),
				Token_MapClose,
			},
		},
		// TODO following doesn't work yet because of type-loss issues when converting away from reflect.Value
		//  (which are in turn blocked from easily resolution because of the tricky detail that map vals are not addressable..).
		//{
		//	title: "wildcard map of literals",
		//	targetFn: func() interface{} {
		//		return &map[string]int{
		//			"a": 1,
		//		}
		//	},
		//	expectSeq: []Token{
		//		Token_MapOpen,
		//		"a", 1,
		//		Token_MapClose,
		//	},
		//},
	}
	for _, tr := range tt {
		if tr.suite == nil {
			tr.suite = &Suite{}
		}
		marshaller := NewMarshaler(tr.suite, tr.targetFn())

		// Run steps.
		var done bool
		var err error
		var tok Token
		for n, expectTok := range tr.expectSeq {
			done, err = marshaller.Step(&tok)
			if !IsTokenEqual(expectTok, tok) {
				t.Errorf("step %d yielded wrong token: expected %s, got %s", n, TokenToString(expectTok), TokenToString(tok))
			}
			if err != nil {
				t.Errorf("step %d (expecting %#v) errored: %s", n, expectTok, err)
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
