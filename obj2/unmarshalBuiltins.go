package obj

import (
	"fmt"
	"reflect"

	. "github.com/polydawn/refmt/tok"
)

type ptrDerefDelegateUnmarshalMachine struct {
	UnmarshalMachine
	peelCount int

	isNil bool
}

func (mach *ptrDerefDelegateUnmarshalMachine) Reset(slab *unmarshalSlab, rv reflect.Value, rt reflect.Type) error {
	mach.isNil = false
	for i := 0; i < mach.peelCount; i++ {
		if rv.IsNil() {
			mach.isNil = true
			return nil
		}
		rv = rv.Elem()
	}
	return mach.UnmarshalMachine.Reset(slab, rv, rt) // REVIEW: this rt should be peeled by here.  do we... ignore the arg and cache it at mach conf time?
}
func (mach *ptrDerefDelegateUnmarshalMachine) Step(driver *UnmarshalDriver, slab *unmarshalSlab, tok *Token) (done bool, err error) {
	if mach.isNil {
		tok.Type = TNull
		return true, nil
	}
	return mach.UnmarshalMachine.Step(driver, slab, tok)
}

type unmarshalMachinePrimitive struct {
	kind reflect.Kind

	rv reflect.Value
}

func (mach *unmarshalMachinePrimitive) Reset(_ *unmarshalSlab, rv reflect.Value, _ reflect.Type) error {
	mach.rv = rv
	return nil
}
func (mach *unmarshalMachinePrimitive) Step(_ *UnmarshalDriver, _ *unmarshalSlab, tok *Token) (done bool, err error) {
	switch mach.kind {
	case reflect.Bool:
		switch tok.Type {
		case TBool:
			mach.rv.SetBool(tok.Bool)
			return true, nil
		default:
			return true, fmt.Errorf("cannot assign %s to bool field", tok)
		}
	case reflect.String:
		switch tok.Type {
		case TString:
			mach.rv.SetString(tok.Str)
			return true, nil
		default:
			return true, fmt.Errorf("cannot assign %s to string field", tok)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch tok.Type {
		case TInt:
			mach.rv.SetInt(tok.Int)
			return true, nil
		case TUint:
			mach.rv.SetInt(int64(tok.Uint)) // todo: overflow check
			return true, nil
		default:
			return true, fmt.Errorf("cannot assign %s to int field", tok)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch tok.Type {
		case TInt:
			if tok.Int >= 0 {
				mach.rv.SetUint(uint64(tok.Int))
				return true, nil
			}
			return true, fmt.Errorf("cannot assign %s to uint field", tok)
		case TUint:
			mach.rv.SetUint(tok.Uint)
			return true, nil
		default:
			return true, fmt.Errorf("cannot assign %s to uint field", tok)
		}
	case reflect.Float32, reflect.Float64:
		switch tok.Type {
		case TFloat64:
			mach.rv.SetFloat(tok.Float64)
			return true, nil
		case TInt:
			mach.rv.SetFloat(float64(tok.Int))
			return true, nil
		case TUint:
			mach.rv.SetFloat(float64(tok.Uint))
			return true, nil
		default:
			return true, fmt.Errorf("cannot assign %s to float field", tok)
		}
	case reflect.Slice: // implicitly bytes; no other slices are "primitve"
		switch tok.Type {
		case TBytes:
			mach.rv.SetBytes(tok.Bytes)
			return true, nil
		default:
			return true, fmt.Errorf("cannot assign %s to []byte field", tok)
		}
	default:
		panic("unhandled")
	}
}
