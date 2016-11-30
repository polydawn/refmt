package xlate

import (
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

func (p TokenPump) Run() {
	// TODO
}

func (p TokenPump) step() {
	// TODO
}
