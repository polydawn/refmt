package again

import (
	"fmt"
	"io"
	"reflect"
)

func jsonToJson(r io.Reader, w io.Writer) {
	turnBothCranks(
		NewJsonDecoder(r),
		NewJsonEncoder(w),
	)
}

func turnBothCranks(tokenSrc TokenSrc, tokenSink TokenSink) error {
	var tok Token
	var srcDone, sinkDone bool
	var err error
	for {
		srcDone, err = tokenSrc.Step(&tok)
		if err != nil {
			return err
		}
		sinkDone, err = tokenSink.Step(&tok)
		if err != nil {
			return err
		}
		if srcDone {
			if sinkDone {
				return nil
			}
			return fmt.Errorf("src at end of item but sink expects more")
		}
	}
}

/*
	Fill with address of primitive (or []byte), or the magic const tokens
	for beginning and ending of maps and arrays.

	Decoder implementations are encouraged to use `util.DecodeBag` to contain
	primitives during decode, then return the address of the relevant
	primitive field from the `DecodeBag` as a `Token`.  This avoids repeated
	pointer allocations.
*/
type Token interface{}

const (
	Token_MapOpen  = '{'
	Token_MapClose = '}'
	Token_ArrOpen  = '['
	Token_ArrClose = ']'
)

type TokenSrc interface {
	Step(fillme *Token) (done bool, err error)
	Reset()
}

type TokenSink interface {
	Step(consume *Token) (done bool, err error)
	Reset()
}

//
// Constructors
//

func NewJsonDecoder(r io.Reader /* optional *JsonSchemaNotes */) TokenSrc  { return nil }
func NewJsonEncoder(w io.Writer /* optional *JsonSchemaNotes */) TokenSink { return nil }

func NewVarTokenizer(v interface{} /* TODO visitmagicks */) TokenSrc { return nil }
func NewVarReceiver(v interface{} /* TODO visitmagicks */) TokenSink { return nil }

type varReceiver struct {
	step func(*Token) (done bool, err error)
	done bool
	err  error
}

// used at initialization to figure out the first step given the type of var
//
func (vr *varReceiver) stepFor(v interface{}) func(*Token) (done bool, err error) {
	switch v.(type) {
	// For total wildcards:
	//  Return a machine that will pick between a literal or `map[string]interface{}`
	//  or `[]interface{}` based on the next token.
	case *interface{}:
		return wildcardStep(v)
	// For single literals:
	//  we have a single machine that handles all these.
	case *string, *[]byte,
		*int, *int8, *int16, *int32, *int64,
		*uint, *uint8, *uint16, *uint32, *uint64:
		dec := &literalDecoderMachine{}
		dec.Reset(v)
		return dec.Step
	// Anything that has real type info:
	//  ... Plaaaay ball!
	default:
		// TODO mustAddressable check goes here.
		if reflect.TypeOf(v).Kind() == reflect.Interface {
			// special path because we can recycle the decoder machines, if they implement resettable.
		}
		// any other concrete type or particular interface:
		//  must have its own visit func defined.
		//  we don't know if it expects to be a map, lit, arr, etc until it takes over.
		//  (the rest of our functions here are the exception: they're half inlined here -- TODO maybe don't be like that; this lookup only makes sense for top level wtf-is-this'es)
		panic("TODO mappersuite lookup")
	}
}

func wildcardStep(target interface{}) func(*Token) (bool, error) {
	return func(tok *Token) (done bool, err error) {
		// If it's a special state, start an object.
		//  (Or, blow up if its a special state that's silly).
		switch *tok {
		case Token_MapOpen:
			// Fill in our wildcard ref with a blank map,
			//  and make a new machine for it; hand off everything.
			target = make(map[string]interface{})
			dec := &wildcardMapDecoderMachine{}
			dec.Reset(target)
			return dec.Step(tok)
		case Token_ArrOpen:
			// TODO same as maps, but with a machine for arrays
			panic("NYI")
		case Token_MapClose:
			return true, fmt.Errorf("unexpected mapClose; expected start of value")
		case Token_ArrClose:
			return true, fmt.Errorf("unexpected arrClose; expected start of value")
		default:
			// If it wasn't the start of composite, shell out to the machine for literals.
			dec := &literalDecoderMachine{}
			dec.Reset(target)
			return dec.Step(tok)
		}
	}
}

