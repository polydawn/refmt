package obj

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/polydawn/refmt/obj2/atlas"
	. "github.com/polydawn/refmt/testutil"
	. "github.com/polydawn/refmt/tok"
	"github.com/polydawn/refmt/tok/fixtures"
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
		valueFn:        func() interface{} { str := "value"; return str },
		marshalResults: &marshalResults{},
		unmarshalResults: []unmarshalResults{
			{title: "into string",
				slotFn:    func() interface{} { var str string; return str },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf("")}},
			{title: "into *string",
				slotFn: func() interface{} { var str string; return &str }},
			{title: "into wildcard",
				slotFn:    func() interface{} { var v interface{}; return v },
				expectErr: fmt.Errorf("unsettable")},
			{title: "into *wildcard",
				slotFn: func() interface{} { var v interface{}; return &v }},
			{title: "into map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return v },
				expectErr: fmt.Errorf("incompatable")},
			{title: "into *map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return &v },
				expectErr: fmt.Errorf("incompatable")},
			{title: "into []iface",
				slotFn:    func() interface{} { var v []interface{}; return v },
				expectErr: fmt.Errorf("incompatable")},
			{title: "into *[]iface",
				slotFn:    func() interface{} { var v []interface{}; return &v },
				expectErr: fmt.Errorf("incompatable")},
		},
	},
}

func TestMarshaller(t *testing.T) {
	// Package all the values from one step into a struct, just so that
	// we can assert on them all at once and make one green checkmark render per step.
	type marshallerStep struct {
		tok Token
		err error
	}
	Convey("Marshaller suite:", t, func() {
		for _, tr := range objFixtures {
			Convey(fmt.Sprintf("%q fixture sequence:", tr.title), func() {
				// Set up marshaller.
				marshaller := NewMarshaler(tr.atlas)
				marshaller.Bind(tr.valueFn())

				Convey("Steps...", func() {
					// Run steps until the marshaller says done or error.
					// For each step, assert the token matches fixtures;
					// when error and expected one, skip token check on that step
					// and finalize with the assertion.
					// If marshaller doesn't stop when we expected it to based
					// on fixture length, let it keep running three more steps
					// so we get that much more debug info.
					var done bool
					var err error
					var tok Token
					expectSteps := len(tr.sequence.Tokens) - 1
					for nStep := 0; nStep < expectSteps+3; nStep++ {
						done, err = marshaller.Step(&tok)
						if err != nil && tr.marshalResults.expectErr != nil {
							Convey("Result (error expected)", func() {
								So(err, ShouldResemble, tr.marshalResults.expectErr)
							})
							return
						}
						So(
							marshallerStep{tok, err},
							ShouldResemble,
							marshallerStep{tr.sequence.Tokens[nStep], nil},
						)
						if done {
							Convey("Result (halted correctly)", func() {
								So(nStep, ShouldEqual, expectSteps)
							})
							return
						}
					}
				})
			})
		}
	})
}

func TestUnmarshaller(t *testing.T) {
	for _, tr := range objFixtures {
		for _, trr := range tr.unmarshalResults {
			t.Skip("wip")

			// Grab slot.
			slot := trr.slotFn()

			// Set up unmarshaller.
			title := strings.Join([]string{tr.title, trr.title}, ", ")
			unmarshaller := NewUnmarshaler(tr.atlas)
			unmarshaller.Bind(slot)

			// Run steps.
			var done bool
			var err error
			for n, tok := range tr.sequence.Tokens {
				done, err = unmarshaller.Step(&tok)
				if err != nil {
					t.Errorf("step %d (inputting %s) errored: %s", n, tok, err)
				}
				if done && n != len(tr.sequence.Tokens)-1 {
					t.Errorf("done early! on step %d out of %d tokens", n, len(tr.sequence.Tokens))
				}
			}
			if !done {
				t.Errorf("still not done after %d tokens!", len(tr.sequence.Tokens))
			}

			// Check resulting object.
			Assert(t, title, tr.valueFn(), slot)
			t.Logf("test %q complete", tr.title)
		}
	}
}
