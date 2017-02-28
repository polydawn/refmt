package obj

import (
	"fmt"

	. "github.com/polydawn/go-xlate/tok"
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

func (m *MarshalMachineLiteral) Reset(_ *marshalSlab, target interface{}) error {
	m.target = target
	return nil
}

func (m MarshalMachineLiteral) Step(_ *MarshalDriver, _ *marshalSlab, tok *Token) (done bool, err error) {
	// Honestly, this entire set of paths does so little work we should think about inlining it
	// into the machine-picker (or earlier) entirely and never allocing or returning a machine.
	switch v2 := m.target.(type) {
	case *bool:
		tok.Type = TBool
		tok.Bool = *v2
	case *string:
		tok.Type = TString
		tok.Str = *v2
	case *[]byte:
		tok.Type = TBytes
		tok.Bytes = *v2
	case *int:
		tok.Type = TInt
		tok.Int = int64(*v2)
	case *int8:
		tok.Type = TInt
		tok.Int = int64(*v2)
	case *int16:
		tok.Type = TInt
		tok.Int = int64(*v2)
	case *int32:
		tok.Type = TInt
		tok.Int = int64(*v2)
	case *int64:
		tok.Type = TInt
		tok.Int = *v2
	case *uint:
		tok.Type = TUint
		tok.Uint = uint64(*v2)
	case *uint8:
		tok.Type = TUint
		tok.Uint = uint64(*v2)
	case *uint16:
		tok.Type = TUint
		tok.Uint = uint64(*v2)
	case *uint32:
		tok.Type = TUint
		tok.Uint = uint64(*v2)
	case *uint64:
		tok.Type = TUint
		tok.Uint = *v2
	case *uintptr:
		tok.Type = TUint
		tok.Uint = uint64(*v2)
	case *float32:
		tok.Type = TFloat64
		tok.Float64 = float64(*v2)
	case *float64:
		tok.Type = TFloat64
		tok.Float64 = *v2
	default:
		panic(fmt.Errorf("cannot marshal unhandled type %T", m.target))
	}
	return true, nil
}