type wildcardMapDecoderMachine struct {
	target map[string]interface{}
	step   func(*Token) (done bool, err error)
	key    string // The key consumed by the prev `step_AcceptKey`.
}

func (dm *wildcardMapDecoderMachine) Reset(target interface{}) {
	dm.target = target.(map[string]interface{})
	dm.step = dm.step_Initial
	dm.key = ""
}

func (dm *wildcardMapDecoderMachine) Step(tok *Token) (done bool, err error) {
	return dm.step(tok)
}

func (dm *wildcardMapDecoderMachine) step_Initial(tok *Token) (done bool, err error) {
	// If it's a special state, start an object.
	//  (Or, blow up if its a special state that's silly).
	switch *tok {
	case Token_MapOpen:
		// Great.  Consumed.
		dm.step = dm.step_AcceptKey
		return false, nil
	case Token_ArrOpen:
		return true, fmt.Errorf("unexpected arrOpen; expected start of map")
	case Token_MapClose:
		return true, fmt.Errorf("unexpected mapClose; expected start of map")
	case Token_ArrClose:
		return true, fmt.Errorf("unexpected arrClose; expected start of map")
	default:
		panic(fmt.Errorf("unexpected literal of type %T; expected start of map", *tok))
	}
}
func (dm *wildcardMapDecoderMachine) step_AcceptKey(tok *Token) (done bool, err error) {
	switch *tok {
	case Token_MapOpen:
		return true, fmt.Errorf("unexpected mapOpen; expected map key")
	case Token_ArrOpen:
		return true, fmt.Errorf("unexpected arrOpen; expected map key")
	case Token_MapClose:
		// no special checks for ends of wildcard map; no such thing as incomplete.
		return true, nil
	case Token_ArrClose:
		return true, fmt.Errorf("unexpected arrClose; expected map key")
	}
	switch k := (*tok).(type) {
	case *string:
		if err = dm.mustAcceptKey(*k); err != nil {
			return true, err
		}
		dm.key = *k
		dm.step = dm.step_AcceptValue
		return false, nil
	default:
		panic(fmt.Errorf("unexpected literal of type %T; expected start of struct", *tok))
	}
}
func (dm *wildcardMapDecoderMachine) mustAcceptKey(k string) error {
	if _, exists := dm.target[k]; exists {
		return fmt.Errorf("repeated key %q", k)
	}
	return nil
}

func (dm *wildcardMapDecoderMachine) step_AcceptValue(tok *Token) (done bool, err error) {
	/*
		driver.Fill(
			tok, // still meant for next person and the real step is to come; we just had to figure out types, here.
			dm.Addr(dm.key),
			dm.step_postValue(), // driver returns to us after the value is done by calling this.
			    // may actually be that we stash that stepfunc, and give driver more general self pointer and Resume func in interface.
		)
	*/
	return false, nil // TODO
}

type literalDecoderMachine struct {
	target interface{}
}

func (dm *literalDecoderMachine) Reset(target interface{}) {
	dm.target = target
}

func (dm *literalDecoderMachine) Step(tok *Token) (done bool, err error) {
	var ok bool
	switch v2 := dm.target.(type) {
	case *string:
		*v2, ok = (*tok).(string)
	case *[]byte:
		panic("TODO")
	case *int, *int8, *int16, *int32, *int64:
		panic("TODO")
	case *uint, *uint8, *uint16, *uint32, *uint64:
		panic("TODO")
	default:
		panic(fmt.Errorf("cannot unmarshall into unhandled type %T", dm.target))
	}
	if ok {
		return true, nil
	}
	return true, fmt.Errorf("unexpected token of type %T, expected literal of type %T", *tok, dm.target)
}

/*
	Suppose we have the following var to unmarshal into:

		var thingy SomeType

	Where SomeType is defined as:

		type SomeType struct {
			AnInt int
			Something interface{}
		}

	The flow of a VarReciever working on this will be something like the following:

		- Begin handling a var of type `SomeType`.
		- Look up the hander for that type info.
		- The handler is accepts the val ref, and returns a step function.
		- The step function is called with the token.
		- [Much work ensues.]
		- If the step function returns done, we return entirely;
		  otherwise we hang onto the next stepFunc, and return.

	The flow of the specific handler for SomeType will look like this:

		- Expect a MapOpen token.
		- Expect a MapKey token.  Return a step func expecting that matching value.
		  - When called with the next token, this step func grabs the ref
		    of the struct field matching the name we were primed with...
		  - And calls dispatch on the whole thing.
		  - (Generally this func looks like it needs {fillingName string, rest},
		    so it can tell what value grab the ref to fill, and decide whether
			to return "expect all done" step.)
		- At any point, it may receive MapClose, which will jump to a check
		  that all fields are either noted as filled (requires sidebar) or
		  are tagged as omitEmpty.
*/

