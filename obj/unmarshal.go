package obj

import (
	. "github.com/polydawn/go-xlate/tok"
)

/*
	Returns a `TokenSink` that will unmarshal tokens into an in-memory value.
*/
func NewUnmarshaler(v interface{} /* TODO visitmagicks */) *UnmarshalDriver {
	d := &UnmarshalDriver{
		step: pickUnmarshalMachine(v),
	}
	return d
}

type UnmarshalDriver struct {
	stack []UnmarshalMachine
	step  UnmarshalMachine
}

type UnmarshalMachine interface {
	Step(*UnmarshalDriver, *Token) (done bool, err error)
}

// for convenience in declaring fields of state machines with internal step funcs
type unmarshalMachineStep func(*UnmarshalDriver, *Token) (done bool, err error)

func (d *UnmarshalDriver) Step(tok *Token) (bool, error) {
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
	Fills `target`,
	first looking up the machine for that type just like it's a new top-level object,
	then pushing the first step with `tok` (the upstream tends to have peeked at it
	in order to decide what to do, but if recursing, it belongs to the next obj),
	then saving this new machine: the driver will then continuing stepping the
	new machine it returns a done status, at which point we'll finally
	"return" by popping back to the last machine on the stack.

	In other words, your UnmarshalMachine calls this when it wants to deal
	with an object, and by the time we call back to your machine again,
	that object will be filled and the stream ready for you to continue.
*/
func (d *UnmarshalDriver) Recurse(tok *Token, target interface{}) error {
	// Push the current machine onto the stack (we'll resume it when the new one is done),
	// and pick a machine to start in on our next item to cover.
	d.stack = append(d.stack, d.step)
	d.step = pickUnmarshalMachine(target) // TODO caller should be able to override this
	// Immediately make a step (we're still the delegate in charge of someone else's step).
	_, err := d.Step(tok)
	return err
}

// Picks an unmarshal machine, returning the custom impls for any
// common/primitive types, and advanced machines where structs get involved.
func pickUnmarshalMachine(v interface{}) UnmarshalMachine {
	switch v2 := v.(type) {
	// For total wildcards:
	//  Return a machine that will pick between a literal or `map[string]interface{}`
	//  or `[]interface{}` based on the next token.
	case *interface{}:
		return newUnmarshalMachineWildcard(v2)
	// For single literals:
	//  we have a single machine that handles all these.
	case *string, *[]byte,
		*int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64:
		return UnmarshalMachineLiteral{v}
	// Anything that has real type info:
	//  Look up what kind of machine to use, based on the type info, and use it.
	default:
		// TODO mustAddressable check goes here.
		panic(ErrUnreachable{"TODO mappersuite lookup"})
	}
}
