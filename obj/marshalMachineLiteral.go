package obj

import (
	"fmt"

	"github.com/polydawn/go-xlate/tok"
)

/*
	A MarshalMachine which yields a single token of some primitive type.
	It supports `int` (all bit widths, and unsigneds), `float`, `bool`,
	`string`, and `[]byte`.

	The `target` slot must be a address of such a primitive type slot.
	Resetting the machine with an an address of other types is incorrect
	usage and will yield a panic during the step func.
*/
type MarshalMachineLiteral struct {
	target interface{}
}

func (m *MarshalMachineLiteral) Reset(_ *slab, target interface{}) error {
	m.target = target
	return nil
}

func (m MarshalMachineLiteral) Step(_ *MarshalDriver, _ *slab, tok *tok.Token) (done bool, err error) {
	// Honestly, this entire set of paths does so little work we should think about inlining it
	// into the machine-picker (or earlier) entirely and never allocing or returning a machine.
	switch v2 := m.target.(type) {
	case *bool:
		*tok = v2
	case *string:
		*tok = v2
	case *[]byte:
		*tok = v2
	case *int8:
		*tok = v2
	case *int16:
		*tok = v2
	case *int32:
		*tok = v2
	case *int64:
		*tok = v2
	case *uint:
		*tok = v2
	case *uint8:
		*tok = v2
	case *uint16:
		*tok = v2
	case *uint32:
		*tok = v2
	case *uint64:
		*tok = v2
	case *uintptr:
		*tok = v2
	case *float32:
		*tok = v2
	case *float64:
		*tok = v2
	case *int:
		*tok = v2
	default:
		panic(fmt.Errorf("cannot marshal unhandled type %T", m.target))
	}
	return true, nil
}
