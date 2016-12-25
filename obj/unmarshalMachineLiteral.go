package obj

import (
	"fmt"

	. "github.com/polydawn/go-xlate/tok"
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

func (m UnmarshalMachineLiteral) Step(_ *UnmarshalDriver, tokp *Token) (done bool, err error) {
	switch tgtp := m.target.(type) {
	case *bool:
		tokc, ok := (*tokp).(*bool)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *string:
		tokc, ok := (*tokp).(*string)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *[]byte:
		tokc, ok := (*tokp).(*[]byte)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *int:
		tokc, ok := (*tokp).(*int)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *int8:
		tokc, ok := (*tokp).(*int8)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *int16:
		tokc, ok := (*tokp).(*int16)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *int32:
		tokc, ok := (*tokp).(*int32)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *int64:
		tokc, ok := (*tokp).(*int64)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *uint:
		tokc, ok := (*tokp).(*uint)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *uint8:
		tokc, ok := (*tokp).(*uint8)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *uint16:
		tokc, ok := (*tokp).(*uint16)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *uint32:
		tokc, ok := (*tokp).(*uint32)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *uint64:
		tokc, ok := (*tokp).(*uint64)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *uintptr:
		tokc, ok := (*tokp).(*uintptr)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *float32:
		tokc, ok := (*tokp).(*float32)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *float64:
		tokc, ok := (*tokp).(*float64)
		if !ok {
			goto err
		}
		*tgtp = *tokc
	case *interface{}:
		// Why a wildcard?  Why not?
		// This machine for literals might end up picked only after the token
		//  comes in, and is aimed at a wildcard spot in memory.
		// In this case, we'll sanity check that the token is a valid type;
		//  then simply pass it through.
		if !IsValidToken(*tokp) {
			return true, fmt.Errorf("invalid token of type %T", *tokp)
		}
		switch tokc := (*tokp).(type) {
		case *bool:
			*tgtp = *tokc
		case *string:
			*tgtp = *tokc
		case *[]byte:
			*tgtp = *tokc
		case *int:
			*tgtp = *tokc
		case *int8:
			*tgtp = *tokc
		case *int16:
			*tgtp = *tokc
		case *int32:
			*tgtp = *tokc
		case *int64:
			*tgtp = *tokc
		case *uint:
			*tgtp = *tokc
		case *uint8:
			*tgtp = *tokc
		case *uint16:
			*tgtp = *tokc
		case *uint32:
			*tgtp = *tokc
		case *uint64:
			*tgtp = *tokc
		case *uintptr:
			*tgtp = *tokc
		case *float32:
			*tgtp = *tokc
		case *float64:
			*tgtp = *tokc
		}
	default:
		panic(fmt.Errorf("cannot unmarshal into unhandled type %T", m.target))
	}
	return true, nil
err:
	return true, fmt.Errorf("unexpected token of type %T, expected literal of type %T", *tokp, m.target)
}
