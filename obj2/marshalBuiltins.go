package obj

import (
	"reflect"

	. "github.com/polydawn/go-xlate/tok"
)

type ptrDerefDelegateMarshalMachine struct {
	MarshalMachine
	peelCount int

	isNil bool
}

func (mach *ptrDerefDelegateMarshalMachine) Reset(slab *marshalSlab, rv reflect.Value, rt reflect.Type) error {
	mach.isNil = false
	for i := 0; i < mach.peelCount; i++ {
		rv = rv.Elem()
		if rv.IsNil() {
			mach.isNil = true
			return nil
		}
	}
	return mach.MarshalMachine.Reset(slab, rv, rt) // REVIEW: this rt should be peeled by here.  do we... ignore the arg and cache it at mach conf time?
}
func (mach *ptrDerefDelegateMarshalMachine) Step(driver *MarshalDriver, slab *marshalSlab, tok *Token) (done bool, err error) {
	if mach.isNil {
		tok.Type = TNull
		return true, nil
	}
	return mach.MarshalMachine.Step(driver, slab, tok)
}
