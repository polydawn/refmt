package fixtures

import (
	"errors"

	. "github.com/polydawn/go-xlate/tok"
)

type Sequence struct {
	Title  string
	Tokens []Token
}

// An array of well-formed token sequences.
var Sequences = []Sequence{
	{"flat string",
		[]Token{
			TokStr("value"),
		},
	},
	{"single row map",
		[]Token{
			Token_MapOpen,
			TokStr("key"),
			TokStr("value"),
			Token_MapClose,
		},
	},
	{"duo row map",
		[]Token{
			Token_MapOpen,
			TokStr("key"),
			TokStr("value"),
			TokStr("k2"),
			TokStr("v2"),
			Token_MapClose,
		},
	},
	{"single entry array",
		[]Token{
			Token_ArrOpen,
			TokStr("value"),
			Token_ArrClose,
		},
	},
	{"duo entry array",
		[]Token{
			Token_ArrOpen,
			TokStr("value"),
			TokStr("v2"),
			Token_ArrClose,
		},
	},
	{"empty map",
		[]Token{
			Token_MapOpen,
			Token_MapClose,
		},
	},
	{"empty array",
		[]Token{
			Token_ArrOpen,
			Token_ArrClose,
		},
	},
	{"array nested in map as non-first and final entry",
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
	},
	{"array nested in map as first and non-final entry",
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
	},
	{"maps nested in array",
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
	},
	{"arrays in arrays in arrays",
		[]Token{
			Token_ArrOpen,
			Token_ArrOpen,
			Token_ArrOpen,
			Token_ArrClose,
			Token_ArrClose,
			Token_ArrClose,
		},
	},
	{"maps nested in maps",
		[]Token{
			Token_MapOpen,
			TokStr("k"),
			Token_MapOpen,
			TokStr("k2"),
			TokStr("v2"),
			Token_MapClose,
			Token_MapClose,
		},
	},
	{"strings needing escape",
		[]Token{
			TokStr("str\nbroken\ttabbed"),
		},
	},
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
