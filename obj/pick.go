package obj

import (
	"fmt"
	"reflect"
)

type Suite struct {
	// Map typeinfo to a factory function for new marshaller state machines.
	// A factory function is needed even though machines are resettable and reusable
	// because we still need more than one machine instance in the case of recursive structures.
	mappings map[reflect.Type]func() MarshalMachine
}

/*
	Folds another behavior dispatch into the suite.

	The `typeHint` parameter is an instance of the type you want dispatch
	to use this machine for.  A zero instance is fine.
	Thus, calls to this method usually resemble the following:

		suite.Add(YourType{}, &SomeMachineImpl{})
*/
func (s *Suite) Add(typeHint interface{}, machFactory func() MarshalMachine) *Suite {
	if s.mappings == nil {
		s.mappings = make(map[reflect.Type]func() MarshalMachine)
	}
	rt := reflect.TypeOf(typeHint)
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	s.mappings[rt] = machFactory
	return s
}

func (s *Suite) pickMarshalMachine(valp interface{}) MarshalMachine {
	mach := s.maybePickMarshalMachine(valp)
	if mach == nil {
		panic(ErrNoHandler{valp})
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

func (s *Suite) marshalMachineForType(val_rt reflect.Type) MarshalMachine {
	mach := s.maybeMarshalMachineForType(val_rt)
	if mach == nil {
		panic(fmt.Errorf("no machine available in suite for type %s", val_rt.Name()))
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
func (s *Suite) maybeMarshalMachineForType(val_rt reflect.Type) MarshalMachine {
	peelCount := 0
	for val_rt.Kind() == reflect.Ptr {
		val_rt = val_rt.Elem()
		peelCount++
	}
	mach := s._maybeMarshalMachineForType(val_rt)
	if mach == nil {
		return nil
	}
	if peelCount > 0 {
		return &ptrDerefDelegateMarshalMachine{mach, peelCount, false}
	}
	return mach
}

func (s *Suite) _maybeMarshalMachineForType(rt reflect.Type) MarshalMachine {
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
		return &MarshalMachineMapWildcard{}
	case reflect.Struct:
		if factory, ok := s.mappings[rt]; ok {
			return factory()
		}
		return nil
	case reflect.Interface:
		panic(ErrUnreachable{"TODO iface"})
	case reflect.Func:
		panic(ErrUnreachable{"TODO func"}) // hey, if we can find it in the suite
	case reflect.Ptr:
		panic(ErrUnreachable{"unreachable: ptrs must already be resolved"})
	default:
		panic(ErrUnreachable{}.Fmt("excursion %s", rt.Kind()))
	}
}
