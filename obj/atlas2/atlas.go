package atlas

import (
	"reflect"
)

type Atlas struct {
	// Map typeinfo to a static description of how that type should be handled.
	// (The internal machinery that will wield this information, and has memory of
	// progress as it does so, is configured using the AtlasEntry, but allocated separately.
	// The machinery is stateful and mutable; the AtlasEntry is not.)
	mappings map[reflect.Type]AtlasEntry
	// todo: others have used 'var rtid uintptr = reflect.ValueOf(rt).Pointer()' -- pointer of the value of the r.T info -- as an index
}

/*
	The AtlasEntry is a declarative roadmap of what we should do for
	marshal and unmarshal of a single object, keyed by type.

	There are a lot of paths your mappings might want to take:

	  - For a struct type, you may simply want to specify some alternate keys, or some to leave out, etc.
	  - For an interface type, you probably want to specify one of our interface muxing strategies
	     with a mapping between enumstr:typeinfo (and, what to do if we get a struct we don't recognize).
	  - For a string, int, or other primitive, you don't need to say anything: defaults will DTRT.
	  - For a typedef'd string, int, or other primitive, you *still* don't need to say anything: but,
	     if you want custom behavior (say, transform the string to an int at the last second, and back again),
		 you can specify transformer functions for that.
	  - For a struct type that you want to turn into a whole different kind (like a string): use
	     those same transform functions.  (You'll no longer need a FieldMap.)
	  - For the most esoteric needs, you can fall all the way back to providing a custom MarshalMachine
	     (but avoid that; it's a lot of work, and one of these other transform methods should suffice).
*/
type AtlasEntry struct {
	// The reflect info of the type this morphism is regarding.
	Type reflect.Type

	// A mapping of fileds in a struct to serial keys.
	// Only valid if `this.Type.Kind() == Struct`.
	StructMap StructMap

	MarshalTransform func(interface{}) interface{}
	// might as well store the extra type info and create it with reflect at suite setup time?

	UnmarshalTransform struct {
		TargetType reflect.Type
		Transform  func(interface{}) interface{}
	}

	// A validation function which will be called for the whole value
	// after unmarshalling reached the end of the object.
	// If it returns an error, the entire unmarshal will error.
	// Not used in marshalling.
	ValidateFn func(v interface{}) error
}

type StructMap struct {
	// An slice of descriptions of each field in the type.
	// Each entry specifies the name by which each field should be referenced
	// when serialized, and defines a way to get an address to the field.
	Fields []StructMapEntry
}

type StructMapEntry struct {
	// The field name; will be emitted as token during marshal, and used for
	// lookup during unmarshal.  Required.
	Name string

	// *One* of the following:

	ReflectRoute []int                         // reflection generates these.
	AddrFunc     func(interface{}) interface{} // custom user function.

	// If true, marshalling will skip this field if its the zero value.
	// (If you need more complex behavior -- for example, a definition of
	// "empty" other than the type's zero value -- this is not for you.
	// Try using a MarshalTransform to make a custom field list dynamically.)
	OmitEmpty bool
}
