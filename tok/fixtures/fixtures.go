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

func init() {
	Sequences = append(Sequences, sequences_Bool...)
	Sequences = append(Sequences, sequences_String...)
	Sequences = append(Sequences, sequences_Map...)
	Sequences = append(Sequences, sequences_Array...)
	Sequences = append(Sequences, sequences_Composite...)
	Sequences = append(Sequences, sequences_Null...)
	Sequences = append(Sequences, sequences_Number...)
	Sequences = append(Sequences, sequences_Bytes...)
	Sequences = append(Sequences, sequences_Tag...)
}

var Sequences []Sequence

// Returns a copy of the sequence with all length info at the start of maps and arrays stripped.
// Use this when testing e.g. json and cbor-in-stream-mode, which doesn't know lengths.
func (s Sequence) SansLengthInfo() Sequence {
	v := Sequence{s.Title, make([]Token, len(s.Tokens))}
	copy(v.Tokens, s.Tokens)
	StompLengths(v.Tokens)
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
