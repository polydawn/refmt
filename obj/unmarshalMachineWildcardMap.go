package obj

import (
	"fmt"

	. "github.com/polydawn/refmt/tok"
)

type UnmarshalMachineMapStringWildcard struct {
	target map[string]interface{}
	step   unmarshalMachineStep
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
	switch tok.Type {
	case TMapOpen:
		// Great.  Consumed.
		m.step = m.step_AcceptKey
		return false, nil
	case TArrOpen:
		return true, fmt.Errorf("unexpected arrOpen; expected start of map")
	case TMapClose:
		return true, fmt.Errorf("unexpected mapClose; expected start of map")
	case TArrClose:
		return true, fmt.Errorf("unexpected arrClose; expected start of map")
	default:
		return true, fmt.Errorf("unexpected token %s; expected start of map", tok)
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
	switch tok.Type {
	case TMapOpen:
		return true, fmt.Errorf("unexpected mapOpen; expected map key")
	case TArrOpen:
		return true, fmt.Errorf("unexpected arrOpen; expected map key")
	case TMapClose:
		// no special checks for ends of wildcard map; no such thing as incomplete.
		return true, nil
	case TArrClose:
		return true, fmt.Errorf("unexpected arrClose; expected map key")
	}
	switch tok.Type {
	case TString:
		if err = m.mustAcceptKey(tok.Str); err != nil {
			return true, err
		}
		m.key = tok.Str
		m.step = m.step_AcceptValue
		return false, nil
	default:
		return true, fmt.Errorf("unexpected token %s; expected key string or end of map", tok)
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
