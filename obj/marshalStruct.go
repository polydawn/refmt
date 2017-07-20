package obj

import (
	"fmt"
	"reflect"

	"github.com/polydawn/refmt/obj/atlas"
	. "github.com/polydawn/refmt/tok"
)

type marshalMachineStructAtlas struct {
	cfg *atlas.StructMap // set on initialization

	rv    reflect.Value
	index int  // Progress marker
	value bool // Progress marker
}

func (mach *marshalMachineStructAtlas) Reset(_ *marshalSlab, rv reflect.Value, _ reflect.Type) error {
	mach.rv = rv
	mach.index = -1
	mach.value = false
	return nil
}

func (mach *marshalMachineStructAtlas) Step(driver *Marshaller, slab *marshalSlab, tok *Token) (done bool, err error) {
	//fmt.Printf("--step on %#v: i=%d/%d v=%v\n", mach.rv, mach.index, len(mach.cfg.Fields), mach.value)
	nEntries := len(mach.cfg.Fields)
	if mach.index < 0 {
		tok.Type = TMapOpen
		tok.Length = nEntries
		mach.index++
		return false, nil
	}
	if mach.index == nEntries {
		tok.Type = TMapClose
		mach.index++
		slab.release()
		return true, nil
	}
	if mach.index > nEntries {
		return true, fmt.Errorf("invalid state: entire struct (%d fields) already consumed", nEntries)
	}

	if mach.value {
		fieldEntry := mach.cfg.Fields[mach.index]
		child_rv := fieldEntry.ReflectRoute.TraverseToValue(mach.rv)
		mach.index++
		mach.value = false
		return false, driver.Recurse(
			tok,
			child_rv,
			fieldEntry.Type,
			slab.requisitionMachine(fieldEntry.Type),
		)
	}
	tok.Type = TString
	tok.Str = mach.cfg.Fields[mach.index].SerialName
	mach.value = true
	if mach.index > 0 {
		slab.release()
	}
	return false, nil
}
