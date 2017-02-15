package obj

import (
	"reflect"
	"testing"

	"github.com/polydawn/go-xlate/obj/atlas"
	. "github.com/polydawn/go-xlate/testutil"
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
	type DD struct {
		A *AA
		F *int
		Z *string
	}
	type DDD struct {
		A **AA
		F *int
		Z *string
	}
	type RR struct {
		R *RR
	}

	tt := []struct {
		title       string
		targetFn    func() interface{} // func returns target, so test source looks like your call param
		suite       *Suite
		expectSeq   []Token
		expectErr   error
		expectPanic error
		errString   string
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
				Add(NN{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
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
				Add(NN{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
					Fields: []atlas.Entry{ /* this should be extraneous */ },
				})).
				Add(AA{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
					Type: reflect.TypeOf(AA{}),
					Fields: []atlas.Entry{
						{Name: "a.y", FieldName: atlas.FieldName{"Y"}},
						{Name: "a.x", FieldName: atlas.FieldName{"X"}},
					},
				})).
				Add(BB{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
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
		{
			title: "struct with fields missing a handler",
			targetFn: func() interface{} {
				return &AA{
					"s",
					BB{},
				}
			},
			suite: (&Suite{}).
				Add(AA{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
					Type: reflect.TypeOf(AA{}),
					Fields: []atlas.Entry{
						{Name: "a.y", FieldName: atlas.FieldName{"Y"}},
						{Name: "a.x", FieldName: atlas.FieldName{"X"}},
					},
				})),
			expectSeq: []Token{
				Token_MapOpen,
				TokStr("a.y"), nil, // last step panics
			},
			expectPanic: ErrNoHandler{},
			errString:   "no machine available in suite for struct of type obj.BB",
		},
		{
			title: "nested structs and ptrs",
			targetFn: func() interface{} {
				f := 2
				return &DD{
					&AA{
						"X",
						BB{"B"},
					},
					&f,
					nil,
				}
			},
			suite: (&Suite{}).
				Add(DD{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
					Type: reflect.TypeOf(DD{}),
					Fields: []atlas.Entry{
						{Name: "1", FieldName: atlas.FieldName{"A"}},
						{Name: "3", FieldName: atlas.FieldName{"Z"}},
						{Name: "2", FieldName: atlas.FieldName{"F"}},
					},
				})).
				Add(AA{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
					Type: reflect.TypeOf(AA{}),
					Fields: []atlas.Entry{
						{Name: "a.y", FieldName: atlas.FieldName{"Y"}},
					},
				})).
				Add(BB{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
					Type: reflect.TypeOf(BB{}),
					Fields: []atlas.Entry{
						{Name: "zee", FieldName: atlas.FieldName{"Z"}},
					},
				})),
			expectSeq: []Token{
				Token_MapOpen,
				TokStr("1"), Token_MapOpen,
				/**/ TokStr("a.y"), Token_MapOpen,
				/**/ /**/ TokStr("zee"), TokStr("B"),
				/**/ Token_MapClose,
				Token_MapClose,
				TokStr("3"), nil,
				TokStr("2"), TokInt(2),
				Token_MapClose,
			},
		},
		{
			title: "nested structs and deep ptrs",
			targetFn: func() interface{} {
				f := 2
				ap := &AA{
					"X",
					BB{"B"},
				}
				return &DDD{
					&ap,
					&f,
					nil,
				}
			},
			suite: (&Suite{}).
				Add(DDD{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
					Type: reflect.TypeOf(DDD{}),
					Fields: []atlas.Entry{
						{Name: "1", FieldName: atlas.FieldName{"A"}},
						{Name: "3", FieldName: atlas.FieldName{"Z"}},
						{Name: "2", FieldName: atlas.FieldName{"F"}},
					},
				})).
				Add(AA{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
					Type: reflect.TypeOf(AA{}),
					Fields: []atlas.Entry{
						{Name: "a.y", FieldName: atlas.FieldName{"Y"}},
					},
				})).
				Add(BB{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
					Type: reflect.TypeOf(BB{}),
					Fields: []atlas.Entry{
						{Name: "zee", FieldName: atlas.FieldName{"Z"}},
					},
				})),
			// should serialize exact same way as previous test:
			expectSeq: []Token{
				Token_MapOpen,
				TokStr("1"), Token_MapOpen,
				/**/ TokStr("a.y"), Token_MapOpen,
				/**/ /**/ TokStr("zee"), TokStr("B"),
				/**/ Token_MapClose,
				Token_MapClose,
				TokStr("3"), nil,
				TokStr("2"), TokInt(2),
				Token_MapClose,
			},
		},
		{
			title: "recursive structures",
			targetFn: func() interface{} {
				return &RR{
					&RR{
						&RR{
							&RR{},
						},
					},
				}
			},
			suite: (&Suite{}).
				Add(RR{}, MarshalMachineStructAtlasFactory(atlas.Atlas{
					Type: reflect.TypeOf(RR{}),
					Fields: []atlas.Entry{
						{Name: "r", FieldName: atlas.FieldName{"R"}},
					},
				})),
			expectSeq: []Token{
				Token_MapOpen,
				/**/ TokStr("r"), Token_MapOpen,
				/**/ /**/ TokStr("r"), Token_MapOpen,
				/**/ /**/ /**/ TokStr("r"), Token_MapOpen,
				/**/ /**/ /**/ /**/ TokStr("r"), nil,
				/**/ /**/ /**/ Token_MapClose,
				/**/ /**/ Token_MapClose,
				/**/ Token_MapClose,
				Token_MapClose,
			},
		},
		{
			title: "map of literals",
			targetFn: func() interface{} {
				return &map[string]int{
					"a": 1,
				}
			},
			expectSeq: []Token{
				Token_MapOpen,
				TokStr("a"), TokInt(1),
				Token_MapClose,
			},
		},
		{
			title: "map of map of literals",
			targetFn: func() interface{} {
				return &map[string]map[string]int{
					"a": map[string]int{
						"b": 2,
					},
				}
			},
			expectSeq: []Token{
				Token_MapOpen,
				TokStr("a"), Token_MapOpen,
				TokStr("b"), TokInt(2),
				Token_MapClose,
				Token_MapClose,
			},
		},
		{
			title: "slice of literals",
			targetFn: func() interface{} {
				return &[]int{1, 2, 3, 4, 5}
			},
			expectSeq: []Token{
				Token_ArrOpen,
				TokInt(1), TokInt(2), TokInt(3), TokInt(4), TokInt(5),
				Token_ArrClose,
			},
		},
	}
	for _, tr := range tt {
		if tr.suite == nil {
			tr.suite = &Suite{}
		}
		err := CapturePanics(func() {
			marshaller := NewMarshaler(tr.suite)
			marshaller.Bind(tr.targetFn())

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
			t.Errorf("test %q failed: expected panic of %T",
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
