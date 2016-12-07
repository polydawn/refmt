package obj

import (
	"fmt"
	"reflect"

	. "github.com/polydawn/go-xlate/tok"
)

// Encodes a slice.
// This machine just wraps the array machine, checking to make sure the value isn't nil.
type MarshalMachineSliceWildcard struct {
	MarshalMachineArrayWildcard
}

func (m *MarshalMachineSliceWildcard) Step(d *MarshalDriver, s *Suite, tok *Token) (done bool, err error) {
	if m.index < 0 {
		if m.target_rv.IsNil() {
			*tok = nil
			return true, nil
		}
	}
	return m.MarshalMachineArrayWildcard.Step(d, s, tok)
}

type MarshalMachineArrayWildcard struct {
	target_rv reflect.Value
	valueMach MarshalMachine
	started   bool
	index     int
	length    int
}

func (m *MarshalMachineArrayWildcard) Reset(s *Suite, valp interface{}) error {
	m.target_rv = reflect.ValueOf(*(valp).(*interface{}))
	m.valueMach = s.marshalMachineForType(m.target_rv.Type().Elem())
	m.index = -1
	m.length = m.target_rv.Len()
	return nil
}

func (m *MarshalMachineArrayWildcard) Step(d *MarshalDriver, s *Suite, tok *Token) (done bool, err error) {
	if m.index < 0 {
		*tok = Token_ArrOpen
		m.index++
		return false, nil
	}
	if m.index == m.length {
		*tok = Token_ArrClose
		m.index++
		return true, nil
	}
	if m.index > m.length {
		return true, fmt.Errorf("invalid state: value already consumed")
	}
	d.Recurse(tok, m.target_rv.Index(m.index).Interface(), m.valueMach)
	m.index++
	return false, nil
}
