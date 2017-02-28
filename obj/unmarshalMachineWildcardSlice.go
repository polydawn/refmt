package obj

import (
	"fmt"

	. "github.com/polydawn/go-xlate/tok"
)

type UnmarshalMachineSliceWildcard struct {
	target *interface{}  // We still need this wildcard form to set into at the end.
	slice  []interface{} // This is our working reference (changes, because `append()`).
	step   unmarshalMachineStep
}

func (m *UnmarshalMachineSliceWildcard) Reset(target interface{}) {
	m.target = target.(*interface{})
	m.step = m.step_Initial
}

func (m *UnmarshalMachineSliceWildcard) Step(vr *UnmarshalDriver, tok *Token) (done bool, err error) {
	return m.step(vr, tok)
}

func (m *UnmarshalMachineSliceWildcard) step_Initial(_ *UnmarshalDriver, tok *Token) (done bool, err error) {
	// If it's a special state, start an object.
	//  (Or, blow up if its a special state that's silly).
	switch tok.Type {
	case TMapOpen:
		return true, fmt.Errorf("unexpected mapOpen; expected start of array")
	case TArrOpen:
		// Great.  Consumed.
		m.step = m.step_AcceptValue
		return false, nil
	case TMapClose:
		return true, fmt.Errorf("unexpected mapClose; expected start of array")
	case TArrClose:
		return true, fmt.Errorf("unexpected arrClose; expected start of array")
	default:
		return true, fmt.Errorf("unexpected token %s; expected start of array", tok)
	}
}

func (m *UnmarshalMachineSliceWildcard) step_AcceptValue(driver *UnmarshalDriver, tok *Token) (done bool, err error) {
	// Either form of open token are valid, but
	// - an arrClose is ours
	// - and a mapClose is clearly invalid.
	switch tok.Type {
	case TMapClose:
		// no special checks for ends of wildcard slice; no such thing as incomplete.
		return false, fmt.Errorf("unexpected mapClose; expected array value or end of array")
	case TArrClose:
		// Finishing step: push our current slice ref all the way to original target.
		*(m.target) = m.slice
		return true, nil
	}
	// Handle it the complex way.
	var v interface{}
	m.slice = append(m.slice, v)
	driver.Recurse(tok, &m.slice[len(m.slice)-1])
	return false, nil
	// Step simply remains `step_AcceptValue` -- arrays don't have much state machine.
}
