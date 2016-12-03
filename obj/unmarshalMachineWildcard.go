package obj

import (
	"fmt"

	. "github.com/polydawn/go-xlate/tok"
)

type UnmarshalMachineWildcard struct {
	target   *interface{}
	delegate UnmarshalMachine // actual machine, once we've demuxed with the first token.
}

func newUnmarshalMachineWildcard(target *interface{}) UnmarshalMachine {
	m := &UnmarshalMachineWildcard{target: target}
	return m
}

func (m *UnmarshalMachineWildcard) Step(driver *UnmarshalDriver, tok *Token) (done bool, err error) {
	if m.delegate == nil {
		return m.step_demux(driver, tok)
	}
	return m.delegate.Step(driver, tok)
}

func (m *UnmarshalMachineWildcard) step_demux(driver *UnmarshalDriver, tok *Token) (done bool, err error) {
	// If it's a special state, start an object.
	//  (Or, blow up if its a special state that's silly).
	switch *tok {
	case Token_MapOpen:
		// Fill in our wildcard ref with a blank map,
		//  and make a new machine for it; hand off everything.
		mp := make(map[string]interface{})
		*(m.target) = mp
		dec := &UnmarshalMachineMapStringWildcard{}
		dec.Reset(mp)
		m.delegate = dec
		return m.delegate.Step(driver, tok)

	case Token_ArrOpen:
		// Similar to maps, but a step more complex: we make a new slot for a *pointer*
		//  to a slice, because slices get new addresses when they grow (whereas by
		//   comparison, maps hide their internal growth).
		dec := &UnmarshalMachineSliceWildcard{}
		// *(m.target) = someSlice // No such step!  Array machine does this at end.
		dec.Reset(m.target)
		m.delegate = dec
		return m.delegate.Step(driver, tok)

	case Token_MapClose:
		return true, fmt.Errorf("unexpected mapClose; expected start of value")

	case Token_ArrClose:
		return true, fmt.Errorf("unexpected arrClose; expected start of value")

	default:
		// If it wasn't the start of composite, shell out to the machine for literals.
		// Don't bother to replace our internal step func because literal machines are never multi-call.
		return UnmarshalMachineLiteral{m.target}.Step(driver, tok)
	}
}
