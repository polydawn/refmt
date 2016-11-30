package obj

import (
	"fmt"

	. "github.com/polydawn/go-xlate/tok"
)

type UnmarshalMachineMapStringWildcard struct {
	target map[string]interface{}
	step   UnmarshalMachine
	key    string      // The key consumed by the prev `step_AcceptKey`.
	tmp    interface{} // A slot to we hand out as a ref to fill during recursions.
}

func (m *UnmarshalMachineMapStringWildcard) Reset(target interface{}) {
	m.target = target.(map[string]interface{})
	m.step = m.step_Initial
	m.key = ""
}

func (m *UnmarshalMachineMapStringWildcard) Step(vr *UnmarshalDriver, tok *Token) (done bool, err error) {
	return m.step(vr, tok)
}

func (m *UnmarshalMachineMapStringWildcard) step_Initial(_ *UnmarshalDriver, tok *Token) (done bool, err error) {
	// If it's a special state, start an object.
	//  (Or, blow up if its a special state that's silly).
	switch *tok {
	case Token_MapOpen:
		// Great.  Consumed.
		m.step = m.step_AcceptKey
		return false, nil
	case Token_ArrOpen:
		return true, fmt.Errorf("unexpected arrOpen; expected start of map")
	case Token_MapClose:
		return true, fmt.Errorf("unexpected mapClose; expected start of map")
	case Token_ArrClose:
		return true, fmt.Errorf("unexpected arrClose; expected start of map")
	default:
		return true, fmt.Errorf("unexpected literal of type %T; expected start of map", *tok)
	}
}

func (m *UnmarshalMachineMapStringWildcard) step_AcceptKey(_ *UnmarshalDriver, tok *Token) (done bool, err error) {
	// First, save any refs from the last value.
	//  (This is fiddly: the delay comes mostly from the handling of slices, which may end up re-allocating
	//   themselves during their decoding.)
	if m.key != "" {
		m.target[m.key] = m.tmp
	}
	// Now switch on tokens.
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
	case string:
		if err = m.mustAcceptKey(k); err != nil {
			return true, err
		}
		m.key = k
		m.step = m.step_AcceptValue
		return false, nil
	default:
		return true, fmt.Errorf("unexpected literal of type %T; expected key string or end of map", *tok)
	}
}

func (m *UnmarshalMachineMapStringWildcard) mustAcceptKey(k string) error {
	if _, exists := m.target[k]; exists {
		return fmt.Errorf("repeated key %q", k)
	}
	return nil
}

func (m *UnmarshalMachineMapStringWildcard) step_AcceptValue(driver *UnmarshalDriver, tok *Token) (done bool, err error) {
	m.step = m.step_AcceptKey
	driver.Recurse(tok, &m.tmp)
	return false, nil
}
