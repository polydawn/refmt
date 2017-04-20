package obj

import (
	"reflect"
	"testing"

	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/polydawn/refmt/testutil"
	. "github.com/polydawn/refmt/tok"
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
	type FF struct {
		I interface{}
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
				Add(NN{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(NN{}),
					Fields: []atlas.Entry{
						{Name: "F", FieldName: atlas.FieldName{"F"}},
						{Name: "X", FieldName: atlas.FieldName{"X"}},
					},
				}}),
			expectSeq: []Token{
				{Type: TMapOpen, Length: 2},
				/**/ TokStr("F"), TokInt(7),
				/**/ TokStr("X"), TokStr("s"),
				{Type: TMapClose},
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
				Add(NN{}, Morphism{Atlas: atlas.Atlas{
					Fields: []atlas.Entry{ /* this should be extraneous */ },
				}}).
				Add(AA{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(AA{}),
					Fields: []atlas.Entry{
						{Name: "a.y", FieldName: atlas.FieldName{"Y"}},
						{Name: "a.x", FieldName: atlas.FieldName{"X"}},
					},
				}}).
				Add(BB{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(BB{}),
					Fields: []atlas.Entry{
						{Name: "zee", FieldName: atlas.FieldName{"Z"}},
					},
				}}),
			expectSeq: []Token{
				{Type: TMapOpen, Length: 2},
				/**/ TokStr("a.y"), {Type: TMapOpen, Length: 1},
				/**/ /**/ TokStr("zee"), TokStr(""),
				/**/ {Type: TMapClose},
				/**/ TokStr("a.x"), TokStr("s"),
				{Type: TMapClose},
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
				Add(AA{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(AA{}),
					Fields: []atlas.Entry{
						{Name: "a.y", FieldName: atlas.FieldName{"Y"}},
						{Name: "a.x", FieldName: atlas.FieldName{"X"}},
					},
				}}),
			expectSeq: []Token{
				{Type: TMapOpen, Length: 2},
				TokStr("a.y"), {Type: TNull}, // last step panics
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
				Add(DD{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(DD{}),
					Fields: []atlas.Entry{
						{Name: "1", FieldName: atlas.FieldName{"A"}},
						{Name: "3", FieldName: atlas.FieldName{"Z"}},
						{Name: "2", FieldName: atlas.FieldName{"F"}},
					},
				}}).
				Add(AA{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(AA{}),
					Fields: []atlas.Entry{
						{Name: "a.y", FieldName: atlas.FieldName{"Y"}},
					},
				}}).
				Add(BB{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(BB{}),
					Fields: []atlas.Entry{
						{Name: "zee", FieldName: atlas.FieldName{"Z"}},
					},
				}}),
			expectSeq: []Token{
				{Type: TMapOpen, Length: 3},
				/**/ TokStr("1"), {Type: TMapOpen, Length: 1},
				/**/ /**/ TokStr("a.y"), {Type: TMapOpen, Length: 1},
				/**/ /**/ /**/ TokStr("zee"), TokStr("B"),
				/**/ /**/ {Type: TMapClose},
				/**/ {Type: TMapClose},
				/**/ TokStr("3"), {Type: TNull},
				/**/ TokStr("2"), TokInt(2),
				{Type: TMapClose},
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
				Add(DDD{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(DDD{}),
					Fields: []atlas.Entry{
						{Name: "1", FieldName: atlas.FieldName{"A"}},
						{Name: "3", FieldName: atlas.FieldName{"Z"}},
						{Name: "2", FieldName: atlas.FieldName{"F"}},
					},
				}}).
				Add(AA{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(AA{}),
					Fields: []atlas.Entry{
						{Name: "a.y", FieldName: atlas.FieldName{"Y"}},
					},
				}}).
				Add(BB{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(BB{}),
					Fields: []atlas.Entry{
						{Name: "zee", FieldName: atlas.FieldName{"Z"}},
					},
				}}),
			// should serialize exact same way as previous test:
			expectSeq: []Token{
				{Type: TMapOpen, Length: 3},
				/**/ TokStr("1"), {Type: TMapOpen, Length: 1},
				/**/ /**/ TokStr("a.y"), {Type: TMapOpen, Length: 1},
				/**/ /**/ /**/ TokStr("zee"), TokStr("B"),
				/**/ /**/ {Type: TMapClose},
				/**/ {Type: TMapClose},
				/**/ TokStr("3"), {Type: TNull},
				/**/ TokStr("2"), TokInt(2),
				{Type: TMapClose},
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
				Add(RR{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(RR{}),
					Fields: []atlas.Entry{
						{Name: "r", FieldName: atlas.FieldName{"R"}},
					},
				}}),
			expectSeq: []Token{
				{Type: TMapOpen, Length: 1},
				/**/ TokStr("r"), {Type: TMapOpen, Length: 1},
				/**/ /**/ TokStr("r"), {Type: TMapOpen, Length: 1},
				/**/ /**/ /**/ TokStr("r"), {Type: TMapOpen, Length: 1},
				/**/ /**/ /**/ /**/ TokStr("r"), {Type: TNull},
				/**/ /**/ /**/ {Type: TMapClose},
				/**/ /**/ {Type: TMapClose},
				/**/ {Type: TMapClose},
				{Type: TMapClose},
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
				{Type: TMapOpen, Length: 1},
				TokStr("a"), TokInt(1),
				{Type: TMapClose},
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
				{Type: TMapOpen, Length: 1},
				/**/ TokStr("a"), {Type: TMapOpen, Length: 1},
				/**/ /**/ TokStr("b"), TokInt(2),
				/**/ {Type: TMapClose},
				{Type: TMapClose},
			},
		},
		{
			title: "slice of literals",
			targetFn: func() interface{} {
				return &[]int{1, 2, 3, 4, 5}
			},
			expectSeq: []Token{
				{Type: TArrOpen, Length: 5},
				TokInt(1), TokInt(2), TokInt(3), TokInt(4), TokInt(5),
				{Type: TArrClose},
			},
		},
		{
			title: "struct containing unset wildcard field",
			targetFn: func() interface{} {
				return &FF{}
			},
			suite: (&Suite{}).
				Add(FF{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(FF{}),
					Fields: []atlas.Entry{
						{Name: "i", FieldName: atlas.FieldName{"I"}},
					},
				}}),
			expectSeq: []Token{
				{Type: TMapOpen, Length: 1},
				TokStr("i"), {Type: TNull},
				{Type: TMapClose},
			},
		},
		{
			title: "struct containing wildcard field set to struct",
			targetFn: func() interface{} {
				return &FF{FF{}}
			},
			suite: (&Suite{}).
				Add(FF{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(FF{}),
					Fields: []atlas.Entry{
						{Name: "i", FieldName: atlas.FieldName{"I"}},
					},
				}}),
			expectSeq: []Token{
				{Type: TMapOpen, Length: 1},
				/**/ TokStr("i"), {Type: TMapOpen, Length: 1},
				/**/ /**/ TokStr("i"), {Type: TNull},
				/**/ {Type: TMapClose},
				{Type: TMapClose},
			},
		},
		{
			title: "struct containing wildcard field set to ptr to struct",
			targetFn: func() interface{} {
				return &FF{&FF{}}
			},
			suite: (&Suite{}).
				Add(FF{}, Morphism{Atlas: atlas.Atlas{
					Type: reflect.TypeOf(FF{}),
					Fields: []atlas.Entry{
						{Name: "i", FieldName: atlas.FieldName{"I"}},
					},
				}}),
			expectSeq: []Token{
				{Type: TMapOpen, Length: 1},
				/**/ TokStr("i"), {Type: TMapOpen, Length: 1},
				/**/ /**/ TokStr("i"), {Type: TNull},
				/**/ {Type: TMapClose},
				{Type: TMapClose},
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
						tr.title, n, expectTok, tok)
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
