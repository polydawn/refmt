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

func (m *MarshalMachineLiteral) Reset(_ *Suite, target interface{}) error {
	m.target = target
	return nil
}

func (m MarshalMachineLiteral) Step(_ *MarshalDriver, _ *Suite, tok *tok.Token) (done bool, err error) {
	// Honestly, this entire set of paths does so little work we should think about inlining it
	// into the machine-picker (or earlier) entirely and never allocing or returning a machine.
	switch v2 := m.target.(type) {
	case *bool,
		*string,
		*[]byte,
		*int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64, *uintptr,
		*float32, *float64:
		*tok = &v2
		return true, nil
	default:
		panic(fmt.Errorf("cannot marshal unhandled type %T", m.target))
	}
}
