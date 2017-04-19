package atlas

import "reflect"

type StructMap struct {
	// An slice of descriptions of each field in the type.
	// Each entry specifies the name by which each field should be referenced
	// when serialized, and defines a way to get an address to the field.
	Fields []StructMapEntry
}

type StructMapEntry struct {
	// The field name; will be emitted as token during marshal, and used for
	// lookup during unmarshal.  Required.
	name string

	// *One* of the following:

	reflectRoute reflectRoute                  // reflection generates these.
	addrFunc     func(interface{}) interface{} // custom user function.

	// If true, marshalling will skip this field if its the zero value.
	omitEmpty bool
}

type reflectRoute []int

func (rr reflectRoute) TraverseToValue(v reflect.Value) reflect.Value {
	for _, i := range rr {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return reflect.Value{}
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return v
}
