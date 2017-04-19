package atlas

import (
	"reflect"
	"strings"
)

/*
	Add a field to the mapping based on its name.

	Given a struct:

		struct{
			X int
			Y struct{ Z int }
		}

	`AddField("X", {"x", ...}) will cause that field to be serialized as key "x";
	`AddField("Y.Z", {"z", ...})` will cause that *nested* field to be serialized
	as key "z" in the same object (e.g. "x" and "z" will be siblings).

	Returns mutated StructMap for convenient call chaining.

	If the fieldName string doesn't map onto the structure type info,
	a panic will be raised.
*/
func (sm *StructMap) AddField(fieldName string, mapping StructMapEntry) *StructMap {
	fieldNameSplit := strings.Split(fieldName, ".")
	// FIXME sigh need the rt obj again... it's all the way up in the AtlasEntry.
	// Can we string together enough builder horseshit to get it down here?
	// Coincidentally, I'd like to be able to define my level of panickyness
	// a line lower at the AtlasBuilder scale.
	rr, err := fieldNameToReflectRoute(nil, fieldNameSplit)
	if err != nil {
		panic(err)
	}
	mapping.reflectRoute = rr
	sm.Fields = append(sm.Fields, mapping)
	return sm
}

func fieldNameToReflectRoute(rt reflect.Type, fieldNameSplit []string) (rr reflectRoute, err error) {
	for _, fn := range fieldNameSplit {
		rf, ok := rt.FieldByName(fn)
		if !ok {
			return nil, ErrStructureMismatch{rt.Name(), "does not have field named " + fn}
		}
		rr = append(rr, rf.Index...)
	}
	return rr, nil
}
