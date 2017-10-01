package obj

import (
	"reflect"

	. "github.com/polydawn/refmt/tok"
)

type unmarshalMachineSliceWildcard struct {
	target_rv reflect.Value
	value_rt  reflect.Type
	valueMach UnmarshalMachine
	step      unmarshalMachineStep
	index     int
}

func (mach *unmarshalMachineSliceWildcard) Reset(slab *unmarshalSlab, rv reflect.Value, rt reflect.Type) error {
	mach.target_rv = rv
	mach.value_rt = rt.Elem()
	mach.valueMach = slab.requisitionMachine(mach.value_rt)
	mach.step = mach.step_Initial
	mach.index = 0
	return nil
}

func (mach *unmarshalMachineSliceWildcard) Step(driver *Unmarshaller, slab *unmarshalSlab, tok *Token) (done bool, err error) {
	return mach.step(driver, slab, tok)
}

func (mach *unmarshalMachineSliceWildcard) step_Initial(_ *Unmarshaller, slab *unmarshalSlab, tok *Token) (done bool, err error) {
	// If it's a special state, start an object.
	//  (Or, blow up if its a special state that's silly).
	switch tok.Type {
	case TMapOpen:
		return true, ErrMalformedTokenStream{tok.Type, "start of array"}
	case TArrOpen:
		// Great.  Consumed.
		mach.step = mach.step_AcceptValue
		// Initialize the slice.
		mach.target_rv.Set(reflect.MakeSlice(mach.target_rv.Type(), 0, 0))
		return false, nil
	case TMapClose:
		return true, ErrMalformedTokenStream{tok.Type, "start of array"}
	case TArrClose:
		return true, ErrMalformedTokenStream{tok.Type, "start of array"}
	case TNull:
		mach.target_rv.Set(reflect.Zero(mach.target_rv.Type()))
		return true, nil
	default:
		return true, ErrMalformedTokenStream{tok.Type, "start of array"}
	}
}

func (mach *unmarshalMachineSliceWildcard) step_AcceptValue(driver *Unmarshaller, slab *unmarshalSlab, tok *Token) (done bool, err error) {
	// Either form of open token are valid, but
	// - an arrClose is ours
	// - and a mapClose is clearly invalid.
	switch tok.Type {
	case TMapClose:
		// no special checks for ends of wildcard slice; no such thing as incomplete.
		return true, ErrMalformedTokenStream{tok.Type, "start of value or end of array"}
	case TArrClose:
		// Finishing step: push our current slice ref all the way to original target.
		// REVIEW does this even require an action anymore? // *(mach.target) = mach.slice
		return true, nil
	}

	// Grow the slice if necessary.
	// FIXME this is ridiculously inefficient, can do much better, this is placeholder quality
	mach.target_rv.Set(reflect.Append(mach.target_rv, reflect.Zero(mach.value_rt)))

	// Recurse on a handle to the next index.
	rv := mach.target_rv.Index(mach.index)
	mach.index++
	return false, driver.Recurse(tok, rv, mach.value_rt, mach.valueMach)
	// Step simply remains `step_AcceptValue` -- arrays don't have much state machine.
}
