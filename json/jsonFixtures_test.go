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

	// Maps
	{"",
		fixtures.SequenceMap["empty map"].SansLengthInfo(),
		`{}`,
		situationEncoding | situationDecoding,
	},

	// Arrays
	{"",
		fixtures.SequenceMap["empty array"].SansLengthInfo(),
		`[]`,
		situationEncoding | situationDecoding,
	},
}
