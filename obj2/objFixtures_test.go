package obj

import (
	"fmt"
	"testing"

	"github.com/polydawn/go-xlate/obj2/atlas"
	. "github.com/polydawn/go-xlate/tok"
	"github.com/polydawn/go-xlate/tok/fixtures"
)

type marshalResults struct {
	expectErr error
	errString string
}
type unmarshalResults struct {
	title string
	// Yields the handle we should give to the unmarshaller to fill.
	// Like `valueFn`, the indirection here is to help
	slotFn    func() interface{}
	expectErr error
	errString string
}

var objFixtures = []struct {
	title string

	// The serial sequence of tokens the value is isomorphic to.
	sequence fixtures.Sequence

	// Yields a value to hand to the marshaller; or, the value we will compare the result against for unmarshal.
	// A func returning a wildcard is used rather than just an `interface{}`, because `&target` conveys very different type information.
	valueFn func() interface{}

	// The suite of mappings to use.
	atlas atlas.Atlas

	// The results expected from marshalling.  If nil: don't run marshal test.
	// (If zero value, the result should be passing and nil errors.)
	marshalResults *marshalResults

	// The results to expect from various unmarshal situations.
	// This is a slice because unmarshal may have different outcomes (usually,
	// erroring vs not) depending on the type of value it was given to populate.
	unmarshalResults []unmarshalResults
}{
	{title: "string literal",
		sequence:       fixtures.SequenceMap["flat string"],
		valueFn:        func() interface{} { str := "value"; return &str },
		marshalResults: &marshalResults{},
		unmarshalResults: []unmarshalResults{
			{title: "into string",
				slotFn:    func() interface{} { var str string; return str },
				expectErr: fmt.Errorf("unsettable")},
			{title: "into string handle",
				slotFn: func() interface{} { var str string; return &str }},
			{title: "into wildcard",
				slotFn:    func() interface{} { var v interface{}; return v },
				expectErr: fmt.Errorf("unsettable")},
			{title: "into wildcard handle",
				slotFn: func() interface{} { var v interface{}; return &v }},
			{title: "into map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return v },
				expectErr: fmt.Errorf("incompatable")},
			{title: "into map[str]iface handle",
				slotFn:    func() interface{} { var v map[string]interface{}; return &v },
				expectErr: fmt.Errorf("incompatable")},
			{title: "into []iface",
				slotFn:    func() interface{} { var v []interface{}; return v },
				expectErr: fmt.Errorf("incompatable")},
			{title: "into []iface handle",
				slotFn:    func() interface{} { var v []interface{}; return &v },
				expectErr: fmt.Errorf("incompatable")},
		},
	},
}

func TestMarshaller(t *testing.T) {
	for _, tr := range objFixtures {
		// Set up marshaller.
		marshaller := NewMarshaler(tr.atlas)
		marshaller.Bind(tr.valueFn())

		// Run steps.
		var done bool
		var err error
		var tok Token
		for n, expectTok := range tr.sequence.Tokens {
			done, err = marshaller.Step(&tok)
			if !IsTokenEqual(expectTok, tok) {
				t.Errorf("test %q failed: step %d yielded wrong token: expected %s, got %s",
					tr.title, n, expectTok, tok)
			}
			if err != nil {
				t.Errorf("test %q failed: step %d (expecting %#v) errored: %s",
					tr.title, n, expectTok, err)
			}
			if done && n != len(tr.sequence.Tokens)-1 {
				t.Errorf("test %q failed: done early! on step %d out of %d tokens",
					tr.title, n, len(tr.sequence.Tokens))
			}
		}
		if !done {
			t.Errorf("test %q failed: still not done after %d tokens!",
				tr.title, len(tr.sequence.Tokens))
		}
		t.Logf("test %q complete", tr.title)
	}
}
