package obj

import (
	"fmt"
	"reflect"

	. "github.com/polydawn/refmt/tok"
)

type unmarshalMachineWildcard struct {
	target_rv reflect.Value
	delegate  UnmarshalMachine // actual machine, once we've demuxed with the first token.
}

func (mach *unmarshalMachineWildcard) Reset(_ *unmarshalSlab, rv reflect.Value, _ reflect.Type) error {
	mach.target_rv = rv
	return nil
}

func (mach *unmarshalMachineWildcard) Step(driver *UnmarshalDriver, slab *unmarshalSlab, tok *Token) (done bool, err error) {
	if mach.delegate == nil {
		return mach.step_demux(driver, slab, tok)
	}
	return mach.delegate.Step(driver, slab, tok)
}

func (mach *unmarshalMachineWildcard) step_demux(driver *UnmarshalDriver, slab *unmarshalSlab, tok *Token) (done bool, err error) {
	// Switch on token type: we may be able to delegate to a primitive machine,
	//  but we may also need to initialize a container type and then hand off.
	switch tok.Type {
	case TMapOpen:
		child := make(map[string]interface{})
		child_rv := reflect.ValueOf(child)
		mach.target_rv.Set(child_rv)
		mach.delegate = &slab.tip().unmarshalMachineMapStringWildcard
		if err := mach.delegate.Reset(slab, child_rv, child_rv.Type()); err != nil {
			return true, err
		}
		return mach.delegate.Step(driver, slab, tok)

	case TArrOpen:
		child := make([]interface{}, 0) // TODO if we have length hint info, use it.
		child_rv := reflect.ValueOf(child)
		mach.target_rv.Set(child_rv) // REVIEW this seems unlikely to be The Way for slices...
		mach.delegate = &slab.tip().unmarshalMachineSliceWildcard
		if err := mach.delegate.Reset(slab, child_rv, child_rv.Type()); err != nil {
			return true, err
		}
		return mach.delegate.Step(driver, slab, tok)

	case TMapClose:
		return true, fmt.Errorf("unexpected mapClose; expected start of value")

	case TArrClose:
		return true, fmt.Errorf("unexpected arrClose; expected start of value")

	case TNull:
		mach.target_rv.Set(reflect.ValueOf(nil))
		return true, nil

	default:
		// If it wasn't the start of composite, shell out to the machine for literals.
		// Don't bother to replace our internal step func because literal machines are never multi-call.
		delegateMach := slab.tip().unmarshalMachinePrimitive
		delegateMach.kind = reflect.Interface
		if err := delegateMach.Reset(slab, mach.target_rv, nil); err != nil {
			return true, err
		}
		return delegateMach.Step(driver, slab, tok)
	}
}
