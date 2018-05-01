package json

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	. "github.com/warpfork/go-wish"

	. "github.com/polydawn/refmt/tok"
	"github.com/polydawn/refmt/tok/fixtures"
)

// note: we still put all tests in one func so we control order.
// this will let us someday refactor all `fixtures.SequenceMap` refs to use a
// func which quietly records which sequences have tests aimed at them, and we
// can read that back at out the end of the tests and use the info to
// proactively warn ourselves when we have unreferenced tok fixtures.

func TestEncoding(t *testing.T) {
	testBoolEncoding(t)
	testStringEncoding(t)
}

func TestDecoding(t *testing.T) {
	testBoolDecoding(t)
	testStringDecoding(t)
}

func checkEncoding(t *testing.T, sequence fixtures.Sequence, expectSerial string, expectErr error) {
	t.Helper()
	outputBuf := &bytes.Buffer{}
	tokenSink := NewEncoder(outputBuf, EncodeOptions{})

	// Run steps, advancing through the token sequence.
	//  If it stops early, just report how many steps in; we Wish on that value.
	//  If it doesn't stop in time, just report that bool; we Wish on that value.
	var nStep int
	var done bool
	var err error
	for _, tok := range sequence.Tokens {
		nStep++
		done, err = tokenSink.Step(&tok)
		if done || err != nil {
			break
		}
	}

	// Assert final result.
	Wish(t, done, ShouldEqual, true)
	Wish(t, nStep, ShouldEqual, len(sequence.Tokens))
	Wish(t, err, ShouldEqual, expectErr)
	Wish(t, outputBuf.String(), ShouldEqual, expectSerial)
}

func checkDecoding(t *testing.T, expectSequence fixtures.Sequence, serial string, expectErr error) {
	t.Helper()
	inputBuf := bytes.NewBufferString(serial)
	tokenSrc := NewDecoder(inputBuf)

	// Run steps, advancing until the decoder reports it's done.
	//  If the decoder keeps yielding more tokens than we expect, that's fine...
	//  we just keep recording them, and we'll diff later.
	//  There's a cutoff when it overshoots by 10 tokens because generally
	//  that indicates we've found some sort of loop bug and 10 extra token
	//  yields is typically enough info to diagnose with.
	var nStep int
	var done bool
	var yield = make([]Token, len(expectSequence.Tokens)+10)
	var err error
	for ; nStep <= len(expectSequence.Tokens)+10; nStep++ {
		done, err = tokenSrc.Step(&yield[nStep])
		if done || err != nil {
			break
		}
	}
	nStep++
	yield = yield[:nStep]

	// Assert final result.
	Wish(t, done, ShouldEqual, true)
	Wish(t, nStep, ShouldEqual, len(expectSequence.Tokens))
	Wish(t, yield, ShouldEqual, expectSequence.Tokens)
	Wish(t, err, ShouldEqual, expectErr)
}

// --------------

var inapplicable = fmt.Errorf("skipme: inapplicable")

var jsonFixtures = []struct {
	title        string
	sequence     fixtures.Sequence
	serial       string
	encodeResult error
	decodeResult error
}{
	// Maps
	{"",
		fixtures.SequenceMap["empty map"].SansLengthInfo(),
		`{}`,
		nil,
		nil,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["empty map"].SansLengthInfo(),
		`{  }`,
		inapplicable,
		nil,
	},
	{"",
		fixtures.SequenceMap["single row map"].SansLengthInfo(),
		`{"key":"value"}`,
		nil,
		nil,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["single row map"].SansLengthInfo(),
		` { "key"  :  "value" } `,
		inapplicable,
		nil,
	},
	{"",
		fixtures.SequenceMap["duo row map"].SansLengthInfo(),
		`{"key":"value","k2":"v2"}`,
		nil,
		nil,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["duo row map"].SansLengthInfo(),
		`{"key":"value",  "k2":"v2"}`,
		inapplicable,
		nil,
	},
	{"decoding with trailing comma",
		fixtures.SequenceMap["duo row map"].SansLengthInfo(),
		`{"key":"value","k2":"v2",}`,
		inapplicable,
		nil,
	},
	{"",
		fixtures.SequenceMap["duo row map alt2"].SansLengthInfo(),
		`{"k2":"v2","key":"value"}`,
		inapplicable,
		nil,
	},

	// Arrays
	{"",
		fixtures.SequenceMap["empty array"].SansLengthInfo(),
		`[]`,
		nil,
		nil,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["empty array"].SansLengthInfo(),
		`  [ ] `,
		inapplicable, nil,
	},
	{"",
		fixtures.SequenceMap["single entry array"].SansLengthInfo(),
		`["value"]`,
		nil,
		nil,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["single entry array"].SansLengthInfo(),
		`  [ "value" ] `,
		inapplicable,
		nil,
	},
	{"",
		fixtures.SequenceMap["duo entry array"].SansLengthInfo(),
		`["value","v2"]`,
		nil,
		nil,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["duo entry array"].SansLengthInfo(),
		`["value",  "v2"]`,
		inapplicable,
		nil,
	},

	// Complex / mixed / nested.
	{"",
		fixtures.SequenceMap["array nested in map as non-first and final entry"].SansLengthInfo(),
		`{"k1":"v1","ke":["oh","whee","wow"]}`,
		nil,
		nil,
	},
	{"",
		fixtures.SequenceMap["array nested in map as first and non-final entry"].SansLengthInfo(),
		`{"ke":["oh","whee","wow"],"k1":"v1"}`,
		nil,
		nil,
	},
	{"",
		fixtures.SequenceMap["maps nested in array"].SansLengthInfo(),
		`[{"k":"v"},"whee",{"k1":"v1"}]`,
		nil,
		nil,
	},
	{"",
		fixtures.SequenceMap["arrays in arrays in arrays"].SansLengthInfo(),
		`[[[]]]`,
		nil,
		nil,
	},
	{"",
		fixtures.SequenceMap["maps nested in maps"].SansLengthInfo(),
		`{"k":{"k2":"v2"}}`,
		nil,
		nil,
	},

	// Errors when decoding invalid inputs!
	{"",
		fixtures.SequenceMap["dangling arr open"].SansLengthInfo().Append(Token{}),
		`[`,
		inapplicable,
		io.EOF, // REVIEW it's probably more explicitly unexpected than that...
	},

	// Numeric.
	{"",
		fixtures.Sequence{"integer zero", []Token{{Type: TInt, Int: 0}}},
		"0",
		nil,
		nil,
	},
	{"",
		fixtures.Sequence{"integer one", []Token{{Type: TInt, Int: 1}}},
		"1",
		nil,
		nil,
	},
	{"",
		fixtures.Sequence{"integer neg 1", []Token{{Type: TInt, Int: -1}}},
		"-1",
		nil,
		nil,
	},
	{"",
		fixtures.Sequence{"integer neg 100", []Token{{Type: TInt, Int: -100}}},
		"-100",
		nil,
		nil,
	},
	{"",
		fixtures.Sequence{"integer 1000000", []Token{{Type: TInt, Int: 1000000}}},
		"1000000",
		nil,
		nil,
	},
	{"",
		fixtures.Sequence{"float 1 e+100", []Token{{Type: TFloat64, Float64: 1.0e+300}}},
		`1e+300`,
		inapplicable, // TODO should support situationEncoding too
		nil,
	},
}
