package fixtures

import (
	"errors"

	. "github.com/polydawn/refmt/tok"
)

type Sequence struct {
	Title  string
	Tokens []Token
}

// An array of well-formed token sequences.
var Sequences = []Sequence{
	// Strings.
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

// Sequences indexed by title.
var SequenceMap map[string]Sequence

func init() {
	SequenceMap = make(map[string]Sequence, len(Sequences))
	for _, v := range Sequences {
		SequenceMap[v.Title] = v
	}
}

// Labels for which way a token sequence fixture is malformed.
type Malformed error

var (
	MalformedUnterminatedArray Malformed = errors.New("MalformedUnterminatedArray")
	MalformedUnterminatedMap   Malformed = errors.New("MalformedUnterminatedMap")
	MalformedNilMapKey         Malformed = errors.New("MalformedNilMapKey")
	MalformedUnbalancedMap     Malformed = errors.New("MalformedUnbalancedMap")
)

// Any array of token sequences that are in some way malformed.
// TokenSinks (i.e. encoders, serializes) should be able to halt and error reasonably on these.
// It's less simple to get TokenSources to emit them, but for decoders comparable inputs should have their own test coverage.
var MalformedSequences = []struct {
	Title     string
	Seq       []Token
	Malformed Malformed
}{}
