package obj

import (
	"testing"

	//. "github.com/polydawn/go-xlate/testutil"
	. "github.com/polydawn/go-xlate/tok"
)

func TestMarshaller(t *testing.T) {
	tt := []struct {
		title     string
		targetFn  func() interface{} // func returns target, so test source looks like your call param
		expectSeq []Token
		expectErr error
	}{
		{
			title:    "simple literal",
			targetFn: func() interface{} { i := 4; return &i },
			expectSeq: []Token{
				4,
			},
		},
	}
	for _, tr := range tt {
		suite := &Suite{}
		marshaller := NewMarshaler(suite, tr.targetFn())

		// Run steps.
		var done bool
		var err error
		var tok Token
		for n, expectTok := range tr.expectSeq {
			done, err = marshaller.Step(&tok)
			if err != nil {
				t.Errorf("step %d (expecting %#v) errored: %s", n, expectTok, err)
			}
			if done && n != len(tr.expectSeq)-1 {
				t.Errorf("done early! on step %d out of %d tokens", n, len(tr.expectSeq))
			}
		}
		if !done {
			t.Errorf("still not done after %d tokens!", len(tr.expectSeq))
		}
		t.Logf("test %q halted correctly and passed", tr.title)
	}
}
