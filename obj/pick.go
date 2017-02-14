package obj

import (
	"fmt"
	"reflect"

	"github.com/polydawn/go-xlate/obj/atlas"
)

type Suite struct {
	// Map typeinfo to a static description of how that type should be handled.
	// (The internal machinery that will wield this information, and has memory of
	// progress as it does so, is configured using the Morphism, but allocated separately.
	// The machinery is stateful and mutable; the Morphism is not.)
	mappings map[reflect.Type]Morphism
}

type Morphism struct {
	Atlas atlas.Atlas // the one kind of object with supported customization at the moment
	// REVIEW: those other funcs in atlas probably belong here, not there.
	// just generally clarify the lines for what should apply for, say, typedef'd ints.
}

/*
	Folds another behavior dispatch into the suite.

	The `typeHint` parameter is an instance of the type you want dispatch
	to use this machine for.  A zero instance is fine.
	Thus, calls to this method usually resemble the following:

		suite.Add(YourType{}, &SomeMachineImpl{})
*/
func (s *Suite) Add(typeHint interface{}, morphism Morphism) *Suite {
	if s.mappings == nil {
		s.mappings = make(map[reflect.Type]Morphism)
	}
	rt := reflect.TypeOf(typeHint)
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	s.mappings[rt] = morphism
	morphism.Atlas.Init()
	return s
}

type slab struct {
	suite *Suite
	rows  []slabRow
}

type slabRow struct {
	ptrDerefDelegateMarshalMachine
	MarshalMachineLiteral
	MarshalMachineMapWildcard
	MarshalMachineSliceWildcard
	MarshalMachineStructAtlas
}

func (s *slab) mustPickMarshalMachine(valp interface{}) MarshalMachine {
	mach := s.pickMarshalMachine(valp)
	if mach == nil {
		panic(ErrNoHandler{valp})
	}
	return mach
}

func (s *slab) mustPickMarshalMachineByType(val_rt reflect.Type) MarshalMachine {
	mach := s.pickMarshalMachineByType(val_rt)
	if mach == nil {
		panic(fmt.Errorf("no machine available in suite for type %s", val_rt.Name()))
	}
	return mach
}

/*
	Picks an unmarshal machine, returning the custom impls for any
	common/primitive types, and advanced machines where structs get involved.

	The argument should be the address of the actual value of interest.

	Returns nil if there is no marshal machine in the suite for this type.
*/
func (s *slab) pickMarshalMachine(valp interface{}) MarshalMachine {
	// TODO : we can use type switches to do some primitives efficiently here
	//  before we turn to the reflective path.
	val_rt := reflect.ValueOf(valp).Elem().Type()
	return s.pickMarshalMachineByType(val_rt)
}

/*
	Like `mustPickMarshalMachine`, but requiring only the reflect type info.
	This is useable when you only have the type info available (rather than an instance);
	this comes up when for example looking up the machine to use for all values
	in a slice based on the slice type info.

	(Using an instance may be able to take faster, non-reflective paths for
	primitive values.)

	In contrast to the method that takes a `valp interface{}`, this type info
	is understood to already be dereferenced.

	Returns nil if there is no marshal machine in the suite for this type.
*/
func (s *slab) pickMarshalMachineByType(val_rt reflect.Type) MarshalMachine {
	peelCount := 0
	for val_rt.Kind() == reflect.Ptr {
		val_rt = val_rt.Elem()
		peelCount++
	}
	mach := s._pickMarshalMachineByType(val_rt)
	if mach == nil {
		return nil
	}
	if peelCount > 0 {
		return &ptrDerefDelegateMarshalMachine{mach, peelCount, false}
	}
	return mach
}

func (s *slab) grow() {
	s.rows = append(s.rows, slabRow{})
}

func (s *slab) _pickMarshalMachineByType(rt reflect.Type) MarshalMachine {
	s.grow()
	off := len(s.rows) - 1
	switch rt.Kind() {
	case reflect.Bool:
		return &s.rows[off].MarshalMachineLiteral
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &s.rows[off].MarshalMachineLiteral
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &s.rows[off].MarshalMachineLiteral
	case reflect.Float32, reflect.Float64:
		return &s.rows[off].MarshalMachineLiteral
	case reflect.String:
		return &s.rows[off].MarshalMachineLiteral
	case reflect.Slice:
		// TODO also bytes should get a special path
		return &s.rows[off].MarshalMachineSliceWildcard
	case reflect.Array:
		return &s.rows[off].MarshalMachineSliceWildcard
	case reflect.Map:
		return &s.rows[off].MarshalMachineMapWildcard
	case reflect.Struct:
		morphism, ok := s.suite.mappings[rt]
		if !ok {
			return nil
		}
		s.rows[off].MarshalMachineStructAtlas.atlas = morphism.Atlas
		return &s.rows[off].MarshalMachineStructAtlas
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
