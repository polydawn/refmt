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

	// --------------------------------------------------------
	// The big escape valves: wanna map to some other kind completely?
	// --------------------------------------------------------

	// Transforms the value we reached by walking (the 'live' value -- which
	// must be of `this.Type`) into another value (the 'serialable' value --
	// which will be of `this.MarshalTransformTargetType`).
	//
	// The target type may be anything, even of a completely different Kind!
	//
	// This transform func runs first, then the resulting value is
	// serialized (by running through the path through Atlas again, so
	// chaining of transform funcs is supported, though not recommended).
	MarshalTransformFunc MarshalTransformFunc
	// The type of value we expect after using the MarshalTransformFunc.
	//
	// The match between transform func and target type should be checked
	// during construction of this AtlasEntry.
	MarshalTransformTargetType reflect.Type

	// Expects a different type (the 'serialable' value -- which will be of
	// 'this.UnmarshalTransformTargetType') than the value we reached by
	// walking (the 'live' value -- which must be of `this.Type`).
	//
	// The target type may be anything, even of a completely different Kind!
	//
	// The unmarshal of that target type will be run first, then the
	// resulting value is fed through this function to produce the real value,
	// which is then placed correctly into bigger mid-unmarshal object tree.
	//
	// For non-primitives, unmarshal of the target type will always target
	// an empty pointer or empty slice, roughly as per if it was
	// operating on a value produced by `TargetType.New()`.
	UnmarshalTransformFunc UnmarshalTransformFunc
	// The type of value we will manufacture an instance of and unmarshal
	// into, then when done provide to the UnmarshalTransformFunc.
	//
	// The match between transform func and target type should be checked
	// during construction of this AtlasEntry.
	UnmarshalTransformTargetType reflect.Type

	// --------------------------------------------------------
	// Standard options for how to map (varies by Kind)
	// --------------------------------------------------------

	// A mapping of fields in a struct to serial keys.
	// Only valid if `this.Type.Kind() == Struct`.
	StructMap StructMap

	// FUTURE: enum-ish primitives, multiplexers for interfaces,
	//  lots of such things will belong here.

	// --------------------------------------------------------
	// Hooks, validate helpers
	// --------------------------------------------------------

	// A validation function which will be called for the whole value
	// after unmarshalling reached the end of the object.
	// If it returns an error, the entire unmarshal will error.
	//
	// Not used in marshalling.
	// Not reachable if an UnmarshalTransform is set.
	ValidateFn func(v interface{}) error
}
