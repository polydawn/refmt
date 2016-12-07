package obj

import (
	"testing"

	"github.com/polydawn/go-xlate/obj/atlas"
	. "github.com/polydawn/go-xlate/testutil"
	. "github.com/polydawn/go-xlate/tok"
)

func TestMarshalMachineStructAtlas(t *testing.T) {
	tt := []struct {
		title  string
		value  interface{}
		atlas  atlas.Atlas
		expect []Token
	}{{
		title: "struct of several primitives",
		value: struct {
			X int
			y int
			z string
		}{
			1, 4, "nine",
		},
		atlas: atlas.Atlas{
			Fields: []atlas.Entry{
				{Name: "x", FieldRoute: []int{0}},
				{Name: "y", FieldRoute: []int{1}},
				{Name: "z", FieldRoute: []int{2}},
			},
		},
		expect: []Token{
			Token_MapOpen,
			"x", 1,
			"y", 4,
			"z", "nine",
			Token_MapClose,
		},
	}}
	t.Skip("not done")
	for _, tr := range tt {
		// Setup
		var tokens []Token
		machine := UnmarshalMachineStructAtlas{
			target: tr.value,
			atlas:  tr.atlas,
		}

		// Run steps.
		var done bool
		var err error
		var tok Token
		for n := 0; n < len(tr.expect); n++ {
			done, err = machine.Step(nil, nil, &tok) // FIXME really need suite objs now
			tokens = append(tokens, tok)
			if err != nil {
				t.Errorf("step %d (yielded %#v) errored: %s", n, tok, err)
			}
			if done && n != len(tr.expect)-1 {
				t.Errorf("done early! on step %d out of %d tokens", n, len(tr.expect))
			}
		}
		if !done {
			t.Errorf("still not done after %d tokens!", len(tr.expect))
		}

		// Check final token stream matches expectations.
		Assert(t, tr.title, tr.expect, tokens)
	}
}