// Returns an atlas so we can use this to build the contin-passing machine without bothering you.
func HandleMe(vreal interface{}) (
	vmediate interface{},
	atl *Atlas,
	after func(), /* closure, already has vreal and vmediate refs */
) {
	return nil, nil, nil
}

type Atlas struct{}

type atlasDecoderMachine struct {
	val      reflect.Value // We're filling this.
	atl      *Atlas        // Our directions.
	step     func(*Token)  // The next step.
	key      string        // The key consumed by the prev `step_AcceptKey`.
	keysDone []string      // List of keys we've completed already (repeats aren't wanted).
}

func NewAtlasDecoderMachine(into reflect.Value, atl *Atlas) *atlasDecoderMachine {
	// TODO this return type should prob have some interface that covers it sufficiently.
	dm := &atlasDecoderMachine{
		atl: atl,
	}
	dm.Reset(into)
	return dm
}

func (dm *atlasDecoderMachine) Reset(into reflect.Value) {
	dm.val = into
	dm.step = dm.step_Initial
	dm.key = ""
	dm.keysDone = dm.keysDone[0:0]
}

func (dm *atlasDecoderMachine) step_Initial(tok *Token) {
	switch *tok {
	case Token_MapOpen:
		// Great.  Consumed.
		dm.step = dm.step_AcceptKey
	case Token_ArrOpen:
		panic("unexpected arrOpen; expected start of struct")
	case Token_MapClose:
		panic("unexpected mapClose; expected start of struct")
	case Token_ArrClose:
		panic("unexpected arrClose; expected start of struct")
	default:
		panic(fmt.Errorf("unexpected literal of type %T; expected start of struct", *tok))
	}
}

func (dm *atlasDecoderMachine) step_AcceptKey(tok *Token) {
	switch *tok {
	case Token_MapOpen:
		panic("unexpected mapOpen; expected map key")
	case Token_ArrOpen:
		panic("unexpected arrOpen; expected map key")
	case Token_MapClose:
		dm.handleEnd()
	case Token_ArrClose:
		panic("unexpected arrClose; expected map key")
	}
	switch k := (*tok).(type) {
	case *string:
		dm.key = *k
		dm.mustAcceptKey(*k)
		//dm.step = dm.step_AcceptValue
		// actually we might wanna just push up our plea now --
		//  this saves us from having to see and forward the token at all,
		//  and makes the pattern of fab-var-filler, ret step func(token) consistent.
		//  if you *really* wanted to implement a breakout for known prims, you could still do that branch here.
		//  HANG ON, nope: keep it in the value step and keep the tok passdown.
		//   do it for parity with arrays, which must have that step
		//   and accept that token during it so they can check for end there.
		/*
			driver.Fill(
				tok, // still meant for next person and the real step is to come; we just had to figure out types, here.
				dm.Addr(dm.key),
				dm.step_postValue(), // driver returns to us after the value is done by calling this.
				    // may actually be that we stash that stepfunc, and give driver more general self pointer and Resume func in interface.
			)
		*/
	default:
		panic(fmt.Errorf("unexpected literal of type %T; expected start of struct", *tok))
	}
}
func (dm *atlasDecoderMachine) mustAcceptKey(k string) {
	for _, x := range dm.keysDone {
		if x == k {
			panic(fmt.Errorf("repeated key %q", k))
		}
	}
	dm.keysDone = append(dm.keysDone, k)
}
func (dm *atlasDecoderMachine) addr(k string) interface{} {
	_ = dm.atl
	return nil // TODO
	// n.b. this is one of the spots where i can't decide if &thing or reflect.Value is better
	// but either way we may want to define a `Slot` type alias to make it readable
}

func (dm *atlasDecoderMachine) step_AcceptValue(tok *Token) {
}

func (dm *atlasDecoderMachine) handleEnd() {
	// TODO check for all filled, etc.  then set terminal states.
}
