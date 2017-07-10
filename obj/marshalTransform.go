package obj

import (
	"reflect"

	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/polydawn/refmt/tok"
)

type marshalMachineTransform struct {
	trFunc   atlas.MarshalTransformFunc
	delegate MarshalMachine
}

func (mach *marshalMachineTransform) Reset(slab *marshalSlab, rv reflect.Value, _ reflect.Type) error {
	tr_rv, err := mach.trFunc(rv)
	if err != nil {
		return err
	}
	return mach.delegate.Reset(slab, tr_rv, tr_rv.Type())
}

func (mach *marshalMachineTransform) Step(driver *Marshaler, slab *marshalSlab, tok *Token) (done bool, err error) {
	return mach.delegate.Step(driver, slab, tok)
}
