package obj

import (
	"fmt"
	"reflect"
)

type Suite struct {
	mappings map[reflect.Type]MarshalMachine
}

func (s *Suite) pickMarshalMachine(valp interface{}) MarshalMachine {
	mach := s.maybePickMarshalMachine(valp)
	if mach == nil {
		panic(fmt.Errorf("no machine available in suite for type %T", valp))
	}
	return mach
}

/*
	Picks an unmarshal machine, returning the custom impls for any
	common/primitive types, and advanced machines where structs get involved.

	The argument should be the address of the actual value of interest.

	Returns nil if there is no marshal machine in the suite for this type.
*/
func (s *Suite) maybePickMarshalMachine(valp interface{}) MarshalMachine {
	// TODO : we can use type switches to do some primitives efficiently here
	//  before we turn to the reflective path.
	val_rt := reflect.ValueOf(valp).Elem().Type()
	return s.maybeMarshalMachineForType(val_rt)
}

func (s *Suite) marshalMachineForType(rt reflect.Type) MarshalMachine {
	mach := s.maybeMarshalMachineForType(rt)
	if mach == nil {
		panic(fmt.Errorf("no machine available in suite for type %s", rt.Name()))
	}
	return mach
}

/*
	Like `pickMarshalMachine`, but requiring only the reflect type info.
	This is useable when you only have the type info available (rather than an instance);
	this comes up when for example looking up the machine to use for all values
	in a slice based on the slice type info.

	(Using an instance may be able to take faster, non-reflective paths for
	primitive values.)

	In contrast to the method that takes a `valp interface{}`, this type info
	is understood to already be dereferenced.

	Returns nil if there is no marshal machine in the suite for this type.
*/
func (s *Suite) maybeMarshalMachineForType(rt reflect.Type) MarshalMachine {
	switch rt.Kind() {
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
		// REVIEW: doing this once, fine.  but unbounded?  questionable.
		return s.maybeMarshalMachineForType(rt.Elem())
	case reflect.Struct:
		return s.mappings[rt]
	case reflect.Interface:
		panic("TODO iface")
	case reflect.Func:
		panic("TODO func") // hey, if we can find it in the suite
	default:
		panic(fmt.Errorf("excursion %s", rt.Kind()))
	}
}
