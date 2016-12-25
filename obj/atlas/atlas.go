/*
	Atlas provides declarative descriptions of how to visit the fields of an object
	(as well as helpful functions to reflect on type declarations and generate
	default atlases that "do the right thing" for your types, following
	familiar conventions like struct tagging).
*/
package atlas

import (
	"fmt"
	"reflect"
)

type Atlas struct {
	// The type this atlas describes.
	Type reflect.Type

	// An slice of descriptions of each field in the type.
	// Each entry specifies the name by which each field should be referenced
	// when serialized, and defines a way to get an address to the field.
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
	FieldRoute FieldRoute                    // autoatlas fills these.
	AddrFunc   func(interface{}) interface{} // custom user function.

	// Optionally, specify exactly what should handle the field value:
	// TODO this is one of {Atlas, func()(Atlas), or TokenSourceMachine|TokenSinkMachine}
	//  the latter is certainly the most correct, but also pretty wicked to export publicly

	// If true, marshalling will skip this field if its the zero value.
	// (If you need more complex behavior -- for example, a definition of
	// "empty" other than the type's zero value -- this is not for you.
	// Try using an AtlasFactory to make a custom field list dynamically.)
	OmitEmpty bool
}

type FieldName []string

type FieldRoute []int

func (atl *Atlas) Init() {
	for i, _ := range atl.Fields {
		atl.Fields[i].init(atl.Type)
	}
}

func (ent *Entry) init(rt reflect.Type) {
	// Validate reference options: only one may be used.
	// If it's a FieldName though, generate a FieldRoute for faster use.
	switch {
	case ent.FieldRoute != nil:
		if ent.FieldName != nil || ent.AddrFunc != nil {
			panic(ErrEntryInvalid{"if FieldRoute is used, no other field selectors may be specified"})
		}
		if len(ent.FieldRoute) == 0 {
			panic(ErrEntryInvalid{"FieldRoute cannot be length zero (would be inf recursion)"})
		}
	case ent.FieldName != nil:
		if ent.FieldRoute != nil || ent.AddrFunc != nil {
			panic(ErrEntryInvalid{"if FieldName is used, no other field selectors may be specified"})
		}
		if len(ent.FieldName) == 0 {
			panic(ErrEntryInvalid{"FieldName cannot be length zero (would be inf recursion)"})
		}
		// transform `FieldName` to a `FieldRoute`.
		for _, fn := range ent.FieldName {
			f, ok := rt.FieldByName(fn)
			if !ok {
				panic(ErrStructureMismatch{rt.Name(), "does not have field named " + fn})
			}
			ent.FieldRoute = append(ent.FieldRoute, f.Index...)
		}
	case ent.AddrFunc != nil:
		if ent.FieldRoute != nil || ent.FieldName != nil {
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
	if ent.FieldRoute == nil {
		panic(fmt.Errorf("atlas.Entry not initialized"))
	}
	// Jump through the defacto first pointer.
	// We're about to check if our traversal will be able to return an addressable field;
	//  but we don't care if the *pointer* itself we have here is addressable.
	v_rv := reflect.ValueOf(v).Elem()
	if !v_rv.CanAddr() {
		panic(fmt.Errorf("values for atlas traversal must be addressable"))
	}
	field_rv := ent.FieldRoute.TraverseToValue(v_rv)
	return field_rv.Addr().Interface()
}

func (fr FieldRoute) TraverseToValue(v reflect.Value) reflect.Value {
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
