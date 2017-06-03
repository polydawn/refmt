package obj

import (
	"fmt"
	"reflect"

	. "github.com/polydawn/refmt/tok"
)

type unmarshalMachineMapStringWildcard struct {
	target_rv reflect.Value
	value_rt  reflect.Type
	valueMach UnmarshalMachine
	step      unmarshalMachineStep
	key_rv    reflect.Value // The key consumed by the prev `step_AcceptKey`.
	tmp_rv    reflect.Value // Addressable handle to a slot for values to unmarshal into.
}

func (mach *unmarshalMachineMapStringWildcard) Reset(slab *unmarshalSlab, rv reflect.Value, rt reflect.Type) error {
	mach.target_rv = rv
	mach.value_rt = rt.Elem()
	mach.valueMach = slab.requisitionMachine(mach.value_rt)
	mach.step = mach.step_Initial
	mach.key_rv = reflect.Value{}
	slot_rv := reflect.New(mach.value_rt)
	mach.tmp_rv = slot_rv.Elem()
	return nil
}

func (mach *unmarshalMachineMapStringWildcard) Step(driver *UnmarshalDriver, slab *unmarshalSlab, tok *Token) (done bool, err error) {
	return mach.step(driver, slab, tok)
}

func (mach *unmarshalMachineMapStringWildcard) step_Initial(_ *UnmarshalDriver, _ *unmarshalSlab, tok *Token) (done bool, err error) {
	// If it's a special state, start an object.
	//  (Or, blow up if its a special state that's silly).
	switch tok.Type {
	case TMapOpen:
		// Great.  Consumed.
		mach.step = mach.step_AcceptKey
		// Initialize the map if it's nil.
		if mach.target_rv.IsNil() {
			mach.target_rv.Set(reflect.MakeMap(mach.target_rv.Type()))
		}
		return false, nil
	case TMapClose:
		return true, fmt.Errorf("unexpected mapClose; expected start of map")
	case TArrClose:
		return true, fmt.Errorf("unexpected arrClose; expected start of map")
	case TArrOpen:
		fallthrough
	default:
		return true, ErrUnmarshalIncongruent{*tok, mach.target_rv}
	}
}

func (mach *unmarshalMachineMapStringWildcard) step_AcceptKey(_ *UnmarshalDriver, _ *unmarshalSlab, tok *Token) (done bool, err error) {
	// First, save any refs from the last value.
	//  (This is fiddly: the delay comes mostly from the handling of slices, which may end up re-allocating
	//   themselves during their decoding.)
	if mach.key_rv != (reflect.Value{}) {
		mach.target_rv.SetMapIndex(mach.key_rv, mach.tmp_rv)
	}
	// Now switch on tokens.
	switch tok.Type {
	case TMapOpen:
		return true, fmt.Errorf("unexpected mapOpen; expected map key")
	case TArrOpen:
		return true, fmt.Errorf("unexpected arrOpen; expected map key")
	case TMapClose:
		// no special checks for ends of wildcard map; no such thing as incomplete.
		return true, nil
	case TArrClose:
		return true, fmt.Errorf("unexpected arrClose; expected map key")
	}
	switch tok.Type {
	case TString:
		key_rv := reflect.ValueOf(tok.Str)
		if err = mach.mustAcceptKey(key_rv); err != nil {
			return true, err
		}
		mach.key_rv = key_rv
		mach.step = mach.step_AcceptValue
		return false, nil
	default:
		return true, fmt.Errorf("unexpected token %s; expected key string or end of map", tok)
	}
}

func (mach *unmarshalMachineMapStringWildcard) mustAcceptKey(key_rv reflect.Value) error {
	if exists := mach.target_rv.MapIndex(key_rv).IsValid(); exists {
		return fmt.Errorf("repeated key %q", key_rv)
	}
	return nil
}

func (mach *unmarshalMachineMapStringWildcard) step_AcceptValue(driver *UnmarshalDriver, slab *unmarshalSlab, tok *Token) (done bool, err error) {
	mach.step = mach.step_AcceptKey
	return false, driver.Recurse(
		tok,
		mach.tmp_rv,
		mach.value_rt,
		mach.valueMach,
	)
}