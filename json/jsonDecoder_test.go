package json

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	. "github.com/polydawn/refmt/tok"
)

func TestJsonDecoder(t *testing.T) {
	// Conjoin step index and token field -- helps us cram a little more info into convey assertion messages
	type step struct {
		tok string
		err error
	}

	Convey("JSON Decoder suite:", t, func() {
		for _, tr := range jsonFixtures {
			// Ignore this row if tagged as inapplicable to decoding.
			if tr.decodeResult == inapplicable {
				continue
			}

			title := tr.sequence.Title
			if tr.title != "" {
				title = strings.Join([]string{tr.sequence.Title, tr.title}, ", ")
			}
			Convey(fmt.Sprintf("%q fixture sequence:", title), func() {
				buf := bytes.NewBuffer([]byte(tr.serial))
				tokenSource := NewDecoder(buf)

				Convey("Steps...", func() {
					var done bool
					var err error
					var nStep int
					expectSteps := len(tr.sequence.Tokens) - 1
					for ; nStep <= expectSteps; nStep++ {
						expectTok := tr.sequence.Tokens[nStep]
						var tok Token
						done, err = tokenSource.Step(&tok)
						So(tok.String(), ShouldResemble, expectTok.String())
						if done || err != nil {
							break
						}
					}
					Convey("Result", FailureContinues, func() {
						So(nStep, ShouldEqual, expectSteps)
						So(err, ShouldEqual, tr.decodeResult)
					})
				})
			})
		}
	})
}
