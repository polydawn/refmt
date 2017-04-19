package obj

import (
	"reflect"

	"github.com/polydawn/go-xlate/obj2/atlas"
)

/*
	A lovely mechanism to stash marshalMachine objects pre-allocated and avoid mallocs.
	Works together with the Atlas: the Atlas says what kind of machinery is needed;
	the marshalSlab "allocates" it and returns it upon your request.
*/
type marshalSlab struct {
	atlas atlas.Atlas
	rows  []marshalSlabRow
}

type marshalSlabRow struct {
	//	ptrDerefDelegateMarshalMachine
	//	marshalMachineLiteral
	//	marshalMachineMapWildcard
	//	marshalMachineSliceWildcard
	//	marshalMachineWildcard
	//	marshalMachineStructAtlas
}

/*
	Return a reference to a machine from the slab.
	*You must release() when done.*

	Errors -- including "no info in Atlas for this type" -- are expressed by
	returning a machine that is a constantly-erroring thunk.
*/
func (slab *marshalSlab) requisitionMachine(forThis reflect.Value) MarshalMachine {
	return nil // TODO
}
