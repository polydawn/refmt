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
func (slab *marshalSlab) requisitionMachine(rt reflect.Type) MarshalMachine {
	// Acquire a row.
	off := len(slab.rows)
	slab.grow()
	row := &slab.rows[off]
	// Flip to rtid.
	rtid := reflect.ValueOf(rt).Pointer()

	// Check primitives first; cheapest (and unoverridable).
	switch rtid {
	case rtid_bool:
		panic("todo")
	case rtid_string:
		panic("todo")
	case rtid_bytes:
		panic("todo")
	case rtid_int:
		panic("todo")
	case rtid_int8:
		panic("todo")
	case rtid_int16:
		panic("todo")
	case rtid_int32:
		panic("todo")
	case rtid_int64:
		panic("todo")
	case rtid_uint:
		panic("todo")
	case rtid_uint8:
		panic("todo")
	case rtid_uint16:
		panic("todo")
	case rtid_uint32:
		panic("todo")
	case rtid_uint64:
		panic("todo")
	case rtid_uintptr:
		panic("todo")
	case rtid_float32:
		panic("todo")
	case rtid_float64:
		panic("todo")
	}

	// Consult atlas second.
	if entry, ok := slab.atlas.Get(rtid); ok {
		_ = entry
		panic("todo")
	}

	// If no specific behavior found, use default behavior based on kind.
	switch rt.Kind() {
	case reflect.Bool:
		panic("todo")
	case reflect.String:
		panic("todo")
	case reflect.Int:
		panic("todo")
	case reflect.Int8:
		panic("todo")
	case reflect.Int16:
		panic("todo")
	case reflect.Int32:
		panic("todo")
	case reflect.Int64:
		panic("todo")
	case reflect.Uint:
		panic("todo")
	case reflect.Uint8:
		panic("todo")
	case reflect.Uint16:
		panic("todo")
	case reflect.Uint32:
		panic("todo")
	case reflect.Uint64:
		panic("todo")
	case reflect.Uintptr:
		panic("todo")
	case reflect.Float32:
		panic("todo")
	case reflect.Float64:
		panic("todo")
	case reflect.Slice:
		// un-typedef'd byte slices were handled already, but a typedef'd one still gets gets treated like a special kind:
		if rt.Elem().Kind() == reflect.Uint8 {
			panic("todo")
		}
		panic("todo")
	case reflect.Array:
		panic("todo")
	case reflect.Map:
		panic("todo")
	case reflect.Struct:
		panic("todo")
	case reflect.Interface:
		panic("todo")
	case reflect.Func:
		panic(fmt.Errorf("functions cannot be marshalled!"))
	case reflect.Ptr:
		panic(fmt.Errorf("unreachable: ptrs must already be resolved"))
	default:
		panic(fmt.Errorf("excursion %s", rt.Kind()))
	}

	// If no joy yet, we're out of ideas: yield an error thunk.
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
