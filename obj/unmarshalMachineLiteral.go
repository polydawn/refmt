package obj

import (
	"fmt"

	"github.com/polydawn/go-xlate/tok"
)

/*
	An UnmarshalMachine which unpacks a single literal of some primitive type.
	It supports `int`, `string`, `bool`, and so on.

	The `target` slot must be a address of such a primitive type, or,
	the address of an `interface{}` slot, which will be filled with whatever
	type of token primitive comes along.
*/
type UnmarshalMachineLiteral struct {
	target interface{}
}

func (m UnmarshalMachineLiteral) Step(_ *UnmarshalDriver, tok *tok.Token) (done bool, err error) {
	var ok bool
	switch v2 := m.target.(type) {
	case *string:
		*v2, ok = (*tok).(string)
	case *[]byte:
		panic("TODO")
	case *int:
		*v2, ok = (*tok).(int)
	case *int8, *int16, *int32, *int64:
		panic("TODO")
	case *uint, *uint8, *uint16, *uint32, *uint64:
		panic("TODO")
	case *interface{}:
		// TODO may want to whitelist tok types here are indeed prim literals as a san check
		*v2 = *tok
		ok = true
	default:
		panic(fmt.Errorf("cannot unmarshal into unhandled type %T", m.target))
	}
	if ok {
		return true, nil
	}
	return true, fmt.Errorf("unexpected token of type %T, expected literal of type %T", *tok, m.target)
}
