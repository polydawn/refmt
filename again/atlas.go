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

func altAtl(v interface{}) map[string]interface{} {
	x := v.(*typeA)
	// the problem here is... you see how this is:
	//  - not easy to mix with the others
	//  - not significantly shorter than much at all
	// the plus sides are small:
	//  - just one func call overhead for the whole struct (major nice)
	//  - just one cast line
	//  - it *is* compile time checked (but so are individual addrfuncs)
	// this also turns out to be orthangonal to anything like derived fields;
	//  those still require before/after funcs with someone holding the temp vars.
	//
	// we could make this an option though, for merge-overriding the rest.
	// would substantially complicate things though.
	return map[string]interface{}{
		"alpha": &(x.Alpha),
		"beta":  &(x.Beta),
		"gamma": &(x.Gamma.Msg),
		"delta": &(x.Delta),
	}
}

// ...
var _ = struct {
	Name      string
	TypeThunk interface{}
	FieldRef  interface{}
}{
	"huh",
	typeA{},
	typeA{}.Gamma.Msg, // no.
	// no, we don't have the technology for this.
	// you'd have to get an addressable instance of the type first,
	// then the path to its field.  THEN we could check the addrs for sameness.
	// but that's:
	//  - not a one-liner, not even close
	//  - pretty heavy magic that's likely to alarm a gopher
	//  - not well defined for anything with pointer hops in the midst of it
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
