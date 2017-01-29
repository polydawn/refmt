package json

import (
	"bytes"
	"testing"

	. "github.com/polydawn/go-xlate/testutil"
	. "github.com/polydawn/go-xlate/tok"
)

func TestJsonSerializer(t *testing.T) {
	tt := []struct {
		title    string
		tokenSeq []Token
		expect   string
	}{
		{
			"flat string",
			[]Token{
				TokStr("value"),
			},
			`"value"`,
		},
		{
			"single row map",
			[]Token{
				Token_MapOpen,
				TokStr("key"),
				TokStr("value"),
				Token_MapClose,
			},
			`{"key":"value"}`,
		},
		{
			"duo row map",
			[]Token{
				Token_MapOpen,
				TokStr("key"),
				TokStr("value"),
				TokStr("k2"),
				TokStr("v2"),
				Token_MapClose,
			},
			`{"key":"value","k2":"v2"}`,
		},
		{
			"single entry array",
			[]Token{
				Token_ArrOpen,
				TokStr("value"),
				Token_ArrClose,
			},
			`["value"]`,
		},
		{
			"duo entry array",
			[]Token{
				Token_ArrOpen,
				TokStr("value"),
				TokStr("v2"),
				Token_ArrClose,
			},
			`["value","v2"]`,
		},
		{
			"empty map",
			[]Token{
				Token_MapOpen,
				Token_MapClose,
			},
			`{}`,
		},
		{
			"empty array",
			[]Token{
				Token_ArrOpen,
				Token_ArrClose,
			},
			`[]`,
		},
		{
			"array nested in map as non-first and final entry",
			[]Token{
				Token_MapOpen,
				TokStr("k1"),
				TokStr("v1"),
				TokStr("ke"),
				Token_ArrOpen,
				TokStr("oh"),
				TokStr("whee"),
				TokStr("wow"),
				Token_ArrClose,
				Token_MapClose,
			},
			`{"k1":"v1","ke":["oh","whee","wow"]}`,
		},
		{
			"array nested in map as first and non-final entry",
			[]Token{
				Token_MapOpen,
				TokStr("ke"),
				Token_ArrOpen,
				TokStr("oh"),
				TokStr("whee"),
				TokStr("wow"),
				Token_ArrClose,
				TokStr("k1"),
				TokStr("v1"),
				Token_MapClose,
			},
			`{"ke":["oh","whee","wow"],"k1":"v1"}`,
		},
		{
			"maps nested in array",
			[]Token{
				Token_ArrOpen,
				Token_MapOpen,
				TokStr("k"),
				TokStr("v"),
				Token_MapClose,
				TokStr("whee"),
				Token_MapOpen,
				TokStr("k1"),
				TokStr("v1"),
				Token_MapClose,
				Token_ArrClose,
			},
			`[{"k":"v"},"whee",{"k1":"v1"}]`,
		},
		{
			"arrays in arrays in arrays",
			[]Token{
				Token_ArrOpen,
				Token_ArrOpen,
				Token_ArrOpen,
				Token_ArrClose,
				Token_ArrClose,
				Token_ArrClose,
			},
			`[[[]]]`,
		},
		{
			"maps nested in maps",
			[]Token{
				Token_MapOpen,
				TokStr("k"),
				Token_MapOpen,
				TokStr("k2"),
				TokStr("v2"),
				Token_MapClose,
				Token_MapClose,
			},
			`{"k":{"k2":"v2"}}`,
		},
		{
			"strings needing escape",
			[]Token{
				TokStr("str\nbroken\ttabbed"),
			},
			`"str\nbroken\ttabbed"`,
		},
	}
	for _, tr := range tt {
		// Set it up.
		buf := &bytes.Buffer{}
		sink := NewSerializer(buf)

		// Run steps.
		var done bool
		var err error
		for n, tok := range tr.tokenSeq {
			done, err = sink.Step(&tok)
			if err != nil {
				t.Errorf("test %q step %d (inputting %#v) errored: %s", tr.title, n, tok, err)
			}
			if done && n != len(tr.tokenSeq)-1 {
				t.Errorf("test %q done early! on step %d out of %d tokens", tr.title, n, len(tr.tokenSeq))
			}
		}
		if !done {
			t.Errorf("test %q still not done after %d tokens!", tr.title, len(tr.tokenSeq))
		}

		// Assert final result.
		Assert(t, tr.title, tr.expect, buf.String())
	}
}
