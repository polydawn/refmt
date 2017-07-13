package json

import (
	"github.com/polydawn/refmt/tok/fixtures"
)

type situation byte

const (
	situationEncoding situation = 0x1
	situationDecoding situation = 0x2
)

var jsonFixtures = []struct {
	title    string
	sequence fixtures.Sequence
	serial   string
	only     situation
}{
	// Strings
	{"",
		fixtures.SequenceMap["empty string"],
		`""`,
		situationEncoding | situationDecoding,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["empty string"].SansLengthInfo(),
		`  "" `,
		situationDecoding,
	},
	{"",
		fixtures.SequenceMap["flat string"],
		`"value"`,
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.SequenceMap["strings needing escape"],
		`"str\nbroken\ttabbed"`,
		situationEncoding | situationDecoding,
	},

	// Maps
	{"",
		fixtures.SequenceMap["empty map"].SansLengthInfo(),
		`{}`,
		situationEncoding | situationDecoding,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["empty map"].SansLengthInfo(),
		`{  }`,
		situationDecoding,
	},
	{"",
		fixtures.SequenceMap["single row map"].SansLengthInfo(),
		`{"key":"value"}`,
		situationEncoding | situationDecoding,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["single row map"].SansLengthInfo(),
		` { "key"  :  "value" } `,
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.SequenceMap["duo row map"].SansLengthInfo(),
		`{"key":"value","k2":"v2"}`,
		situationEncoding | situationDecoding,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["duo row map"].SansLengthInfo(),
		`{"key":"value",  "k2":"v2"}`,
		situationEncoding | situationDecoding,
	},
	{"decoding with trailing comma",
		fixtures.SequenceMap["duo row map"].SansLengthInfo(),
		`{"key":"value","k2":"v2",}`,
		situationEncoding | situationDecoding,
	},
	{"",
		fixtures.SequenceMap["duo row map alt2"].SansLengthInfo(),
		`{"k2":"v2","key":"value"}`,
		situationDecoding,
	},

	// Arrays
	{"",
		fixtures.SequenceMap["empty array"].SansLengthInfo(),
		`[]`,
		situationEncoding | situationDecoding,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["empty array"].SansLengthInfo(),
		`  [ ] `,
		situationDecoding,
	},
	{"",
		fixtures.SequenceMap["single entry array"].SansLengthInfo(),
		`["value"]`,
		situationEncoding | situationDecoding,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["single entry array"].SansLengthInfo(),
		`  [ "value" ] `,
		situationDecoding,
	},
	{"",
		fixtures.SequenceMap["duo entry array"].SansLengthInfo(),
		`["value","v2"]`,
		situationEncoding | situationDecoding,
	},
	{"decoding with extra whitespace",
		fixtures.SequenceMap["duo entry array"].SansLengthInfo(),
		`["value",  "v2"]`,
		situationDecoding,
	},

	// Complex / mixed / nested.
	{"",
		fixtures.SequenceMap["array nested in map as non-first and final entry"].SansLengthInfo(),
		`{"k1":"v1","ke":["oh","whee","wow"]}`,
		situationEncoding | situationDecoding,
	},
}
