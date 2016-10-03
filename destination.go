package xlate

/*
	Destination is the interface a MappingFunc pushes values into as it walks
	the MappingFunc's input object.

	Destination can be imagined having implementations like a json encoder,
	a cbor encoder, or even a *decoding* tool that's turning around and pushing
	the values back onto another structure in memory based on the field names.

	MappingFuncs given object inputs which are pushed into a Destination forms
	an object serialization pipeline;
	A tokenizer pushing chunks of an input stream into a Destination can form
	an object deserialization pipeline.
*/
type Destination interface {
	OpenMap()
	WriteMapKey(k string)
	CloseMap()

	OpenArray()
	CloseArray()

	// add write funcs for bare types?  almost certainly
	// after any of these, you're in a terminal state.

	WriteNull()
	WriteString(string)

	// also essential to have a way to say "skipme" -- this is distinct from writing a nil
}

// DESIGN: destinations are almost entirely leaves.
// Most of the functions write a single leaf, and further calls would be invalid.
// Only two modes are transitions to modes that accept multiple further calls.
// Can we do something to make that state machine better represented in calling code?
// Unclear.  Haven't thought of any ways to do so without a proliferation of types --
// a type for each of the states seems a net loss for a rapidly comprehensible API.
//
// todo: May want a helper mixin for destination implementations to track this state machine stuff
// and raise the errors coherently, though, if going this route.

// DESIGN: nesting Destinations to do schema verification should be a thing.
// For example, nesting Destinations to filter for subsets of values, while
// simultaenously punting to another Destination that parses a wider range of values.

/*
	Error raised when calling any of the output functions on `Destination` when
	the destination has already been fed a leaf state.

	(For example, it's invalid to `WriteString` twice in a row if serializing
	a single valid json object: you'd need to open an array first, then write
	a series of strings, then close the array.)
*/
type ErrAlreadyLeaf struct {
	After Breadcrumb
}

func (ErrAlreadyLeaf) Error() string { return "ErrAlreadyLeaf" }

/*
	Error raised when a Mapper invokes a MappingFunc to handle a value, and the
	Destination had not reached a valid terminal/leaf state when the MappingFunc returns.
*/
type ErrReturnedWithoutLeaf struct {
	During Breadcrumb
}

func (ErrReturnedWithoutLeaf) Error() string { return "ErrReturnedWithoutLeaf" }

/*
	An xpath-like object representing where we are in the object traversal.

	Mostly for use in error messages.
*/
type Breadcrumb []string
