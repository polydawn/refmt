package obj

import (
	. "github.com/polydawn/go-xlate/tok"
)

/*
	Returns a `TokenSource` that will walk over structures in memory,
	emitting tokens representing values and fields as it visits them.
*/
func NewMarshaler(v interface{} /* TODO visitmagicks */) *MarshalDriver {
	d := &MarshalDriver{
		step: pickMarshalMachine(v),
	}
	return d
}

type MarshalDriver struct {
	stack []MarshalMachine
	step  MarshalMachine
}

type MarshalMachine interface {
	Step(*MarshalDriver, *Token) (done bool, err error)
}

// for convenience in declaring fields of state machines with internal step funcs
type marshalMachineStep func(*MarshalDriver, *Token) (done bool, err error)

func (d *MarshalDriver) Step(tok *Token) (bool, error) {
	done, err := d.step.Step(d, tok)
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

/*
	Traverses `target`,
	first looking up the machine for that type just like it's a new top-level object,
	then pushing the first step with `tok` (the upstream tends to have peeked at it
	in order to decide what to do, but if recursing, it belongs to the next obj),
	then saving this new machine: the driver will then continuing stepping the
	new machine it returns a done status, at which point we'll finally
	"return" by popping back to the last machine on the stack.

	In other words, your MarshalMachine calls this when it wants to deal
	with an object, and by the time we call back to your machine again,
	that object will be traversed and the stream ready for you to continue.
*/
func (d *MarshalDriver) Recurse(tok *Token, target interface{}) error {
	// Push the current machine onto the stack (we'll resume it when the new one is done),
	// and pick a machine to start in on our next item to cover.
	d.stack = append(d.stack, d.step)
	d.step = pickMarshalMachine(target) // TODO caller should be able to override this
	// Immediately make a step (we're still the delegate in charge of someone else's step).
	_, err := d.Step(tok)
	return err
}
