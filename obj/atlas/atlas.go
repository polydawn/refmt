/*
	Atlas provides declarative descriptions of how to visit the fields of an object
	(as well as helpful functions to reflect on type declarations and generate
	default atlases that "do the right thing" for your types, following
	familiar conventions like struct tagging).
*/
package atlas

import "reflect"

type Atlas struct {
	Fields []Entry

	// A validation function which will be called for the whole value
	// after unmarshalling reached the end of the object.
	// If it returns an error, the entire unmarshal will error.
	// Not used in marshalling.
	ValidateFn func(v interface{}) error

	// If set, will be called after unmarshalling reached the end of the
	// object, and given a list of keys that appeared, in order of appearance.
	// This may be useful for knowing if a field was explicitly set to the zero
	// value vs simply unspecified, or for recording the order for later use
	// (e.g. so it can be serialized out again later in the same stable order).
	// Not used in marshalling.
	RecordFn func([]string)
}

type Entry struct {
	// The field name; will be emitted as token during marshal, and used for
	// lookup during unmarshal.  Required.
	Name string

	// *One* of the following:

	FieldName  FieldName                     // look up the fields by string name.
	fieldRoute fieldRoute                    // autoatlas fills these.
	AddrFunc   func(interface{}) interface{} // custom user function.

	// If true, marshalling will skip this field if its the zero value.
	// (If you need more complex behavior -- for example, a definition of
	// "empty" other than the type's zero value -- this is not for you.
	// Try using an AtlasFactory to make a custom field list dynamically.)
	OmitEmpty bool
}

type FieldName []string

type fieldRoute []int

func (ent *Entry) init() {
	// Validate reference options: only one may be used.
	// If it's a FieldName though, generate a fieldRoute for faster use.
	switch {
	case ent.fieldRoute != nil:
		if ent.FieldName != nil || ent.AddrFunc != nil {
			panic(ErrEntryInvalid{"if fieldRoute is used, no other field selectors may be specified"})
		}
	case ent.FieldName != nil:
		if ent.fieldRoute != nil || ent.AddrFunc != nil {
			panic(ErrEntryInvalid{"if FieldName is used, no other field selectors may be specified"})
		}
		// TODO transform `FieldName` to a `fieldRoute`
		// FIXME needs type info to reflect on, which isn't currently at hand
	case ent.AddrFunc != nil:
		if ent.fieldRoute != nil || ent.FieldName != nil {
			panic(ErrEntryInvalid{"if AddrFunc is used, no other field selectors may be specified"})
		}
	default:
		panic(ErrEntryInvalid{"one field selector must be specified"})
	}
}

/*
	Returns a reference to a field.
	(If the field is type `T`, the returned `interface{}` contains a `*T`.)
*/
func (ent Entry) Grab(v interface{}) interface{} {
	if ent.AddrFunc != nil {
		return ent.AddrFunc(v)
	}
	return ent.fieldRoute.TraverseToValue(reflect.ValueOf(v)).Interface()
}

func (fr fieldRoute) TraverseToValue(v reflect.Value) reflect.Value {
	for _, i := range fr {
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
