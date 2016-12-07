package obj

import (
	"reflect"
)

/*
	Picks an unmarshal machine, returning the custom impls for any
	common/primitive types, and advanced machines where structs get involved.
*/
func pickMarshalMachine(valp interface{}) MarshalMachine {
	// future: if we wanted to support a function on a type that indicates custom behavior,
	//  this would be the place to do that check, via `val_rt.Implements(markerType)`.

	val_rt := reflect.TypeOf(*(valp).(*interface{}))
	switch val_rt.Kind() {
	case reflect.Bool:
		return &MarshalMachineLiteral{}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &MarshalMachineLiteral{}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &MarshalMachineLiteral{}
	case reflect.Float32, reflect.Float64:
		return &MarshalMachineLiteral{}
	case reflect.String:
		return &MarshalMachineLiteral{}
	case reflect.Slice:
		// TODO also bytes should get a special path
		return &MarshalMachineSliceWildcard{}
	case reflect.Array:
		return &MarshalMachineSliceWildcard{}
	case reflect.Map:
		// Consider (for efficiency in happy paths):
		//		switch v2 := v.(type) {
		//		case map[string]interface{}:
		//			_ = v2
		//			return nil // TODO special
		//		default:
		//			return &MarshalMachineMapWildcard{}
		//		}
		// but, it's not clear how we'd cache this.
		// possibly we should attach this whole method to the Driver,
		//  so it can have state for cache.
		return &MarshalMachineMapWildcard{}
	case reflect.Ptr:
		panic("TODO")
	case reflect.Struct:
		panic("TODO")
	case reflect.Interface:
		panic("TODO")
	default:
		panic("TODO")
	}
}

type Suite struct {
	mappings map[reflect.Type]MarshalMachine // i want to know if i can reuse them!  i hope i can.  why wouldn't i?
}
