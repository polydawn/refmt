package obj

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/polydawn/refmt/obj2/atlas"
	. "github.com/polydawn/refmt/tok"
	"github.com/polydawn/refmt/tok/fixtures"
)

var skipMe = fmt.Errorf("skipme")

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

type tObjStr struct {
	X string
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
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(interface{}(nil))}},
			{title: "into *wildcard",
				slotFn: func() interface{} { var v interface{}; return &v }},
			{title: "into map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return v },
				expectErr: skipMe},
			{title: "into *map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return &v },
				expectErr: skipMe},
			{title: "into []iface",
				slotFn:    func() interface{} { var v []interface{}; return v },
				expectErr: skipMe},
			{title: "into *[]iface",
				slotFn:    func() interface{} { var v []interface{}; return &v },
				expectErr: skipMe},
		},
	},
	{title: "object with one string field, with atlas entry",
		sequence: fixtures.SequenceMap["single row map"],
		valueFn:  func() interface{} { return tObjStr{"value"} },
		atlas: atlas.MustBuild(
			atlas.BuildEntry(tObjStr{}).StructMap().
				AddField("X", atlas.StructMapEntry{SerialName: "key"}).
				Complete(),
		),
		marshalResults: &marshalResults{},
		unmarshalResults: []unmarshalResults{
			{title: "into string",
				slotFn:    func() interface{} { var str string; return str },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf("")}},
			{title: "into *string",
				slotFn:    func() interface{} { var str string; return &str },
				expectErr: ErrUnmarshalIncongruent{Token{Type: TMapOpen, Length: 1}, reflect.ValueOf("")}},
			{title: "into wildcard",
				slotFn:    func() interface{} { var v interface{}; return v },
				expectErr: ErrInvalidUnmarshalTarget{reflect.TypeOf(interface{}(nil))}},
			{title: "into *wildcard",
				slotFn:    func() interface{} { var v interface{}; return &v },
				expectErr: skipMe},
			{title: "into map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return v },
				expectErr: skipMe},
			{title: "into *map[str]iface",
				slotFn:    func() interface{} { var v map[string]interface{}; return &v },
				expectErr: skipMe},
			{title: "into []iface",
				slotFn:    func() interface{} { var v []interface{}; return v },
				expectErr: skipMe},
			{title: "into *[]iface",
				slotFn:    func() interface{} { var v []interface{}; return &v },
				expectErr: skipMe},
		},
	},
}

func TestMarshaller(t *testing.T) {
	// Package all the values from one step into a struct, just so that
	// we can assert on them all at once and make one green checkmark render per step.
	// Stringify the token first so extraneous fields in the union are hidden.
	type step struct {
		tok string
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
								So(err.Error(), ShouldResemble, tr.marshalResults.expectErr.Error())
							})
							return
						}
						if nStep <= expectSteps {
							So(
								step{tok.String(), err},
								ShouldResemble,
								step{tr.sequence.Tokens[nStep].String(), nil},
							)
						} else {
							So(
								step{tok.String(), err},
								ShouldResemble,
								step{Token{}.String(), fmt.Errorf("overshoot")},
							)
						}
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
	// Package all the values from one step into a struct, just so that
	// we can assert on them all at once and make one green checkmark render per step.
	// Stringify the token first so extraneous fields in the union are hidden.
	type step struct {
		tok  string
		err  error
		done bool
	}

	Convey("Unmarshaller suite:", t, func() {
		for _, tr := range objFixtures {
			Convey(fmt.Sprintf("%q fixture sequence:", tr.title), func() {
				for _, trr := range tr.unmarshalResults {
					maybe := Convey
					if trr.expectErr == skipMe {
						maybe = SkipConvey
					}
					maybe(fmt.Sprintf("targetting %s:", trr.title), func() {
						// Grab slot.
						slot := trr.slotFn()

						// Set up unmarshaller.
						unmarshaller := NewUnmarshaler(tr.atlas)
						err := unmarshaller.Bind(slot)
						if err != nil && trr.expectErr != nil {
							Convey("Result (error expected)", func() {
								So(err.Error(), ShouldResemble, trr.expectErr.Error())
							})
							return
						}

						Convey("Steps...", func() {
							// Run steps.
							// This is less complicated than the marshaller test
							// because we know exactly when we'll run out of them.
							var done bool
							var err error
							expectSteps := len(tr.sequence.Tokens) - 1
							for nStep, tok := range tr.sequence.Tokens {
								done, err = unmarshaller.Step(&tok)
								if err != nil && trr.expectErr != nil {
									Convey("Result (error expected)", func() {
										So(err.Error(), ShouldResemble, trr.expectErr.Error())
									})
									return
								}
								if nStep == expectSteps {
									So(
										step{tok.String(), err, done},
										ShouldResemble,
										step{tr.sequence.Tokens[nStep].String(), nil, true},
									)
								} else {
									So(
										step{tok.String(), err, done},
										ShouldResemble,
										step{Token{}.String(), nil, false},
									)
								}
							}

							Convey("Result", func() {
								// Get value back out.  Some reflection required to get around pointers.
								v := reflect.ValueOf(slot).Elem().Interface()
								So(v, ShouldResemble, tr.valueFn())
							})
						})
					})
				}
			})
		}
	})
}
