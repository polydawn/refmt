package obj

import (
	"github.com/polydawn/go-xlate/tok"
)

/*
	Returns a `TokenSink` that will unmarshal tokens into an in-memory value.
*/
func NewUnmarshaler(v interface{} /* TODO visitmagicks */) *unmarshalDriver {
	d := &unmarshalDriver{}
	d.stack = []UnmarshalMachine{
	// todo fixme etc
	}
	return d
}

type unmarshalDriver struct {
	stack []UnmarshalMachine
	top   UnmarshalMachine
}

func (d *unmarshalDriver) Step(consume tok.Token) (done bool, err error) {
	return false, nil
}

type UnmarshalMachine func(*unmarshalDriver, tok.Token) (done bool, err error)
