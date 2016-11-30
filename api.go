package xlate

import (
	"github.com/polydawn/go-xlate/tok"
)

type TokenSource interface {
	Step(fillme *tok.Token) (done bool, err error)
}

type TokenSink interface {
	Step(consume tok.Token) (done bool, err error)
}

type TokenPump struct {
	TokenSource
	TokenSink
}
