package xlate

import (
	"fmt"

	. "github.com/polydawn/go-xlate/tok"
)

type TokenSource interface {
	Step(fillme *Token) (done bool, err error)
}

type TokenSink interface {
	Step(consume *Token) (done bool, err error)
}

type TokenPump struct {
	TokenSource
	TokenSink
}

func (p TokenPump) Run() error {
	var tok Token
	var srcDone, sinkDone bool
	var err error
	for {
		srcDone, err = p.TokenSource.Step(&tok)
		if err != nil {
			return err
		}
		sinkDone, err = p.TokenSink.Step(&tok)
		if err != nil {
			return err
		}
		if srcDone {
			if sinkDone {
				return nil
			}
			return fmt.Errorf("src at end of item but sink expects more")
		}
	}
}
