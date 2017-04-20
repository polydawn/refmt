package obj

import (
	"reflect"

	"github.com/polydawn/refmt/obj/atlas"
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
