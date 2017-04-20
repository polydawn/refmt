package obj

import (
	"reflect"

	"github.com/polydawn/refmt/obj2/atlas"
	. "github.com/polydawn/refmt/tok"
)

/*
	Allocates the machinery for treating an in-memory object like a `TokenSink`.
	This machinery will walk over values,	using received tokens to fill in
	fields as it visits them.

	Initialization must be finished by calling `Bind` to set the value to visit;
	after this, the `Step` function is ready to be pumped.
	Subsequent calls to `Bind` do a full reset, leaving `Step` ready to call
	again and making all of the machinery reusable without re-allocating.
*/
func NewUnmarshaler(atl atlas.Atlas) *UnmarshalDriver {
	d := &UnmarshalDriver{
		unmarshalSlab: unmarshalSlab{
			atlas: atl,
			rows:  make([]unmarshalSlabRow, 0, 10),
		},
		stack: make([]UnmarshalMachine, 0, 10),
	}
	return d
}

func (d *UnmarshalDriver) Bind(v interface{}) error {
	d.stack = d.stack[0:0]
	d.unmarshalSlab.rows = d.unmarshalSlab.rows[0:0]
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		err := ErrInvalidUnmarshalTarget{reflect.TypeOf(v)}
		d.step = &errThunkUnmarshalMachine{err}
		return err
	}
	rt := rv.Type()
	d.step = d.unmarshalSlab.requisitionMachine(rt)
	return d.step.Reset(&d.unmarshalSlab, rv, rt)
}

type UnmarshalDriver struct {
	unmarshalSlab unmarshalSlab
	stack         []UnmarshalMachine
	step          UnmarshalMachine
}

type UnmarshalMachine interface {
	Reset(*unmarshalSlab, reflect.Value, reflect.Type) error
	Step(*UnmarshalDriver, *unmarshalSlab, *Token) (done bool, err error)
}

func (d *UnmarshalDriver) Step(tok *Token) (bool, error) {
	done, err := d.step.Step(d, &d.unmarshalSlab, tok)
	// If the step errored: out, entirely.
	if err != nil {
		return true, err
	}
	// If the step wasn't done, return same status.
	if !done {
		return false, nil
	}
	// If it WAS done, pop next, or if stack empty, we're entirely done.
	nSteps := len(d.stack) - 1
	if nSteps == -1 {
		return true, nil // that's all folks
	}
	d.step = d.stack[nSteps]
	d.stack = d.stack[0:nSteps]
	return false, nil
}
