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

func (m *MarshalMachineSliceWildcard) Step(d *MarshalDriver, s *marshalSlab, tok *Token) (done bool, err error) {
	if m.index < 0 {
		if m.target_rv.IsNil() {
			tok.Type = TNull
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

func (m *MarshalMachineArrayWildcard) Reset(s *marshalSlab, valp interface{}) error {
	m.target_rv = reflect.ValueOf(valp).Elem()
	m.valueMach = s.mustPickMarshalMachineByType(m.target_rv.Type().Elem())
	m.index = -1
	m.length = m.target_rv.Len()
	return nil
}

func (m *MarshalMachineArrayWildcard) Step(d *MarshalDriver, s *marshalSlab, tok *Token) (done bool, err error) {
	if m.index < 0 {
		tok.Type = TArrOpen
		tok.Length = m.target_rv.Len()
		m.index++
		return false, nil
	}
	if m.index == m.length {
		tok.Type = TArrClose
		m.index++
		s.release()
		return true, nil
	}
	if m.index > m.length {
		return true, fmt.Errorf("invalid state: value already consumed")
	}
	d.Recurse(tok, m.target_rv.Index(m.index).Addr().Interface(), m.valueMach)
	m.index++
	return false, nil
}
