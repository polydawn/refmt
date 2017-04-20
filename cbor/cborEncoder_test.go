package cbor

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/polydawn/refmt/testutil"
)

func TestCborEncoder(t *testing.T) {
	tt := cborFixtures
	// Loop over test table.
	for _, tr := range tt {
		// Skip if not tagged for encoding.
		if tr.only&situationEncoding == 0 {
			continue
		}

		// Set it up.
		title := tr.sequence.Title
		if tr.title != "" {
			title = strings.Join([]string{tr.sequence.Title, tr.title}, ", ")
		}
		buf := &bytes.Buffer{}
		tokenSink := NewEncoder(buf)

		// Run steps.
		var done bool
		var err error
		for n, tok := range tr.sequence.Tokens {
			done, err = tokenSink.Step(&tok)
			if err != nil {
				t.Errorf("test %q step %d (inputting %s) errored: %s", title, n, tok, err)
			}
			if done && n != len(tr.sequence.Tokens)-1 {
				t.Errorf("test %q done early! on index=%d out of %d tokens", title, n, len(tr.sequence.Tokens))
			}
		}
		if !done {
			t.Errorf("test %q still not done after %d tokens!", title, len(tr.sequence.Tokens))
		}

		// Assert final result.
		Assert(t, title, tr.serial, buf.Bytes())

		t.Logf("test %q --- done", title)
	}
}
