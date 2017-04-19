package obj

import (
	"fmt"
	"reflect"

	"github.com/polydawn/go-xlate/obj2/atlas"
	. "github.com/polydawn/go-xlate/tok"
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

	errThunkMarshalMachine
}

/*
	Return a reference to a machine from the slab.
	*You must release() when done.*

	Errors -- including "no info in Atlas for this type" -- are expressed by
	returning a machine that is a constantly-erroring thunk.
*/
func (slab *marshalSlab) requisitionMachine(forThis reflect.Type) MarshalMachine {
	off := len(slab.rows)
	slab.grow()
	row := &slab.rows[off]
	m := &row.errThunkMarshalMachine
	m.err = fmt.Errorf("no machine found")
	return m
}

func (s *marshalSlab) grow() {
	s.rows = append(s.rows, marshalSlabRow{})
}

func (s *marshalSlab) release() {
	s.rows = s.rows[0 : len(s.rows)-1]
}

type errThunkMarshalMachine struct {
	err error
}

func (m *errThunkMarshalMachine) Reset(_ *marshalSlab, _ reflect.Value, _ reflect.Type) error {
	return m.err
}
func (m *errThunkMarshalMachine) Step(d *MarshalDriver, s *marshalSlab, tok *Token) (done bool, err error) {
	return true, m.err
}
