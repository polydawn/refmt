package again

/*
	So the cladistics are like this:

	  - TokenSource/TokenSink
	    - impl: json
		- impl: cbor
		- impl: [... other serial forms ...]
		- impl: var
		  - uses: VarVisitStepFunc -- but rarely shown; very powerful but very ugly primitive

	  - VarVisitStepFunc interface
	    - impl: Atlas
		- impl: [custom]
		- impl: atlas-gen

	Atlases are the defacto most reasonable way to define fields in a type to visit:
	in what order, and by what names.
*/

/*
	Two things are needed, and several things desirable:

	  - required: the ability to list fields (for serializing).
	  - required: the ability to map a name to a memory location (for both directions).
	  - optional: metadata about fields like "don't emit token if zero-valued" (or predicate? too complicated).

	It's possible to write a custom VarVisitStepFunc for decoding, but
	generally way more work than you want to sign up for.
	(Operating correctly for arbitrarily ordered tokens is a handful.)
*/

type typeA struct {
	Alpha string
	Beta  string
	Gamma typeB
	Delta int
}

type typeB struct {
	Msg string
}

var atlForTypeA = []AtlasField{
	{Name: "alpha", AddrFunc: func(v interface{}) interface{} { return &(v.(*typeA).Alpha) }},
	{Name: "beta",
		fieldRoute: []int{1}},
	{Name: "gamma",
		FieldName: FieldName{"Gamma", "Msg"}},
	{Name: "delta",
		FieldName: FieldName{"Delta"}},
}

type Atlas interface {
	Fields() []string
	Addr(fieldName string) interface{}
}

type AtlasField struct {
	Name string

	// *One* of the following:

	FieldName  FieldName                     // look up the fields by string name.
	fieldRoute []int                         // autoatlas fills these.
	AddrFunc   func(interface{}) interface{} // custom user function.

	// Behavioral options:

	OmitEmpty bool
}

type FieldName []string
