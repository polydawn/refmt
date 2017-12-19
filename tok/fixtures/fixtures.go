/*
	Token stream test fixtures.
	This is a public package because it is used by tests in the `json`, `cbor`, and `obj` packages.
	It should not be seen in the imports outside of testing.
*/
package fixtures

import (
	. "github.com/polydawn/refmt/tok"
)

type Sequence struct {
	Title  string
	Tokens []Token
}

// An array of well-formed token sequences.
var Sequences = []Sequence{
	// Booleans.
	{"true",
		[]Token{
			{Type: TBool, Bool: true},
		},
	},
	{"false",
		[]Token{
			{Type: TBool, Bool: false},
		},
	},

	// Strings.
	{"empty string",
		[]Token{
			TokStr(""),
		},
	},
	{"flat string",
		[]Token{
			TokStr("value"),
		},
	},
	{"strings needing escape",
		[]Token{
			TokStr("str\nbroken\ttabbed"),
		},
	},

	// Maps.
	{"empty map",
		[]Token{
			{Type: TMapOpen, Length: 0},
			{Type: TMapClose},
		},
	},
	{"single row map",
		[]Token{
			{Type: TMapOpen, Length: 1},
			TokStr("key"),
			TokStr("value"),
			{Type: TMapClose},
		},
	},
	{"duo row map",
		[]Token{
			{Type: TMapOpen, Length: 2},
			TokStr("key"),
			TokStr("value"),
			TokStr("k2"),
			TokStr("v2"),
			{Type: TMapClose},
		},
	},
	{"duo row map alt2",
		// same as previous, but map entries in a different order -- useful to test that unmarshaller can accept that (or, reject non-canonical orders!)
		[]Token{
			{Type: TMapOpen, Length: 2},
			TokStr("k2"),
			TokStr("v2"),
			TokStr("key"),
			TokStr("value"),
			{Type: TMapClose},
		},
	},
	{"quad map default order",
		[]Token{
			{Type: TMapOpen, Length: 4},
			TokStr("1"), TokStr("1"),
			TokStr("b"), TokStr("2"),
			TokStr("bc"), TokStr("3"),
			TokStr("d"), TokStr("4"),
			{Type: TMapClose},
		},
	},
	{"quad map rfc7049 order",
		[]Token{
			{Type: TMapOpen, Length: 4},
			TokStr("1"), TokStr("1"),
			TokStr("b"), TokStr("2"),
			TokStr("d"), TokStr("3"),
			TokStr("bc"), TokStr("4"),
			{Type: TMapClose},
		},
	},

	// Arrays.
	{"empty array",
		[]Token{
			{Type: TArrOpen, Length: 0},
			{Type: TArrClose},
		},
	},
	{"single entry array",
		[]Token{
			{Type: TArrOpen, Length: 1},
			TokStr("value"),
			{Type: TArrClose},
		},
	},
	{"duo entry array",
		[]Token{
			{Type: TArrOpen, Length: 2},
			TokStr("value"),
			TokStr("v2"),
			{Type: TArrClose},
		},
	},

	// Complex / mixed / nested.
	{"array nested in map as non-first and final entry",
		[]Token{
			{Type: TMapOpen, Length: 2},
			TokStr("k1"),
			TokStr("v1"),
			TokStr("ke"),
			{Type: TArrOpen, Length: 3},
			TokStr("oh"),
			TokStr("whee"),
			TokStr("wow"),
			{Type: TArrClose},
			{Type: TMapClose},
		},
	},
	{"array nested in map as first and non-final entry",
		[]Token{
			{Type: TMapOpen, Length: 2},
			TokStr("ke"),
			{Type: TArrOpen, Length: 3},
			TokStr("oh"),
			TokStr("whee"),
			TokStr("wow"),
			{Type: TArrClose},
			TokStr("k1"),
			TokStr("v1"),
			{Type: TMapClose},
		},
	},
	{"maps nested in array",
		[]Token{
			{Type: TArrOpen, Length: 3},
			{Type: TMapOpen, Length: 1},
			TokStr("k"),
			TokStr("v"),
			{Type: TMapClose},
			TokStr("whee"),
			{Type: TMapOpen, Length: 1},
			TokStr("k1"),
			TokStr("v1"),
			{Type: TMapClose},
			{Type: TArrClose},
		},
	},
	{"arrays in arrays in arrays",
		[]Token{
			{Type: TArrOpen, Length: 1},
			{Type: TArrOpen, Length: 1},
			{Type: TArrOpen, Length: 0},
			{Type: TArrClose},
			{Type: TArrClose},
			{Type: TArrClose},
		},
	},
	{"maps nested in maps",
		[]Token{
			{Type: TMapOpen, Length: 1},
			TokStr("k"),
			{Type: TMapOpen, Length: 1},
			TokStr("k2"),
			TokStr("v2"),
			{Type: TMapClose},
			{Type: TMapClose},
		},
	},
	{"maps nested in maps with mixed nulls",
		[]Token{
			{Type: TMapOpen, Length: 2},
			TokStr("k"),
			{Type: TMapOpen, Length: 1},
			TokStr("k2"),
			TokStr("v2"),
			{Type: TMapClose},
			TokStr("k2"),
			{Type: TNull},
			{Type: TMapClose},
		},
	},
	{"map[str][]map[str]int",
		// this one is primarily for the objmapper tests
		[]Token{
			{Type: TMapOpen, Length: 1},
			TokStr("k"),
			{Type: TArrOpen, Length: 2},
			{Type: TMapOpen, Length: 1},
			TokStr("k2"),
			TokInt(1),
			{Type: TMapClose},
			{Type: TMapOpen, Length: 1},
			TokStr("k2"),
			TokInt(2),
			{Type: TMapClose},
			{Type: TArrClose},
			{Type: TMapClose},
		},
	},

	// Empty and null and null-at-depth.
	{"empty",
		[]Token{},
	},
	{"null",
		[]Token{
			{Type: TNull},
		},
	},
	{"null in array",
		[]Token{
			{Type: TArrOpen, Length: 1},
			{Type: TNull},
			{Type: TArrClose},
		},
	},
	{"null in map",
		[]Token{
			{Type: TMapOpen, Length: 1},
			TokStr("k"),
			{Type: TNull},
			{Type: TMapClose},
		},
	},
	{"null in array in array",
		[]Token{
			{Type: TArrOpen, Length: 1},
			{Type: TArrOpen, Length: 1},
			{Type: TNull},
			{Type: TArrClose},
			{Type: TArrClose},
		},
	},
	{"null in middle of array",
		[]Token{
			{Type: TArrOpen, Length: 5},
			TokStr("one"),
			{Type: TNull},
			TokStr("three"),
			{Type: TNull},
			TokStr("five"),
			{Type: TArrClose},
		},
	},

	// Numbers.
	// Warning: surprisingly contentious topic.
	// CBOR can't distinguish between positive numbers and unsigned;
	// JSON can't generally distinguish much of anything from anything, and is
	// subject to disasterous issues around floating point precision.
	//
	// Commented out because they're functionally useless -- packages define their own vagueries.
	// Except some which are used in the objmapping fixtures.
	//	{"integer zero", []Token{{Type: TInt, Int: 0}}},
	{"integer one", []Token{{Type: TInt, Int: 1}}},
	//	{"integer neg one", []Token{{Type: TInt, Int: -1}}},

	// Byte strings.
	// Warning: contentious topic.
	// JSON can't clearly represent binary types, and must use string transforms.
	{"short byte array",
		[]Token{
			{Type: TBytes, Bytes: []byte(`value`)}, // Note 'Length' field not used; would be redundant.
		},
	},
	{"long zero byte array",
		[]Token{
			{Type: TBytes, Bytes: make([]byte, 400)},
		},
	},

	// Tags.
	// Warning: contentious topic.
	// This is basically a CBOR-specific feature.
	// We also baked some support for it into the obj traversers, though that
	// of course also does not make much sense except used in combo with CBOR.
	{"tagged object",
		[]Token{
			{Type: TMapOpen, Length: 1, Tagged: true, Tag: 50},
			{Type: TString, Str: "k"},
			{Type: TString, Str: "v"},
			{Type: TMapClose},
		},
	},
	{"tagged string",
		[]Token{
			{Type: TString, Str: "wahoo", Tagged: true, Tag: 50},
		},
	},
	{"array with mixed tagged values",
		[]Token{
			{Type: TArrOpen, Length: 2},
			{Type: TUint, Uint: 400, Tagged: true, Tag: 40},
			{Type: TString, Str: "500", Tagged: true, Tag: 50},
			{Type: TArrClose},
		},
	},
	{"object with deeper tagged values",
		[]Token{
			{Type: TMapOpen, Length: 5},
			{Type: TString, Str: "k1"}, {Type: TString, Str: "500", Tagged: true, Tag: 50},
			{Type: TString, Str: "k2"}, {Type: TString, Str: "untagged"},
			{Type: TString, Str: "k3"}, {Type: TString, Str: "600", Tagged: true, Tag: 60},
			{Type: TString, Str: "k4"}, {Type: TArrOpen, Length: 2},
			/**/ {Type: TString, Str: "asdf", Tagged: true, Tag: 50},
			/**/ {Type: TString, Str: "qwer", Tagged: true, Tag: 50},
			/**/ {Type: TArrClose},
			{Type: TString, Str: "k5"}, {Type: TString, Str: "505", Tagged: true, Tag: 50},
			{Type: TMapClose},
		},
	},

	// Partial sequences!
	// Decoders may emit these before hitting an error (like EOF, or invalid following serial token).
	// Encoders may consume these, but ending after them would be an unexpected end of sequence error.
	{"dangling arr open",
		[]Token{
			{Type: TArrOpen, Length: 1},
		},
	},
}

// Returns a copy of the sequence with all length info at the start of maps and arrays stripped.
// Use this when testing e.g. json and cbor-in-stream-mode, which doesn't know lengths.
func (s Sequence) SansLengthInfo() Sequence {
	v := Sequence{s.Title, make([]Token, len(s.Tokens))}
	copy(v.Tokens, s.Tokens)
	for i := range v.Tokens {
		v.Tokens[i].Length = -1
	}
	return v
}

// Returns a copy of the sequence with the given token appened.
// This is mostly useful to test failure modes, like
// appending an invalid token at the end so decoder lengths match up.
func (s Sequence) Append(tok Token) Sequence {
	v := Sequence{s.Title, make([]Token, len(s.Tokens)+1)}
	copy(v.Tokens, s.Tokens)
	v.Tokens[len(s.Tokens)] = tok
	return v
}

// Sequences indexed by title.
var SequenceMap map[string]Sequence

func init() {
	SequenceMap = make(map[string]Sequence, len(Sequences))
	for _, v := range Sequences {
		SequenceMap[v.Title] = v
	}
}
