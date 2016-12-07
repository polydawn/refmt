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

func (m *MarshalMachineSliceWildcard) Step(d *MarshalDriver, tok *Token) (done bool, err error) {
	if m.index < 0 {
		if m.target_rv.IsNil() {
			*tok = nil
			return true, nil
		}
	}
	return m.MarshalMachineArrayWildcard.Step(d, tok)
}

type MarshalMachineArrayWildcard struct {
	target_rv reflect.Value
	started   bool
	index     int
	length    int
}

func (m *MarshalMachineArrayWildcard) Reset(valp interface{}) {
	m.target_rv = reflect.ValueOf(*(valp).(*interface{}))
	// TODO we should choose an encoder machine here!  and this is why drivers must allow that dictation.
	// Primitives should be iterable in the step here directly, too.
	// Somewhat bizarrely, this means... reset should take a driver?!  either that or it has to go in the first step.
	m.index = -1
	m.length = m.target_rv.Len()
}

func (m *MarshalMachineArrayWildcard) Step(d *MarshalDriver, tok *Token) (done bool, err error) {
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
	d.Recurse(tok, m.target_rv.Index(m.index))
	m.index++
	return false, nil
}
