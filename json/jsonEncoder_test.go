package json

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/polydawn/refmt/testutil"
	"github.com/polydawn/refmt/tok/fixtures"
)

func TestJsonEncoder(t *testing.T) {
	tt := []struct {
		title    string
		sequence fixtures.Sequence
		expect   string
	}{
		{"",
			fixtures.SequenceMap["flat string"],
			`"value"`,
		},
		{"",
			fixtures.SequenceMap["single row map"],
			`{"key":"value"}`,
		},
		{"",
			fixtures.SequenceMap["duo row map"],
			`{"key":"value","k2":"v2"}`,
		},
		{"",
			fixtures.SequenceMap["single entry array"],
			`["value"]`,
		},
		{"",
			fixtures.SequenceMap["duo entry array"],
			`["value","v2"]`,
		},
		{"",
			fixtures.SequenceMap["empty map"],
			`{}`,
		},
		{"",
			fixtures.SequenceMap["empty array"],
			`[]`,
		},
		{"",
			fixtures.SequenceMap["array nested in map as non-first and final entry"],
			`{"k1":"v1","ke":["oh","whee","wow"]}`,
		},
		{"",
			fixtures.SequenceMap["array nested in map as first and non-final entry"],
			`{"ke":["oh","whee","wow"],"k1":"v1"}`,
		},
		{"",
			fixtures.SequenceMap["maps nested in array"],
			`[{"k":"v"},"whee",{"k1":"v1"}]`,
		},
		{"",
			fixtures.SequenceMap["arrays in arrays in arrays"],
			`[[[]]]`,
		},
		{"",
			fixtures.SequenceMap["maps nested in maps"],
			`{"k":{"k2":"v2"}}`,
		},
		{"",
			fixtures.SequenceMap["strings needing escape"],
			`"str\nbroken\ttabbed"`,
		},
	}
	for _, tr := range tt {
		// Set it up.
		title := tr.sequence.Title
		if tr.title != "" {
			title = strings.Join([]string{tr.sequence.Title, tr.title}, ", ")
		}
		buf := &bytes.Buffer{}
		sink := NewEncoder(buf)

		// Run steps.
		var done bool
		var err error
		for n, tok := range tr.sequence.Tokens {
			done, err = sink.Step(&tok)
			if err != nil {
				t.Errorf("test %q step %d (inputting %#v) errored: %s", title, n, tok, err)
			}
			if done && n != len(tr.sequence.Tokens)-1 {
				t.Errorf("test %q done early! on step %d out of %d tokens", title, n, len(tr.sequence.Tokens))
			}
		}
		if !done {
			t.Errorf("test %q still not done after %d tokens!", title, len(tr.sequence.Tokens))
		}

		// Assert final result.
		Assert(t, title, tr.expect, buf.String())
	}
}
