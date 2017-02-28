package cbor

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/polydawn/go-xlate/tok"
	"github.com/polydawn/go-xlate/tok/fixtures"
)

func bcat(bss ...[]byte) []byte {
	l := 0
	for _, bs := range bss {
		l += len(bs)
	}
	rbs := make([]byte, 0, l)
	for _, bs := range bss {
		rbs = append(rbs, bs...)
	}
	return rbs
}

func b(b byte) []byte { return []byte{b} }

func TestCborDecoder(t *testing.T) {
	tt := []struct {
		title    string
		sequence fixtures.Sequence
		input    []byte
	}{
		{"",
			fixtures.SequenceMap["flat string"],
			bcat(b(0x60+5), []byte(`value`)),
		},
		{"indefinite length string (single actual hunk)",
			fixtures.SequenceMap["flat string"],
			bcat(b(0x7f), b(0x60+5), []byte(`value`), b(0xff)),
		},
		{"indefinite length string (multiple hunks)",
			fixtures.SequenceMap["flat string"],
			bcat(b(0x7f), b(0x60+2), []byte(`va`), b(0x60+3), []byte(`lue`), b(0xff)),
		},
		{"",
			fixtures.SequenceMap["strings needing escape"],
			bcat(b(0x60+17), []byte("str\nbroken\ttabbed")),
		},

		{"",
			fixtures.SequenceMap["empty map"],
			bcat(b(0xa0)),
		},
		{"indefinite length",
			fixtures.SequenceMap["empty map"],
			bcat(b(0xbf), b(0xff)),
		},
		{"",
			fixtures.SequenceMap["single row map"],
			bcat(b(0xa0+1),
				b(0x60+3), []byte(`key`), b(0x60+5), []byte(`value`),
			),
		},
		{"indefinite length",
			fixtures.SequenceMap["single row map"],
			bcat(b(0xbf),
				b(0x60+3), []byte(`key`), b(0x60+5), []byte(`value`),
				b(0xff),
			),
		},
		{"",
			fixtures.SequenceMap["duo row map"],
			bcat(b(0xa0+2),
				b(0x60+3), []byte(`key`), b(0x60+5), []byte(`value`),
				b(0x60+2), []byte(`k2`), b(0x60+2), []byte(`v2`),
			),
		},
		{"indefinite length",
			fixtures.SequenceMap["duo row map"],
			bcat(b(0xbf),
				b(0x60+3), []byte(`key`), b(0x60+5), []byte(`value`),
				b(0x60+2), []byte(`k2`), b(0x60+2), []byte(`v2`),
				b(0xff),
			),
		},

		{"",
			fixtures.SequenceMap["empty array"],
			bcat(b(0x80)),
		},
		{"indefinite length",
			fixtures.SequenceMap["empty array"],
			bcat(b(0x9f), b(0xff)),
		},
		{"",
			fixtures.SequenceMap["single entry array"],
			bcat(b(0x80+1),
				b(0x60+5), []byte(`value`),
			),
		},
		{"indefinite length",
			fixtures.SequenceMap["single entry array"],
			bcat(b(0x9f),
				b(0x60+5), []byte(`value`),
				b(0xff),
			),
		},
		{"indefinite length with nested indef string",
			fixtures.SequenceMap["single entry array"],
			bcat(b(0x9f),
				bcat(b(0x7f), b(0x60+5), []byte(`value`), b(0xff)),
				b(0xff),
			),
		},
		{"",
			fixtures.SequenceMap["duo entry array"],
			bcat(b(0x80+2),
				b(0x60+5), []byte(`value`),
				b(0x60+2), []byte(`v2`),
			),
		},
		{"indefinite length",
			fixtures.SequenceMap["duo entry array"],
			bcat(b(0x9f),
				b(0x60+5), []byte(`value`),
				b(0x60+2), []byte(`v2`),
				b(0xff),
			),
		},

		{"all indefinite length",
			fixtures.SequenceMap["array nested in map as non-first and final entry"],
			bcat(b(0xbf),
				b(0x60+2), []byte(`k1`), b(0x60+2), []byte(`v1`),
				b(0x60+2), []byte(`ke`), bcat(b(0x9f),
					b(0x60+2), []byte(`oh`),
					b(0x60+4), []byte(`whee`),
					b(0x60+3), []byte(`wow`),
					b(0xff),
				),
				b(0xff),
			),
		},
		{"all indefinite length",
			fixtures.SequenceMap["array nested in map as first and non-final entry"],
			bcat(b(0xbf),
				b(0x60+2), []byte(`ke`), bcat(b(0x9f),
					b(0x60+2), []byte(`oh`),
					b(0x60+4), []byte(`whee`),
					b(0x60+3), []byte(`wow`),
					b(0xff),
				),
				b(0x60+2), []byte(`k1`), b(0x60+2), []byte(`v1`),
				b(0xff),
			),
		},
		{"all indefinite length",
			fixtures.SequenceMap["maps nested in array"],
			bcat(b(0x9f),
				bcat(b(0xbf),
					b(0x60+1), []byte(`k`), b(0x60+1), []byte(`v`),
					b(0xff),
				),
				b(0x60+4), []byte(`whee`),
				bcat(b(0xbf),
					b(0x60+2), []byte(`k1`), b(0x60+2), []byte(`v1`),
					b(0xff),
				),
				b(0xff),
			),
		},
		{"all indefinite length",
			fixtures.SequenceMap["arrays in arrays in arrays"],
			bcat(b(0x9f), b(0x9f), b(0x9f), b(0xff), b(0xff), b(0xff)),
		},
		{"all indefinite length",
			fixtures.SequenceMap["maps nested in maps"],
			bcat(b(0xbf),
				b(0x60+1), []byte(`k`), bcat(b(0xbf),
					b(0x60+2), []byte(`k2`), b(0x60+2), []byte(`v2`),
					b(0xff),
				),
				b(0xff),
			),
		},
	}
	for _, tr := range tt {
		// Set it up.
		title := tr.sequence.Title
		if tr.title != "" {
			title = strings.Join([]string{tr.sequence.Title, tr.title}, ", ")
		}
		buf := bytes.NewBuffer(tr.input)
		tokenSource := NewDecoder(buf)

		// Run steps.
		var done bool
		var err error
		var tok Token
		for n, expectTok := range tr.sequence.Tokens {
			done, err = tokenSource.Step(&tok)
			if err != nil {
				t.Errorf("test %q step %d (inputting %#v) errored: %s", title, n, expectTok, err)
			}
			if !IsTokenEqual(expectTok, tok) {
				t.Errorf("test %q failed: step %d yielded wrong token: expected %s, got %s",
					title, n, TokenToString(expectTok), TokenToString(tok))
			}
			if done && n != len(tr.sequence.Tokens)-1 {
				t.Errorf("test %q done early! on index=%d out of %d tokens", title, n, len(tr.sequence.Tokens))
			}
		}
		if !done {
			t.Errorf("test %q still not done after %d tokens!", title, len(tr.sequence.Tokens))
		}

		t.Logf("test %q --- done", title)
	}
}
