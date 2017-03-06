package cbor

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/polydawn/go-xlate/tok"
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
	tt := cborFixtures
	for _, tr := range tt {
		// Set it up.
		title := tr.sequence.Title
		if tr.title != "" {
			title = strings.Join([]string{tr.sequence.Title, tr.title}, ", ")
		}
		buf := bytes.NewBuffer(tr.serial)
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
					title, n, expectTok, tok)
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
