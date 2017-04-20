package obj

import (
	"fmt"

	. "github.com/polydawn/refmt/tok"
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

func (m UnmarshalMachineLiteral) Step(_ *UnmarshalDriver, tok *Token) (done bool, err error) {
	switch tgtp := m.target.(type) {
	case *bool:
		switch tok.Type {
		case TBool:
			*tgtp = tok.Bool
		default:
			goto errNoMatch
		}
	case *string:
		switch tok.Type {
		case TString:
			*tgtp = tok.Str
		default:
			goto errNoMatch
		}
	case *[]byte:
		switch tok.Type {
		case TBytes:
			*tgtp = tok.Bytes
		default:
			goto errNoMatch
		}
	case *int:
		switch tok.Type {
		case TInt:
			*tgtp = int(tok.Int) // TODO overflow check
		case TUint:
			*tgtp = int(tok.Uint) // TODO overflow check
		default:
			goto errNoMatch
		}
	case *int8:
		switch tok.Type {
		case TInt:
			*tgtp = int8(tok.Int) // TODO overflow check
		case TUint:
			*tgtp = int8(tok.Uint) // TODO overflow check
		default:
			goto errNoMatch
		}
	case *int16:
		switch tok.Type {
		case TInt:
			*tgtp = int16(tok.Int) // TODO overflow check
		case TUint:
			*tgtp = int16(tok.Uint) // TODO overflow check
		default:
			goto errNoMatch
		}
	case *int32:
		switch tok.Type {
		case TInt:
			*tgtp = int32(tok.Int) // TODO overflow check
		case TUint:
			*tgtp = int32(tok.Uint) // TODO overflow check
		default:
			goto errNoMatch
		}
	case *int64:
		switch tok.Type {
		case TInt:
			*tgtp = tok.Int
		case TUint:
			*tgtp = int64(tok.Uint) // TODO overflow check
		default:
			goto errNoMatch
		}
	case *uint:
		switch tok.Type {
		case TInt:
			*tgtp = uint(tok.Int) // TODO overflow check
		case TUint:
			*tgtp = uint(tok.Uint) // TODO overflow check
		default:
			goto errNoMatch
		}
	case *uint8:
		switch tok.Type {
		case TInt:
			*tgtp = uint8(tok.Int) // TODO overflow check
		case TUint:
			*tgtp = uint8(tok.Uint) // TODO overflow check
		default:
			goto errNoMatch
		}
	case *uint16:
		switch tok.Type {
		case TInt:
			*tgtp = uint16(tok.Int) // TODO overflow check
		case TUint:
			*tgtp = uint16(tok.Uint) // TODO overflow check
		default:
			goto errNoMatch
		}
	case *uint32:
		switch tok.Type {
		case TInt:
			*tgtp = uint32(tok.Int) // TODO overflow check
		case TUint:
			*tgtp = uint32(tok.Uint) // TODO overflow check
		default:
			goto errNoMatch
		}
	case *uint64:
		switch tok.Type {
		case TInt:
			*tgtp = uint64(tok.Int)
		case TUint:
			*tgtp = tok.Uint
		default:
			goto errNoMatch
		}
	case *uintptr:
		switch tok.Type {
		case TInt:
			*tgtp = uintptr(tok.Int)
		case TUint:
			*tgtp = uintptr(tok.Uint)
		default:
			goto errNoMatch
		}
	case *float32:
		switch tok.Type {
		case TFloat64:
			*tgtp = float32(tok.Float64) // TODO overflow check
		case TInt:
			*tgtp = float32(tok.Int) // TODO overflow check
		case TUint:
			*tgtp = float32(tok.Uint) // TODO overflow check
		default:
			goto errNoMatch
		}
	case *float64:
		switch tok.Type {
		case TFloat64:
			*tgtp = tok.Float64
		case TInt:
			*tgtp = float64(tok.Int) // TODO overflow check
		case TUint:
			*tgtp = float64(tok.Uint) // TODO overflow check
		default:
			goto errNoMatch
		}
	case *interface{}:
		// Why a wildcard?  Why not?
		// This machine for literals might end up picked only after the token
		//  comes in, and is aimed at a wildcard spot in memory.
		// In this case, we'll sanity check that the token is a valid type;
		//  then simply pass it through.
		switch tok.Type {
		case TString:
			*tgtp = tok.Str
		case TBytes:
			*tgtp = tok.Bytes
		case TBool:
			*tgtp = tok.Bool
		case TInt:
			*tgtp = tok.Int
		case TUint:
			*tgtp = tok.Uint
		case TFloat64:
			*tgtp = tok.Float64
		default:
			return true, fmt.Errorf("unexpected token %s, expected literal of any type", tok)
		}
		return true, nil
	default:
		panic(fmt.Errorf("cannot unmarshal into unhandled type %T", m.target))
	}
	return true, nil
errNoMatch:
	return true, fmt.Errorf("unexpected token of type %s, expected literal of type %T", tok.Type, m.target)
}
